package logger

import (
	"context"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/sirupsen/logrus"
)

type LogrusLogger struct {
	l *logrus.Logger
}

func NewTextLogrusLogger(w io.Writer, logLevel string) *LogrusLogger {
	log := logrus.New()
	switch strings.ToUpper(logLevel) {
	case "INFO":
		log.SetLevel(logrus.InfoLevel)
	case "DEBUG":
		log.SetLevel(logrus.DebugLevel)
	}
	log.SetOutput(w)
	log.SetFormatter(&logrus.TextFormatter{})
	return &LogrusLogger{log}
}

func NewJsonLogrusLogger(w io.Writer, logLevel string) *LogrusLogger {
	log := logrus.New()
	switch strings.ToUpper(logLevel) {
	case "INFO":
		log.SetLevel(logrus.InfoLevel)
	case "DEBUG":
		log.SetLevel(logrus.DebugLevel)
	default:
		log.SetLevel(logrus.DebugLevel)
	}
	log.SetOutput(w)
	log.SetFormatter(&logrus.JSONFormatter{})

	return &LogrusLogger{log}
}

func (log *LogrusLogger) Info(message string, opt ...any) {
	params := log.parseLogrusOpt(opt...)
	log.l.WithFields(params).Info(message)
}

func (log *LogrusLogger) Debug(message string, opt ...any) {
	params := log.parseLogrusOpt(opt...)
	log.l.WithFields(params).Debug(message)
}

func (log *LogrusLogger) Warn(message string, opt ...any) {
	params := log.parseLogrusOpt(opt...)
	log.l.WithFields(params).Warn(message)
}

func (log *LogrusLogger) Error(message string, opt ...any) {
	params := log.parseLogrusOpt(opt...)
	log.l.WithFields(params).Error(message)
}

func (log *LogrusLogger) Log(_ context.Context, lvl int, message string, fields ...any) {
	// adapter log level from grpclog(slog) using slog level log
	var logrusLevel logrus.Level
	switch LogLevel(lvl) {
	case LevelDebug:
		logrusLevel = logrus.DebugLevel
	case LevelInfo:
		logrusLevel = logrus.InfoLevel
	case LevelWarn:
		logrusLevel = logrus.WarnLevel
	case LevelError:
		logrusLevel = logrus.ErrorLevel
	}

	logrusFields := make(logrus.Fields)

	for i := 0; i < len(fields); i += 2 {
		key := fmt.Sprintf("%v", fields[i])
		value := fmt.Sprintf("%v", fields[i+1])
		logrusFields[key] = value
	}

	log.l.WithFields(logrusFields).Log(logrusLevel, message)
}

func (log *LogrusLogger) Sync() error {
	return nil
}

func (log *LogrusLogger) parseLogrusOpt(opt ...any) logrus.Fields {
	params := make(logrus.Fields)
	for k, v := range opt {
		switch val := v.(type) {
		case map[string]any:
			// TODO: if use map any maybe using this func
			// logrus.Fields(val)
			for key, value := range val {
				params[key] = value
			}
		case error:
			params["error"] = val
		default:
			params["param_"+strconv.Itoa(k)] = val
		}
	}
	return params
}
