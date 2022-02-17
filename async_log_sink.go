package zlog

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"net/url"
	"sync"
	"sync/atomic"

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
	msgChans   chan []byte
	writer     *WriteCloseFlusher
	ctx        context.Context
	cancel     context.CancelFunc
	wg         sync.WaitGroup
}

// AsyncLoggerSink 定义工厂函数
func AsyncLoggerSink(url *url.URL) (sink zap.Sink, err error) {
	jack := &lumberjack.Logger{
		Filename:   getLogFilePath(&defaultLogOptions),
		Compress:   false,
		LocalTime:  true,
		MaxSize:    maxFileSize,
		MaxBackups: maxBackups,
		MaxAge:     maxAge,
	}
	bw := bufio.NewWriterSize(jack, defaultLogOptions.bufioSize)
	wc := &WriteCloseFlusher{
		Writer:  bw,
		Flusher: bw,
		Closer:  jack,
	}
	c := &AsyncLogSink{
		writer:   wc,
		msgChans: make(chan []byte, maxChanSize),
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

	if !defaultLogOptions.overflow {
		c.msgChans <- cp
	} else {
		select {
		case c.msgChans <- cp:
		default:
			failCounts := atomic.AddUint64(&c.failCounts, 1)
			if failCounts%100 == 0 {
				c.msgChans <- addField(failCounts, "blockNums", cp)
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
		if !closed {
			select {
			case msg = <-c.msgChans:
			case <-c.ctx.Done():
				closed = true
			}
		} else {
			select {
			case msg = <-c.msgChans:
			default:
				c.writer.Flush()
				return
			}
		}

		if len(msg) > 0 {
			c.writer.Write([]byte(msg))
			msg = nil
		}

		if len(c.msgChans) == 0 {
			c.writer.Flush()
		}
	}
}
