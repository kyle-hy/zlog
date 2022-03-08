package zlog

import (
	"os"
	"runtime"
	"strings"

	"go.uber.org/zap/zapcore"
)

// ext 获取函数名(同文件路径扩展名)
func ext(path string) string {
	for i := len(path) - 1; i >= 0 && !os.IsPathSeparator(path[i]); i-- {
		if path[i] == '.' {
			return path[i+1:]
		}
	}
	return ""
}

// callerEncoder will add caller to log. format is "filename:lineNum:funcName", e.g:"zaplog/zaplog_test.go:15:zaplog.TestNewLogger"
func callerEncoder(caller zapcore.EntryCaller, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(strings.Join([]string{caller.TrimmedPath(), ext(runtime.FuncForPC(caller.PC).Name())}, ":"))
}
