package logger

import (
	"database/internal/database/config"
	"fmt"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

//go:generate go run go.uber.org/mock/mockgen -package logger -destination mock.go -source log.go Logger
type Logger interface {
	Infof(template string, args ...any)
	Errorf(template string, args ...any)
	Warnf(template string, args ...any)
	Debug(msg string, fields ...zap.Field)
}

type Log struct {
	base *zap.SugaredLogger
	zap  *zap.Logger
}

func (l Log) Debug(msg string, fields ...zap.Field) {
	l.zap.Debug(msg, fields...)
}

func (l Log) Errorf(template string, args ...any) {
	l.base.Errorf(template, args...)
}

func (l Log) Warnf(template string, args ...any) {
	l.base.Warnf(template, args...)
}

func (l Log) Infof(template string, args ...any) {
	l.base.Infof(template, args...)
}

func NewLogger(conf *config.Config) (Logger, error) {
	if conf == nil {
		return nil, fmt.Errorf("config should be initialized")
	}

	zapConfig := zap.NewDevelopmentConfig()

	zapConfig.Level = zap.NewAtomicLevelAt(
		mapConfigLevel(conf.Logging.Level),
	)

	l, err := zapConfig.Build()
	if err != nil {
		return nil, fmt.Errorf("[logger] failed to build logger: %w", err)
	}

	return &Log{zap: l, base: l.Sugar()}, nil
}

func mapConfigLevel(level string) zapcore.Level {
	switch level {
	case config.DefaultLogLevel:
		return zapcore.WarnLevel
	case config.LogLevelDev:
		return zapcore.DebugLevel
	default:
		return zapcore.WarnLevel
	}
}
