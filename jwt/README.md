# JWT Module

The JWT module provides comprehensive JSON Web Token management with HMAC-SHA256 signing, token generation, validation, and refresh mechanisms for secure authentication.

## Features

- **Token Generation**: Create access and refresh tokens with custom claims
- **Token Validation**: Secure token parsing and validation
- **Refresh Mechanism**: Secure token refresh workflow
- **Custom Claims**: Flexible claim structure with metadata support
- **HMAC-SHA256 Signing**: Industry-standard token signing
- **Configurable Expiry**: Separate expiry times for access and refresh tokens
- **Security Best Practices**: Built-in security measures and validation

## Installation

```bash
go get github.com/saipulimdn/gopackkit/jwt
```

## Quick Start

### Basic Usage

```go
package main

import (
    "log"
    "time"
    
    "github.com/saipulimdn/gopackkit/jwt"
)

func main() {
    // Create JWT manager with configuration
    config := jwt.Config{
        SecretKey:            "your-super-secret-key-here",
        AccessTokenExpiry:    15 * time.Minute,
        RefreshTokenExpiry:   7 * 24 * time.Hour, // 7 days
        Issuer:              "your-app",
    }
    
    jwtManager := jwt.New(config)
    
    // Create user claims
    claims := jwt.Claims{
        UserID: "user123",
        Email:  "user@example.com",
        Role:   "admin",
        Metadata: map[string]interface{}{
            "department": "engineering",
            "level":      "senior",
        },
    }
    
    // Generate tokens
    tokens, err := jwtManager.GenerateTokens(claims)
    if err != nil {
        log.Fatal("Failed to generate tokens:", err)
    }
    
    log.Printf("Access Token: %s", tokens.AccessToken)
    log.Printf("Refresh Token: %s", tokens.RefreshToken)
    log.Printf("Expires At: %v", tokens.ExpiresAt)
    
    // Validate access token
    validatedClaims, err := jwtManager.ValidateAccessToken(tokens.AccessToken)
    if err != nil {
        log.Fatal("Token validation failed:", err)
    }
    
    log.Printf("Validated User ID: %s", validatedClaims.UserID)
    log.Printf("Validated Email: %s", validatedClaims.Email)
    log.Printf("Validated Role: %s", validatedClaims.Role)
}
```

## Configuration

```go
type Config struct {
    SecretKey            string        // Secret key for signing tokens
    AccessTokenExpiry    time.Duration // Access token expiry duration
    RefreshTokenExpiry   time.Duration // Refresh token expiry duration
    Issuer              string        // Token issuer
}
```

### Default Configuration

```go
config := jwt.Config{
    SecretKey:            "change-this-secret-key",
    AccessTokenExpiry:    15 * time.Minute,
    RefreshTokenExpiry:   7 * 24 * time.Hour,
    Issuer:              "gopackkit",
}
```

### Environment Variables

Configure JWT using environment variables:

```bash
export JWT_SECRET_KEY=your-super-secret-key-here
export JWT_ACCESS_TOKEN_EXPIRY=15m
export JWT_REFRESH_TOKEN_EXPIRY=168h  # 7 days
export JWT_ISSUER=your-app-name
```

## Token Management

### Generate Tokens

```go
package main

import (
    "github.com/saipulimdn/gopackkit/jwt"
)

func generateUserTokens(jwtManager *jwt.Manager, userID, email, role string) (*jwt.TokenPair, error) {
    claims := jwt.Claims{
        UserID: userID,
        Email:  email,
        Role:   role,
        Metadata: map[string]interface{}{
            "login_time": time.Now().Unix(),
            "ip_address": "192.168.1.1",
        },
    }
    
    tokens, err := jwtManager.GenerateTokens(claims)
    if err != nil {
        return nil, err
    }
    
    return tokens, nil
}

func main() {
    config := jwt.Config{
        SecretKey:            "your-secret-key",
        AccessTokenExpiry:    30 * time.Minute,
        RefreshTokenExpiry:   30 * 24 * time.Hour, // 30 days
        Issuer:              "my-app",
    }
    
    jwtManager := jwt.New(config)
    
    // Generate tokens for a user
    tokens, err := generateUserTokens(jwtManager, "user456", "user@example.com", "user")
    if err != nil {
        log.Fatal("Failed to generate tokens:", err)
    }
    
    log.Printf("Generated tokens for user: %+v", tokens)
}
```

### Validate Access Token

```go
package main

import (
    "fmt"
    "github.com/saipulimdn/gopackkit/jwt"
)

func validateAndGetUser(jwtManager *jwt.Manager, tokenString string) (*jwt.Claims, error) {
    // Validate the token
    claims, err := jwtManager.ValidateAccessToken(tokenString)
    if err != nil {
        return nil, fmt.Errorf("invalid token: %w", err)
    }
    
    // Additional validation if needed
    if claims.UserID == "" {
        return nil, fmt.Errorf("invalid user ID in token")
    }
    
    if claims.Email == "" {
        return nil, fmt.Errorf("invalid email in token")
    }
    
    return claims, nil
}

func main() {
    jwtManager := jwt.New(jwt.Config{
        SecretKey: "your-secret-key",
        Issuer:    "my-app",
    })
    
    // Example token (you would get this from request headers)
    tokenString := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
    
    claims, err := validateAndGetUser(jwtManager, tokenString)
    if err != nil {
        log.Printf("Token validation failed: %v", err)
        return
    }
    
    log.Printf("Valid token for user: %s (%s)", claims.UserID, claims.Email)
    log.Printf("User role: %s", claims.Role)
    
    // Access metadata
    if loginTime, ok := claims.Metadata["login_time"]; ok {
        log.Printf("User logged in at: %v", loginTime)
    }
}
```

### Refresh Tokens

```go
package main

import (
    "github.com/saipulimdn/gopackkit/jwt"
)

func refreshUserTokens(jwtManager *jwt.Manager, refreshToken string) (*jwt.TokenPair, error) {
    // Refresh tokens using the refresh token
    newTokens, err := jwtManager.RefreshTokens(refreshToken)
    if err != nil {
        return nil, err
    }
    
    return newTokens, nil
}

func main() {
    jwtManager := jwt.New(jwt.Config{
        SecretKey:            "your-secret-key",
        AccessTokenExpiry:    15 * time.Minute,
        RefreshTokenExpiry:   7 * 24 * time.Hour,
        Issuer:              "my-app",
    })
    
    // Example refresh token (you would get this from client)
    refreshToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
    
    newTokens, err := refreshUserTokens(jwtManager, refreshToken)
    if err != nil {
        log.Printf("Token refresh failed: %v", err)
        return
    }
    
    log.Printf("Refreshed tokens successfully")
    log.Printf("New Access Token: %s", newTokens.AccessToken)
    log.Printf("New Refresh Token: %s", newTokens.RefreshToken)
}
```

## HTTP Integration

### Authentication Middleware

```go
package main

import (
    "context"
    "net/http"
    "strings"
    
    "github.com/saipulimdn/gopackkit/jwt"
)

type contextKey string

const userContextKey contextKey = "user"

func authMiddleware(jwtManager *jwt.Manager) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            // Get Authorization header
            authHeader := r.Header.Get("Authorization")
            if authHeader == "" {
                http.Error(w, "Authorization header required", http.StatusUnauthorized)
                return
            }
            
            // Check Bearer token format
            parts := strings.Split(authHeader, " ")
            if len(parts) != 2 || parts[0] != "Bearer" {
                http.Error(w, "Invalid authorization header format", http.StatusUnauthorized)
                return
            }
            
            tokenString := parts[1]
            
            // Validate token
            claims, err := jwtManager.ValidateAccessToken(tokenString)
            if err != nil {
                http.Error(w, "Invalid token", http.StatusUnauthorized)
                return
            }
            
            // Add user to context
            ctx := context.WithValue(r.Context(), userContextKey, claims)
            next.ServeHTTP(w, r.WithContext(ctx))
        })
    }
}

func getUserFromContext(ctx context.Context) (*jwt.Claims, bool) {
    user, ok := ctx.Value(userContextKey).(*jwt.Claims)
    return user, ok
}

func protectedHandler(w http.ResponseWriter, r *http.Request) {
    user, ok := getUserFromContext(r.Context())
    if !ok {
        http.Error(w, "User not found in context", http.StatusInternalServerError)
        return
    }
    
    response := map[string]interface{}{
        "message": "Protected endpoint accessed successfully",
        "user_id": user.UserID,
        "email":   user.Email,
        "role":    user.Role,
    }
    
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(response)
}

func main() {
    jwtManager := jwt.New(jwt.Config{
        SecretKey: "your-secret-key",
        Issuer:    "my-app",
    })
    
    // Setup routes with middleware
    mux := http.NewServeMux()
    
    // Protected route
    protectedRoute := authMiddleware(jwtManager)(http.HandlerFunc(protectedHandler))
    mux.Handle("/protected", protectedRoute)
    
    log.Println("Server starting on :8080")
    log.Fatal(http.ListenAndServe(":8080", mux))
}
```

### Login Endpoint

```go
package main

import (
    "encoding/json"
    "net/http"
    "time"
    
    "github.com/saipulimdn/gopackkit/jwt"
)

type LoginRequest struct {
    Email    string `json:"email"`
    Password string `json:"password"`
}

type LoginResponse struct {
    AccessToken  string    `json:"access_token"`
    RefreshToken string    `json:"refresh_token"`
    TokenType    string    `json:"token_type"`
    ExpiresAt    time.Time `json:"expires_at"`
    User         UserInfo  `json:"user"`
}

type UserInfo struct {
    ID    string `json:"id"`
    Email string `json:"email"`
    Role  string `json:"role"`
}

func loginHandler(jwtManager *jwt.Manager) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        if r.Method != http.MethodPost {
            http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
            return
        }
        
        var req LoginRequest
        if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
            http.Error(w, "Invalid JSON", http.StatusBadRequest)
            return
        }
        
        // Validate credentials (implement your authentication logic)
        user, err := authenticateUser(req.Email, req.Password)
        if err != nil {
            http.Error(w, "Invalid credentials", http.StatusUnauthorized)
            return
        }
        
        // Create JWT claims
        claims := jwt.Claims{
            UserID: user.ID,
            Email:  user.Email,
            Role:   user.Role,
            Metadata: map[string]interface{}{
                "login_time": time.Now().Unix(),
                "user_agent": r.UserAgent(),
            },
        }
        
        // Generate tokens
        tokens, err := jwtManager.GenerateTokens(claims)
        if err != nil {
            http.Error(w, "Failed to generate tokens", http.StatusInternalServerError)
            return
        }
        
        // Prepare response
        response := LoginResponse{
            AccessToken:  tokens.AccessToken,
            RefreshToken: tokens.RefreshToken,
            TokenType:    "Bearer",
            ExpiresAt:    tokens.ExpiresAt,
            User: UserInfo{
                ID:    user.ID,
                Email: user.Email,
                Role:  user.Role,
            },
        }
        
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(response)
    }
}

// Mock authentication function
func authenticateUser(email, password string) (*UserInfo, error) {
    // Implement your authentication logic here
    // This is just a mock implementation
    if email == "user@example.com" && password == "password123" {
        return &UserInfo{
            ID:    "user123",
            Email: email,
            Role:  "user",
        }, nil
    }
    
    return nil, errors.New("invalid credentials")
}
```

### Token Refresh Endpoint

```go
package main

import (
    "encoding/json"
    "net/http"
    
    "github.com/saipulimdn/gopackkit/jwt"
)

type RefreshRequest struct {
    RefreshToken string `json:"refresh_token"`
}

type RefreshResponse struct {
    AccessToken  string    `json:"access_token"`
    RefreshToken string    `json:"refresh_token"`
    TokenType    string    `json:"token_type"`
    ExpiresAt    time.Time `json:"expires_at"`
}

func refreshHandler(jwtManager *jwt.Manager) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        if r.Method != http.MethodPost {
            http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
            return
        }
        
        var req RefreshRequest
        if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
            http.Error(w, "Invalid JSON", http.StatusBadRequest)
            return
        }
        
        if req.RefreshToken == "" {
            http.Error(w, "Refresh token required", http.StatusBadRequest)
            return
        }
        
        // Refresh tokens
        newTokens, err := jwtManager.RefreshTokens(req.RefreshToken)
        if err != nil {
            http.Error(w, "Invalid refresh token", http.StatusUnauthorized)
            return
        }
        
        // Prepare response
        response := RefreshResponse{
            AccessToken:  newTokens.AccessToken,
            RefreshToken: newTokens.RefreshToken,
            TokenType:    "Bearer",
            ExpiresAt:    newTokens.ExpiresAt,
        }
        
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(response)
    }
}

func main() {
    jwtManager := jwt.New(jwt.Config{
        SecretKey:            "your-secret-key",
        AccessTokenExpiry:    15 * time.Minute,
        RefreshTokenExpiry:   7 * 24 * time.Hour,
        Issuer:              "my-app",
    })
    
    http.HandleFunc("/login", loginHandler(jwtManager))
    http.HandleFunc("/refresh", refreshHandler(jwtManager))
    
    log.Println("Auth server starting on :8080")
    log.Fatal(http.ListenAndServe(":8080", nil))
}
```

## Advanced Usage

### Role-Based Access Control

```go
package main

import (
    "context"
    "net/http"
    
    "github.com/saipulimdn/gopackkit/jwt"
)

func requireRole(jwtManager *jwt.Manager, requiredRole string) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return authMiddleware(jwtManager)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            user, ok := getUserFromContext(r.Context())
            if !ok {
                http.Error(w, "User not found", http.StatusInternalServerError)
                return
            }
            
            if user.Role != requiredRole {
                http.Error(w, "Insufficient permissions", http.StatusForbidden)
                return
            }
            
            next.ServeHTTP(w, r)
        }))
    }
}

func requireAnyRole(jwtManager *jwt.Manager, roles ...string) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return authMiddleware(jwtManager)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            user, ok := getUserFromContext(r.Context())
            if !ok {
                http.Error(w, "User not found", http.StatusInternalServerError)
                return
            }
            
            hasRole := false
            for _, role := range roles {
                if user.Role == role {
                    hasRole = true
                    break
                }
            }
            
            if !hasRole {
                http.Error(w, "Insufficient permissions", http.StatusForbidden)
                return
            }
            
            next.ServeHTTP(w, r)
        }))
    }
}

// Usage example
func setupRoutes(jwtManager *jwt.Manager) *http.ServeMux {
    mux := http.NewServeMux()
    
    // Admin only endpoint
    adminHandler := requireRole(jwtManager, "admin")(http.HandlerFunc(adminOnlyHandler))
    mux.Handle("/admin", adminHandler)
    
    // User or Admin endpoint
    userHandler := requireAnyRole(jwtManager, "user", "admin")(http.HandlerFunc(userHandler))
    mux.Handle("/user", userHandler)
    
    return mux
}
```

### Custom Claims Validation

```go
package main

import (
    "errors"
    "time"
    
    "github.com/saipulimdn/gopackkit/jwt"
)

type CustomClaims struct {
    jwt.Claims
    Department string   `json:"department"`
    Permissions []string `json:"permissions"`
    LastLogin   int64    `json:"last_login"`
}

func validateCustomClaims(jwtManager *jwt.Manager, tokenString string) (*CustomClaims, error) {
    // First validate the standard claims
    standardClaims, err := jwtManager.ValidateAccessToken(tokenString)
    if err != nil {
        return nil, err
    }
    
    // Create custom claims structure
    customClaims := &CustomClaims{
        Claims: *standardClaims,
    }
    
    // Extract custom fields from metadata
    if dept, ok := standardClaims.Metadata["department"].(string); ok {
        customClaims.Department = dept
    }
    
    if perms, ok := standardClaims.Metadata["permissions"].([]interface{}); ok {
        for _, perm := range perms {
            if permStr, ok := perm.(string); ok {
                customClaims.Permissions = append(customClaims.Permissions, permStr)
            }
        }
    }
    
    if lastLogin, ok := standardClaims.Metadata["last_login"].(float64); ok {
        customClaims.LastLogin = int64(lastLogin)
    }
    
    // Custom validation logic
    if customClaims.Department == "" {
        return nil, errors.New("department is required")
    }
    
    if len(customClaims.Permissions) == 0 {
        return nil, errors.New("user must have at least one permission")
    }
    
    // Check if last login is too old (example: 30 days)
    if time.Now().Unix()-customClaims.LastLogin > 30*24*60*60 {
        return nil, errors.New("session expired due to inactivity")
    }
    
    return customClaims, nil
}

func generateCustomTokens(jwtManager *jwt.Manager, userID, email, role, department string, permissions []string) (*jwt.TokenPair, error) {
    claims := jwt.Claims{
        UserID: userID,
        Email:  email,
        Role:   role,
        Metadata: map[string]interface{}{
            "department":  department,
            "permissions": permissions,
            "last_login":  time.Now().Unix(),
        },
    }
    
    return jwtManager.GenerateTokens(claims)
}
```

## Security Best Practices

### Secure Token Storage

```go
package main

import (
    "net/http"
    "time"
)

func setSecureTokenCookie(w http.ResponseWriter, name, token string, expires time.Time) {
    cookie := &http.Cookie{
        Name:     name,
        Value:    token,
        Expires:  expires,
        HttpOnly: true,  // Prevent XSS attacks
        Secure:   true,  // HTTPS only
        SameSite: http.SameSiteStrictMode, // CSRF protection
        Path:     "/",
    }
    
    http.SetCookie(w, cookie)
}

func clearTokenCookie(w http.ResponseWriter, name string) {
    cookie := &http.Cookie{
        Name:     name,
        Value:    "",
        Expires:  time.Now().Add(-time.Hour),
        HttpOnly: true,
        Secure:   true,
        SameSite: http.SameSiteStrictMode,
        Path:     "/",
    }
    
    http.SetCookie(w, cookie)
}

// Usage in login handler
func secureLoginHandler(jwtManager *jwt.Manager) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        // ... authentication logic ...
        
        tokens, err := jwtManager.GenerateTokens(claims)
        if err != nil {
            http.Error(w, "Failed to generate tokens", http.StatusInternalServerError)
            return
        }
        
        // Set secure cookies
        setSecureTokenCookie(w, "access_token", tokens.AccessToken, tokens.ExpiresAt)
        setSecureTokenCookie(w, "refresh_token", tokens.RefreshToken, time.Now().Add(7*24*time.Hour))
        
        // Return response without tokens in body for added security
        response := map[string]interface{}{
            "message": "Login successful",
            "user":    userInfo,
        }
        
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(response)
    }
}
```

### Token Blacklisting

```go
package main

import (
    "sync"
    "time"
    
    "github.com/saipulimdn/gopackkit/jwt"
)

type TokenBlacklist struct {
    tokens map[string]time.Time
    mutex  sync.RWMutex
}

func NewTokenBlacklist() *TokenBlacklist {
    bl := &TokenBlacklist{
        tokens: make(map[string]time.Time),
    }
    
    // Start cleanup goroutine
    go bl.cleanup()
    
    return bl
}

func (bl *TokenBlacklist) Add(tokenID string, expiry time.Time) {
    bl.mutex.Lock()
    defer bl.mutex.Unlock()
    bl.tokens[tokenID] = expiry
}

func (bl *TokenBlacklist) IsBlacklisted(tokenID string) bool {
    bl.mutex.RLock()
    defer bl.mutex.RUnlock()
    
    expiry, exists := bl.tokens[tokenID]
    if !exists {
        return false
    }
    
    // Check if token has expired
    if time.Now().After(expiry) {
        go func() {
            bl.mutex.Lock()
            delete(bl.tokens, tokenID)
            bl.mutex.Unlock()
        }()
        return false
    }
    
    return true
}

func (bl *TokenBlacklist) cleanup() {
    ticker := time.NewTicker(1 * time.Hour)
    defer ticker.Stop()
    
    for range ticker.C {
        bl.mutex.Lock()
        now := time.Now()
        for tokenID, expiry := range bl.tokens {
            if now.After(expiry) {
                delete(bl.tokens, tokenID)
            }
        }
        bl.mutex.Unlock()
    }
}

// Enhanced validation with blacklist check
func validateWithBlacklist(jwtManager *jwt.Manager, blacklist *TokenBlacklist, tokenString string) (*jwt.Claims, error) {
    claims, err := jwtManager.ValidateAccessToken(tokenString)
    if err != nil {
        return nil, err
    }
    
    // Extract token ID from claims (you might need to add this to your Claims struct)
    if tokenID, ok := claims.Metadata["jti"].(string); ok {
        if blacklist.IsBlacklisted(tokenID) {
            return nil, errors.New("token has been revoked")
        }
    }
    
    return claims, nil
}
```

## Testing

### Unit Tests

```go
package main

import (
    "testing"
    "time"
    
    "github.com/saipulimdn/gopackkit/jwt"
)

func TestJWTGeneration(t *testing.T) {
    config := jwt.Config{
        SecretKey:            "test-secret-key",
        AccessTokenExpiry:    15 * time.Minute,
        RefreshTokenExpiry:   24 * time.Hour,
        Issuer:              "test-app",
    }
    
    jwtManager := jwt.New(config)
    
    claims := jwt.Claims{
        UserID: "test-user",
        Email:  "test@example.com",
        Role:   "user",
    }
    
    tokens, err := jwtManager.GenerateTokens(claims)
    if err != nil {
        t.Fatalf("Failed to generate tokens: %v", err)
    }
    
    if tokens.AccessToken == "" {
        t.Error("Access token should not be empty")
    }
    
    if tokens.RefreshToken == "" {
        t.Error("Refresh token should not be empty")
    }
    
    if tokens.ExpiresAt.IsZero() {
        t.Error("ExpiresAt should not be zero")
    }
}

func TestJWTValidation(t *testing.T) {
    config := jwt.Config{
        SecretKey: "test-secret-key",
        Issuer:    "test-app",
    }
    
    jwtManager := jwt.New(config)
    
    claims := jwt.Claims{
        UserID: "test-user",
        Email:  "test@example.com",
        Role:   "admin",
    }
    
    tokens, err := jwtManager.GenerateTokens(claims)
    if err != nil {
        t.Fatalf("Failed to generate tokens: %v", err)
    }
    
    // Validate access token
    validatedClaims, err := jwtManager.ValidateAccessToken(tokens.AccessToken)
    if err != nil {
        t.Fatalf("Failed to validate token: %v", err)
    }
    
    if validatedClaims.UserID != claims.UserID {
        t.Errorf("Expected UserID %s, got %s", claims.UserID, validatedClaims.UserID)
    }
    
    if validatedClaims.Email != claims.Email {
        t.Errorf("Expected Email %s, got %s", claims.Email, validatedClaims.Email)
    }
    
    if validatedClaims.Role != claims.Role {
        t.Errorf("Expected Role %s, got %s", claims.Role, validatedClaims.Role)
    }
}

func TestJWTRefresh(t *testing.T) {
    config := jwt.Config{
        SecretKey:            "test-secret-key",
        AccessTokenExpiry:    15 * time.Minute,
        RefreshTokenExpiry:   24 * time.Hour,
        Issuer:              "test-app",
    }
    
    jwtManager := jwt.New(config)
    
    claims := jwt.Claims{
        UserID: "test-user",
        Email:  "test@example.com",
        Role:   "user",
    }
    
    // Generate initial tokens
    originalTokens, err := jwtManager.GenerateTokens(claims)
    if err != nil {
        t.Fatalf("Failed to generate original tokens: %v", err)
    }
    
    // Refresh tokens
    newTokens, err := jwtManager.RefreshTokens(originalTokens.RefreshToken)
    if err != nil {
        t.Fatalf("Failed to refresh tokens: %v", err)
    }
    
    if newTokens.AccessToken == originalTokens.AccessToken {
        t.Error("New access token should be different from original")
    }
    
    // Validate new access token
    validatedClaims, err := jwtManager.ValidateAccessToken(newTokens.AccessToken)
    if err != nil {
        t.Fatalf("Failed to validate new access token: %v", err)
    }
    
    if validatedClaims.UserID != claims.UserID {
        t.Errorf("Expected UserID %s, got %s", claims.UserID, validatedClaims.UserID)
    }
}

func TestInvalidToken(t *testing.T) {
    config := jwt.Config{
        SecretKey: "test-secret-key",
        Issuer:    "test-app",
    }
    
    jwtManager := jwt.New(config)
    
    // Test with invalid token
    _, err := jwtManager.ValidateAccessToken("invalid.token.here")
    if err == nil {
        t.Error("Expected error for invalid token")
    }
    
    // Test with empty token
    _, err = jwtManager.ValidateAccessToken("")
    if err == nil {
        t.Error("Expected error for empty token")
    }
}
```

## Best Practices

1. **Use Strong Secret Keys**: Use cryptographically secure random keys (at least 256 bits)
2. **Short Access Token Expiry**: Keep access tokens short-lived (15-30 minutes)
3. **Secure Storage**: Store refresh tokens securely (HTTP-only cookies)
4. **Token Rotation**: Rotate refresh tokens on each use
5. **Validate All Claims**: Always validate all token claims including expiry
6. **Use HTTPS**: Always use HTTPS in production
7. **Implement Token Revocation**: Maintain a blacklist for revoked tokens
8. **Monitor Token Usage**: Log and monitor token generation and validation

## Troubleshooting

### Common Issues

1. **Invalid Signature**: Check that the secret key matches between generation and validation
2. **Token Expired**: Verify token expiry times and system clock synchronization
3. **Invalid Issuer**: Ensure issuer claim matches expected value
4. **Malformed Token**: Validate token format and encoding
5. **Clock Skew**: Account for small time differences between systems

### Debug Mode

```go
package main

import (
    "encoding/base64"
    "encoding/json"
    "fmt"
    "strings"
    
    "github.com/saipulimdn/gopackkit/jwt"
)

func debugToken(tokenString string) {
    fmt.Println("=== JWT Debug Information ===")
    
    parts := strings.Split(tokenString, ".")
    if len(parts) != 3 {
        fmt.Println("Invalid token format")
        return
    }
    
    // Decode header
    headerBytes, err := base64.RawURLEncoding.DecodeString(parts[0])
    if err != nil {
        fmt.Printf("Failed to decode header: %v\n", err)
        return
    }
    
    fmt.Printf("Header: %s\n", string(headerBytes))
    
    // Decode payload
    payloadBytes, err := base64.RawURLEncoding.DecodeString(parts[1])
    if err != nil {
        fmt.Printf("Failed to decode payload: %v\n", err)
        return
    }
    
    fmt.Printf("Payload: %s\n", string(payloadBytes))
    
    // Parse claims
    var claims map[string]interface{}
    if err := json.Unmarshal(payloadBytes, &claims); err != nil {
        fmt.Printf("Failed to parse claims: %v\n", err)
        return
    }
    
    fmt.Println("\nClaims:")
    for key, value := range claims {
        fmt.Printf("  %s: %v\n", key, value)
    }
}
```
