package logger

import (
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// zapLogger wraps zap.Logger to implement our Logger interface
type zapLogger struct {
	logger *zap.Logger
	sugar  *zap.SugaredLogger
}

// newZapLogger creates a new zap-based logger
func newZapLogger(config Config) Logger {
	// Create encoder config
	var encoderConfig zapcore.EncoderConfig
	if config.Format == "json" {
		encoderConfig = zap.NewProductionEncoderConfig()
	} else {
		encoderConfig = zap.NewDevelopmentEncoderConfig()
		encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	}
	
	// Create encoder
	var encoder zapcore.Encoder
	if config.Format == "json" {
		encoder = zapcore.NewJSONEncoder(encoderConfig)
	} else {
		encoder = zapcore.NewConsoleEncoder(encoderConfig)
	}
	
	// Create writer syncer
	writer := zapcore.AddSync(getWriter(config))
	
	// Create core
	core := zapcore.NewCore(encoder, writer, parseZapLevel(config.Level))
	
	// Create logger
	logger := zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))
	
	return &zapLogger{
		logger: logger,
		sugar:  logger.Sugar(),
	}
}

func (z *zapLogger) Debug(msg string, fields ...interface{}) {
	z.sugar.Debugw(msg, fields...)
}

func (z *zapLogger) Info(msg string, fields ...interface{}) {
	z.sugar.Infow(msg, fields...)
}

func (z *zapLogger) Warn(msg string, fields ...interface{}) {
	z.sugar.Warnw(msg, fields...)
}

func (z *zapLogger) Error(msg string, fields ...interface{}) {
	z.sugar.Errorw(msg, fields...)
}

func (z *zapLogger) Fatal(msg string, fields ...interface{}) {
	z.sugar.Fatalw(msg, fields...)
}

func (z *zapLogger) WithField(key string, value interface{}) Logger {
	newLogger := z.logger.With(zap.Any(key, value))
	return &zapLogger{
		logger: newLogger,
		sugar:  newLogger.Sugar(),
	}
}

func (z *zapLogger) WithFields(fields map[string]interface{}) Logger {
	zapFields := make([]zap.Field, 0, len(fields))
	for key, value := range fields {
		zapFields = append(zapFields, zap.Any(key, value))
	}
	
	newLogger := z.logger.With(zapFields...)
	return &zapLogger{
		logger: newLogger,
		sugar:  newLogger.Sugar(),
	}
}

// parseZapLevel converts string level to zap level
func parseZapLevel(level string) zapcore.Level {
	switch strings.ToLower(level) {
	case "debug":
		return zap.DebugLevel
	case "info":
		return zap.InfoLevel
	case "warn", "warning":
		return zap.WarnLevel
	case "error":
		return zap.ErrorLevel
	case "fatal":
		return zap.FatalLevel
	default:
		return zap.InfoLevel
	}
}
