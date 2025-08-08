package logger

import (
	"io"
	"os"
)

// Logger interface defines common logging methods
type Logger interface {
	Debug(msg string, fields ...interface{})
	Info(msg string, fields ...interface{})
	Warn(msg string, fields ...interface{})
	Error(msg string, fields ...interface{})
	Fatal(msg string, fields ...interface{})
	WithField(key string, value interface{}) Logger
	WithFields(fields map[string]interface{}) Logger
}

// Config holds logger configuration
type Config struct {
	Level    string `json:"level" yaml:"level" env:"LOG_LEVEL" default:"info"`
	Format   string `json:"format" yaml:"format" env:"LOG_FORMAT" default:"text"`
	Output   string `json:"output" yaml:"output" env:"LOG_OUTPUT" default:"stdout"`
	Backend  string `json:"backend" yaml:"backend" env:"LOG_BACKEND" default:"logrus"`
	Filename string `json:"filename" yaml:"filename" env:"LOG_FILENAME"`
}

// LogLevel represents log levels
type LogLevel string

const (
	DebugLevel LogLevel = "debug"
	InfoLevel  LogLevel = "info"
	WarnLevel  LogLevel = "warn"
	ErrorLevel LogLevel = "error"
	FatalLevel LogLevel = "fatal"
)

// LogFormat represents log output formats
type LogFormat string

const (
	TextFormat LogFormat = "text"
	JSONFormat LogFormat = "json"
)

// LogBackend represents logging backends
type LogBackend string

const (
	LogrusBackend LogBackend = "logrus"
	ZapBackend    LogBackend = "zap"
)

// New creates a logger with default configuration
func New() Logger {
	return NewWithConfig(Config{
		Level:   "info",
		Format:  "text",
		Output:  "stdout",
		Backend: "logrus",
	})
}

// NewWithConfig creates a logger with custom configuration
func NewWithConfig(config Config) Logger {
	switch LogBackend(config.Backend) {
	case ZapBackend:
		return newZapLogger(config)
	case LogrusBackend:
		return newLogrusLogger(config)
	default:
		return newSimpleLogger(config)
	}
}

// getWriter returns the appropriate writer based on output configuration
func getWriter(config Config) io.Writer {
	switch config.Output {
	case "stderr":
		return os.Stderr
	case "file":
		if config.Filename != "" {
			file, err := os.OpenFile(config.Filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
			if err != nil {
				// Fallback to stdout if file can't be opened
				return os.Stdout
			}
			return file
		}
		return os.Stdout
	default:
		return os.Stdout
	}
}
