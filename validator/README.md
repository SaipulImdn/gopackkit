# Validator Module

The Validator module provides secure data validation with a focus on preventing common security vulnerabilities. It offers safe validation rules without regex injection risks and includes specialized validators for common data types.

## Features

- **Security-First Design**: No regex injection vulnerabilities
- **Safe Email Validation**: Basic email checking without complex regex
- **Phone Number Validation**: Length-based validation (10-15 digits)
- **Alphanumeric Validation**: Character-by-character validation
- **Special Character Detection**: Safe special character filtering
- **Required Field Validation**: Comprehensive null/empty checking
- **Struct Tag Support**: Easy integration with Go structs

## Installation

```bash
go get github.com/saipulimdn/gopackkit/validator
```

## Quick Start

### Basic Usage

```go
package main

import (
    "fmt"
    "log"
    
    "github.com/saipulimdn/gopackkit/validator"
)

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
        Description: "Simple description without special chars",
    }
    
    err := v.Struct(user)
    if err != nil {
        log.Printf("Validation errors: %v", err)
        return
    }
    
    fmt.Println("Validation successful!")
}
```

### Individual Field Validation

```go
package main

import (
    "fmt"
    "github.com/saipulimdn/gopackkit/validator"
)

func main() {
    v := validator.New()
    
    // Validate individual fields
    email := "user@example.com"
    if err := v.Var(email, "required,email_safe"); err != nil {
        fmt.Printf("Email validation failed: %v\n", err)
    } else {
        fmt.Println("Email is valid")
    }
    
    phone := "1234567890"
    if err := v.Var(phone, "required,phone"); err != nil {
        fmt.Printf("Phone validation failed: %v\n", err)
    } else {
        fmt.Println("Phone is valid")
    }
    
    username := "user123"
    if err := v.Var(username, "required,alphanumeric"); err != nil {
        fmt.Printf("Username validation failed: %v\n", err)
    } else {
        fmt.Println("Username is valid")
    }
}
```

## Validation Rules

### Required Validation

```go
type User struct {
    Name  string `validate:"required"`
    Email string `validate:"required"`
    Age   int    `validate:"required"` // Checks for zero value
}

// Examples
user1 := User{Name: "John", Email: "john@example.com", Age: 25} // Valid
user2 := User{Name: "", Email: "john@example.com", Age: 25}    // Invalid: Name is empty
user3 := User{Name: "John", Email: "", Age: 25}                // Invalid: Email is empty
user4 := User{Name: "John", Email: "john@example.com", Age: 0} // Invalid: Age is zero
```

### Email Safe Validation

Safe email validation that checks for the presence of `@` symbol without complex regex:

```go
type Contact struct {
    Email string `validate:"required,email_safe"`
}

// Valid emails
validEmails := []string{
    "user@example.com",
    "test.email@domain.org",
    "user+tag@example.co.uk",
    "user123@test-domain.com",
}

// Invalid emails
invalidEmails := []string{
    "invalid-email",     // No @ symbol
    "@example.com",      // Missing local part
    "user@",             // Missing domain
    "",                  // Empty string
    "user@@example.com", // Multiple @ symbols
}
```

### Phone Validation

Validates phone numbers with 10-15 digits:

```go
type Contact struct {
    Phone string `validate:"required,phone"`
}

// Valid phone numbers
validPhones := []string{
    "1234567890",      // 10 digits
    "12345678901",     // 11 digits
    "123456789012345", // 15 digits
}

// Invalid phone numbers
invalidPhones := []string{
    "123456789",       // Too short (9 digits)
    "1234567890123456", // Too long (16 digits)
    "12345abcde",      // Contains letters
    "123-456-7890",    // Contains special chars
    "",                // Empty string
}
```

### Alphanumeric Validation

Checks if string contains only letters and numbers:

```go
type User struct {
    Username string `validate:"required,alphanumeric"`
}

// Valid usernames
validUsernames := []string{
    "user123",
    "JohnDoe",
    "testuser",
    "User123ABC",
}

// Invalid usernames
invalidUsernames := []string{
    "user_123",     // Contains underscore
    "user-name",    // Contains hyphen
    "user@name",    // Contains @ symbol
    "user name",    // Contains space
    "user!",        // Contains exclamation
}
```

### No Special Characters Validation

Ensures string doesn't contain potentially dangerous special characters:

```go
type Content struct {
    Description string `validate:"no_special_chars"`
}

// Valid descriptions
validDescriptions := []string{
    "Simple text description",
    "Text with numbers 123",
    "Description with spaces",
    "",  // Empty is allowed
}

// Invalid descriptions (contain special chars)
invalidDescriptions := []string{
    "Text with <script>",     // Contains < >
    "SQL injection'; DROP",  // Contains ; '
    "XSS attack &lt;script", // Contains & <
    "Command injection `ls`", // Contains ` `
    "Path traversal ../../../", // Contains special chars
}
```

## Advanced Usage

### Custom Validation Messages

```go
package main

import (
    "errors"
    "fmt"
    "strings"
    
    "github.com/saipulimdn/gopackkit/validator"
)

type User struct {
    Name     string `validate:"required" json:"name"`
    Email    string `validate:"required,email_safe" json:"email"`
    Username string `validate:"required,alphanumeric" json:"username"`
}

func validateUserWithCustomMessages(user User) map[string]string {
    v := validator.New()
    err := v.Struct(user)
    
    if err == nil {
        return nil
    }
    
    errors := make(map[string]string)
    
    // Parse validation errors and create custom messages
    errStr := err.Error()
    
    if strings.Contains(errStr, "Name") && strings.Contains(errStr, "required") {
        errors["name"] = "Name is required and cannot be empty"
    }
    
    if strings.Contains(errStr, "Email") && strings.Contains(errStr, "required") {
        errors["email"] = "Email address is required"
    } else if strings.Contains(errStr, "Email") && strings.Contains(errStr, "email_safe") {
        errors["email"] = "Please enter a valid email address"
    }
    
    if strings.Contains(errStr, "Username") && strings.Contains(errStr, "required") {
        errors["username"] = "Username is required"
    } else if strings.Contains(errStr, "Username") && strings.Contains(errStr, "alphanumeric") {
        errors["username"] = "Username can only contain letters and numbers"
    }
    
    return errors
}

func main() {
    user := User{
        Name:     "",
        Email:    "invalid-email",
        Username: "user_name!",
    }
    
    if errors := validateUserWithCustomMessages(user); errors != nil {
        fmt.Println("Validation errors:")
        for field, message := range errors {
            fmt.Printf("  %s: %s\n", field, message)
        }
    }
}
```

### Conditional Validation

```go
package main

import (
    "github.com/saipulimdn/gopackkit/validator"
)

type UserProfile struct {
    UserType    string `validate:"required"`
    Email       string `validate:"required,email_safe"`
    Phone       string // Conditional validation
    CompanyName string // Only required for business users
    Website     string // Only for business users
}

func validateUserProfile(profile UserProfile) error {
    v := validator.New()
    
    // Basic validation
    if err := v.Struct(profile); err != nil {
        return err
    }
    
    // Conditional validation based on user type
    if profile.UserType == "business" {
        // Business users must have company name
        if err := v.Var(profile.CompanyName, "required"); err != nil {
            return errors.New("company name is required for business users")
        }
        
        // Validate website if provided
        if profile.Website != "" {
            if err := v.Var(profile.Website, "no_special_chars"); err != nil {
                return errors.New("website URL contains invalid characters")
            }
        }
    }
    
    // Phone is optional but must be valid if provided
    if profile.Phone != "" {
        if err := v.Var(profile.Phone, "phone"); err != nil {
            return errors.New("phone number format is invalid")
        }
    }
    
    return nil
}

func main() {
    // Business user example
    businessUser := UserProfile{
        UserType:    "business",
        Email:       "business@example.com",
        Phone:       "1234567890",
        CompanyName: "Example Corp",
        Website:     "https://example.com",
    }
    
    if err := validateUserProfile(businessUser); err != nil {
        log.Printf("Business user validation failed: %v", err)
    } else {
        fmt.Println("Business user validation successful")
    }
    
    // Regular user example
    regularUser := UserProfile{
        UserType: "regular",
        Email:    "user@example.com",
        Phone:    "9876543210",
    }
    
    if err := validateUserProfile(regularUser); err != nil {
        log.Printf("Regular user validation failed: %v", err)
    } else {
        fmt.Println("Regular user validation successful")
    }
}
```

## API Integration Examples

### HTTP Handler with Validation

```go
package main

import (
    "encoding/json"
    "net/http"
    
    "github.com/saipulimdn/gopackkit/validator"
)

type CreateUserRequest struct {
    Name     string `json:"name" validate:"required"`
    Email    string `json:"email" validate:"required,email_safe"`
    Username string `json:"username" validate:"required,alphanumeric"`
    Phone    string `json:"phone" validate:"phone"`
}

type ValidationError struct {
    Field   string `json:"field"`
    Message string `json:"message"`
}

type ErrorResponse struct {
    Error   string            `json:"error"`
    Details []ValidationError `json:"details,omitempty"`
}

func parseValidationErrors(err error) []ValidationError {
    var errors []ValidationError
    
    // Simple error parsing - in real apps, you might want more sophisticated parsing
    errStr := err.Error()
    
    // This is a simplified example - you'd implement proper error parsing
    errors = append(errors, ValidationError{
        Field:   "validation",
        Message: errStr,
    })
    
    return errors
}

func createUserHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }
    
    var req CreateUserRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        w.WriteHeader(http.StatusBadRequest)
        json.NewEncoder(w).Encode(ErrorResponse{
            Error: "Invalid JSON format",
        })
        return
    }
    
    // Validate request
    v := validator.New()
    if err := v.Struct(req); err != nil {
        w.WriteHeader(http.StatusBadRequest)
        json.NewEncoder(w).Encode(ErrorResponse{
            Error:   "Validation failed",
            Details: parseValidationErrors(err),
        })
        return
    }
    
    // Process valid request
    user := map[string]interface{}{
        "id":       "user123",
        "name":     req.Name,
        "email":    req.Email,
        "username": req.Username,
        "phone":    req.Phone,
    }
    
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(map[string]interface{}{
        "success": true,
        "user":    user,
    })
}

func main() {
    http.HandleFunc("/users", createUserHandler)
    
    log.Println("Server starting on :8080")
    log.Fatal(http.ListenAndServe(":8080", nil))
}
```

### Middleware for Request Validation

```go
package main

import (
    "context"
    "encoding/json"
    "net/http"
    "reflect"
    
    "github.com/saipulimdn/gopackkit/validator"
)

type ValidatedRequest interface {
    Validate() error
}

func ValidationMiddleware(next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        // Skip validation for GET requests
        if r.Method == http.MethodGet {
            next(w, r)
            return
        }
        
        // Get request type from context (you'd set this based on your routing)
        requestType := r.Context().Value("requestType")
        if requestType == nil {
            next(w, r)
            return
        }
        
        // Create new instance of request type
        reqValue := reflect.New(requestType.(reflect.Type)).Interface()
        
        // Decode JSON
        if err := json.NewDecoder(r.Body).Decode(reqValue); err != nil {
            http.Error(w, "Invalid JSON", http.StatusBadRequest)
            return
        }
        
        // Validate using struct tags
        v := validator.New()
        if err := v.Struct(reqValue); err != nil {
            w.WriteHeader(http.StatusBadRequest)
            json.NewEncoder(w).Encode(map[string]string{
                "error": "Validation failed: " + err.Error(),
            })
            return
        }
        
        // Add validated request to context
        ctx := context.WithValue(r.Context(), "validatedRequest", reqValue)
        next(w, r.WithContext(ctx))
    }
}
```

## Security Features

### XSS Prevention

```go
package main

import (
    "html"
    "strings"
    
    "github.com/saipulimdn/gopackkit/validator"
)

type SafeContent struct {
    Title       string `validate:"required,no_special_chars"`
    Description string `validate:"no_special_chars"`
    Content     string // Will be sanitized separately
}

func sanitizeInput(input string) string {
    // HTML escape
    escaped := html.EscapeString(input)
    
    // Remove potentially dangerous patterns
    dangerous := []string{
        "<script",
        "</script>",
        "javascript:",
        "onload=",
        "onerror=",
        "onclick=",
    }
    
    cleaned := escaped
    for _, pattern := range dangerous {
        cleaned = strings.ReplaceAll(strings.ToLower(cleaned), pattern, "")
    }
    
    return cleaned
}

func validateAndSanitize(content SafeContent) (SafeContent, error) {
    v := validator.New()
    
    // Validate struct
    if err := v.Struct(content); err != nil {
        return content, err
    }
    
    // Sanitize content that allows HTML
    content.Content = sanitizeInput(content.Content)
    
    return content, nil
}
```

### SQL Injection Prevention

```go
package main

import (
    "strings"
    "github.com/saipulimdn/gopackkit/validator"
)

type SearchQuery struct {
    Query    string `validate:"required,no_special_chars"`
    Category string `validate:"alphanumeric"`
    SortBy   string `validate:"alphanumeric"`
}

func validateSearchQuery(query SearchQuery) error {
    v := validator.New()
    
    // Basic validation
    if err := v.Struct(query); err != nil {
        return err
    }
    
    // Additional SQL injection checks
    sqlKeywords := []string{
        "drop", "delete", "update", "insert", "create",
        "alter", "exec", "execute", "union", "select",
    }
    
    queryLower := strings.ToLower(query.Query)
    for _, keyword := range sqlKeywords {
        if strings.Contains(queryLower, keyword) {
            return errors.New("query contains potentially dangerous SQL keywords")
        }
    }
    
    return nil
}
```

## Performance Considerations

### Batch Validation

```go
package main

import (
    "fmt"
    "sync"
    
    "github.com/saipulimdn/gopackkit/validator"
)

type User struct {
    Name  string `validate:"required"`
    Email string `validate:"required,email_safe"`
}

func validateUsersBatch(users []User) map[int]error {
    v := validator.New()
    results := make(map[int]error)
    var mu sync.Mutex
    var wg sync.WaitGroup
    
    // Validate users concurrently
    for i, user := range users {
        wg.Add(1)
        go func(index int, u User) {
            defer wg.Done()
            
            err := v.Struct(u)
            
            mu.Lock()
            if err != nil {
                results[index] = err
            }
            mu.Unlock()
        }(i, user)
    }
    
    wg.Wait()
    return results
}

func main() {
    users := []User{
        {Name: "John", Email: "john@example.com"},
        {Name: "", Email: "invalid"},
        {Name: "Jane", Email: "jane@example.com"},
    }
    
    results := validateUsersBatch(users)
    
    for index, err := range results {
        fmt.Printf("User %d validation failed: %v\n", index, err)
    }
}
```

## Testing

### Unit Tests

```go
package main

import (
    "testing"
    
    "github.com/saipulimdn/gopackkit/validator"
)

func TestEmailValidation(t *testing.T) {
    v := validator.New()
    
    validEmails := []string{
        "test@example.com",
        "user.name@domain.org",
        "user+tag@example.co.uk",
    }
    
    invalidEmails := []string{
        "invalid-email",
        "@example.com",
        "user@",
        "",
    }
    
    for _, email := range validEmails {
        err := v.Var(email, "email_safe")
        if err != nil {
            t.Errorf("Expected %s to be valid, got error: %v", email, err)
        }
    }
    
    for _, email := range invalidEmails {
        err := v.Var(email, "email_safe")
        if err == nil {
            t.Errorf("Expected %s to be invalid, but validation passed", email)
        }
    }
}

func TestPhoneValidation(t *testing.T) {
    v := validator.New()
    
    validPhones := []string{
        "1234567890",      // 10 digits
        "12345678901",     // 11 digits
        "123456789012345", // 15 digits
    }
    
    invalidPhones := []string{
        "123456789",       // Too short
        "1234567890123456", // Too long
        "12345abcde",      // Contains letters
        "123-456-7890",    // Contains special chars
    }
    
    for _, phone := range validPhones {
        err := v.Var(phone, "phone")
        if err != nil {
            t.Errorf("Expected %s to be valid, got error: %v", phone, err)
        }
    }
    
    for _, phone := range invalidPhones {
        err := v.Var(phone, "phone")
        if err == nil {
            t.Errorf("Expected %s to be invalid, but validation passed", phone)
        }
    }
}

func TestAlphanumericValidation(t *testing.T) {
    v := validator.New()
    
    validInputs := []string{
        "user123",
        "JohnDoe",
        "testuser",
        "ABC123def",
    }
    
    invalidInputs := []string{
        "user_123",
        "user-name",
        "user@name",
        "user name",
        "user!",
    }
    
    for _, input := range validInputs {
        err := v.Var(input, "alphanumeric")
        if err != nil {
            t.Errorf("Expected %s to be valid, got error: %v", input, err)
        }
    }
    
    for _, input := range invalidInputs {
        err := v.Var(input, "alphanumeric")
        if err == nil {
            t.Errorf("Expected %s to be invalid, but validation passed", input)
        }
    }
}
```

## Best Practices

1. **Always Validate User Input**: Never trust data from external sources
2. **Use Struct Tags**: Leverage validation tags for clean, declarative validation
3. **Combine with Sanitization**: Use validation alongside input sanitization
4. **Handle Errors Gracefully**: Provide meaningful error messages to users
5. **Validate Early**: Perform validation at the API boundary
6. **Security First**: Use secure validation rules to prevent attacks
7. **Test Thoroughly**: Write comprehensive tests for all validation scenarios

## Troubleshooting

### Common Issues

1. **False Positives**: Adjust validation rules if legitimate input is rejected
2. **Performance Issues**: Use batch validation for large datasets
3. **Custom Requirements**: Implement additional validation logic as needed
4. **Error Handling**: Ensure proper error message handling in your application

### Debug Mode

```go
package main

import (
    "fmt"
    "reflect"
    
    "github.com/saipulimdn/gopackkit/validator"
)

func debugValidation(data interface{}) {
    v := validator.New()
    
    fmt.Println("=== Validation Debug ===")
    
    // Get struct info
    val := reflect.ValueOf(data)
    typ := reflect.TypeOf(data)
    
    if val.Kind() == reflect.Ptr {
        val = val.Elem()
        typ = typ.Elem()
    }
    
    // Show field validation rules
    for i := 0; i < val.NumField(); i++ {
        field := typ.Field(i)
        value := val.Field(i)
        tag := field.Tag.Get("validate")
        
        fmt.Printf("Field: %s\n", field.Name)
        fmt.Printf("  Value: %v\n", value.Interface())
        fmt.Printf("  Rules: %s\n", tag)
        
        // Test individual field
        if tag != "" {
            err := v.Var(value.Interface(), tag)
            if err != nil {
                fmt.Printf("  Error: %v\n", err)
            } else {
                fmt.Printf("  Status: Valid\n")
            }
        }
        fmt.Println()
    }
    
    // Test entire struct
    err := v.Struct(data)
    if err != nil {
        fmt.Printf("Overall validation failed: %v\n", err)
    } else {
        fmt.Println("Overall validation: Passed")
    }
}
```
