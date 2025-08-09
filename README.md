# GoPackKit

GoPackKit is a collection of reusable utility modules for Go applications. This toolkit provides common functionality including logging, HTTP client, configuration management, validation, object storage, authentication, and secure password handling.

## Installation

```bash
go get github.com/saipulimdn/gopackkit
```

## Modules

GoPackKit provides the following modules, each with comprehensive documentation:

### Core Utilities
- **[Logger](logger/)** - Multi-backend logging with Simple, Logrus, and Zap support
- **[HTTP Client](httpclient/)** - HTTP client with retry mechanism and timeout configuration
- **[Config](config/)** - Environment variable loader with struct tag support
- **[Validator](validator/)** - Secure data validation with XSS/injection prevention

### Authentication & Security  
- **[JWT](jwt/)** - JSON Web Token management with HMAC-SHA256 signing
- **[Password](password/)** - Secure password hashing and validation using bcrypt

### Storage & Communication
- **[MinIO](minio/)** - MinIO client for object storage with presigned URL support
- **[gRPC](grpc/)** - gRPC client and server with comprehensive configuration

## Quick Start

Each module can be used independently. Here's a simple example using multiple modules:

```go
package main

import (
    "log"
    "time"
    
    "github.com/saipulimdn/gopackkit/config"
    "github.com/saipulimdn/gopackkit/logger"
    "github.com/saipulimdn/gopackkit/password"
    "github.com/saipulimdn/gopackkit/jwt"
)

type AppConfig struct {
    Port       int    `env:"PORT" default:"8080"`
    LogLevel   string `env:"LOG_LEVEL" default:"info"`
    JWTSecret  string `env:"JWT_SECRET" default:"your-secret-key"`
}

func main() {
    // Load configuration
    var cfg AppConfig
    if err := config.Load(&cfg); err != nil {
        log.Fatal("Failed to load config:", err)
    }
    
    // Setup logger
    logger := logger.NewLogrus(logger.LogrusConfig{
        Level:  cfg.LogLevel,
        Format: "json",
    })
    
    // Setup password manager
    pm := password.New()
    
    // Setup JWT manager
    jwtManager := jwt.New(jwt.Config{
        SecretKey:         cfg.JWTSecret,
        AccessTokenExpiry: 15 * time.Minute,
        Issuer:           "gopackkit-example",
    })
    
    logger.Info("Application started successfully",
        "port", cfg.Port,
        "log_level", cfg.LogLevel,
    )
    
    // Example usage
    testPassword := "MySecurePassword123!"
    
    // Hash password
    hashedPassword, err := pm.Hash(testPassword)
    if err != nil {
        logger.Error("Failed to hash password", "error", err)
        return
    }
    
    // Generate JWT token
    claims := jwt.Claims{
        UserID: "user123",
        Email:  "user@example.com",
        Role:   "user",
    }
    
    tokens, err := jwtManager.GenerateTokens(claims)
    if err != nil {
        logger.Error("Failed to generate tokens", "error", err)
        return
    }
    
    logger.Info("Successfully processed user",
        "password_hash", hashedPassword.Hash,
        "access_token", tokens.AccessToken,
    )
}
```

## Module Overview

### [Logger](logger/)
Multi-backend logging system supporting Simple, Logrus, and Zap with:
- Configurable log levels and formats
- JSON and text output support
- Environment variable configuration
- Thread-safe operations

### [HTTP Client](httpclient/)
Robust HTTP client with enterprise features:
- Automatic retry with exponential backoff
- Configurable timeouts and connection pooling
- JSON request/response handling
- Comprehensive error handling

### [Config](config/)
Environment variable configuration loader:
- Struct tag-based mapping
- Type conversion support
- Default value handling
- Nested struct support

### [Validator](validator/)
Security-focused data validation:
- XSS and injection prevention
- Safe email and phone validation
- Alphanumeric checking
- No regex injection vulnerabilities

### [JWT](jwt/)
JSON Web Token management:
- HMAC-SHA256 signing
- Access and refresh token support
- Custom claims handling
- Token validation and refresh

### [Password](password/)
Secure password operations:
- Bcrypt hashing with configurable cost
- Password strength validation
- Random password generation
- Hash migration support

### [MinIO](minio/)
Object storage client:
- Presigned URL generation (GET, PUT, POST)
- Object upload/download operations
- Bucket management
- TLS support

### [gRPC](grpc/)
gRPC client and server implementation:
- Connection management and health checking
- TLS/SSL support
- Keep-alive configuration
- Graceful shutdown handling

## Features

### Security First
- **Input Validation**: Comprehensive validation without regex injection risks
- **Secure Hashing**: Bcrypt with configurable cost factors
- **JWT Security**: HMAC-SHA256 signing with proper claim validation
- **TLS Support**: Full TLS/SSL support across all network modules

### Production Ready
- **Error Handling**: Comprehensive error handling and logging
- **Configuration**: Environment variable support across all modules
- **Performance**: Optimized for production workloads
- **Monitoring**: Built-in health checks and status monitoring

### Developer Friendly
- **Easy Integration**: Simple APIs with sensible defaults
- **Comprehensive Documentation**: Detailed examples and use cases
- **Testing Support**: Built-in testing utilities and examples
- **Flexible Configuration**: Customizable settings for different environments

## Installation & Usage

### Individual Modules

You can import and use individual modules as needed:

```go
// Logger module
import "github.com/saipulimdn/gopackkit/logger"

// HTTP Client module  
import "github.com/saipulimdn/gopackkit/httpclient"

// Config module
import "github.com/saipulimdn/gopackkit/config"

// And so on...
```

### Full Toolkit

Or import the entire toolkit:

```go
import (
    "github.com/saipulimdn/gopackkit/logger"
    "github.com/saipulimdn/gopackkit/httpclient"
    "github.com/saipulimdn/gopackkit/config"
    "github.com/saipulimdn/gopackkit/validator"
    "github.com/saipulimdn/gopackkit/jwt"
    "github.com/saipulimdn/gopackkit/password"
    "github.com/saipulimdn/gopackkit/minio"
    "github.com/saipulimdn/gopackkit/grpc"
)
```

## Examples

### Complete Web Application

```go
package main

import (
    "encoding/json"
    "net/http"
    "time"
    
    "github.com/saipulimdn/gopackkit/config"
    "github.com/saipulimdn/gopackkit/jwt"
    "github.com/saipulimdn/gopackkit/logger"
    "github.com/saipulimdn/gopackkit/password"
    "github.com/saipulimdn/gopackkit/validator"
    "github.com/saipulimdn/gopackkit/httpclient"
)

type Config struct {
    Port      int    `env:"PORT" default:"8080"`
    JWTSecret string `env:"JWT_SECRET" default:"your-secret-key"`
    LogLevel  string `env:"LOG_LEVEL" default:"info"`
}

type LoginRequest struct {
    Email    string `json:"email" validate:"required,email_safe"`
    Password string `json:"password" validate:"required"`
}

type User struct {
    ID       string `json:"id"`
    Email    string `json:"email"`
    Password string `json:"-"`
}

func main() {
    // Load configuration
    var cfg Config
    if err := config.Load(&cfg); err != nil {
        panic(err)
    }
    
    // Setup components
    log := logger.NewLogrus(logger.LogrusConfig{Level: cfg.LogLevel, Format: "json"})
    pm := password.New()
    jwtManager := jwt.New(jwt.Config{SecretKey: cfg.JWTSecret, Issuer: "myapp"})
    validator := validator.New()
    httpClient := httpclient.New()
    
    // Login handler
    http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
        var req LoginRequest
        if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
            http.Error(w, "Invalid JSON", http.StatusBadRequest)
            return
        }
        
        // Validate request
        if err := validator.Struct(req); err != nil {
            http.Error(w, err.Error(), http.StatusBadRequest)
            return
        }
        
        // Authenticate user (simplified)
        user := User{ID: "user123", Email: req.Email}
        storedHash := "$2a$12$..." // from database
        
        if err := pm.Verify(req.Password, storedHash); err != nil {
            log.Warn("Invalid login attempt", "email", req.Email)
            http.Error(w, "Invalid credentials", http.StatusUnauthorized)
            return
        }
        
        // Generate JWT tokens
        claims := jwt.Claims{UserID: user.ID, Email: user.Email, Role: "user"}
        tokens, err := jwtManager.GenerateTokens(claims)
        if err != nil {
            log.Error("Failed to generate tokens", "error", err)
            http.Error(w, "Internal error", http.StatusInternalServerError)
            return
        }
        
        log.Info("User logged in", "user_id", user.ID, "email", user.Email)
        json.NewEncoder(w).Encode(tokens)
    })
    
    log.Info("Server starting", "port", cfg.Port)
    http.ListenAndServe(fmt.Sprintf(":%d", cfg.Port), nil)
}
```

### Microservice with gRPC

```go
package main

import (
    "context"
    "log"
    "net"
    
    "github.com/saipulimdn/gopackkit/grpc"
    "github.com/saipulimdn/gopackkit/logger"
    "github.com/saipulimdn/gopackkit/config"
    "github.com/saipulimdn/gopackkit/validator"
)

type ServiceConfig struct {
    GRPCPort  int    `env:"GRPC_PORT" default:"9000"`
    HTTPPort  int    `env:"HTTP_PORT" default:"8080"`
    LogLevel  string `env:"LOG_LEVEL" default:"info"`
}

func main() {
    var cfg ServiceConfig
    config.Load(&cfg)
    
    log := logger.NewZap(logger.ZapConfig{Level: cfg.LogLevel})
    validator := validator.New()
    
    // Create gRPC server
    grpcServer, err := grpc.NewServerWithConfig(grpc.ServerConfig{
        Port: cfg.GRPCPort,
        EnableReflection: true,
    })
    if err != nil {
        log.Fatal("Failed to create gRPC server", "error", err)
    }
    
    // Register services and start server
    log.Info("Starting microservice", "grpc_port", cfg.GRPCPort)
    if err := grpcServer.Start(); err != nil {
        log.Fatal("Failed to start gRPC server", "error", err)
    }
}
```

### File Upload Service with MinIO

```go
package main

import (
    "encoding/json"
    "net/http"
    "time"
    
    "github.com/saipulimdn/gopackkit/minio"
    "github.com/saipulimdn/gopackkit/logger"
    "github.com/saipulimdn/gopackkit/validator"
)

type UploadRequest struct {
    Filename string `json:"filename" validate:"required"`
    Bucket   string `json:"bucket" validate:"required,alphanumeric"`
}

func main() {
    log := logger.NewSimple()
    validator := validator.New()
    
    // Setup MinIO client
    minioClient, err := minio.New(minio.Config{
        Endpoint:        "localhost:9000",
        AccessKeyID:     "minioadmin",
        SecretAccessKey: "minioadmin",
        UseSSL:          false,
    })
    if err != nil {
        log.Fatal("Failed to create MinIO client", "error", err)
    }
    
    http.HandleFunc("/upload-url", func(w http.ResponseWriter, r *http.Request) {
        var req UploadRequest
        json.NewDecoder(r.Body).Decode(&req)
        
        if err := validator.Struct(req); err != nil {
            http.Error(w, err.Error(), http.StatusBadRequest)
            return
        }
        
        // Generate presigned upload URL
        uploadURL, err := minioClient.GeneratePresignedPutURL(
            req.Bucket, req.Filename, 1*time.Hour)
        if err != nil {
            http.Error(w, "Failed to generate URL", http.StatusInternalServerError)
            return
        }
        
        json.NewEncoder(w).Encode(map[string]string{"upload_url": uploadURL})
    })
    
    log.Info("File upload service started on :8080")
    http.ListenAndServe(":8080", nil)
}
```

## Environment Variables

All modules support environment variable configuration:

```bash
# Logger
export LOG_LEVEL=info
export LOG_FORMAT=json

# HTTP Client  
export HTTP_CLIENT_TIMEOUT=30s
export HTTP_CLIENT_MAX_RETRIES=3

# JWT
export JWT_SECRET_KEY=your-super-secret-key
export JWT_ACCESS_TOKEN_EXPIRY=15m
export JWT_REFRESH_TOKEN_EXPIRY=7d

# Password
export PASSWORD_MIN_LENGTH=12
export PASSWORD_BCRYPT_COST=14

# MinIO
export MINIO_ENDPOINT=localhost:9000
export MINIO_ACCESS_KEY_ID=minioadmin
export MINIO_SECRET_ACCESS_KEY=minioadmin

# gRPC
export GRPC_CLIENT_HOST=api.example.com
export GRPC_CLIENT_PORT=9000
export GRPC_SERVER_PORT=9090
```

## Testing

Each module includes comprehensive tests. Run tests for all modules:

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run tests for specific module
go test ./logger
go test ./httpclient
go test ./config
```

## Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Make your changes with tests
4. Ensure all tests pass (`go test ./...`)
5. Commit your changes (`git commit -m 'Add amazing feature'`)
6. Push to the branch (`git push origin feature/amazing-feature`)
7. Open a Pull Request

### Development Guidelines

- Write comprehensive tests for new features
- Update documentation for API changes
- Follow Go best practices and conventions
- Ensure backwards compatibility when possible
- Add examples for new functionality

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Support

- **Documentation**: Each module has detailed README with examples
- **Issues**: Report bugs and request features via GitHub Issues
- **Discussions**: Use GitHub Discussions for questions and community support

## Changelog

See [CHANGELOG.md](CHANGELOG.md) for release notes and version history.
