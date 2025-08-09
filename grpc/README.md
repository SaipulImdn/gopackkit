# gRPC Module

The gRPC module provides both client and server implementations for gRPC communication with comprehensive configuration options, connection management, and middleware support.

## Features

### Client Features
- **Automatic Connection Management**: Connection pooling and health checking
- **Configurable Timeouts**: Request and connection timeout settings
- **Keep-Alive Support**: Configurable keep-alive parameters
- **TLS/SSL Support**: Secure connections with certificate management
- **Retry Logic**: Built-in retry mechanism with exponential backoff
- **Connection State Monitoring**: Real-time connection state tracking

### Server Features
- **Graceful Shutdown**: Proper server shutdown handling
- **TLS/SSL Support**: Secure server configuration
- **Interceptors**: Request/response middleware support
- **Health Checks**: Built-in health check service
- **Connection Limits**: Configurable connection management
- **Reflection Support**: gRPC reflection for development

## Installation

```bash
go get github.com/saipulimdn/gopackkit/grpc
```

## Quick Start

### gRPC Client

```go
package main

import (
    "context"
    "log"
    "time"
    
    "github.com/saipulimdn/gopackkit/grpc"
    "google.golang.org/grpc"
)

func main() {
    // Create client with default configuration
    client, err := grpc.NewClient()
    if err != nil {
        log.Fatal("Failed to create gRPC client:", err)
    }
    defer client.Close()
    
    // Get the connection for service calls
    conn := client.GetConnection()
    
    // Use connection with your gRPC service
    // serviceClient := yourpb.NewYourServiceClient(conn)
    
    log.Println("gRPC client connected successfully")
}
```

### gRPC Server

```go
package main

import (
    "log"
    
    "github.com/saipulimdn/gopackkit/grpc"
)

func main() {
    // Create server with default configuration
    server, err := grpc.NewServer()
    if err != nil {
        log.Fatal("Failed to create gRPC server:", err)
    }
    
    // Register your services here
    // yourpb.RegisterYourServiceServer(server.GetServer(), &yourServiceImpl{})
    
    // Start server
    log.Println("Starting gRPC server on :9000")
    if err := server.Start(); err != nil {
        log.Fatal("Failed to start server:", err)
    }
}
```

## Client Configuration

### Default Configuration

```go
client, err := grpc.NewClient()
```

### Custom Configuration

```go
package main

import (
    "time"
    "github.com/saipulimdn/gopackkit/grpc"
)

func main() {
    config := grpc.ClientConfig{
        Host:                   "api.example.com",
        Port:                   443,
        MaxRecvMsgSize:        8 * 1024 * 1024, // 8MB
        MaxSendMsgSize:        8 * 1024 * 1024, // 8MB
        KeepAliveTime:         30 * time.Second,
        KeepAliveTimeout:      5 * time.Second,
        KeepAliveWithoutStream: true,
        ConnectionTimeout:     15 * time.Second,
        EnableTLS:             true,
        InsecureTLS:          false,
        ServerNameOverride:   "api.example.com",
        Block:                 true,
    }
    
    client, err := grpc.NewClientWithConfig(config)
    if err != nil {
        log.Fatal("Failed to create client:", err)
    }
    defer client.Close()
}
```

### Client Configuration Options

```go
type ClientConfig struct {
    Host                     string        // Server host
    Port                     int           // Server port
    MaxRecvMsgSize          int           // Maximum receive message size
    MaxSendMsgSize          int           // Maximum send message size
    KeepAliveTime           time.Duration // Keep-alive time
    KeepAliveTimeout        time.Duration // Keep-alive timeout
    KeepAliveWithoutStream  bool          // Keep-alive without stream
    ConnectionTimeout       time.Duration // Connection timeout
    EnableTLS               bool          // Enable TLS
    InsecureTLS            bool          // Skip TLS verification
    ServerNameOverride     string        // Server name override
    CertFile               string        // Certificate file path
    Block                  bool          // Block until connection ready
}
```

## Server Configuration

### Default Configuration

```go
server, err := grpc.NewServer()
```

### Custom Configuration

```go
package main

import (
    "time"
    "github.com/saipulimdn/gopackkit/grpc"
)

func main() {
    config := grpc.ServerConfig{
        Host:                    "0.0.0.0",
        Port:                    9090,
        MaxRecvMsgSize:         8 * 1024 * 1024,
        MaxSendMsgSize:         8 * 1024 * 1024,
        ConnectionTimeout:      10 * time.Second,
        KeepAliveTime:          30 * time.Second,
        KeepAliveTimeout:       5 * time.Second,
        KeepAliveMinTime:       5 * time.Second,
        KeepAlivePermitWithoutStream: false,
        MaxConnectionIdle:      15 * time.Minute,
        MaxConnectionAge:       30 * time.Minute,
        MaxConnectionAgeGrace:  5 * time.Minute,
        EnableTLS:              true,
        CertFile:              "/path/to/cert.pem",
        KeyFile:               "/path/to/key.pem",
        EnableReflection:       true,
    }
    
    server, err := grpc.NewServerWithConfig(config)
    if err != nil {
        log.Fatal("Failed to create server:", err)
    }
    
    // Start server
    if err := server.Start(); err != nil {
        log.Fatal("Failed to start server:", err)
    }
}
```

## Advanced Usage

### Client with Service Implementation

```go
package main

import (
    "context"
    "log"
    "time"
    
    "github.com/saipulimdn/gopackkit/grpc"
    "google.golang.org/grpc/codes"
    "google.golang.org/grpc/status"
)

// Assuming you have a proto-generated service
// import "your/proto/package"

type UserServiceClient struct {
    client *grpc.Client
    // serviceClient yourpb.UserServiceClient
}

func NewUserServiceClient(host string, port int) (*UserServiceClient, error) {
    config := grpc.ClientConfig{
        Host:              host,
        Port:              port,
        ConnectionTimeout: 10 * time.Second,
        EnableTLS:        false,
    }
    
    client, err := grpc.NewClientWithConfig(config)
    if err != nil {
        return nil, err
    }
    
    // Wait for connection to be ready
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()
    
    if err := client.WaitForConnection(ctx); err != nil {
        client.Close()
        return nil, err
    }
    
    return &UserServiceClient{
        client: client,
        // serviceClient: yourpb.NewUserServiceClient(client.GetConnection()),
    }, nil
}

func (c *UserServiceClient) GetUser(ctx context.Context, userID string) error {
    // Check connection health
    if err := c.client.HealthCheck(); err != nil {
        return status.Error(codes.Unavailable, "client not connected")
    }
    
    // Make gRPC call
    // req := &yourpb.GetUserRequest{UserId: userID}
    // resp, err := c.serviceClient.GetUser(ctx, req)
    // if err != nil {
    //     return err
    // }
    
    log.Printf("Successfully retrieved user: %s", userID)
    return nil
}

func (c *UserServiceClient) Close() error {
    return c.client.Close()
}

func main() {
    client, err := NewUserServiceClient("localhost", 9000)
    if err != nil {
        log.Fatal("Failed to create user service client:", err)
    }
    defer client.Close()
    
    ctx := context.Background()
    err = client.GetUser(ctx, "user123")
    if err != nil {
        log.Fatal("Failed to get user:", err)
    }
}
```

### Server with Service Implementation

```go
package main

import (
    "context"
    "log"
    "os"
    "os/signal"
    "syscall"
    
    "github.com/saipulimdn/gopackkit/grpc"
    "google.golang.org/grpc/codes"
    "google.golang.org/grpc/status"
)

// Assuming you have a proto-generated service
// import "your/proto/package"

type UserServiceServer struct {
    // yourpb.UnimplementedUserServiceServer
}

func (s *UserServiceServer) GetUser(ctx context.Context, req interface{}) (interface{}, error) {
    // Implementation example
    // req := req.(*yourpb.GetUserRequest)
    
    // Validate request
    // if req.UserId == "" {
    //     return nil, status.Error(codes.InvalidArgument, "user ID is required")
    // }
    
    // Business logic here
    log.Printf("Getting user with ID: %v", req)
    
    // Return response
    // return &yourpb.GetUserResponse{
    //     User: &yourpb.User{
    //         Id: req.UserId,
    //         Name: "John Doe",
    //         Email: "john@example.com",
    //     },
    // }, nil
    
    return nil, nil
}

func main() {
    // Create server
    config := grpc.ServerConfig{
        Port:             9000,
        EnableReflection: true,
    }
    
    server, err := grpc.NewServerWithConfig(config)
    if err != nil {
        log.Fatal("Failed to create gRPC server:", err)
    }
    
    // Register service
    userService := &UserServiceServer{}
    // yourpb.RegisterUserServiceServer(server.GetServer(), userService)
    
    // Handle graceful shutdown
    go func() {
        sigChan := make(chan os.Signal, 1)
        signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
        <-sigChan
        
        log.Println("Shutting down gRPC server...")
        server.Stop()
    }()
    
    // Start server
    log.Printf("Starting gRPC server on port %d", config.Port)
    if err := server.Start(); err != nil {
        log.Fatal("Failed to start server:", err)
    }
}
```

### TLS Configuration

#### Client with TLS

```go
package main

import (
    "github.com/saipulimdn/gopackkit/grpc"
)

func main() {
    // TLS with certificate file
    config := grpc.ClientConfig{
        Host:               "secure-api.example.com",
        Port:               443,
        EnableTLS:          true,
        CertFile:          "/path/to/ca-cert.pem",
        ServerNameOverride: "secure-api.example.com",
    }
    
    client, err := grpc.NewClientWithConfig(config)
    if err != nil {
        log.Fatal("Failed to create TLS client:", err)
    }
    defer client.Close()
    
    // For development/testing with self-signed certificates
    devConfig := grpc.ClientConfig{
        Host:        "localhost",
        Port:        9443,
        EnableTLS:   true,
        InsecureTLS: true, // Skip certificate verification
    }
    
    devClient, err := grpc.NewClientWithConfig(devConfig)
    if err != nil {
        log.Fatal("Failed to create dev TLS client:", err)
    }
    defer devClient.Close()
}
```

#### Server with TLS

```go
package main

import (
    "github.com/saipulimdn/gopackkit/grpc"
)

func main() {
    config := grpc.ServerConfig{
        Port:      9443,
        EnableTLS: true,
        CertFile:  "/path/to/server-cert.pem",
        KeyFile:   "/path/to/server-key.pem",
    }
    
    server, err := grpc.NewServerWithConfig(config)
    if err != nil {
        log.Fatal("Failed to create TLS server:", err)
    }
    
    log.Println("Starting secure gRPC server on :9443")
    if err := server.Start(); err != nil {
        log.Fatal("Failed to start secure server:", err)
    }
}
```

## Connection Management

### Client Connection Monitoring

```go
package main

import (
    "context"
    "log"
    "time"
    
    "github.com/saipulimdn/gopackkit/grpc"
    "google.golang.org/grpc/connectivity"
)

func monitorConnection(client *grpc.Client) {
    ticker := time.NewTicker(5 * time.Second)
    defer ticker.Stop()
    
    for {
        select {
        case <-ticker.C:
            if client.IsConnected() {
                log.Printf("Connection status: Connected to %s", client.GetAddress())
            } else {
                log.Printf("Connection status: Disconnected from %s", client.GetAddress())
                
                // Attempt to reconnect
                ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
                err := client.WaitForConnection(ctx)
                cancel()
                
                if err != nil {
                    log.Printf("Reconnection failed: %v", err)
                } else {
                    log.Println("Reconnected successfully")
                }
            }
            
            // Perform health check
            if err := client.HealthCheck(); err != nil {
                log.Printf("Health check failed: %v", err)
            }
        }
    }
}

func main() {
    client, err := grpc.NewClient()
    if err != nil {
        log.Fatal("Failed to create client:", err)
    }
    defer client.Close()
    
    // Start connection monitoring
    go monitorConnection(client)
    
    // Your application logic here
    select {} // Keep running
}
```

## Environment Variables

Configure gRPC client and server using environment variables:

### Client Environment Variables

```bash
# Connection settings
export GRPC_CLIENT_HOST=api.example.com
export GRPC_CLIENT_PORT=9000
export GRPC_CLIENT_CONNECTION_TIMEOUT=15s
export GRPC_CLIENT_BLOCK=true

# Message size limits
export GRPC_CLIENT_MAX_RECV_MSG_SIZE=8388608  # 8MB
export GRPC_CLIENT_MAX_SEND_MSG_SIZE=8388608  # 8MB

# Keep-alive settings
export GRPC_CLIENT_KEEP_ALIVE_TIME=30s
export GRPC_CLIENT_KEEP_ALIVE_TIMEOUT=5s
export GRPC_CLIENT_KEEP_ALIVE_WITHOUT_STREAM=true

# TLS settings
export GRPC_CLIENT_ENABLE_TLS=true
export GRPC_CLIENT_INSECURE_TLS=false
export GRPC_CLIENT_SERVER_NAME_OVERRIDE=api.example.com
export GRPC_CLIENT_CERT_FILE=/path/to/cert.pem
```

### Server Environment Variables

```bash
# Server settings
export GRPC_SERVER_HOST=0.0.0.0
export GRPC_SERVER_PORT=9000
export GRPC_SERVER_CONNECTION_TIMEOUT=10s

# Message size limits
export GRPC_SERVER_MAX_RECV_MSG_SIZE=8388608  # 8MB
export GRPC_SERVER_MAX_SEND_MSG_SIZE=8388608  # 8MB

# Keep-alive settings
export GRPC_SERVER_KEEP_ALIVE_TIME=30s
export GRPC_SERVER_KEEP_ALIVE_TIMEOUT=5s
export GRPC_SERVER_KEEP_ALIVE_MIN_TIME=5s
export GRPC_SERVER_KEEP_ALIVE_PERMIT_WITHOUT_STREAM=false

# Connection management
export GRPC_SERVER_MAX_CONNECTION_IDLE=15m
export GRPC_SERVER_MAX_CONNECTION_AGE=30m
export GRPC_SERVER_MAX_CONNECTION_AGE_GRACE=5m

# TLS settings
export GRPC_SERVER_ENABLE_TLS=false
export GRPC_SERVER_CERT_FILE=/path/to/cert.pem
export GRPC_SERVER_KEY_FILE=/path/to/key.pem

# Development settings
export GRPC_SERVER_ENABLE_REFLECTION=true
```

## Examples

### Microservice Communication

```go
package main

import (
    "context"
    "log"
    "time"
    
    "github.com/saipulimdn/gopackkit/grpc"
)

type ServiceRegistry struct {
    userService    *grpc.Client
    paymentService *grpc.Client
    orderService   *grpc.Client
}

func NewServiceRegistry() (*ServiceRegistry, error) {
    // User Service Client
    userClient, err := grpc.NewClientWithConfig(grpc.ClientConfig{
        Host: "user-service",
        Port: 9001,
        ConnectionTimeout: 10 * time.Second,
    })
    if err != nil {
        return nil, err
    }
    
    // Payment Service Client
    paymentClient, err := grpc.NewClientWithConfig(grpc.ClientConfig{
        Host: "payment-service",
        Port: 9002,
        ConnectionTimeout: 10 * time.Second,
    })
    if err != nil {
        userClient.Close()
        return nil, err
    }
    
    // Order Service Client
    orderClient, err := grpc.NewClientWithConfig(grpc.ClientConfig{
        Host: "order-service",
        Port: 9003,
        ConnectionTimeout: 10 * time.Second,
    })
    if err != nil {
        userClient.Close()
        paymentClient.Close()
        return nil, err
    }
    
    return &ServiceRegistry{
        userService:    userClient,
        paymentService: paymentClient,
        orderService:   orderClient,
    }, nil
}

func (s *ServiceRegistry) Close() {
    s.userService.Close()
    s.paymentService.Close()
    s.orderService.Close()
}

func (s *ServiceRegistry) HealthCheck() error {
    services := map[string]*grpc.Client{
        "user-service":    s.userService,
        "payment-service": s.paymentService,
        "order-service":   s.orderService,
    }
    
    for name, client := range services {
        if err := client.HealthCheck(); err != nil {
            log.Printf("Health check failed for %s: %v", name, err)
            return err
        }
    }
    
    return nil
}

func main() {
    registry, err := NewServiceRegistry()
    if err != nil {
        log.Fatal("Failed to create service registry:", err)
    }
    defer registry.Close()
    
    // Periodic health checks
    ticker := time.NewTicker(30 * time.Second)
    defer ticker.Stop()
    
    for {
        select {
        case <-ticker.C:
            if err := registry.HealthCheck(); err != nil {
                log.Printf("Service health check failed: %v", err)
            } else {
                log.Println("All services healthy")
            }
        }
    }
}
```

## Best Practices

### Client Best Practices

1. **Reuse Connections**: Create one client per service and reuse it
2. **Connection Pooling**: Use default connection pooling for better performance
3. **Health Monitoring**: Implement regular health checks
4. **Graceful Degradation**: Handle connection failures gracefully
5. **Timeout Management**: Set appropriate timeouts for operations
6. **TLS in Production**: Always use TLS in production environments

### Server Best Practices

1. **Graceful Shutdown**: Implement proper shutdown handling
2. **Resource Limits**: Set appropriate message size limits
3. **Connection Management**: Configure keep-alive and connection limits
4. **Error Handling**: Return appropriate gRPC status codes
5. **Monitoring**: Implement metrics and health checks
6. **Security**: Use TLS and proper authentication

## Testing

### Unit Testing

```go
package main

import (
    "testing"
    "time"
    
    "github.com/saipulimdn/gopackkit/grpc"
)

func TestClientCreation(t *testing.T) {
    config := grpc.ClientConfig{
        Host:              "localhost",
        Port:              9000,
        ConnectionTimeout: 5 * time.Second,
        Block:             false, // Don't block in tests
    }
    
    client, err := grpc.NewClientWithConfig(config)
    if err != nil {
        t.Fatalf("Failed to create client: %v", err)
    }
    defer client.Close()
    
    if client.GetAddress() != "localhost:9000" {
        t.Errorf("Expected address localhost:9000, got %s", client.GetAddress())
    }
}

func TestServerCreation(t *testing.T) {
    config := grpc.ServerConfig{
        Port: 0, // Use random port for testing
    }
    
    server, err := grpc.NewServerWithConfig(config)
    if err != nil {
        t.Fatalf("Failed to create server: %v", err)
    }
    
    // Test server configuration
    if server.GetAddress() == "" {
        t.Error("Server address should not be empty")
    }
}
```

## Troubleshooting

### Common Issues

1. **Connection Refused**: Check if server is running and port is correct
2. **TLS Handshake Failed**: Verify certificate configuration
3. **Context Deadline Exceeded**: Increase timeout values
4. **Message Size Exceeded**: Increase MaxRecvMsgSize/MaxSendMsgSize
5. **Connection Drops**: Configure keep-alive settings

### Debug Mode

```go
// Enable gRPC logging for debugging
import "google.golang.org/grpc/grpclog"

func init() {
    grpclog.SetLoggerV2(grpclog.NewLoggerV2(os.Stdout, os.Stdout, os.Stderr))
}
```

### Health Check Integration

```go
package main

import (
    "context"
    "net/http"
    
    "github.com/saipulimdn/gopackkit/grpc"
)

func healthHandler(client *grpc.Client) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        if err := client.HealthCheck(); err != nil {
            w.WriteHeader(http.StatusServiceUnavailable)
            w.Write([]byte("gRPC client unhealthy: " + err.Error()))
            return
        }
        
        w.WriteHeader(http.StatusOK)
        w.Write([]byte("gRPC client healthy"))
    }
}

func main() {
    client, err := grpc.NewClient()
    if err != nil {
        log.Fatal("Failed to create gRPC client:", err)
    }
    defer client.Close()
    
    // Health check endpoint
    http.HandleFunc("/health", healthHandler(client))
    
    log.Println("Health check server starting on :8080")
    http.ListenAndServe(":8080", nil)
}
```
