# Password Module

The Password module provides secure password hashing, validation, and management using bcrypt with comprehensive password strength checking and security best practices.

## Features

- **Secure Bcrypt Hashing**: Industry-standard password hashing with configurable cost
- **Password Strength Validation**: Comprehensive strength checking with scoring system
- **Random Password Generation**: Cryptographically secure password generation
- **Hash Migration Detection**: Detect and handle different bcrypt costs
- **Configurable Rules**: Customizable validation rules for different security requirements
- **Security Best Practices**: Built-in protection against common password vulnerabilities

## Installation

```bash
go get github.com/saipulimdn/gopackkit/password
```

## Quick Start

### Basic Usage

```go
package main

import (
    "log"
    
    "github.com/saipulimdn/gopackkit/password"
)

func main() {
    // Create password manager with default configuration
    pm := password.New()
    
    // Hash a password
    plainPassword := "MySecurePassword123!"
    hashedPassword, err := pm.Hash(plainPassword)
    if err != nil {
        log.Fatal("Failed to hash password:", err)
    }
    
    log.Printf("Original: %s", plainPassword)
    log.Printf("Hashed: %s", hashedPassword.Hash)
    log.Printf("Cost: %d", hashedPassword.Cost)
    
    // Verify password
    err = pm.Verify(plainPassword, hashedPassword.Hash)
    if err != nil {
        log.Fatal("Password verification failed:", err)
    }
    
    log.Println("Password verified successfully!")
    
    // Validate password strength
    validation := pm.Validate(plainPassword)
    log.Printf("Password strength: %s (Score: %d/10)", validation.Strength, validation.Score)
    log.Printf("Valid: %t", validation.Valid)
    
    if len(validation.Issues) > 0 {
        log.Println("Issues found:")
        for _, issue := range validation.Issues {
            log.Printf("  - %s", issue)
        }
    }
}
```

### Custom Configuration

```go
package main

import (
    "github.com/saipulimdn/gopackkit/password"
)

func main() {
    // Create custom configuration
    config := password.Config{
        MinLength:      12,   // Minimum 12 characters
        MaxLength:      128,  // Maximum 128 characters
        RequireUpper:   true, // Require uppercase letters
        RequireLower:   true, // Require lowercase letters
        RequireDigit:   true, // Require digits
        RequireSpecial: true, // Require special characters
        BcryptCost:     14,   // Higher cost for better security
    }
    
    pm := password.NewWithConfig(config)
    
    // Test password with custom rules
    testPassword := "MyVerySecurePassword123!"
    
    validation := pm.Validate(testPassword)
    log.Printf("Password validation result: %+v", validation)
    
    if validation.Valid {
        hashedPassword, err := pm.Hash(testPassword)
        if err != nil {
            log.Fatal("Failed to hash password:", err)
        }
        
        log.Printf("Password hashed successfully with cost %d", hashedPassword.Cost)
    }
}
```

## Configuration

```go
type Config struct {
    MinLength      int  // Minimum password length
    MaxLength      int  // Maximum password length
    RequireUpper   bool // Require uppercase letters
    RequireLower   bool // Require lowercase letters
    RequireDigit   bool // Require numeric digits
    RequireSpecial bool // Require special characters
    BcryptCost     int  // Bcrypt cost factor (4-31)
}
```

### Default Configuration

```go
config := password.Config{
    MinLength:      8,
    MaxLength:      64,
    RequireUpper:   true,
    RequireLower:   true,
    RequireDigit:   true,
    RequireSpecial: false,
    BcryptCost:     12,
}
```

### Environment Variables

Configure password requirements using environment variables:

```bash
export PASSWORD_MIN_LENGTH=12
export PASSWORD_MAX_LENGTH=128
export PASSWORD_REQUIRE_UPPER=true
export PASSWORD_REQUIRE_LOWER=true
export PASSWORD_REQUIRE_DIGIT=true
export PASSWORD_REQUIRE_SPECIAL=true
export PASSWORD_BCRYPT_COST=14
```

## Password Operations

### Hash Password

```go
package main

import (
    "github.com/saipulimdn/gopackkit/password"
)

func hashUserPassword(pm *password.Manager, plainPassword string) (*password.HashedPassword, error) {
    // Validate password before hashing
    validation := pm.Validate(plainPassword)
    if !validation.Valid {
        return nil, fmt.Errorf("password validation failed: %v", validation.Issues)
    }
    
    // Hash the password
    hashedPassword, err := pm.Hash(plainPassword)
    if err != nil {
        return nil, fmt.Errorf("failed to hash password: %w", err)
    }
    
    log.Printf("Password hashed successfully")
    log.Printf("  Hash: %s", hashedPassword.Hash)
    log.Printf("  Cost: %d", hashedPassword.Cost)
    log.Printf("  Created: %v", hashedPassword.CreatedAt)
    
    return hashedPassword, nil
}

func main() {
    pm := password.New()
    
    userPassword := "UserSecurePassword123!"
    
    hashedPassword, err := hashUserPassword(pm, userPassword)
    if err != nil {
        log.Fatal("Password hashing failed:", err)
    }
    
    // Store hashedPassword.Hash in your database
    fmt.Printf("Store this hash in database: %s\n", hashedPassword.Hash)
}
```

### Verify Password

```go
package main

import (
    "github.com/saipulimdn/gopackkit/password"
)

func authenticateUser(pm *password.Manager, inputPassword, storedHash string) error {
    // Verify password against stored hash
    err := pm.Verify(inputPassword, storedHash)
    if err != nil {
        log.Printf("Authentication failed for user: %v", err)
        return fmt.Errorf("invalid credentials")
    }
    
    log.Println("User authenticated successfully")
    
    // Check if password needs rehashing (cost upgrade)
    if pm.NeedsRehash(storedHash) {
        log.Println("Password hash needs upgrade")
        
        // Rehash with current cost
        newHash, err := pm.Hash(inputPassword)
        if err != nil {
            log.Printf("Failed to rehash password: %v", err)
        } else {
            log.Printf("New hash generated: %s", newHash.Hash)
            // Update database with newHash.Hash
        }
    }
    
    return nil
}

func main() {
    pm := password.New()
    
    // Simulate stored hash from database
    storedHash := "$2a$12$LQv3c1yqBWVHxkd0LHAkCOYz6TtxMQJqhN8/LewdBPj1xq/DAMOvu"
    inputPassword := "UserSecurePassword123!"
    
    err := authenticateUser(pm, inputPassword, storedHash)
    if err != nil {
        log.Fatal("Authentication failed:", err)
    }
    
    log.Println("User login successful")
}
```

### Validate Password Strength

```go
package main

import (
    "github.com/saipulimdn/gopackkit/password"
)

func checkPasswordStrength(pm *password.Manager, password string) {
    validation := pm.Validate(password)
    
    fmt.Printf("Password: %s\n", password)
    fmt.Printf("Valid: %t\n", validation.Valid)
    fmt.Printf("Strength: %s\n", validation.Strength)
    fmt.Printf("Score: %d/10\n", validation.Score)
    
    if len(validation.Issues) > 0 {
        fmt.Println("Issues:")
        for _, issue := range validation.Issues {
            fmt.Printf("  - %s\n", issue)
        }
    }
    
    fmt.Println("Requirements met:")
    for requirement, met := range validation.Requirements {
        status := "✗"
        if met {
            status = "✓"
        }
        fmt.Printf("  %s %s\n", status, requirement)
    }
    fmt.Println()
}

func main() {
    pm := password.New()
    
    // Test various passwords
    passwords := []string{
        "123456",                    // Very weak
        "password",                  // Weak
        "Password123",               // Fair
        "MySecurePassword123",       // Good
        "MyV3ryS3cur3P@ssw0rd!",    // Strong
        "C0mpl3x!P@ssw0rd#W1th$p3c1@lCh@rs&N0s", // Very strong
    }
    
    for _, pwd := range passwords {
        checkPasswordStrength(pm, pwd)
    }
}
```

### Generate Random Password

```go
package main

import (
    "github.com/saipulimdn/gopackkit/password"
)

func generateSecurePassword(pm *password.Manager, length int) (string, error) {
    // Generate random password
    randomPassword, err := pm.GenerateRandomPassword(length)
    if err != nil {
        return "", fmt.Errorf("failed to generate password: %w", err)
    }
    
    // Validate generated password
    validation := pm.Validate(randomPassword)
    if !validation.Valid {
        // This should rarely happen with properly generated passwords
        log.Printf("Generated password failed validation: %v", validation.Issues)
        // Try again or adjust generation logic
        return generateSecurePassword(pm, length)
    }
    
    log.Printf("Generated password: %s", randomPassword)
    log.Printf("Strength: %s (Score: %d)", validation.Strength, validation.Score)
    
    return randomPassword, nil
}

func generatePasswordsForUsers(pm *password.Manager) {
    // Generate passwords for different use cases
    useCases := []struct {
        name   string
        length int
    }{
        {"Temporary Password", 12},
        {"Strong User Password", 16},
        {"Admin Password", 20},
        {"API Key Password", 32},
    }
    
    for _, useCase := range useCases {
        password, err := generateSecurePassword(pm, useCase.length)
        if err != nil {
            log.Printf("Failed to generate %s: %v", useCase.name, err)
            continue
        }
        
        fmt.Printf("%s (%d chars): %s\n", useCase.name, useCase.length, password)
    }
}

func main() {
    pm := password.New()
    
    generatePasswordsForUsers(pm)
}
```

## HTTP Integration

### User Registration Handler

```go
package main

import (
    "encoding/json"
    "net/http"
    
    "github.com/saipulimdn/gopackkit/password"
)

type RegisterRequest struct {
    Email    string `json:"email"`
    Password string `json:"password"`
    Name     string `json:"name"`
}

type RegisterResponse struct {
    Success bool   `json:"success"`
    Message string `json:"message"`
    UserID  string `json:"user_id,omitempty"`
}

type PasswordError struct {
    Field   string   `json:"field"`
    Issues  []string `json:"issues"`
    Score   int      `json:"score"`
    Strength string  `json:"strength"`
}

func registerHandler(pm *password.Manager) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        if r.Method != http.MethodPost {
            http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
            return
        }
        
        var req RegisterRequest
        if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
            http.Error(w, "Invalid JSON", http.StatusBadRequest)
            return
        }
        
        // Validate password strength
        validation := pm.Validate(req.Password)
        if !validation.Valid {
            passwordError := PasswordError{
                Field:    "password",
                Issues:   validation.Issues,
                Score:    validation.Score,
                Strength: validation.Strength,
            }
            
            w.WriteHeader(http.StatusBadRequest)
            json.NewEncoder(w).Encode(map[string]interface{}{
                "success": false,
                "error":   "Password does not meet requirements",
                "details": passwordError,
            })
            return
        }
        
        // Hash password
        hashedPassword, err := pm.Hash(req.Password)
        if err != nil {
            log.Printf("Failed to hash password: %v", err)
            http.Error(w, "Internal server error", http.StatusInternalServerError)
            return
        }
        
        // Save user to database (mock implementation)
        userID := saveUser(req.Email, req.Name, hashedPassword.Hash)
        
        // Return success response
        response := RegisterResponse{
            Success: true,
            Message: "User registered successfully",
            UserID:  userID,
        }
        
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusCreated)
        json.NewEncoder(w).Encode(response)
    }
}

// Mock function to save user
func saveUser(email, name, passwordHash string) string {
    // In real application, save to database
    log.Printf("Saving user: %s, %s, %s", email, name, passwordHash)
    return "user_" + email // Mock user ID
}
```

### Login Handler with Password Verification

```go
package main

import (
    "encoding/json"
    "net/http"
    "time"
    
    "github.com/saipulimdn/gopackkit/password"
)

type LoginRequest struct {
    Email    string `json:"email"`
    Password string `json:"password"`
}

type LoginResponse struct {
    Success   bool   `json:"success"`
    Message   string `json:"message"`
    Token     string `json:"token,omitempty"`
    ExpiresAt string `json:"expires_at,omitempty"`
}

type User struct {
    ID           string
    Email        string
    Name         string
    PasswordHash string
}

func loginHandler(pm *password.Manager) http.HandlerFunc {
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
        
        // Get user from database (mock implementation)
        user, err := getUserByEmail(req.Email)
        if err != nil {
            w.WriteHeader(http.StatusUnauthorized)
            json.NewEncoder(w).Encode(LoginResponse{
                Success: false,
                Message: "Invalid credentials",
            })
            return
        }
        
        // Verify password
        err = pm.Verify(req.Password, user.PasswordHash)
        if err != nil {
            log.Printf("Password verification failed for user %s: %v", req.Email, err)
            w.WriteHeader(http.StatusUnauthorized)
            json.NewEncoder(w).Encode(LoginResponse{
                Success: false,
                Message: "Invalid credentials",
            })
            return
        }
        
        // Check if password needs rehashing
        if pm.NeedsRehash(user.PasswordHash) {
            log.Printf("Password needs rehashing for user %s", req.Email)
            go func() {
                // Rehash password in background
                newHash, err := pm.Hash(req.Password)
                if err != nil {
                    log.Printf("Failed to rehash password for user %s: %v", req.Email, err)
                    return
                }
                
                // Update user's password hash in database
                updateUserPasswordHash(user.ID, newHash.Hash)
            }()
        }
        
        // Generate session token (mock implementation)
        token, expiresAt := generateSessionToken(user.ID)
        
        response := LoginResponse{
            Success:   true,
            Message:   "Login successful",
            Token:     token,
            ExpiresAt: expiresAt.Format(time.RFC3339),
        }
        
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(response)
    }
}

// Mock functions
func getUserByEmail(email string) (*User, error) {
    // In real application, query database
    if email == "user@example.com" {
        return &User{
            ID:           "user123",
            Email:        email,
            Name:         "Test User",
            PasswordHash: "$2a$12$LQv3c1yqBWVHxkd0LHAkCOYz6TtxMQJqhN8/LewdBPj1xq/DAMOvu",
        }, nil
    }
    return nil, errors.New("user not found")
}

func updateUserPasswordHash(userID, newHash string) {
    log.Printf("Updated password hash for user %s", userID)
}

func generateSessionToken(userID string) (string, time.Time) {
    // Mock token generation
    expiresAt := time.Now().Add(24 * time.Hour)
    return "mock_token_" + userID, expiresAt
}
```

### Password Change Handler

```go
package main

import (
    "encoding/json"
    "net/http"
    
    "github.com/saipulimdn/gopackkit/password"
)

type ChangePasswordRequest struct {
    CurrentPassword string `json:"current_password"`
    NewPassword     string `json:"new_password"`
}

type ChangePasswordResponse struct {
    Success bool   `json:"success"`
    Message string `json:"message"`
}

func changePasswordHandler(pm *password.Manager) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        if r.Method != http.MethodPost {
            http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
            return
        }
        
        // Get user from authentication context (simplified)
        userID := r.Header.Get("X-User-ID") // In real app, get from JWT token
        if userID == "" {
            http.Error(w, "Unauthorized", http.StatusUnauthorized)
            return
        }
        
        var req ChangePasswordRequest
        if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
            http.Error(w, "Invalid JSON", http.StatusBadRequest)
            return
        }
        
        // Get current user
        user, err := getUserByID(userID)
        if err != nil {
            http.Error(w, "User not found", http.StatusNotFound)
            return
        }
        
        // Verify current password
        err = pm.Verify(req.CurrentPassword, user.PasswordHash)
        if err != nil {
            w.WriteHeader(http.StatusBadRequest)
            json.NewEncoder(w).Encode(ChangePasswordResponse{
                Success: false,
                Message: "Current password is incorrect",
            })
            return
        }
        
        // Validate new password
        validation := pm.Validate(req.NewPassword)
        if !validation.Valid {
            w.WriteHeader(http.StatusBadRequest)
            json.NewEncoder(w).Encode(map[string]interface{}{
                "success": false,
                "message": "New password does not meet requirements",
                "issues":  validation.Issues,
                "score":   validation.Score,
            })
            return
        }
        
        // Check if new password is different from current
        err = pm.Verify(req.NewPassword, user.PasswordHash)
        if err == nil {
            w.WriteHeader(http.StatusBadRequest)
            json.NewEncoder(w).Encode(ChangePasswordResponse{
                Success: false,
                Message: "New password must be different from current password",
            })
            return
        }
        
        // Hash new password
        newHash, err := pm.Hash(req.NewPassword)
        if err != nil {
            log.Printf("Failed to hash new password: %v", err)
            http.Error(w, "Internal server error", http.StatusInternalServerError)
            return
        }
        
        // Update password in database
        err = updateUserPasswordHash(userID, newHash.Hash)
        if err != nil {
            log.Printf("Failed to update password in database: %v", err)
            http.Error(w, "Internal server error", http.StatusInternalServerError)
            return
        }
        
        // Return success response
        response := ChangePasswordResponse{
            Success: true,
            Message: "Password changed successfully",
        }
        
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(response)
        
        log.Printf("Password changed successfully for user %s", userID)
    }
}

func getUserByID(userID string) (*User, error) {
    // Mock implementation
    return &User{
        ID:           userID,
        Email:        "user@example.com",
        PasswordHash: "$2a$12$LQv3c1yqBWVHxkd0LHAkCOYz6TtxMQJqhN8/LewdBPj1xq/DAMOvu",
    }, nil
}
```

## Password Strength Levels

The module provides five levels of password strength:

### Weak (0-2 points)
- Basic passwords that don't meet most requirements
- Examples: "123456", "password", "qwerty"

### Fair (3-4 points)
- Passwords that meet some requirements but lack complexity
- Examples: "Password1", "hello123"

### Good (5-6 points)
- Passwords that meet most requirements
- Examples: "MyPassword123", "SecurePass1"

### Strong (7-8 points)
- Passwords that meet all requirements
- Examples: "MySecurePassword123", "Str0ngP@ssw0rd"

### Very Strong (9+ points)
- Excellent passwords with high complexity
- Examples: "MyV3ryS3cur3P@ssw0rd!", "C0mpl3x!P@ssw0rd#2023"

## Advanced Usage

### Password Policy Enforcement

```go
package main

import (
    "errors"
    "time"
    
    "github.com/saipulimdn/gopackkit/password"
)

type PasswordPolicy struct {
    pm                *password.Manager
    minScore          int
    preventReuse      bool
    maxAge            time.Duration
    passwordHistory   map[string][]string // userID -> password hashes
    passwordAges      map[string]time.Time // userID -> last change time
}

func NewPasswordPolicy(config password.Config) *PasswordPolicy {
    return &PasswordPolicy{
        pm:              password.NewWithConfig(config),
        minScore:        6, // Require "Good" or better
        preventReuse:    true,
        maxAge:          90 * 24 * time.Hour, // 90 days
        passwordHistory: make(map[string][]string),
        passwordAges:    make(map[string]time.Time),
    }
}

func (pp *PasswordPolicy) ValidatePassword(userID, newPassword string) error {
    // Check password strength
    validation := pp.pm.Validate(newPassword)
    if !validation.Valid {
        return fmt.Errorf("password validation failed: %v", validation.Issues)
    }
    
    if validation.Score < pp.minScore {
        return fmt.Errorf("password too weak (score: %d, required: %d)", validation.Score, pp.minScore)
    }
    
    // Check password reuse
    if pp.preventReuse {
        history, exists := pp.passwordHistory[userID]
        if exists {
            for _, oldHash := range history {
                if pp.pm.Verify(newPassword, oldHash) == nil {
                    return errors.New("password has been used recently")
                }
            }
        }
    }
    
    return nil
}

func (pp *PasswordPolicy) SetPassword(userID, newPassword string) error {
    // Validate against policy
    if err := pp.ValidatePassword(userID, newPassword); err != nil {
        return err
    }
    
    // Hash password
    hashedPassword, err := pp.pm.Hash(newPassword)
    if err != nil {
        return err
    }
    
    // Update password history
    if pp.preventReuse {
        history := pp.passwordHistory[userID]
        history = append(history, hashedPassword.Hash)
        
        // Keep only last 5 passwords
        if len(history) > 5 {
            history = history[len(history)-5:]
        }
        
        pp.passwordHistory[userID] = history
    }
    
    // Update password age
    pp.passwordAges[userID] = time.Now()
    
    return nil
}

func (pp *PasswordPolicy) IsPasswordExpired(userID string) bool {
    lastChange, exists := pp.passwordAges[userID]
    if !exists {
        return true // No password set
    }
    
    return time.Since(lastChange) > pp.maxAge
}

func (pp *PasswordPolicy) GetPasswordAge(userID string) time.Duration {
    lastChange, exists := pp.passwordAges[userID]
    if !exists {
        return 0
    }
    
    return time.Since(lastChange)
}
```

### Bulk Password Operations

```go
package main

import (
    "sync"
    
    "github.com/saipulimdn/gopackkit/password"
)

type BulkPasswordProcessor struct {
    pm          *password.Manager
    concurrency int
}

func NewBulkPasswordProcessor(config password.Config, concurrency int) *BulkPasswordProcessor {
    return &BulkPasswordProcessor{
        pm:          password.NewWithConfig(config),
        concurrency: concurrency,
    }
}

type PasswordJob struct {
    UserID      string
    Password    string
    Result      chan PasswordResult
}

type PasswordResult struct {
    UserID string
    Hash   string
    Error  error
}

func (bpp *BulkPasswordProcessor) ProcessPasswords(jobs []PasswordJob) {
    jobChan := make(chan PasswordJob, len(jobs))
    
    // Start workers
    var wg sync.WaitGroup
    for i := 0; i < bpp.concurrency; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            for job := range jobChan {
                result := PasswordResult{UserID: job.UserID}
                
                // Validate password
                validation := bpp.pm.Validate(job.Password)
                if !validation.Valid {
                    result.Error = fmt.Errorf("validation failed: %v", validation.Issues)
                } else {
                    // Hash password
                    hashedPassword, err := bpp.pm.Hash(job.Password)
                    if err != nil {
                        result.Error = err
                    } else {
                        result.Hash = hashedPassword.Hash
                    }
                }
                
                job.Result <- result
            }
        }()
    }
    
    // Send jobs
    for _, job := range jobs {
        jobChan <- job
    }
    close(jobChan)
    
    // Wait for completion
    wg.Wait()
}

// Usage example
func bulkPasswordExample() {
    processor := NewBulkPasswordProcessor(password.Config{
        BcryptCost: 12,
    }, 4) // 4 concurrent workers
    
    // Create jobs
    jobs := []PasswordJob{
        {UserID: "user1", Password: "SecurePassword123!", Result: make(chan PasswordResult, 1)},
        {UserID: "user2", Password: "AnotherPassword456!", Result: make(chan PasswordResult, 1)},
        {UserID: "user3", Password: "ThirdPassword789!", Result: make(chan PasswordResult, 1)},
    }
    
    // Process passwords
    go processor.ProcessPasswords(jobs)
    
    // Collect results
    for _, job := range jobs {
        result := <-job.Result
        if result.Error != nil {
            log.Printf("Failed to process password for %s: %v", result.UserID, result.Error)
        } else {
            log.Printf("Password hashed for %s: %s", result.UserID, result.Hash)
        }
    }
}
```

## Testing

### Unit Tests

```go
package main

import (
    "strings"
    "testing"
    
    "github.com/saipulimdn/gopackkit/password"
)

func TestPasswordHashing(t *testing.T) {
    pm := password.New()
    
    testPassword := "TestPassword123!"
    
    // Test hashing
    hashedPassword, err := pm.Hash(testPassword)
    if err != nil {
        t.Fatalf("Failed to hash password: %v", err)
    }
    
    if hashedPassword.Hash == "" {
        t.Error("Hash should not be empty")
    }
    
    if hashedPassword.Cost == 0 {
        t.Error("Cost should not be zero")
    }
    
    // Test verification
    err = pm.Verify(testPassword, hashedPassword.Hash)
    if err != nil {
        t.Errorf("Password verification failed: %v", err)
    }
    
    // Test wrong password
    err = pm.Verify("WrongPassword", hashedPassword.Hash)
    if err == nil {
        t.Error("Expected error for wrong password")
    }
}

func TestPasswordValidation(t *testing.T) {
    pm := password.New()
    
    testCases := []struct {
        password      string
        expectedValid bool
        expectedScore int
    }{
        {"123456", false, 1},
        {"password", false, 2},
        {"Password123", true, 5},
        {"MySecurePassword123!", true, 8},
        {"", false, 0},
        {"a", false, 1},
    }
    
    for _, tc := range testCases {
        validation := pm.Validate(tc.password)
        
        if validation.Valid != tc.expectedValid {
            t.Errorf("Password %s: expected valid=%t, got valid=%t", 
                tc.password, tc.expectedValid, validation.Valid)
        }
        
        if validation.Score < tc.expectedScore-1 || validation.Score > tc.expectedScore+1 {
            t.Errorf("Password %s: expected score around %d, got %d", 
                tc.password, tc.expectedScore, validation.Score)
        }
    }
}

func TestRandomPasswordGeneration(t *testing.T) {
    pm := password.New()
    
    lengths := []int{12, 16, 20, 32}
    
    for _, length := range lengths {
        password, err := pm.GenerateRandomPassword(length)
        if err != nil {
            t.Errorf("Failed to generate password of length %d: %v", length, err)
            continue
        }
        
        if len(password) != length {
            t.Errorf("Expected password length %d, got %d", length, len(password))
        }
        
        // Validate generated password
        validation := pm.Validate(password)
        if !validation.Valid {
            t.Errorf("Generated password failed validation: %v", validation.Issues)
        }
    }
}

func TestPasswordStrengthLevels(t *testing.T) {
    pm := password.New()
    
    strengthTests := []struct {
        password         string
        expectedStrength string
    }{
        {"123456", "Weak"},
        {"password123", "Fair"},
        {"Password123", "Good"},
        {"MySecurePassword123", "Strong"},
        {"MyV3ryS3cur3P@ssw0rd!", "Very Strong"},
    }
    
    for _, test := range strengthTests {
        validation := pm.Validate(test.password)
        if !strings.Contains(validation.Strength, test.expectedStrength) {
            t.Errorf("Password %s: expected strength %s, got %s", 
                test.password, test.expectedStrength, validation.Strength)
        }
    }
}
```

## Best Practices

1. **Use Strong Bcrypt Cost**: Use cost 12 or higher for production
2. **Validate Before Hashing**: Always validate password strength before hashing
3. **Handle Errors Properly**: Provide helpful error messages without exposing security details
4. **Implement Password Policies**: Use comprehensive password policies for enterprise applications
5. **Monitor for Breaches**: Regularly check passwords against known breach databases
6. **Use HTTPS**: Always transmit passwords over secure connections
7. **Rate Limiting**: Implement rate limiting for authentication attempts
8. **Secure Storage**: Store hashed passwords securely in your database

## Troubleshooting

### Common Issues

1. **High CPU Usage**: Reduce bcrypt cost for better performance vs security trade-off
2. **Validation Too Strict**: Adjust password requirements based on user feedback
3. **Random Generation Fails**: Ensure system has sufficient entropy
4. **Hash Verification Slow**: Consider implementing caching for frequently accessed passwords
5. **Memory Usage**: Monitor memory usage with high bcrypt costs

### Performance Tuning

```go
package main

import (
    "time"
    
    "github.com/saipulimdn/gopackkit/password"
)

func benchmarkBcryptCost() {
    costs := []int{10, 11, 12, 13, 14}
    testPassword := "BenchmarkPassword123!"
    
    for _, cost := range costs {
        config := password.Config{BcryptCost: cost}
        pm := password.NewWithConfig(config)
        
        start := time.Now()
        _, err := pm.Hash(testPassword)
        duration := time.Since(start)
        
        if err != nil {
            log.Printf("Cost %d failed: %v", cost, err)
            continue
        }
        
        log.Printf("Cost %d: %v", cost, duration)
    }
}
```
