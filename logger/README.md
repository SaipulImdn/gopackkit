# Logger Module

Logger module menyediakan multiple backend logging dengan interface yang unified. Module ini mendukung Simple logger, Logrus, dan Zap dengan konfigurasi yang fleksibel.

## Features

- **Multiple Backends**: Simple, Logrus, Zap
- **Configurable Log Levels**: Debug, Info, Warn, Error, Fatal, Panic
- **Multiple Output Formats**: JSON, Text, Console
- **Thread-Safe Operations**: Aman digunakan di concurrent environment
- **Easy Configuration**: Konfigurasi melalui struct dengan sensible defaults

## Installation

```bash
go get github.com/saipulimdn/gopackkit/logger
```

## Quick Start

### Simple Logger

```go
package main

import "github.com/saipulimdn/gopackkit/logger"

func main() {
    log := logger.NewSimple()
    
    log.Info("Application started")
    log.Warn("This is a warning")
    log.Error("This is an error")
}
```

### Logrus Logger

```go
package main

import "github.com/saipulimdn/gopackkit/logger"

func main() {
    // Default configuration
    log := logger.NewLogrus(logger.DefaultLogrusConfig())
    
    // Custom configuration
    config := logger.LogrusConfig{
        Level:  "debug",
        Format: "json",
    }
    log := logger.NewLogrus(config)
    
    log.Debug("Debug message")
    log.Info("Info message")
    log.Error("Error message")
}
```

### Zap Logger

```go
package main

import "github.com/saipulimdn/gopackkit/logger"

func main() {
    // Default configuration
    log := logger.NewZap(logger.DefaultZapConfig())
    
    // Custom configuration
    config := logger.ZapConfig{
        Level:       "info",
        Development: true,
        Encoding:    "console",
    }
    log := logger.NewZap(config)
    
    log.Info("Application started")
    log.Warn("Warning message")
    log.Error("Error occurred")
}
```

## Configuration

### Logrus Configuration

```go
type LogrusConfig struct {
    Level  string `json:"level" yaml:"level" env:"LOG_LEVEL" default:"info"`
    Format string `json:"format" yaml:"format" env:"LOG_FORMAT" default:"json"`
}
```

**Supported Levels**: `panic`, `fatal`, `error`, `warn`, `info`, `debug`, `trace`
**Supported Formats**: `json`, `text`

### Zap Configuration

```go
type ZapConfig struct {
    Level       string `json:"level" yaml:"level" env:"LOG_LEVEL" default:"info"`
    Development bool   `json:"development" yaml:"development" env:"LOG_DEVELOPMENT" default:"false"`
    Encoding    string `json:"encoding" yaml:"encoding" env:"LOG_ENCODING" default:"json"`
}
```

**Supported Levels**: `debug`, `info`, `warn`, `error`, `dpanic`, `panic`, `fatal`
**Supported Encodings**: `json`, `console`

## Methods

All loggers implement the common interface:

```go
type Logger interface {
    Debug(msg string, keysAndValues ...interface{})
    Info(msg string, keysAndValues ...interface{})
    Warn(msg string, keysAndValues ...interface{})
    Error(msg string, keysAndValues ...interface{})
    Fatal(msg string, keysAndValues ...interface{})
    Panic(msg string, keysAndValues ...interface{})
}
```

### Usage with Structured Logging

```go
log := logger.NewLogrus(logger.LogrusConfig{
    Level:  "info",
    Format: "json",
})

// Simple message
log.Info("User created")

// With structured data
log.Info("User created", 
    "user_id", "12345",
    "email", "user@example.com",
    "role", "admin",
)

// With error
log.Error("Database connection failed",
    "error", err,
    "host", "localhost:5432",
    "retries", 3,
)
```

## Environment Variables

Anda dapat mengkonfigurasi logger menggunakan environment variables:

```bash
# Log level
export LOG_LEVEL=debug

# Log format (logrus)
export LOG_FORMAT=json

# Development mode (zap)
export LOG_DEVELOPMENT=true

# Encoding (zap)
export LOG_ENCODING=console
```

## Examples

### Web API Logging

```go
package main

import (
    "net/http"
    "time"
    
    "github.com/saipulimdn/gopackkit/logger"
)

func main() {
    log := logger.NewLogrus(logger.LogrusConfig{
        Level:  "info",
        Format: "json",
    })
    
    // Middleware untuk logging request
    http.HandleFunc("/api/users", func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()
        
        // Log incoming request
        log.Info("Incoming request",
            "method", r.Method,
            "path", r.URL.Path,
            "remote_addr", r.RemoteAddr,
            "user_agent", r.UserAgent(),
        )
        
        // Handle request
        w.WriteHeader(http.StatusOK)
        w.Write([]byte(`{"message": "success"}`))
        
        // Log response
        duration := time.Since(start)
        log.Info("Request completed",
            "method", r.Method,
            "path", r.URL.Path,
            "status", 200,
            "duration_ms", duration.Milliseconds(),
        )
    })
    
    log.Info("Server starting", "port", 8080)
    http.ListenAndServe(":8080", nil)
}
```

### Error Logging with Context

```go
package main

import (
    "errors"
    "github.com/saipulimdn/gopackkit/logger"
)

func processUser(userID string) error {
    log := logger.NewZap(logger.ZapConfig{
        Level:       "debug",
        Development: true,
        Encoding:    "console",
    })
    
    log.Debug("Processing user", "user_id", userID)
    
    // Simulate some processing
    if userID == "" {
        err := errors.New("user ID cannot be empty")
        log.Error("Invalid user ID",
            "error", err,
            "user_id", userID,
        )
        return err
    }
    
    // Simulate database error
    if userID == "error" {
        err := errors.New("database connection failed")
        log.Error("Database error occurred",
            "error", err,
            "user_id", userID,
            "operation", "user_lookup",
        )
        return err
    }
    
    log.Info("User processed successfully", "user_id", userID)
    return nil
}

func main() {
    log := logger.NewSimple()
    
    users := []string{"123", "", "456", "error", "789"}
    
    for _, userID := range users {
        if err := processUser(userID); err != nil {
            log.Error("Failed to process user",
                "user_id", userID,
                "error", err,
            )
        }
    }
}
```

### Configuration from Environment

```go
package main

import (
    "os"
    "github.com/saipulimdn/gopackkit/logger"
)

func main() {
    // Get configuration from environment
    logLevel := os.Getenv("LOG_LEVEL")
    if logLevel == "" {
        logLevel = "info"
    }
    
    logFormat := os.Getenv("LOG_FORMAT")
    if logFormat == "" {
        logFormat = "json"
    }
    
    config := logger.LogrusConfig{
        Level:  logLevel,
        Format: logFormat,
    }
    
    log := logger.NewLogrus(config)
    
    log.Info("Logger configured",
        "level", config.Level,
        "format", config.Format,
    )
}
```

## Performance

### Benchmarks

- **Simple Logger**: Paling cepat untuk development
- **Logrus**: Balance antara performance dan features
- **Zap**: Paling performant untuk production

### Memory Usage

- **Simple Logger**: Minimal memory footprint
- **Logrus**: Moderate memory usage
- **Zap**: Optimized untuk minimal allocations

## Best Practices

1. **Gunakan structured logging** dengan key-value pairs
2. **Set log level yang appropriate** untuk environment (debug untuk dev, info untuk prod)
3. **Include context information** seperti request ID, user ID
4. **Avoid logging sensitive information** seperti passwords, tokens
5. **Use consistent key names** across aplikasi
6. **Log errors dengan sufficient context** untuk debugging

## Testing

```go
package main

import (
    "testing"
    "github.com/saipulimdn/gopackkit/logger"
)

func TestLogger(t *testing.T) {
    log := logger.NewSimple()
    
    // Test tidak akan panic
    log.Info("Test message")
    log.Debug("Debug message")
    log.Warn("Warning message")
    
    // Test dengan parameters
    log.Info("Test with params",
        "key1", "value1",
        "key2", 123,
    )
}
```

## Troubleshooting

### Common Issues

1. **Log level tidak berubah**: Pastikan environment variable `LOG_LEVEL` di-set dengan benar
2. **Format tidak sesuai**: Check `LOG_FORMAT` environment variable
3. **Performance issues**: Gunakan Zap untuk high-throughput applications

### Debug Mode

```go
// Enable debug mode untuk melihat semua log messages
config := logger.LogrusConfig{
    Level:  "debug",
    Format: "text", // More readable untuk debugging
}
log := logger.NewLogrus(config)
```
