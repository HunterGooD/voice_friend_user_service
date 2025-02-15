package logger

import "context"

type Logger interface {
	Info(message string, opt ...any)
	Debug(message string, opt ...any)
	Warn(message string, opt ...any)
	Error(message string, opt ...any)
	Log(ctx context.Context, lvl int, message string, fields ...any)
	Sync() error
}

type LogLevel int

const (
	LevelDebug LogLevel = -4
	LevelInfo  LogLevel = 0
	LevelWarn  LogLevel = 4
	LevelError LogLevel = 8
)
