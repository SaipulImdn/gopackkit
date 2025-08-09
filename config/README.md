# Config Module

Config module menyediakan functionality untuk loading environment variables ke dalam Go structs dengan menggunakan struct tags. Module ini mendukung type conversion, default values, dan validation.

## Features

- **Struct Tag Support**: Automatic mapping menggunakan `env` tags
- **Default Values**: Built-in default value support
- **Type Conversion**: Automatic conversion untuk berbagai Go types
- **Environment Loading**: Load dari environment variables
- **Validation**: Basic validation untuk required fields
- **Nested Structs**: Support untuk nested configuration

## Installation

```bash
go get github.com/saipulimdn/gopackkit/config
```

## Quick Start

### Basic Usage

```go
package main

import (
    "fmt"
    "log"
    "time"
    
    "github.com/saipulimdn/gopackkit/config"
)

type AppConfig struct {
    Port        int           `env:"PORT" default:"8080"`
    DatabaseURL string        `env:"DATABASE_URL" default:"localhost:5432"`
    Debug       bool          `env:"DEBUG" default:"false"`
    Timeout     time.Duration `env:"TIMEOUT" default:"30s"`
    SecretKey   string        `env:"SECRET_KEY"`
}

func main() {
    var cfg AppConfig
    
    // Load configuration from environment
    err := config.Load(&cfg)
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("Port: %d\n", cfg.Port)
    fmt.Printf("Database URL: %s\n", cfg.DatabaseURL)
    fmt.Printf("Debug: %t\n", cfg.Debug)
    fmt.Printf("Timeout: %v\n", cfg.Timeout)
    fmt.Printf("Secret Key: %s\n", cfg.SecretKey)
}
```

### Environment Variables

```bash
# Set environment variables
export PORT=9000
export DATABASE_URL=postgresql://user:pass@localhost:5432/db
export DEBUG=true
export TIMEOUT=60s
export SECRET_KEY=your-secret-key

# Run application
go run main.go
```

## Supported Types

Config module mendukung automatic type conversion untuk:

### Basic Types

```go
type Config struct {
    // String
    Name        string `env:"APP_NAME" default:"MyApp"`
    
    // Integers
    Port        int    `env:"PORT" default:"8080"`
    MaxConns    int32  `env:"MAX_CONNECTIONS" default:"100"`
    BufferSize  int64  `env:"BUFFER_SIZE" default:"1024"`
    
    // Unsigned Integers
    Workers     uint   `env:"WORKERS" default:"4"`
    Retries     uint32 `env:"RETRIES" default:"3"`
    
    // Floating Point
    Rate        float32 `env:"RATE" default:"0.5"`
    Precision   float64 `env:"PRECISION" default:"0.001"`
    
    // Boolean
    Debug       bool `env:"DEBUG" default:"false"`
    EnableTLS   bool `env:"ENABLE_TLS" default:"true"`
    
    // Duration
    Timeout     time.Duration `env:"TIMEOUT" default:"30s"`
    Interval    time.Duration `env:"INTERVAL" default:"1m"`
}
```

### Duration Format

Duration values support Go's standard duration format:

```bash
export TIMEOUT=30s      # 30 seconds
export INTERVAL=5m      # 5 minutes
export DEADLINE=2h      # 2 hours
export RETENTION=24h    # 24 hours
export CLEANUP=7d       # 7 days (if using custom parser)
```

## Advanced Usage

### Nested Structs

```go
type DatabaseConfig struct {
    Host         string        `env:"DB_HOST" default:"localhost"`
    Port         int           `env:"DB_PORT" default:"5432"`
    Username     string        `env:"DB_USERNAME" default:"postgres"`
    Password     string        `env:"DB_PASSWORD"`
    Database     string        `env:"DB_NAME" default:"myapp"`
    SSLMode      string        `env:"DB_SSL_MODE" default:"disable"`
    MaxConns     int           `env:"DB_MAX_CONNECTIONS" default:"25"`
    ConnTimeout  time.Duration `env:"DB_CONN_TIMEOUT" default:"10s"`
}

type RedisConfig struct {
    Host         string        `env:"REDIS_HOST" default:"localhost"`
    Port         int           `env:"REDIS_PORT" default:"6379"`
    Password     string        `env:"REDIS_PASSWORD"`
    Database     int           `env:"REDIS_DB" default:"0"`
    PoolSize     int           `env:"REDIS_POOL_SIZE" default:"10"`
    IdleTimeout  time.Duration `env:"REDIS_IDLE_TIMEOUT" default:"5m"`
}

type AppConfig struct {
    // Application settings
    Port         int    `env:"PORT" default:"8080"`
    Environment  string `env:"ENVIRONMENT" default:"development"`
    LogLevel     string `env:"LOG_LEVEL" default:"info"`
    
    // Nested configurations
    Database     DatabaseConfig
    Redis        RedisConfig
    
    // Security
    JWTSecret    string        `env:"JWT_SECRET"`
    SessionTTL   time.Duration `env:"SESSION_TTL" default:"24h"`
}

func main() {
    var cfg AppConfig
    
    err := config.Load(&cfg)
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("App Port: %d\n", cfg.Port)
    fmt.Printf("DB Host: %s\n", cfg.Database.Host)
    fmt.Printf("Redis Host: %s\n", cfg.Redis.Host)
}
```

### Environment File Example

```bash
# .env file
PORT=8080
ENVIRONMENT=production
LOG_LEVEL=info

# Database configuration
DB_HOST=db.example.com
DB_PORT=5432
DB_USERNAME=myapp
DB_PASSWORD=secretpassword
DB_NAME=myapp_production
DB_SSL_MODE=require
DB_MAX_CONNECTIONS=50
DB_CONN_TIMEOUT=15s

# Redis configuration
REDIS_HOST=redis.example.com
REDIS_PORT=6379
REDIS_PASSWORD=redispassword
REDIS_DB=0
REDIS_POOL_SIZE=20
REDIS_IDLE_TIMEOUT=10m

# Security
JWT_SECRET=your-super-secret-jwt-key
SESSION_TTL=48h
```

### Validation

```go
package main

import (
    "errors"
    "fmt"
    "log"
    
    "github.com/saipulimdn/gopackkit/config"
)

type AppConfig struct {
    Port        int    `env:"PORT" default:"8080"`
    DatabaseURL string `env:"DATABASE_URL"`
    APIKey      string `env:"API_KEY"`
    Environment string `env:"ENVIRONMENT" default:"development"`
}

func (c *AppConfig) Validate() error {
    if c.Port < 1 || c.Port > 65535 {
        return errors.New("port must be between 1 and 65535")
    }
    
    if c.Environment == "production" && c.DatabaseURL == "" {
        return errors.New("DATABASE_URL is required in production")
    }
    
    if c.Environment == "production" && c.APIKey == "" {
        return errors.New("API_KEY is required in production")
    }
    
    return nil
}

func main() {
    var cfg AppConfig
    
    // Load configuration
    err := config.Load(&cfg)
    if err != nil {
        log.Fatal("Failed to load config:", err)
    }
    
    // Validate configuration
    err = cfg.Validate()
    if err != nil {
        log.Fatal("Configuration validation failed:", err)
    }
    
    fmt.Println("Configuration loaded and validated successfully")
}
```

## Examples

### Web Server Configuration

```go
package main

import (
    "fmt"
    "log"
    "net/http"
    "time"
    
    "github.com/saipulimdn/gopackkit/config"
)

type ServerConfig struct {
    // Server settings
    Host            string        `env:"HOST" default:"0.0.0.0"`
    Port            int           `env:"PORT" default:"8080"`
    ReadTimeout     time.Duration `env:"READ_TIMEOUT" default:"30s"`
    WriteTimeout    time.Duration `env:"WRITE_TIMEOUT" default:"30s"`
    IdleTimeout     time.Duration `env:"IDLE_TIMEOUT" default:"60s"`
    ShutdownTimeout time.Duration `env:"SHUTDOWN_TIMEOUT" default:"10s"`
    
    // TLS settings
    EnableTLS       bool   `env:"ENABLE_TLS" default:"false"`
    CertFile        string `env:"CERT_FILE"`
    KeyFile         string `env:"KEY_FILE"`
    
    // CORS settings
    AllowOrigins    string `env:"ALLOW_ORIGINS" default:"*"`
    AllowMethods    string `env:"ALLOW_METHODS" default:"GET,POST,PUT,DELETE"`
    AllowHeaders    string `env:"ALLOW_HEADERS" default:"Content-Type,Authorization"`
}

func main() {
    var cfg ServerConfig
    
    err := config.Load(&cfg)
    if err != nil {
        log.Fatal("Failed to load configuration:", err)
    }
    
    // Create HTTP server
    server := &http.Server{
        Addr:         fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
        ReadTimeout:  cfg.ReadTimeout,
        WriteTimeout: cfg.WriteTimeout,
        IdleTimeout:  cfg.IdleTimeout,
    }
    
    // Setup routes
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusOK)
        w.Write([]byte("Hello, World!"))
    })
    
    log.Printf("Server starting on %s", server.Addr)
    
    // Start server
    if cfg.EnableTLS {
        log.Fatal(server.ListenAndServeTLS(cfg.CertFile, cfg.KeyFile))
    } else {
        log.Fatal(server.ListenAndServe())
    }
}
```

### Database Configuration

```go
package main

import (
    "database/sql"
    "fmt"
    "log"
    "time"
    
    "github.com/saipulimdn/gopackkit/config"
    _ "github.com/lib/pq" // PostgreSQL driver
)

type DatabaseConfig struct {
    Driver          string        `env:"DB_DRIVER" default:"postgres"`
    Host            string        `env:"DB_HOST" default:"localhost"`
    Port            int           `env:"DB_PORT" default:"5432"`
    Username        string        `env:"DB_USERNAME" default:"postgres"`
    Password        string        `env:"DB_PASSWORD"`
    Database        string        `env:"DB_NAME" default:"myapp"`
    SSLMode         string        `env:"DB_SSL_MODE" default:"disable"`
    MaxOpenConns    int           `env:"DB_MAX_OPEN_CONNS" default:"25"`
    MaxIdleConns    int           `env:"DB_MAX_IDLE_CONNS" default:"5"`
    ConnMaxLifetime time.Duration `env:"DB_CONN_MAX_LIFETIME" default:"1h"`
    ConnMaxIdleTime time.Duration `env:"DB_CONN_MAX_IDLE_TIME" default:"30m"`
}

func (c *DatabaseConfig) ConnectionString() string {
    return fmt.Sprintf(
        "host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
        c.Host, c.Port, c.Username, c.Password, c.Database, c.SSLMode,
    )
}

func main() {
    var cfg DatabaseConfig
    
    err := config.Load(&cfg)
    if err != nil {
        log.Fatal("Failed to load config:", err)
    }
    
    // Connect to database
    db, err := sql.Open(cfg.Driver, cfg.ConnectionString())
    if err != nil {
        log.Fatal("Failed to connect to database:", err)
    }
    defer db.Close()
    
    // Configure connection pool
    db.SetMaxOpenConns(cfg.MaxOpenConns)
    db.SetMaxIdleConns(cfg.MaxIdleConns)
    db.SetConnMaxLifetime(cfg.ConnMaxLifetime)
    db.SetConnMaxIdleTime(cfg.ConnMaxIdleTime)
    
    // Test connection
    err = db.Ping()
    if err != nil {
        log.Fatal("Failed to ping database:", err)
    }
    
    log.Println("Database connected successfully")
}
```

### Microservice Configuration

```go
package main

import (
    "log"
    "time"
    
    "github.com/saipulimdn/gopackkit/config"
)

type ServiceConfig struct {
    // Service identity
    ServiceName    string `env:"SERVICE_NAME" default:"my-service"`
    ServiceVersion string `env:"SERVICE_VERSION" default:"1.0.0"`
    Environment    string `env:"ENVIRONMENT" default:"development"`
    
    // HTTP Server
    HTTPPort       int           `env:"HTTP_PORT" default:"8080"`
    HTTPTimeout    time.Duration `env:"HTTP_TIMEOUT" default:"30s"`
    
    // gRPC Server
    GRPCPort       int           `env:"GRPC_PORT" default:"9090"`
    GRPCTimeout    time.Duration `env:"GRPC_TIMEOUT" default:"30s"`
    
    // Database
    DatabaseURL    string `env:"DATABASE_URL"`
    
    // Message Queue
    RedisURL       string `env:"REDIS_URL" default:"redis://localhost:6379"`
    QueueName      string `env:"QUEUE_NAME" default:"default"`
    
    // External APIs
    UserServiceURL    string        `env:"USER_SERVICE_URL" default:"http://localhost:8081"`
    PaymentServiceURL string        `env:"PAYMENT_SERVICE_URL" default:"http://localhost:8082"`
    APITimeout        time.Duration `env:"API_TIMEOUT" default:"10s"`
    
    // Observability
    LogLevel          string `env:"LOG_LEVEL" default:"info"`
    MetricsPort       int    `env:"METRICS_PORT" default:"9000"`
    TracingEndpoint   string `env:"TRACING_ENDPOINT"`
    
    // Security
    JWTSecret         string        `env:"JWT_SECRET"`
    TokenExpiry       time.Duration `env:"TOKEN_EXPIRY" default:"1h"`
    RateLimitRPS      int           `env:"RATE_LIMIT_RPS" default:"100"`
}

func main() {
    var cfg ServiceConfig
    
    err := config.Load(&cfg)
    if err != nil {
        log.Fatal("Failed to load configuration:", err)
    }
    
    log.Printf("Starting %s v%s in %s environment", 
        cfg.ServiceName, cfg.ServiceVersion, cfg.Environment)
    log.Printf("HTTP Server: :%d", cfg.HTTPPort)
    log.Printf("gRPC Server: :%d", cfg.GRPCPort)
    log.Printf("Metrics: :%d", cfg.MetricsPort)
    
    // Initialize and start services...
}
```

## Environment Variables Best Practices

### Naming Convention

```bash
# Use consistent prefixes for related settings
APP_NAME=myapp
APP_VERSION=1.0.0
APP_ENVIRONMENT=production

# Group by functionality
DB_HOST=localhost
DB_PORT=5432
DB_USERNAME=postgres

REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=secret

LOG_LEVEL=info
LOG_FORMAT=json

# Use descriptive names
HTTP_READ_TIMEOUT=30s
HTTP_WRITE_TIMEOUT=30s
JWT_ACCESS_TOKEN_EXPIRY=15m
JWT_REFRESH_TOKEN_EXPIRY=7d
```

### Production vs Development

```bash
# Development (.env.development)
ENVIRONMENT=development
LOG_LEVEL=debug
DEBUG=true
DB_HOST=localhost
ENABLE_TLS=false

# Production (.env.production)
ENVIRONMENT=production
LOG_LEVEL=info
DEBUG=false
DB_HOST=prod-db.example.com
ENABLE_TLS=true
```

## Testing

```go
package main

import (
    "os"
    "testing"
    "time"
    
    "github.com/saipulimdn/gopackkit/config"
)

type TestConfig struct {
    Port    int           `env:"TEST_PORT" default:"8080"`
    Timeout time.Duration `env:"TEST_TIMEOUT" default:"30s"`
    Debug   bool          `env:"TEST_DEBUG" default:"false"`
}

func TestConfigLoad(t *testing.T) {
    // Set test environment variables
    os.Setenv("TEST_PORT", "9000")
    os.Setenv("TEST_TIMEOUT", "60s")
    os.Setenv("TEST_DEBUG", "true")
    
    defer func() {
        // Clean up
        os.Unsetenv("TEST_PORT")
        os.Unsetenv("TEST_TIMEOUT")
        os.Unsetenv("TEST_DEBUG")
    }()
    
    var cfg TestConfig
    err := config.Load(&cfg)
    if err != nil {
        t.Fatalf("Failed to load config: %v", err)
    }
    
    // Test values
    if cfg.Port != 9000 {
        t.Errorf("Expected port 9000, got %d", cfg.Port)
    }
    
    if cfg.Timeout != 60*time.Second {
        t.Errorf("Expected timeout 60s, got %v", cfg.Timeout)
    }
    
    if !cfg.Debug {
        t.Errorf("Expected debug true, got %t", cfg.Debug)
    }
}

func TestConfigDefaults(t *testing.T) {
    var cfg TestConfig
    err := config.Load(&cfg)
    if err != nil {
        t.Fatalf("Failed to load config: %v", err)
    }
    
    // Test default values
    if cfg.Port != 8080 {
        t.Errorf("Expected default port 8080, got %d", cfg.Port)
    }
    
    if cfg.Timeout != 30*time.Second {
        t.Errorf("Expected default timeout 30s, got %v", cfg.Timeout)
    }
    
    if cfg.Debug {
        t.Errorf("Expected default debug false, got %t", cfg.Debug)
    }
}
```

## Troubleshooting

### Common Issues

1. **Type conversion errors**: Pastikan format environment variable sesuai dengan Go type
2. **Missing required fields**: Set environment variables atau provide default values
3. **Duration parsing errors**: Gunakan Go duration format (1s, 1m, 1h)
4. **Nested struct issues**: Pastikan environment variable names unik

### Debug Mode

```go
package main

import (
    "fmt"
    "os"
    "reflect"
    
    "github.com/saipulimdn/gopackkit/config"
)

func debugConfig(cfg interface{}) {
    v := reflect.ValueOf(cfg).Elem()
    t := reflect.TypeOf(cfg).Elem()
    
    fmt.Println("=== Configuration Debug ===")
    for i := 0; i < v.NumField(); i++ {
        field := t.Field(i)
        value := v.Field(i)
        envTag := field.Tag.Get("env")
        defaultTag := field.Tag.Get("default")
        
        envValue := os.Getenv(envTag)
        
        fmt.Printf("%s:\n", field.Name)
        fmt.Printf("  Env Var: %s\n", envTag)
        fmt.Printf("  Env Value: %s\n", envValue)
        fmt.Printf("  Default: %s\n", defaultTag)
        fmt.Printf("  Final Value: %v\n", value.Interface())
        fmt.Println()
    }
}

func main() {
    var cfg AppConfig
    
    err := config.Load(&cfg)
    if err != nil {
        log.Fatal(err)
    }
    
    debugConfig(&cfg)
}
```
