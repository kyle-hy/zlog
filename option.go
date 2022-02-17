package zlog

import (
	"fmt"
	"os"
	"path"
)

const (
	defaultLogPath = "./log/%s/%s.log"
)

// Options 属性
type Options struct {
	testEnv   bool                   // 测试环境日志级别为debug
	logPath   string                 // 日志路径
	withGID   bool                   // 打印协程id
	stdout    bool                   // 日志同时打印到标准输出
	overflow  bool                   // 日志缓存管道溢出则丢弃日志
	rotate    bool                   // 是否使用lumberjack滚动日志
	bufioSize int                    // 写文件io的缓存大小
	fields    map[string]interface{} // 日志默认附加的字段
}

var defaultOptions = Options{
	testEnv:   true,
	withGID:   false,
	overflow:  false,
	rotate:    true,
	bufioSize: 1024 * 8,
}

// 由于日志文件配套工具有相关限制，故不提供灵活的文件路径
func getLogFilePath(opt *Options) string {
	if len(opt.logPath) == 0 {
		return fmt.Sprintf(defaultLogPath, processName(), processName())
	}
	return opt.logPath
}

func processName() string {
	return path.Base(os.Args[0])

}

// Option 属性选项
type Option func(*Options)

// Overflow 设置日志缓存管道溢出后是否丢弃
func Overflow(discard bool) Option {
	return func(o *Options) {
		o.overflow = discard
	}
}

// WithGID 打印协程ID
func WithGID(withGID bool) Option {
	return func(o *Options) {
		o.withGID = withGID
	}
}

// BufioSize bufio缓存的大小, 默认1024*8
func BufioSize(bufioSize int) Option {
	return func(o *Options) {
		o.bufioSize = bufioSize
	}
}

// LogPath 日志文件路径
func LogPath(logPath string) Option {
	return func(o *Options) {
		o.logPath = logPath
	}
}

// Rotate 滚动日志
func Rotate(rotate bool) Option {
	return func(o *Options) {
		o.rotate = rotate
	}
}

// WithFields 所有日志都附带的字段
func WithFields(fields map[string]interface{}) Option {
	return func(o *Options) {
		o.fields = fields
	}
}

// Stdout 日志打印到标准输出
func Stdout(stdout bool) Option {
	return func(o *Options) {
		o.stdout = stdout
	}
}

// TestEnv 是否是测试环境
func TestEnv(testEnv bool) Option {
	return func(o *Options) {
		o.testEnv = testEnv
	}
}
