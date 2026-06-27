package logger

import (
	"fmt"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	log     *zap.Logger
	Sugared *zap.SugaredLogger
)

type Config struct {
	Level      string // debug, info, warn, error
	Format     string // json, text
	FilePath   string
	MaxSize    int
	MaxBackups int
	MaxAge     int
	Compress   bool
	Console    bool
}

func Init(cfg *Config) error {
	var level zapcore.Level
	switch cfg.Level {
	case "debug":
		level = zapcore.DebugLevel
	case "info":
		level = zapcore.InfoLevel
	case "warn":
		level = zapcore.WarnLevel
	case "error":
		level = zapcore.ErrorLevel
	default:
		level = zapcore.InfoLevel
	}

	var encoderConfig zapcore.EncoderConfig
	if cfg.Format == "text" {
		encoderConfig = zapcore.EncoderConfig{
			TimeKey:        "timestamp",
			LevelKey:       "level",
			NameKey:        "logger",
			MessageKey:     "message",
			StacktraceKey:  "stacktrace",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.CapitalColorLevelEncoder,
			EncodeTime:     zapcore.ISO8601TimeEncoder,
			EncodeDuration: zapcore.StringDurationEncoder,
		}
	} else {
		encoderConfig = zapcore.EncoderConfig{
			TimeKey:        "timestamp",
			LevelKey:       "level",
			NameKey:        "logger",
			MessageKey:     "message",
			StacktraceKey:  "stacktrace",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.LowercaseLevelEncoder,
			EncodeTime:     zapcore.ISO8601TimeEncoder,
			EncodeDuration: zapcore.StringDurationEncoder,
		}
	}

	var encoder zapcore.Encoder
	if cfg.Format == "text" {
		encoder = zapcore.NewConsoleEncoder(encoderConfig)
	} else {
		encoder = zapcore.NewJSONEncoder(encoderConfig)
	}

	var cores []zapcore.Core

	if cfg.Console {
		cores = append(cores, zapcore.NewCore(encoder, zapcore.AddSync(os.Stdout), level))
	}

	if cfg.FilePath != "" {
		fileWriter := zapcore.AddSync(&lumberjack.Logger{
			Filename:   cfg.FilePath,
			MaxSize:    cfg.MaxSize,
			MaxBackups: cfg.MaxBackups,
			MaxAge:     cfg.MaxAge,
			Compress:   cfg.Compress,
		})
		cores = append(cores, zapcore.NewCore(encoder, fileWriter, level))
	}

	var core zapcore.Core
	if len(cores) == 0 {
		core = zapcore.NewCore(encoder, zapcore.AddSync(os.Stdout), level)
	} else if len(cores) == 1 {
		core = cores[0]
	} else {
		core = zapcore.NewTee(cores...)
	}

	log = zap.New(core, zap.AddStacktrace(zapcore.ErrorLevel))
	Sugared = log.Sugar()

	return nil
}

func newFileRotator(cfg *Config) *os.File {
	return os.Stdout
}

func isZapFields(args []interface{}) bool {
	for _, arg := range args {
		if _, ok := arg.(zap.Field); !ok {
			return false
		}
	}
	return true
}

func argsToFields(args []interface{}) []zap.Field {
	fields := make([]zap.Field, 0, len(args))
	for _, arg := range args {
		if f, ok := arg.(zap.Field); ok {
			fields = append(fields, f)
		}
	}
	return fields
}

func Debug(msg string, args ...interface{}) {
	if len(args) == 0 {
		log.Debug(msg)
		return
	}
	if isZapFields(args) {
		log.Debug(msg, argsToFields(args)...)
	} else {
		log.Debug(fmt.Sprintf(msg, args...))
	}
}

func Info(msg string, args ...interface{}) {
	if len(args) == 0 {
		log.Info(msg)
		return
	}
	if isZapFields(args) {
		log.Info(msg, argsToFields(args)...)
	} else {
		log.Info(fmt.Sprintf(msg, args...))
	}
}

func Warn(msg string, args ...interface{}) {
	if len(args) == 0 {
		log.Warn(msg)
		return
	}
	if isZapFields(args) {
		log.Warn(msg, argsToFields(args)...)
	} else {
		log.Warn(fmt.Sprintf(msg, args...))
	}
}

func Error(msg string, args ...interface{}) {
	if len(args) == 0 {
		log.Error(msg)
		return
	}
	if isZapFields(args) {
		log.Error(msg, argsToFields(args)...)
	} else {
		log.Error(fmt.Sprintf(msg, args...))
	}
}

func Fatal(msg string, args ...interface{}) {
	if len(args) == 0 {
		log.Fatal(msg)
		return
	}
	if isZapFields(args) {
		log.Fatal(msg, argsToFields(args)...)
	} else {
		log.Fatal(fmt.Sprintf(msg, args...))
	}
}

func Debugf(template string, args ...interface{}) {
	Sugared.Debugf(template, args...)
}

func Infof(template string, args ...interface{}) {
	Sugared.Infof(template, args...)
}

func Warnf(template string, args ...interface{}) {
	Sugared.Warnf(template, args...)
}

func Errorf(template string, args ...interface{}) {
	Sugared.Errorf(template, args...)
}

func Fatalf(template string, args ...interface{}) {
	Sugared.Fatalf(template, args...)
}

func Sync() {
	if log != nil {
		log.Sync()
	}
}
