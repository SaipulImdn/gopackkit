# GoPackKit

GoPackKit adalah kumpulan utility modules untuk Go yang dapat digunakan kembali di berbagai project. Module ini menyediakan fungsionalitas umum seperti logging, HTTP client, konfigurasi, validasi, MinIO integration, JWT token management, dan password hashing.

## Installation

```bash
go get github.com/saipulimdn/gopackkit
```

## Modules

- [Logger](#logger) - Multi-backend logging (Simple, Logrus, Zap)
- [HTTP Client](#http-client) - HTTP client dengan retry mechanism
- [Config](#config) - Environment variable loader
- [Validator](#validator) - Data validation dengan security focus
- [MinIO](#minio) - MinIO client untuk object storage
- [JWT](#jwt) - JWT token management
- [Password](#password) - Password hashing dan validation

---

## Logger

Logger module menyediakan multiple backend logging dengan interface yang unified.

### Features

- Multiple backends: Simple, Logrus, Zap
- Configurable log levels dan formats
- JSON dan text output formats
- Thread-safe operations

### Usage

```go
import "github.com/saipulimdn/gopackkit/logger"

// Simple Logger
simpleLogger := logger.NewSimple()
simpleLogger.Info("This is an info message")

// Logrus Logger
logrusLogger := logger.NewLogrus(logger.LogrusConfig{
    Level:  "info",
    Format: "json",
})
logrusLogger.Error("This is an error message")

// Zap Logger
zapLogger := logger.NewZap(logger.ZapConfig{
    Level:       "debug",
    Development: true,
    Encoding:    "console",
})
zapLogger.Debug("This is a debug message")
```

### Configuration

```go
// Logrus Configuration
config := logger.LogrusConfig{
    Level:  "info",        // panic, fatal, error, warn, info, debug, trace
    Format: "json",        // json, text
}

// Zap Configuration
config := logger.ZapConfig{
    Level:       "info",   // debug, info, warn, error, dpanic, panic, fatal
    Development: false,    // true for development mode
    Encoding:    "json",   // json, console
}
```

---

## HTTP Client

HTTP client dengan retry mechanism dan configurable timeout.

### Features

- Automatic retry dengan exponential backoff
- Configurable timeout dan retry attempts
- Request/Response middleware support
- Error handling yang robust

### Usage

```go
import "github.com/saipulimdn/gopackkit/httpclient"

// Default client
client := httpclient.New()

// Custom configuration
config := httpclient.Config{
    Timeout:      30 * time.Second,
    MaxRetries:   3,
    RetryDelay:   1 * time.Second,
    MaxRetryDelay: 10 * time.Second,
}
client := httpclient.NewWithConfig(config)

// Make requests
resp, err := client.Get("https://api.example.com/data")
if err != nil {
    log.Fatal(err)
}
defer resp.Body.Close()

// POST with JSON
data := map[string]interface{}{
    "name": "John Doe",
    "email": "john@example.com",
}
resp, err := client.PostJSON("https://api.example.com/users", data)
```

### Methods

- `Get(url string) (*http.Response, error)`
- `Post(url, contentType string, body io.Reader) (*http.Response, error)`
- `PostJSON(url string, data interface{}) (*http.Response, error)`
- `Put(url, contentType string, body io.Reader) (*http.Response, error)`
- `Delete(url string) (*http.Response, error)`
- `Do(req *http.Request) (*http.Response, error)`

---

## Config

Environment variable loader dengan struct tag support.

### Features

- Automatic environment variable mapping
- Default values support
- Type conversion (string, int, bool, duration)
- Struct tag based configuration

### Usage

```go
import "github.com/saipulimdn/gopackkit/config"

type AppConfig struct {
    Port        int           `env:"PORT" default:"8080"`
    DatabaseURL string        `env:"DATABASE_URL" default:"localhost:5432"`
    Debug       bool          `env:"DEBUG" default:"false"`
    Timeout     time.Duration `env:"TIMEOUT" default:"30s"`
    SecretKey   string        `env:"SECRET_KEY"`
}

func main() {
    var cfg AppConfig
    
    // Load from environment variables
    err := config.Load(&cfg)
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("Server running on port: %d\n", cfg.Port)
    fmt.Printf("Database URL: %s\n", cfg.DatabaseURL)
    fmt.Printf("Debug mode: %t\n", cfg.Debug)
}
```

### Supported Types

- `string`
- `int`, `int8`, `int16`, `int32`, `int64`
- `uint`, `uint8`, `uint16`, `uint32`, `uint64`
- `float32`, `float64`
- `bool`
- `time.Duration`

---

## Validator

Data validation dengan security focus untuk menghindari vulnerability.

### Features

- Secure validation rules (no regex injection)
- Phone number validation (10-15 digits)
- Safe email validation (basic @ check)
- Alphanumeric dan no special chars validation
- Required field validation

### Usage

```go
import "github.com/saipulimdn/gopackkit/validator"

type User struct {
    Name        string `validate:"required"`
    Email       string `validate:"required,email_safe"`
    Phone       string `validate:"required,phone"`
    Username    string `validate:"required,alphanumeric"`
    Description string `validate:"no_special_chars"`
}

func main() {
    v := validator.New()
    
    user := User{
        Name:        "John Doe",
        Email:       "john@example.com",
        Phone:       "1234567890",
        Username:    "johndoe123",
        Description: "Simple description",
    }
    
    err := v.Struct(user)
    if err != nil {
        fmt.Printf("Validation errors: %v\n", err)
    }
}
```

### Validation Rules

- `required` - Field tidak boleh kosong
- `email_safe` - Email format basic (mengandung @)
- `phone` - Nomor telepon 10-15 digit
- `alphanumeric` - Hanya huruf dan angka
- `no_special_chars` - Tidak mengandung karakter special

---

## MinIO

MinIO client untuk object storage operations.

### Features

- Presigned URL generation (GET, PUT, POST)
- Object upload dan download
- Bucket operations
- Configurable expiration times

### Usage

```go
import "github.com/saipulimdn/gopackkit/minio"

// Configuration
config := minio.Config{
    Endpoint:        "localhost:9000",
    AccessKeyID:     "minioadmin",
    SecretAccessKey: "minioadmin",
    UseSSL:          false,
    Region:          "us-east-1",
}

client, err := minio.New(config)
if err != nil {
    log.Fatal(err)
}

// Generate presigned URL for upload
uploadURL, err := client.GeneratePresignedPutURL("my-bucket", "my-object.jpg", 1*time.Hour)
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Upload URL: %s\n", uploadURL)

// Generate presigned URL for download
downloadURL, err := client.GeneratePresignedGetURL("my-bucket", "my-object.jpg", 1*time.Hour)
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Download URL: %s\n", downloadURL)

// Upload object
err = client.UploadObject("my-bucket", "my-object.txt", "./local-file.txt")
if err != nil {
    log.Fatal(err)
}
```

### Methods

- `GeneratePresignedGetURL(bucket, object string, expiry time.Duration) (string, error)`
- `GeneratePresignedPutURL(bucket, object string, expiry time.Duration) (string, error)`
- `GeneratePresignedPostURL(bucket, object string, expiry time.Duration) (map[string]string, error)`
- `UploadObject(bucket, object, filePath string) error`
- `DownloadObject(bucket, object, filePath string) error`
- `DeleteObject(bucket, object string) error`
- `ObjectExists(bucket, object string) (bool, error)`

---

## JWT

JWT token management dengan HMAC signing.

### Features

- Token generation dengan custom claims
- Token validation dan parsing
- Refresh token mechanism
- Configurable expiration times
- HMAC-SHA256 signing

### Usage

```go
import "github.com/saipulimdn/gopackkit/jwt"

// Configuration
config := jwt.Config{
    SecretKey:            "your-secret-key",
    AccessTokenExpiry:    15 * time.Minute,
    RefreshTokenExpiry:   7 * 24 * time.Hour,
    Issuer:              "your-app",
}

jwtManager := jwt.New(config)

// Generate token
claims := jwt.Claims{
    UserID: "user123",
    Email:  "user@example.com",
    Role:   "admin",
}

tokens, err := jwtManager.GenerateTokens(claims)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Access Token: %s\n", tokens.AccessToken)
fmt.Printf("Refresh Token: %s\n", tokens.RefreshToken)

// Validate token
validClaims, err := jwtManager.ValidateAccessToken(tokens.AccessToken)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("User ID: %s\n", validClaims.UserID)
fmt.Printf("Email: %s\n", validClaims.Email)

// Refresh tokens
newTokens, err := jwtManager.RefreshTokens(tokens.RefreshToken)
if err != nil {
    log.Fatal(err)
}
```

### Token Claims

```go
type Claims struct {
    UserID   string                 `json:"user_id"`
    Email    string                 `json:"email"`
    Role     string                 `json:"role"`
    Metadata map[string]interface{} `json:"metadata,omitempty"`
}
```

---

## Password

Password hashing dan validation menggunakan bcrypt.

### Features

- Secure bcrypt hashing
- Password strength validation
- Configurable validation rules
- Random password generation
- Hash migration detection

### Usage

```go
import "github.com/saipulimdn/gopackkit/password"

// Default configuration
pm := password.New()

// Hash password
hashedPassword, err := pm.Hash("mySecurePassword123")
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Hash: %s\n", hashedPassword.Hash)

// Verify password
err = pm.Verify("mySecurePassword123", hashedPassword.Hash)
if err != nil {
    log.Fatal("Password verification failed:", err)
}

fmt.Println("Password verified successfully!")

// Validate password strength
validation := pm.Validate("mySecurePassword123")
fmt.Printf("Valid: %t\n", validation.Valid)
fmt.Printf("Strength: %s\n", validation.Strength)
fmt.Printf("Score: %d\n", validation.Score)

// Generate random password
randomPassword, err := pm.GenerateRandomPassword(16)
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Generated password: %s\n", randomPassword)
```

### Custom Configuration

```go
config := password.Config{
    MinLength:      12,
    MaxLength:      64,
    RequireUpper:   true,
    RequireLower:   true,
    RequireDigit:   true,
    RequireSpecial: true,
    BcryptCost:     14,
}

pm := password.NewWithConfig(config)
```

### Password Strength Levels

- **Weak** (0-2 points): Basic passwords
- **Fair** (3-4 points): Some requirements met
- **Good** (5-6 points): Most requirements met
- **Strong** (7-8 points): All requirements met
- **Very Strong** (9+ points): Excellent passwords

---

## Examples

### Complete Web API Example

```go
package main

import (
    "encoding/json"
    "fmt"
    "net/http"
    "time"
    
    "github.com/saipulimdn/gopackkit/config"
    "github.com/saipulimdn/gopackkit/jwt"
    "github.com/saipulimdn/gopackkit/logger"
    "github.com/saipulimdn/gopackkit/password"
    "github.com/saipulimdn/gopackkit/validator"
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

func main() {
    // Load configuration
    var cfg Config
    if err := config.Load(&cfg); err != nil {
        panic(err)
    }
    
    // Setup logger
    log := logger.NewLogrus(logger.LogrusConfig{
        Level:  cfg.LogLevel,
        Format: "json",
    })
    
    // Setup JWT
    jwtManager := jwt.New(jwt.Config{
        SecretKey:         cfg.JWTSecret,
        AccessTokenExpiry: 15 * time.Minute,
        RefreshTokenExpiry: 7 * 24 * time.Hour,
        Issuer:           "my-app",
    })
    
    // Setup password manager
    pm := password.New()
    
    // Setup validator
    v := validator.New()
    
    // Login handler
    http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
        var req LoginRequest
        if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
            http.Error(w, "Invalid JSON", http.StatusBadRequest)
            return
        }
        
        // Validate request
        if err := v.Struct(req); err != nil {
            http.Error(w, err.Error(), http.StatusBadRequest)
            return
        }
        
        // In real app, get user from database
        storedHash := "$2a$12$..." // from database
        
        // Verify password
        if err := pm.Verify(req.Password, storedHash); err != nil {
            log.Warn("Invalid login attempt", "email", req.Email)
            http.Error(w, "Invalid credentials", http.StatusUnauthorized)
            return
        }
        
        // Generate tokens
        claims := jwt.Claims{
            UserID: "user123",
            Email:  req.Email,
            Role:   "user",
        }
        
        tokens, err := jwtManager.GenerateTokens(claims)
        if err != nil {
            log.Error("Failed to generate tokens", "error", err)
            http.Error(w, "Internal server error", http.StatusInternalServerError)
            return
        }
        
        // Return tokens
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(tokens)
        
        log.Info("User logged in successfully", "email", req.Email)
    })
    
    log.Info("Server starting", "port", cfg.Port)
    http.ListenAndServe(fmt.Sprintf(":%d", cfg.Port), nil)
}
```

## Security Features

### Validator Security
- **No Regex Injection**: Menggunakan character-by-character validation
- **Safe Email Check**: Basic @ validation tanpa complex regex
- **Phone Validation**: Length-based validation (10-15 digits)

### Password Security
- **Bcrypt Hashing**: Industry standard algorithm
- **Configurable Cost**: Adjustable untuk future-proofing
- **Secure Random**: Menggunakan crypto/rand
- **Pattern Detection**: Deteksi weak patterns

### JWT Security
- **HMAC-SHA256**: Secure signing algorithm
- **Configurable Expiry**: Access dan refresh token expiry
- **Token Refresh**: Secure token refresh mechanism

## Contributing

1. Fork the repository
2. Create feature branch (`git checkout -b feature/amazing-feature`)
3. Commit changes (`git commit -m 'Add amazing feature'`)
4. Push to branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Support

Jika ada pertanyaan atau issue, silakan buat issue di GitHub repository ini.
