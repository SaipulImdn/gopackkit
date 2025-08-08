# GoPackKit

A collection of reusable Go utilities and packages for common development needs.

## Installation

```bash
go get github.com/saipulimdn/gopackkit
```

## Packages

### Logger
Standard logger with multiple backends and configurable output formats.

```go
import "github.com/saipulimdn/gopackkit/logger"

// Simple usage
log := logger.New()
log.Info("Hello World")

// With configuration
config := logger.Config{
    Level:   "info",
    Format:  "json",
    Backend: "zap",
}
log := logger.NewWithConfig(config)
log.Error("Something went wrong", "error", err)
```

### HTTP Client
HTTP client with retry mechanism and middleware support.

```go
import "github.com/saipulimdn/gopackkit/httpclient"

client := httpclient.New()
resp, err := client.Get("https://api.example.com/data")
```

### Configuration
Environment variable and configuration management.

```go
import "github.com/saipulimdn/gopackkit/config"

type AppConfig struct {
    Port     int    `env:"PORT" default:"8080"`
    Database string `env:"DATABASE_URL"`
}

var cfg AppConfig
err := config.Load(&cfg)
```

### Validator
Data validation with struct tag support.

```go
import "github.com/saipulimdn/gopackkit/validator"

type User struct {
    Email string `validate:"required,email"`
    Age   int    `validate:"min=18,max=100"`
}

user := User{Email: "test@example.com", Age: 25}
err := validator.Validate(user)
```

### JWT Token
JWT token generation, validation, and management.

```go
import "github.com/saipulimdn/gopackkit/jwt"

config := jwt.Config{
    SecretKey:       "your-secret-key",
    AccessTokenTTL:  15 * time.Minute,
    RefreshTokenTTL: 7 * 24 * time.Hour,
    Issuer:          "your-app",
}

jwtManager, err := jwt.New(config)

// Generate token pair
tokenPair, err := jwtManager.GenerateTokenPair("user123", "john", "john@example.com", []string{"user"}, nil)

// Validate token
tokenInfo, err := jwtManager.ValidateToken(tokenPair.AccessToken)

// Refresh token
newTokenPair, err := jwtManager.RefreshToken(tokenPair.RefreshToken)
```

### MinIO Client
MinIO client with presigned URL generation and object management.

```go
import "github.com/saipulimdn/gopackkit/minio"

config := minio.Config{
    Endpoint:        "localhost:9000",
    AccessKeyID:     "minioadmin",
    SecretAccessKey: "minioadmin",
    UseSSL:          false,
}

client, err := minio.New(config)

// Generate presigned GET URL
getURL, err := client.GetPresignedURL(ctx, "bucket", "object.jpg", nil)

// Generate presigned PUT URL  
putURL, err := client.PutPresignedURL(ctx, "bucket", "object.jpg", 
    minio.CustomPresignedURLOptions(time.Hour))

// Delete object
err = client.DeleteObject(ctx, "bucket", "object.jpg", nil)
```

## License

MIT License
