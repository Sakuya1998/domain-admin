package logger

import (
	"domain-admin/pkg/config"
	"os"
	"path/filepath"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var Log *zap.Logger

func InitLogger(cfg config.LogConfig) {
	level := zap.InfoLevel
	if l, err := zap.ParseAtomicLevel(cfg.Level); err == nil {
		level = l.Level()
	}

	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		FunctionKey:    zapcore.OmitKey,
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	if cfg.Development {
		encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	}

	var encoder zapcore.Encoder
	if cfg.Format == "console" {
		encoder = zapcore.NewConsoleEncoder(encoderConfig)
	} else {
		encoder = zapcore.NewJSONEncoder(encoderConfig)
	}

	var cores []zapcore.Core

	if cfg.Output == "file" || cfg.Output == "both" {
		dir := filepath.Dir(cfg.FilePath)
		_ = os.MkdirAll(dir, 0755)

		writeSyncer := zapcore.AddSync(&lumberjack.Logger{
			Filename:   cfg.FilePath,
			MaxSize:    cfg.MaxSize,
			MaxBackups: cfg.MaxBackups,
			MaxAge:     cfg.MaxAge,
			Compress:   cfg.Compress,
		})
		cores = append(cores, zapcore.NewCore(encoder, writeSyncer, level))
	}

	if cfg.Output == "console" || cfg.Output == "both" {
		cores = append(cores, zapcore.NewCore(encoder, zapcore.AddSync(os.Stdout), level))
	}

	core := zapcore.NewTee(cores...)
	options := []zap.Option{zap.AddCaller(), zap.AddCallerSkip(1)}
	if cfg.Development {
		options = append(options, zap.Development())
	}

	Log = zap.New(core, options...)
}

// Sync ensures all buffered logs are flushed
func Sync() {
	if Log != nil {
		_ = Log.Sync()
	}
}

// ---- 结构化日志，接受 zap.Field ----

func Debug(msg string, fields ...zap.Field) {
	Log.Debug(msg, fields...)
}

func Info(msg string, fields ...zap.Field) {
	Log.Info(msg, fields...)
}

func Warn(msg string, fields ...zap.Field) {
	Log.Warn(msg, fields...)
}

func Error(msg string, fields ...zap.Field) {
	Log.Error(msg, fields...)
}

func Fatal(msg string, fields ...zap.Field) {
	Log.Fatal(msg, fields...)
}

// ---- 结构化日志，接受 map[string]interface{}，无需外部引用 zap ----

func DebugKV(msg string, kvs map[string]interface{}) {
	Log.Debug(msg, mapToFields(kvs)...)
}

func InfoKV(msg string, kvs map[string]interface{}) {
	Log.Info(msg, mapToFields(kvs)...)
}

func WarnKV(msg string, kvs map[string]interface{}) {
	Log.Warn(msg, mapToFields(kvs)...)
}

func ErrorKV(msg string, kvs map[string]interface{}) {
	Log.Error(msg, mapToFields(kvs)...)
}

func FatalKV(msg string, kvs map[string]interface{}) {
	Log.Fatal(msg, mapToFields(kvs)...)
}

func mapToFields(kvs map[string]interface{}) []zap.Field {
	fields := make([]zap.Field, 0, len(kvs))
	for k, v := range kvs {
		fields = append(fields, zap.Any(k, v))
	}
	return fields
}

// ---- 格式化日志，类似 fmt.Printf ----

func Debugf(format string, args ...interface{}) {
	Log.Sugar().Debugf(format, args...)
}

func Infof(format string, args ...interface{}) {
	Log.Sugar().Infof(format, args...)
}

func Warnf(format string, args ...interface{}) {
	Log.Sugar().Warnf(format, args...)
}

func Errorf(format string, args ...interface{}) {
	Log.Sugar().Errorf(format, args...)
}

func Fatalf(format string, args ...interface{}) {
	Log.Sugar().Fatalf(format, args...)
}

// ---- 其他辅助 ----

// WithOptions clones the logger with the given options
func WithOptions(opts ...zap.Option) *zap.Logger {
	return Log.WithOptions(opts...)
}

// With creates a child logger and adds structured context
func With(fields ...zap.Field) *zap.Logger {
	return Log.With(fields...)
}
