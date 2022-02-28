package zlog

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"net/url"
	"os"
	"path/filepath"
	"sync"
	"sync/atomic"
	"time"

	"github.com/kyle-hy/zlog/chanmgr"
	"go.uber.org/zap"
	"gopkg.in/natefinch/lumberjack.v2"
)

const (
	maxChanSize = 256 * 1024
	maxFileSize = 4 * 1024 // 4GBytes
	maxBackups  = 10
	maxAge      = 7
)

// Flusher .
type Flusher interface {
	Flush() error
}

// WriteCloseFlusher .
type WriteCloseFlusher struct {
	io.Writer
	io.Closer
	Flusher
}

// AsyncLogSink 定义一个结构体
type AsyncLogSink struct {
	closed     bool
	failCounts uint64
	chanMgr    *chanmgr.ChanMgr
	writer     *WriteCloseFlusher
	ctx        context.Context
	cancel     context.CancelFunc
	wg         sync.WaitGroup
}

// AsyncLoggerSink 定义工厂函数
func AsyncLoggerSink(url *url.URL) (sink zap.Sink, err error) {
	var writer io.WriteCloser
	filePath := getLogFilePath(&defaultOptions)

	if defaultOptions.rotate {
		writer = &lumberjack.Logger{
			Filename:   filePath,
			Compress:   true,
			LocalTime:  true,
			MaxSize:    maxFileSize,
			MaxBackups: maxBackups,
			MaxAge:     maxAge,
		}
	} else {
		err := os.MkdirAll(filepath.Dir(filePath), 0755)
		if err != nil {
			return nil, err
		}
		openFlag := os.O_CREATE | os.O_WRONLY | os.O_APPEND // 使用第三方程序rotate的话，用append模式打开，否则会形成空洞的大文件
		writer, err = os.OpenFile(filePath, openFlag, os.FileMode(0644))
		if err != nil {
			return nil, err
		}
	}

	bw := bufio.NewWriterSize(writer, defaultOptions.bufioSize)
	wc := &WriteCloseFlusher{
		Writer:  bw,
		Flusher: bw,
		Closer:  writer,
	}
	c := &AsyncLogSink{
		writer:  wc,
		chanMgr: chanmgr.NewChanMgr(256, maxChanSize/256),
	}

	c.ctx, c.cancel = context.WithCancel(context.Background())
	c.wg.Add(1)
	go func() {
		c.loop()
		c.wg.Done()
	}()

	return c, nil
}

// Sync 定义Sync方法以实现Sink接口
func (c *AsyncLogSink) Sync() error {
	c.Close()
	return nil
}

// Close 定义Close方法以实现Sink接口
func (c *AsyncLogSink) Close() error {
	time.Sleep(time.Millisecond * 2) // 短暂等待日志写入管道

	if c.closed {
		return nil
	}
	c.closed = true
	if c.cancel != nil {
		c.cancel()
	}
	c.wg.Wait() // wait until all msgs have been consumed
	c.writer.Close()
	return nil
}

// 定义Write方法以实现Sink接口
func (c *AsyncLogSink) Write(p []byte) (n int, err error) {
	// zap框架复用切片p参数,需要拷贝否则错乱
	cp := make([]byte, len(p))
	copy(cp, p)

	msgChan, _ := c.chanMgr.NextWrite()
	if !defaultOptions.overflow {
		msgChan <- cp
	} else {
		select {
		case msgChan <- cp:
		default:
			failCounts := atomic.AddUint64(&c.failCounts, 1)
			if failCounts%100 == 0 {
				msgChan <- addField(failCounts, "blockNums", cp)
			}
		}
	}
	return len(p), nil
}

func addField(failCounts uint64, name string, msg []byte) []byte {
	b := bytes.TrimSuffix(msg, []byte("}\n"))
	b = append(b, []byte(fmt.Sprintf(",\"%s\":%d}\n", name, failCounts))...)
	return b
}

func (c *AsyncLogSink) loop() {
	defer func() {
		recover()
	}()

	var msg []byte
	closed := false

	for {
		msgChan, idx := c.chanMgr.NextRead()
		if !closed {
			select {
			case msg = <-msgChan:
			case <-c.ctx.Done():
				closed = true
			}
		} else {
			select {
			case msg = <-msgChan:
			default:
				c.writer.Flush()
				return
			}
		}

		if len(msg) > 0 {
			c.writer.Write([]byte(msg))
			msg = nil
		}

		if c.chanMgr.Len(idx+1) == 0 {
			c.writer.Flush()
		}
	}
}
