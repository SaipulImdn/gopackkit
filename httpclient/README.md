# HTTP Client Module

HTTP Client module menyediakan HTTP client dengan retry mechanism, timeout configuration, dan middleware support. Module ini dirancang untuk memberikan reliability dan ease of use dalam melakukan HTTP requests.

## Features

- **Automatic Retry**: Exponential backoff dengan configurable attempts
- **Configurable Timeout**: Request timeout yang dapat disesuaikan
- **Multiple HTTP Methods**: GET, POST, PUT, DELETE support
- **JSON Support**: Built-in JSON request/response handling
- **Error Handling**: Comprehensive error handling dengan retry logic
- **Middleware Support**: Request/Response middleware chain
- **Connection Pooling**: Reusable HTTP connections

## Installation

```bash
go get github.com/saipulimdn/gopackkit/httpclient
```

## Quick Start

### Basic Usage

```go
package main

import (
    "fmt"
    "log"
    
    "github.com/saipulimdn/gopackkit/httpclient"
)

func main() {
    // Create default client
    client := httpclient.New()
    
    // Make GET request
    resp, err := client.Get("https://api.github.com/users/octocat")
    if err != nil {
        log.Fatal(err)
    }
    defer resp.Body.Close()
    
    fmt.Printf("Status: %s\n", resp.Status)
}
```

### Custom Configuration

```go
package main

import (
    "time"
    "github.com/saipulimdn/gopackkit/httpclient"
)

func main() {
    config := httpclient.Config{
        Timeout:       30 * time.Second,
        MaxRetries:    5,
        RetryDelay:    2 * time.Second,
        MaxRetryDelay: 30 * time.Second,
    }
    
    client := httpclient.NewWithConfig(config)
    
    resp, err := client.Get("https://api.example.com/data")
    if err != nil {
        log.Fatal(err)
    }
    defer resp.Body.Close()
}
```

## Configuration

```go
type Config struct {
    Timeout       time.Duration // Request timeout
    MaxRetries    int           // Maximum retry attempts
    RetryDelay    time.Duration // Initial retry delay
    MaxRetryDelay time.Duration // Maximum retry delay
}
```

### Default Configuration

```go
config := httpclient.Config{
    Timeout:       10 * time.Second,
    MaxRetries:    3,
    RetryDelay:    1 * time.Second,
    MaxRetryDelay: 10 * time.Second,
}
```

## Methods

### GET Request

```go
// Simple GET
resp, err := client.Get("https://api.example.com/users")

// GET with query parameters
url := "https://api.example.com/users?page=1&limit=10"
resp, err := client.Get(url)
```

### POST Request

```go
// POST with form data
data := strings.NewReader("name=John&email=john@example.com")
resp, err := client.Post(
    "https://api.example.com/users",
    "application/x-www-form-urlencoded",
    data,
)

// POST with JSON (recommended)
userData := map[string]interface{}{
    "name":  "John Doe",
    "email": "john@example.com",
    "age":   30,
}
resp, err := client.PostJSON("https://api.example.com/users", userData)
```

### PUT Request

```go
// PUT with JSON data
updateData := map[string]interface{}{
    "name": "John Updated",
    "age":  31,
}
jsonData, _ := json.Marshal(updateData)
resp, err := client.Put(
    "https://api.example.com/users/123",
    "application/json",
    bytes.NewReader(jsonData),
)
```

### DELETE Request

```go
// Simple DELETE
resp, err := client.Delete("https://api.example.com/users/123")
```

### Custom Request

```go
// Create custom request
req, err := http.NewRequest("PATCH", "https://api.example.com/users/123", body)
if err != nil {
    log.Fatal(err)
}

// Add headers
req.Header.Set("Authorization", "Bearer token123")
req.Header.Set("Content-Type", "application/json")

// Execute request
resp, err := client.Do(req)
```

## Advanced Usage

### JSON Response Handling

```go
package main

import (
    "encoding/json"
    "log"
    
    "github.com/saipulimdn/gopackkit/httpclient"
)

type User struct {
    ID    int    `json:"id"`
    Name  string `json:"name"`
    Email string `json:"email"`
}

func main() {
    client := httpclient.New()
    
    resp, err := client.Get("https://api.example.com/users/1")
    if err != nil {
        log.Fatal(err)
    }
    defer resp.Body.Close()
    
    var user User
    if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
        log.Fatal(err)
    }
    
    log.Printf("User: %+v\n", user)
}
```

### Error Handling

```go
package main

import (
    "fmt"
    "net/http"
    
    "github.com/saipulimdn/gopackkit/httpclient"
)

func main() {
    client := httpclient.New()
    
    resp, err := client.Get("https://api.example.com/users")
    if err != nil {
        log.Printf("Request failed: %v\n", err)
        return
    }
    defer resp.Body.Close()
    
    // Check HTTP status
    if resp.StatusCode != http.StatusOK {
        fmt.Printf("HTTP Error: %s\n", resp.Status)
        return
    }
    
    // Process successful response
    fmt.Println("Request successful")
}
```

### Retry Logic Example

```go
package main

import (
    "time"
    "github.com/saipulimdn/gopackkit/httpclient"
)

func main() {
    // Configure aggressive retry
    config := httpclient.Config{
        Timeout:       5 * time.Second,
        MaxRetries:    5,           // Will retry up to 5 times
        RetryDelay:    1 * time.Second,  // Start with 1 second delay
        MaxRetryDelay: 16 * time.Second, // Max delay between retries
    }
    
    client := httpclient.NewWithConfig(config)
    
    // This will automatically retry on failure
    resp, err := client.Get("https://unreliable-api.example.com/data")
    if err != nil {
        log.Printf("Failed after all retries: %v\n", err)
        return
    }
    defer resp.Body.Close()
    
    log.Println("Request successful")
}
```

## Environment Variables

Anda dapat mengkonfigurasi HTTP client menggunakan environment variables:

```bash
# Request timeout
export HTTP_CLIENT_TIMEOUT=30s

# Maximum retries
export HTTP_CLIENT_MAX_RETRIES=5

# Retry delay
export HTTP_CLIENT_RETRY_DELAY=2s

# Maximum retry delay
export HTTP_CLIENT_MAX_RETRY_DELAY=30s
```

## Examples

### API Client Implementation

```go
package main

import (
    "encoding/json"
    "fmt"
    "time"
    
    "github.com/saipulimdn/gopackkit/httpclient"
)

type APIClient struct {
    client  *httpclient.Client
    baseURL string
    apiKey  string
}

type Response struct {
    Data    interface{} `json:"data"`
    Message string      `json:"message"`
    Status  string      `json:"status"`
}

func NewAPIClient(baseURL, apiKey string) *APIClient {
    config := httpclient.Config{
        Timeout:       15 * time.Second,
        MaxRetries:    3,
        RetryDelay:    2 * time.Second,
        MaxRetryDelay: 10 * time.Second,
    }
    
    return &APIClient{
        client:  httpclient.NewWithConfig(config),
        baseURL: baseURL,
        apiKey:  apiKey,
    }
}

func (api *APIClient) GetUser(userID string) (*User, error) {
    url := fmt.Sprintf("%s/users/%s", api.baseURL, userID)
    
    req, err := http.NewRequest("GET", url, nil)
    if err != nil {
        return nil, err
    }
    
    // Add authentication
    req.Header.Set("Authorization", "Bearer "+api.apiKey)
    req.Header.Set("Content-Type", "application/json")
    
    resp, err := api.client.Do(req)
    if err != nil {
        return nil, fmt.Errorf("request failed: %w", err)
    }
    defer resp.Body.Close()
    
    if resp.StatusCode != http.StatusOK {
        return nil, fmt.Errorf("API error: %s", resp.Status)
    }
    
    var response Response
    if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
        return nil, fmt.Errorf("decode error: %w", err)
    }
    
    // Convert response.Data to User struct
    userBytes, _ := json.Marshal(response.Data)
    var user User
    json.Unmarshal(userBytes, &user)
    
    return &user, nil
}

func (api *APIClient) CreateUser(user User) (*User, error) {
    url := fmt.Sprintf("%s/users", api.baseURL)
    
    resp, err := api.client.PostJSON(url, user)
    if err != nil {
        return nil, fmt.Errorf("request failed: %w", err)
    }
    defer resp.Body.Close()
    
    if resp.StatusCode != http.StatusCreated {
        return nil, fmt.Errorf("API error: %s", resp.Status)
    }
    
    var response Response
    if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
        return nil, fmt.Errorf("decode error: %w", err)
    }
    
    userBytes, _ := json.Marshal(response.Data)
    var createdUser User
    json.Unmarshal(userBytes, &createdUser)
    
    return &createdUser, nil
}

func main() {
    client := NewAPIClient("https://api.example.com", "your-api-key")
    
    // Get user
    user, err := client.GetUser("123")
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("User: %+v\n", user)
    
    // Create user
    newUser := User{
        Name:  "Jane Doe",
        Email: "jane@example.com",
    }
    
    createdUser, err := client.CreateUser(newUser)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Created User: %+v\n", createdUser)
}
```

### File Upload Example

```go
package main

import (
    "bytes"
    "io"
    "mime/multipart"
    "os"
    
    "github.com/saipulimdn/gopackkit/httpclient"
)

func uploadFile(client *httpclient.Client, filePath string) error {
    // Open file
    file, err := os.Open(filePath)
    if err != nil {
        return err
    }
    defer file.Close()
    
    // Create multipart form
    var body bytes.Buffer
    writer := multipart.NewWriter(&body)
    
    // Add file field
    part, err := writer.CreateFormFile("file", filePath)
    if err != nil {
        return err
    }
    
    // Copy file content
    _, err = io.Copy(part, file)
    if err != nil {
        return err
    }
    
    // Add other form fields
    writer.WriteField("description", "File upload")
    
    // Close writer
    err = writer.Close()
    if err != nil {
        return err
    }
    
    // Make request
    resp, err := client.Post(
        "https://api.example.com/upload",
        writer.FormDataContentType(),
        &body,
    )
    if err != nil {
        return err
    }
    defer resp.Body.Close()
    
    if resp.StatusCode != http.StatusOK {
        return fmt.Errorf("upload failed: %s", resp.Status)
    }
    
    fmt.Println("File uploaded successfully")
    return nil
}

func main() {
    client := httpclient.New()
    
    err := uploadFile(client, "/path/to/file.jpg")
    if err != nil {
        log.Fatal(err)
    }
}
```

## Best Practices

1. **Set appropriate timeouts** berdasarkan expected response time
2. **Configure retry logic** untuk network failures
3. **Handle HTTP status codes** dengan proper error handling
4. **Use JSON methods** untuk API communications
5. **Add authentication headers** secara konsisten
6. **Close response bodies** untuk prevent memory leaks
7. **Use context** untuk request cancellation (custom requests)

## Performance Tips

1. **Reuse client instances** - HTTP clients maintain connection pools
2. **Set reasonable timeouts** - Avoid hanging requests
3. **Configure retry delays** - Prevent overwhelming servers
4. **Use connection pooling** - Default Go HTTP client behavior

## Testing

```go
package main

import (
    "net/http"
    "net/http/httptest"
    "testing"
    
    "github.com/saipulimdn/gopackkit/httpclient"
)

func TestHTTPClient(t *testing.T) {
    // Create test server
    server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusOK)
        w.Write([]byte(`{"message": "success"}`))
    }))
    defer server.Close()
    
    // Test client
    client := httpclient.New()
    
    resp, err := client.Get(server.URL)
    if err != nil {
        t.Fatalf("Request failed: %v", err)
    }
    defer resp.Body.Close()
    
    if resp.StatusCode != http.StatusOK {
        t.Errorf("Expected status 200, got %d", resp.StatusCode)
    }
}
```

## Troubleshooting

### Common Issues

1. **Timeout errors**: Increase timeout configuration
2. **Connection refused**: Check target server availability
3. **TLS errors**: Verify SSL certificates
4. **Retry exhaustion**: Check server health and increase retry limits

### Debug Mode

```go
// Enable verbose logging untuk debugging
import "net/http/httputil"

// Dump request/response for debugging
req, _ := http.NewRequest("GET", url, nil)
reqDump, _ := httputil.DumpRequestOut(req, true)
fmt.Printf("Request: %s\n", reqDump)

resp, err := client.Do(req)
if err == nil {
    respDump, _ := httputil.DumpResponse(resp, true)
    fmt.Printf("Response: %s\n", respDump)
}
```
