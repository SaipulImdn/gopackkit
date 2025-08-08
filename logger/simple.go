package logger

import (
	"fmt"
	"log"
	"strings"
)

// simpleLogger is a basic logger implementation using Go's standard log package
type simpleLogger struct {
	logger *log.Logger
	level  int
	fields map[string]interface{}
}

const (
	levelDebug = iota
	levelInfo
	levelWarn
	levelError
	levelFatal
)

// newSimpleLogger creates a new simple logger
func newSimpleLogger(config Config) Logger {
	writer := getWriter(config)
	level := parseSimpleLevel(config.Level)

	var logger *log.Logger
	if config.Format == "json" {
		logger = log.New(writer, "", 0)
	} else {
		logger = log.New(writer, "", log.LstdFlags)
	}

	return &simpleLogger{
		logger: logger,
		level:  level,
		fields: make(map[string]interface{}),
	}
}

func (s *simpleLogger) Debug(msg string, fields ...interface{}) {
	if s.level <= levelDebug {
		s.log("DEBUG", msg, fields...)
	}
}

func (s *simpleLogger) Info(msg string, fields ...interface{}) {
	if s.level <= levelInfo {
		s.log("INFO", msg, fields...)
	}
}

func (s *simpleLogger) Warn(msg string, fields ...interface{}) {
	if s.level <= levelWarn {
		s.log("WARN", msg, fields...)
	}
}

func (s *simpleLogger) Error(msg string, fields ...interface{}) {
	if s.level <= levelError {
		s.log("ERROR", msg, fields...)
	}
}

func (s *simpleLogger) Fatal(msg string, fields ...interface{}) {
	s.log("FATAL", msg, fields...)
	log.Fatal(msg)
}

func (s *simpleLogger) WithField(key string, value interface{}) Logger {
	newFields := make(map[string]interface{})
	for k, v := range s.fields {
		newFields[k] = v
	}
	newFields[key] = value

	return &simpleLogger{
		logger: s.logger,
		level:  s.level,
		fields: newFields,
	}
}

func (s *simpleLogger) WithFields(fields map[string]interface{}) Logger {
	newFields := make(map[string]interface{})
	for k, v := range s.fields {
		newFields[k] = v
	}
	for k, v := range fields {
		newFields[k] = v
	}

	return &simpleLogger{
		logger: s.logger,
		level:  s.level,
		fields: newFields,
	}
}

func (s *simpleLogger) log(level, msg string, fields ...interface{}) {
	// Combine existing fields with new fields
	allFields := make(map[string]interface{})
	for k, v := range s.fields {
		allFields[k] = v
	}

	// Add new fields from parameters
	for i := 0; i < len(fields)-1; i += 2 {
		if key, ok := fields[i].(string); ok && i+1 < len(fields) {
			allFields[key] = fields[i+1]
		}
	}

	// Build log message
	var fieldsStr string
	if len(allFields) > 0 {
		var parts []string
		for k, v := range allFields {
			parts = append(parts, fmt.Sprintf("%s=%v", k, v))
		}
		fieldsStr = " " + strings.Join(parts, " ")
	}

	logMsg := fmt.Sprintf("[%s] %s%s", level, msg, fieldsStr)
	s.logger.Output(3, logMsg)
}

// parseSimpleLevel converts string level to simple logger level
func parseSimpleLevel(level string) int {
	switch strings.ToLower(level) {
	case "debug":
		return levelDebug
	case "info":
		return levelInfo
	case "warn", "warning":
		return levelWarn
	case "error":
		return levelError
	case "fatal":
		return levelFatal
	default:
		return levelInfo
	}
}
