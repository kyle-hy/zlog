package logger

import (
	"fmt"
	"os"

	"github.com/kyle-hy/zlog"
	"github.com/v2pro/plz/gls"
)

func logFormat(template string, fmtArgs []interface{}) string {
	// Format with Sprint, Sprintf, or neither.
	msg := template
	if msg == "" && len(fmtArgs) > 0 {
		msg = fmt.Sprint(fmtArgs...)
	} else if msg != "" && len(fmtArgs) > 0 {
		msg = fmt.Sprintf(template, fmtArgs...)
	}
	return msg
}

// Debugf logs a message at DebugLevel.
func Debugf(template string, fmtArgs ...interface{}) {
	if zlog.InnerLog() != nil {
		zlog.InnerLog().Debug(logFormat(template, fmtArgs), zlog.GoID(gls.GoID()))
	} else {
		fmt.Printf("log not init. "+template, fmtArgs)
	}
}

// Infof logs a message at InfoLevel.
func Infof(template string, fmtArgs ...interface{}) {
	if zlog.InnerLog() != nil {
		zlog.InnerLog().Info(logFormat(template, fmtArgs), zlog.GoID(gls.GoID()))
	} else {
		fmt.Printf("log not init. "+template, fmtArgs)
	}
}

// Warnf logs a message at WarnLevel.
// at the log site, as well as any fields accumulated on the logger.
func Warnf(template string, fmtArgs ...interface{}) {
	if zlog.InnerLog() != nil {
		zlog.InnerLog().Warn(logFormat(template, fmtArgs), zlog.GoID(gls.GoID()))
	} else {
		fmt.Printf("log not init. "+template, fmtArgs)
	}
}

// Errorf logs a message at ErrorLevel.
// at the log site, as well as any fields accumulated on the logger.
func Errorf(template string, fmtArgs ...interface{}) {
	if zlog.InnerLog() != nil {
		zlog.InnerLog().Error(logFormat(template, fmtArgs), zlog.GoID(gls.GoID()))
	} else {
		fmt.Printf("log not init. "+template, fmtArgs)
	}
}

// PanicAsyncf logs a message at ErrorLevel and flush to file.
// The logger then closed and panics, even if logging at PanicLevel is disabled.
func PanicAsyncf(template string, fmtArgs ...interface{}) {
	if zlog.InnerLog() != nil {
		msg := logFormat(template, fmtArgs)
		zlog.InnerLog().Error("panic:"+msg, zlog.GoID(gls.GoID()))
		zlog.InnerLog().Sync()
		panic(msg)
	} else {
		fmt.Printf("log not init. "+template, fmtArgs)
	}
}

// FatalAsyncf logs a message at FatalLevel and flush to file.
// The logger then calls os.Exit(1), even if logging at FatalLevel is disabled.
func FatalAsyncf(template string, fmtArgs ...interface{}) {
	if zlog.InnerLog() != nil {
		zlog.InnerLog().Error("fatal:"+logFormat(template, fmtArgs), zlog.GoID(gls.GoID()))
		zlog.InnerLog().Sync()
		os.Exit(1)
	} else {
		fmt.Printf("log not init. "+template, fmtArgs)
	}
}

// DPanicf logs a message at DPanicLevel.
// If the logger is in development mode, it then panics (DPanic means
// "development panic"). This is useful for catching errors that are
// recoverable, but shouldn't ever happen.
func DPanicf(template string, fmtArgs ...interface{}) {
	if zlog.InnerLog() != nil {
		zlog.InnerLog().DPanic("panic:"+logFormat(template, fmtArgs), zlog.GoID(gls.GoID()))
	} else {
		fmt.Printf("log not init. "+template, fmtArgs)
	}
}

// Panicf logs a message at PanicLevel.
// The logger then panics, even if logging at PanicLevel is disabled.
func Panicf(template string, fmtArgs ...interface{}) {
	if zlog.InnerLog() != nil {
		zlog.InnerLog().Panic("panic:"+logFormat(template, fmtArgs), zlog.GoID(gls.GoID()))
	} else {
		fmt.Printf("log not init. "+template, fmtArgs)
	}
}

// Fatalf logs a message at FatalLevel.
// The logger then calls os.Exit(1), even if logging at FatalLevel is disabled.
func Fatalf(template string, fmtArgs ...interface{}) {
	if zlog.InnerLog() != nil {
		zlog.InnerLog().Fatal("fatal:"+logFormat(template, fmtArgs), zlog.GoID(gls.GoID()))
	} else {
		fmt.Printf("log not init. "+template, fmtArgs)
	}
}

// Debug logs a message at DebugLevel. The message includes any fields passed
// at the log site, as well as any fields accumulated on the logger.
func Debug(msg ...interface{}) {
	if zlog.InnerLog() != nil {
		zlog.InnerLog().Debug(logFormat("", msg), zlog.GoID(gls.GoID()))
	} else {
		fmt.Println("log not init. msg:", msg)
	}
}

// Info logs a message at InfoLevel. The message includes any fields passed
// at the log site, as well as any fields accumulated on the logger.
func Info(msg ...interface{}) {
	if zlog.InnerLog() != nil {
		zlog.InnerLog().Info(logFormat("", msg), zlog.GoID(gls.GoID()))
	} else {
		fmt.Println("log not init. msg:", msg)
	}
}

// Warn logs a message at WarnLevel. The message includes any fields passed
// at the log site, as well as any fields accumulated on the logger.
func Warn(msg ...interface{}) {
	if zlog.InnerLog() != nil {
		zlog.InnerLog().Warn(logFormat("", msg), zlog.GoID(gls.GoID()))
	} else {
		fmt.Println("log not init. msg:", msg)
	}
}

// Error logs a message at ErrorLevel. The message includes any fields passed
// at the log site, as well as any fields accumulated on the logger.
func Error(msg ...interface{}) {
	if zlog.InnerLog() != nil {
		zlog.InnerLog().Error(logFormat("", msg), zlog.GoID(gls.GoID()))
	} else {
		fmt.Println("log not init. msg:", msg)
	}
}

// PanicAsync logs a message at ErrorLevel and flush to file. The message includes any fields passed
// at the log site, as well as any fields accumulated on the logger.
// The logger then closed and panics, even if logging at PanicLevel is disabled.
func PanicAsync(msg ...interface{}) {
	if zlog.InnerLog() != nil {
		zlog.InnerLog().Error("panic:"+logFormat("", msg), zlog.GoID(gls.GoID()))
		zlog.InnerLog().Sync()
		panic(msg)
	} else {
		fmt.Println("log not init. msg:", msg)
	}
}

// FatalAsync logs a message at FatalLevel and flush to file. The message includes any fields passed
// at the log site, as well as any fields accumulated on the logger.
//
// The logger then calls os.Exit(1), even if logging at FatalLevel is disabled.
func FatalAsync(msg ...interface{}) {
	if zlog.InnerLog() != nil {
		zlog.InnerLog().Error("fatal:"+logFormat("", msg), zlog.GoID(gls.GoID()))
		zlog.InnerLog().Sync()
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
func DPanic(msg ...interface{}) {
	if zlog.InnerLog() != nil {
		zlog.InnerLog().DPanic(logFormat("", msg), zlog.GoID(gls.GoID()))
	} else {
		fmt.Println("log not init. msg:", msg)
	}
}

// Panic logs a message at PanicLevel. The message includes any fields passed
// at the log site, as well as any fields accumulated on the logger.
//
// The logger then panics, even if logging at PanicLevel is disabled.
func Panic(msg ...interface{}) {
	if zlog.InnerLog() != nil {
		zlog.InnerLog().Panic(logFormat("", msg), zlog.GoID(gls.GoID()))
	} else {
		fmt.Println("log not init. msg:", msg)
	}
}

// Fatal logs a message at FatalLevel. The message includes any fields passed
// at the log site, as well as any fields accumulated on the logger.
//
// The logger then calls os.Exit(1), even if logging at FatalLevel is disabled.
func Fatal(msg ...interface{}) {
	if zlog.InnerLog() != nil {
		zlog.InnerLog().Fatal(logFormat("", msg), zlog.GoID(gls.GoID()))
	} else {
		fmt.Println("log not init. msg:", msg)
	}
}
