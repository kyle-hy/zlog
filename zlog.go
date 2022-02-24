package zlog

import (
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/v2pro/plz/gls"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	appInnerLog *zap.Logger
	initOnce    sync.Once
)

// GetLogger  获取 appInnerLog
func GetLogger() *zap.Logger {
	return appInnerLog
}

// InitLog 初始化日志
func InitLog(opts ...Option) error {
	closeLog()

	for _, opt := range opts {
		opt(&defaultOptions)
	}

	innerLog, err := newLogger(&defaultOptions)
	if err != nil {
		return err
	}
	appInnerLog = innerLog
	return nil
}

func closeLog() error {
	if appInnerLog != nil {
		return appInnerLog.Sync()
	}
	return nil
}

func epochFullTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006-01-02 15:04:05"))
}

// newLogger 初始化日志
func newLogger(opt *Options) (*zap.Logger, error) {
	logLevel := zap.NewAtomicLevelAt(zap.DebugLevel)
	if opt.testEnv {
		logLevel = zap.NewAtomicLevelAt(zap.InfoLevel)
	}

	// 将AsyncLoggerSink工厂函数注册到zap中, 自定义协议名为 AsyncLog
	if err := zap.RegisterSink("AsyncLog", AsyncLoggerSink); err != nil {
		fmt.Println(err)
		return nil, err
	}
	outPaths := []string{"AsyncLog://127.0.0.1"}
	if opt.stdout {
		outPaths = append(outPaths, "stdout")
	}

	var zc = zap.Config{
		Level:             logLevel,
		Development:       false,
		DisableCaller:     false,
		DisableStacktrace: false,
		Sampling:          nil,
		Encoding:          "json",
		EncoderConfig: zapcore.EncoderConfig{
			MessageKey:     "msg",
			LevelKey:       "level",
			TimeKey:        "time",
			NameKey:        "name",
			CallerKey:      "caller",
			StacktraceKey:  "stack",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.LowercaseLevelEncoder,
			EncodeTime:     epochFullTimeEncoder, // EncodeTime: zapcore.ISO8601TimeEncoder,
			EncodeDuration: zapcore.StringDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
			EncodeName:     zapcore.FullNameEncoder,
		},
		OutputPaths:      outPaths,
		ErrorOutputPaths: []string{"stderr"},
		InitialFields:    opt.fields,
	}

	return zc.Build(zap.AddCallerSkip(1))
}

// Debug logs a message at DebugLevel. The message includes any fields passed
// at the log site, as well as any fields accumulated on the logger.
func Debug(msg string, fields ...zapcore.Field) {
	if appInnerLog != nil {
		appInnerLog.Debug(msg, addGoID(fields)...)
	} else {
		fmt.Println("log not init. msg:", msg)
	}
}

// Info logs a message at InfoLevel. The message includes any fields passed
// at the log site, as well as any fields accumulated on the logger.
func Info(msg string, fields ...zapcore.Field) {
	if appInnerLog != nil {
		appInnerLog.Info(msg, addGoID(fields)...)
	} else {
		fmt.Println("log not init. msg:", msg)
	}
}

// Warn logs a message at WarnLevel. The message includes any fields passed
// at the log site, as well as any fields accumulated on the logger.
func Warn(msg string, fields ...zapcore.Field) {
	if appInnerLog != nil {
		appInnerLog.Warn(msg, addGoID(fields)...)
	} else {
		fmt.Println("log not init. msg:", msg)
	}
}

// Error logs a message at ErrorLevel. The message includes any fields passed
// at the log site, as well as any fields accumulated on the logger.
func Error(msg string, fields ...zapcore.Field) {
	if appInnerLog != nil {
		appInnerLog.Error(msg, addGoID(fields)...)
	} else {
		fmt.Println("log not init. msg:", msg)
	}
}

// PanicAsync logs a message at ErrorLevel and flush to file. The message includes any fields passed
// at the log site, as well as any fields accumulated on the logger.
// The logger then closed and panics, even if logging at PanicLevel is disabled.
func PanicAsync(msg string, fields ...zapcore.Field) {
	if appInnerLog != nil {
		appInnerLog.Error("panic:"+msg, addGoID(fields)...)
		appInnerLog.Sync()
		panic(msg)
	} else {
		fmt.Println("log not init. msg:", msg)
	}
}

// FatalAsync logs a message at FatalLevel and flush to file. The message includes any fields passed
// at the log site, as well as any fields accumulated on the logger.
//
// The logger then calls os.Exit(1), even if logging at FatalLevel is disabled.
func FatalAsync(msg string, fields ...zapcore.Field) {
	if appInnerLog != nil {
		appInnerLog.Error("fatal:"+msg, addGoID(fields)...)
		appInnerLog.Sync()
		os.Exit(1)
	} else {
		fmt.Println("log not init. msg:", msg)
	}
}

// DPanic logs a message at DPanicLevel. The message includes any fields passed
// at the log site, as well as any fields accumulated on the logger.
//
// If the logger is in development mode, it then panics (DPanic means
// "development panic"). This is useful for catching errors that are
// recoverable, but shouldn't ever happen.
func DPanic(msg string, fields ...zapcore.Field) {
	if appInnerLog != nil {
		appInnerLog.DPanic(msg, addGoID(fields)...)
	} else {
		fmt.Println("log not init. msg:", msg)
	}
}

// Panic logs a message at PanicLevel. The message includes any fields passed
// at the log site, as well as any fields accumulated on the logger.
//
// The logger then panics, even if logging at PanicLevel is disabled.
func Panic(msg string, fields ...zapcore.Field) {
	if appInnerLog != nil {
		appInnerLog.Panic(msg, addGoID(fields)...)
	} else {
		fmt.Println("log not init. msg:", msg)
	}
}

// Fatal logs a message at FatalLevel. The message includes any fields passed
// at the log site, as well as any fields accumulated on the logger.
//
// The logger then calls os.Exit(1), even if logging at FatalLevel is disabled.
func Fatal(msg string, fields ...zapcore.Field) {
	if appInnerLog != nil {
		appInnerLog.Fatal(msg, addGoID(fields)...)
	} else {
		fmt.Println("log not init. msg:", msg)
	}
}

// Sync calls the underlying syslogCore's Sync method, flushing any buffered log
// entries. Applications should take care to call Sync before exiting.
func Sync() error {
	Info("log closed")
	if appInnerLog != nil {
		return appInnerLog.Sync()
	}
	return nil
}

// LogLevelEnable returns true if the given level is at or above this level.
func LogLevelEnable(level zapcore.Level) bool {
	return appInnerLog.Core().Enabled(level)
}

func addGoID(fields []zapcore.Field) []zapcore.Field {
	if defaultOptions.withGID {
		return append(fields, GoID(gls.GoID()))
	}
	return fields
}
