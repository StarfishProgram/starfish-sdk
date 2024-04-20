package starfish_sdk

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// LogConfig 日志配置
type LogConfig struct {
	Level string `toml:"level"` // 级别 : debug, info, warn, error
}

// ILog Log
type ILog interface {
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
	// Ins 获取zap实例
	Ins() *zap.Logger
	// AddCallerSkip 跳过栈帧
	AddCallerSkip(skip int) ILog
}

type log struct{ ins *zap.Logger }

func (l *log) Debug(args ...interface{}) {
	l.ins.Sugar().Debugln(args...)
}

func (l *log) Info(args ...interface{}) {
	l.ins.Sugar().Infoln(args...)
}

func (l *log) Warn(args ...interface{}) {
	l.ins.Sugar().Warnln(args...)
}

func (l *log) Error(args ...interface{}) {
	l.ins.Sugar().Errorln(args...)
}

func (l *log) Panic(args ...interface{}) {
	l.ins.Sugar().Panicln(args...)
}

func (l *log) Debugf(format string, args ...interface{}) {
	l.ins.Sugar().Debugf(format, args...)
}

func (l *log) Infof(format string, args ...interface{}) {
	l.ins.Sugar().Infof(format, args...)
}

func (l *log) Warnf(format string, args ...interface{}) {
	l.ins.Sugar().Warnf(format, args...)
}

func (l *log) Errorf(format string, args ...interface{}) {
	l.ins.Sugar().Errorf(format, args...)
}

func (l *log) Panicf(format string, args ...interface{}) {
	l.ins.Sugar().Panicf(format, args...)
}

func (l *log) Ins() *zap.Logger {
	return l.ins
}

func (l *log) AddCallerSkip(skip int) ILog {
	_ins := l.ins.WithOptions(zap.AddCallerSkip(skip))
	return &log{_ins}
}

var logIns ILog

// Log 获取日志
func Log() ILog {
	return logIns
}

// InitLog 初始化日志
func InitLog(config *LogConfig) ILog {
	level, err := zapcore.ParseLevel(config.Level)
	if err != nil {
		level = zapcore.DebugLevel
	}
	atom := zap.NewAtomicLevelAt(level)
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.TimeKey = "time"
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	logger := zap.New(
		zapcore.NewCore(
			zapcore.NewConsoleEncoder(encoderConfig),
			zapcore.Lock(os.Stdout),
			atom,
		),
		zap.AddCaller(),
		zap.AddCallerSkip(1),
	)
	ins := &log{logger}
	logIns = ins
	return ins
}
