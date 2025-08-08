package logger

import (
	"strings"

	"github.com/sirupsen/logrus"
)

// logrusLogger wraps logrus.Logger to implement our Logger interface
type logrusLogger struct {
	logger *logrus.Logger
	entry  *logrus.Entry
}

// newLogrusLogger creates a new logrus-based logger
func newLogrusLogger(config Config) Logger {
	logger := logrus.New()

	// Set level
	logger.SetLevel(parseLogrusLevel(config.Level))

	// Set formatter
	if config.Format == string(JSONFormat) {
		logger.SetFormatter(&logrus.JSONFormatter{})
	} else {
		logger.SetFormatter(&logrus.TextFormatter{
			FullTimestamp: true,
		})
	}

	// Set output
	logger.SetOutput(getWriter(config))

	return &logrusLogger{
		logger: logger,
		entry:  logger.WithFields(logrus.Fields{}),
	}
}

func (l *logrusLogger) Debug(msg string, fields ...interface{}) {
	l.logWithFields(l.entry.Debug, msg, fields...)
}

func (l *logrusLogger) Info(msg string, fields ...interface{}) {
	l.logWithFields(l.entry.Info, msg, fields...)
}

func (l *logrusLogger) Warn(msg string, fields ...interface{}) {
	l.logWithFields(l.entry.Warn, msg, fields...)
}

func (l *logrusLogger) Error(msg string, fields ...interface{}) {
	l.logWithFields(l.entry.Error, msg, fields...)
}

func (l *logrusLogger) Fatal(msg string, fields ...interface{}) {
	l.logWithFields(l.entry.Fatal, msg, fields...)
}

func (l *logrusLogger) WithField(key string, value interface{}) Logger {
	return &logrusLogger{
		logger: l.logger,
		entry:  l.entry.WithField(key, value),
	}
}

func (l *logrusLogger) WithFields(fields map[string]interface{}) Logger {
	return &logrusLogger{
		logger: l.logger,
		entry:  l.entry.WithFields(fields),
	}
}

// logWithFields handles the key-value pairs and logs the message
func (l *logrusLogger) logWithFields(logFunc func(...interface{}), msg string, fields ...interface{}) {
	if len(fields) == 0 {
		logFunc(msg)
		return
	}

	// Convert fields to logrus.Fields
	logrusFields := make(logrus.Fields)
	for i := 0; i < len(fields)-1; i += 2 {
		if key, ok := fields[i].(string); ok && i+1 < len(fields) {
			logrusFields[key] = fields[i+1]
		}
	}

	if len(logrusFields) > 0 {
		l.entry.WithFields(logrusFields).Log(l.logger.Level, msg)
	} else {
		logFunc(msg)
	}
}

// parseLogrusLevel converts string level to logrus level
func parseLogrusLevel(level string) logrus.Level {
	switch strings.ToLower(level) {
	case string(DebugLevel):
		return logrus.DebugLevel
	case string(InfoLevel):
		return logrus.InfoLevel
	case "warn", "warning":
		return logrus.WarnLevel
	case string(ErrorLevel):
		return logrus.ErrorLevel
	case string(FatalLevel):
		return logrus.FatalLevel
	default:
		return logrus.InfoLevel
	}
}
