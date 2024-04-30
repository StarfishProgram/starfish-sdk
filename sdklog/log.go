package sdklog

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Config 日志配置
type Config struct {
	Level string `toml:"level"` // 级别 : debug, info, warn, error
}

type _Log struct{ ins *zap.Logger }

func (l *_Log) Debug(args ...interface{}) {
	l.ins.Sugar().Debugln(args...)
}

func (l *_Log) Info(args ...interface{}) {
	l.ins.Sugar().Infoln(args...)
}

func (l *_Log) Warn(args ...interface{}) {
	l.ins.Sugar().Warnln(args...)
}

func (l *_Log) Error(args ...interface{}) {
	l.ins.Sugar().Errorln(args...)
}

func (l *_Log) Panic(args ...interface{}) {
	l.ins.Sugar().Panicln(args...)
}

func (l *_Log) Debugf(format string, args ...interface{}) {
	l.ins.Sugar().Debugf(format, args...)
}

func (l *_Log) Infof(format string, args ...interface{}) {
	l.ins.Sugar().Infof(format, args...)
}

func (l *_Log) Warnf(format string, args ...interface{}) {
	l.ins.Sugar().Warnf(format, args...)
}

func (l *_Log) Errorf(format string, args ...interface{}) {
	l.ins.Sugar().Errorf(format, args...)
}

func (l *_Log) Panicf(format string, args ...interface{}) {
	l.ins.Sugar().Panicf(format, args...)
}

func (l *_Log) AddCallerSkip(skip int) Log {
	_ins := l.ins.WithOptions(zap.AddCallerSkip(skip))
	return &_Log{_ins}
}

var ins Log

// Log 获取日志
func Ins() Log {
	return ins
}

// Debug 调试日志
func Debug(args ...interface{}) {
	ins.AddCallerSkip(1).Debug(args...)
}

// Info 信息日志
func Info(args ...interface{}) {
	ins.AddCallerSkip(1).Info(args...)
}

// Warn 警告日志
func Warn(args ...interface{}) {
	ins.AddCallerSkip(1).Warn(args...)
}

// Error 错误日志
func Error(args ...interface{}) {
	ins.AddCallerSkip(1).Error(args...)
}

// Panic 恐慌日志
func Panic(args ...interface{}) {
	ins.AddCallerSkip(1).Panic(args...)
}

// Debugf 调试日志
func Debugf(format string, args ...interface{}) {
	ins.AddCallerSkip(1).Debugf(format, args...)
}

// Infof 信息日志
func Infof(format string, args ...interface{}) {
	ins.AddCallerSkip(1).Infof(format, args...)
}

// Warnf 警告日志
func Warnf(format string, args ...interface{}) {
	ins.AddCallerSkip(1).Warnf(format, args...)
}

// Errorf 错误日志
func Errorf(format string, args ...interface{}) {
	ins.AddCallerSkip(1).Errorf(format, args...)
}

// Panicf 恐慌日志
func Panicf(format string, args ...interface{}) {
	ins.AddCallerSkip(1).Panicf(format, args...)
}

// AddCallerSkip 跳过栈帧
func AddCallerSkip(skip int) Log {
	return ins.AddCallerSkip(1).AddCallerSkip(skip)
}

// Init 初始化日志
func Init(config *Config) {
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
	ins = &_Log{logger}
}
