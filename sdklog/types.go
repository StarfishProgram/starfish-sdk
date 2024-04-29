package sdklog

// Log Log
type Log interface {
	// Debug 调试日志
	Debug(args ...interface{})
	// Info 信息日志
	Info(args ...interface{})
	// Warn 警告日志
	Warn(args ...interface{})
	// Error 错误日志
	Error(args ...interface{})
	// Panic 恐慌日志
	Panic(args ...interface{})
	// Debugf 调试日志
	Debugf(format string, args ...interface{})
	// Infof 信息日志
	Infof(format string, args ...interface{})
	// Warnf 警告日志
	Warnf(format string, args ...interface{})
	// Errorf 错误日志
	Errorf(format string, args ...interface{})
	// Panicf 恐慌日志
	Panicf(format string, args ...interface{})
	// AddCallerSkip 跳过栈帧
	AddCallerSkip(skip int) Log
}
