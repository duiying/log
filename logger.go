package log

import (
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
)

type Logger struct {
	*zap.Logger
}

func NewLogger(cfg *Config) *Logger {
	var allCore []zapcore.Core

	// 使用 lumberjack 对日志进行切割
	hook := lumberjack.Logger{
		Filename:   cfg.Filename,
		MaxSize:    cfg.MaxSize,
		MaxBackups: cfg.MaxBackups,
		MaxAge:     cfg.MaxAge,
		Compress:   cfg.Compress,
		LocalTime:  true,
	}

	encoderConfig := getEncoderConfig()

	ll := zap.NewAtomicLevelAt(zap.DebugLevel)

	fileWriter := zapcore.AddSync(&hook)
	allCore = append(allCore, zapcore.NewCore(zapcore.NewJSONEncoder(encoderConfig), fileWriter, ll))
	if !cfg.Online {
		c := zap.NewDevelopmentEncoderConfig()
		c.EncodeLevel = zapcore.CapitalColorLevelEncoder
		c.EncodeCaller = zapcore.ShortCallerEncoder
		allCore = append(allCore, zapcore.NewCore(zapcore.NewConsoleEncoder(c), zapcore.Lock(os.Stdout), ll))
	}

	core := zapcore.NewTee(allCore...)
	opts := []zap.Option{
		zap.AddCaller(),
		zap.AddCallerSkip(1),
		// Error 及以上级别的日志打印堆栈信息
		zap.AddStacktrace(zap.ErrorLevel),
	}

	l := zap.New(core,
		opts...,
	)

	lg := &Logger{
		Logger: l,
	}

	zap.RedirectStdLog(l)
	return lg
}

func getEncoderConfig() zapcore.EncoderConfig {
	return zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeName:     zapcore.FullNameEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
}

func (l *Logger) encodeKeyValues(keysAndValues []any) (fields []zap.Field) {
	kvLen := len(keysAndValues)
	if kvLen == 0 {
		return
	}
	if kvLen&1 == 1 {
		l.Logger.WithOptions(zap.AddCaller(), zap.AddCallerSkip(1)).Error("keysAndValues length is not even")
	}
	fields = make([]zap.Field, 0, kvLen>>1)
	kvLen--
	for i := 0; i < kvLen; i += 2 {
		key, ok := keysAndValues[i].(string)
		if !ok {
			l.Logger.WithOptions(zap.AddCaller(), zap.AddCallerSkip(1)).Error(fmt.Sprintf("key is not string: %v", keysAndValues[i]))
			continue
		}
		fields = append(fields, zap.Any(key, keysAndValues[i+1]))
	}
	return
}

func (l *Logger) With(keysAndValues ...any) *Logger {
	fs := l.encodeKeyValues(keysAndValues)
	cloneL := *l
	cloneL.Logger = l.Logger.With(fs...)
	return &cloneL
}

func (l *Logger) Debug(msg string, keysAndValues ...any) {
	l.Logger.Debug(msg, l.encodeKeyValues(keysAndValues)...)
}

func (l *Logger) Info(msg string, keysAndValues ...any) {
	l.Logger.Info(msg, l.encodeKeyValues(keysAndValues)...)
}

func (l *Logger) Warn(msg string, keysAndValues ...any) {
	l.Logger.Warn(msg, l.encodeKeyValues(keysAndValues)...)
}

func (l *Logger) Error(msg string, keysAndValues ...any) {
	l.Logger.Error(msg, l.encodeKeyValues(keysAndValues)...)
}

func (l *Logger) DPanic(msg string, keysAndValues ...any) {
	l.Logger.DPanic(msg, l.encodeKeyValues(keysAndValues)...)
}

func (l *Logger) Panic(msg string, keysAndValues ...any) {
	l.Logger.Panic(msg, l.encodeKeyValues(keysAndValues)...)
}

func (l *Logger) Fatal(msg string, keysAndValues ...any) {
	l.Logger.Fatal(msg, l.encodeKeyValues(keysAndValues)...)
}
