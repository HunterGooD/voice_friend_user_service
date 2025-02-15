package logger

import (
	"context"
	"strconv"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type ZapLogger struct {
	l *zap.Logger
}

func NewZapLogger() *ZapLogger {
	config := zap.NewProductionConfig()
	config.EncoderConfig.TimeKey = "timestamp"
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	l, _ := config.Build()
	return &ZapLogger{l}
}

func (z *ZapLogger) Info(message string, opt ...any) {
	fields := z.parseZapOpt(opt)
	z.l.Info(message, fields...)
}

func (z *ZapLogger) Debug(message string, opt ...any) {
	fields := z.parseZapOpt(opt)
	z.l.Debug(message, fields...)
}

func (z *ZapLogger) Warn(message string, opt ...any) {
	fields := z.parseZapOpt(opt)
	z.l.Warn(message, fields...)
}

func (z *ZapLogger) Error(message string, opt ...any) {
	fields := z.parseZapOpt(opt)
	z.l.Error(message, fields...)
}

func (z *ZapLogger) Sync() error {
	return z.l.Sync()
}

func (z *ZapLogger) Log(ctx context.Context, lvl int, message string, fields ...any) {
	sugar := z.l.Sugar()
	var zapLevel zapcore.Level
	switch LogLevel(lvl) {
	case LevelDebug:
		zapLevel = zapcore.DebugLevel
	case LevelInfo:
		zapLevel = zapcore.InfoLevel
	case LevelWarn:
		zapLevel = zapcore.WarnLevel
	case LevelError:
		zapLevel = zapcore.ErrorLevel
	}
	sugar.Logw(zapLevel, message, fields...)
}

func (z *ZapLogger) parseZapOpt(opt ...any) []zap.Field {
	params := make([]zap.Field, 0, len(opt))
	for k, v := range opt {
		switch val := v.(type) {
		case map[string]any:
			for key, value := range val {
				params = append(params, zap.Any(key, value))
			}
		case error:
			params = append(params, zap.Error(val))
		default:
			params = append(params, zap.Any("param_"+strconv.Itoa(k), val))
		}
	}
	return params
}
