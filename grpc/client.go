package grpc

import (
	"context"
	"crypto/tls"
	"fmt"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/status"
)

// Client represents a gRPC client instance
type Client struct {
    conn   *grpc.ClientConn
    config ClientConfig
}

// ClientConfig holds gRPC client configuration
type ClientConfig struct {
    Host                     string        `json:"host" yaml:"host" env:"GRPC_CLIENT_HOST" default:"localhost"`
    Port                     int           `json:"port" yaml:"port" env:"GRPC_CLIENT_PORT" default:"9000"`
    MaxRecvMsgSize           int           `json:"max_recv_msg_size" yaml:"max_recv_msg_size" env:"GRPC_CLIENT_MAX_RECV_MSG_SIZE" default:"4194304"`     // 4MB
    MaxSendMsgSize           int           `json:"max_send_msg_size" yaml:"max_send_msg_size" env:"GRPC_CLIENT_MAX_SEND_MSG_SIZE" default:"4194304"`     // 4MB
    KeepAliveTime            time.Duration `json:"keep_alive_time" yaml:"keep_alive_time" env:"GRPC_CLIENT_KEEP_ALIVE_TIME" default:"30s"`
    KeepAliveTimeout         time.Duration `json:"keep_alive_timeout" yaml:"keep_alive_timeout" env:"GRPC_CLIENT_KEEP_ALIVE_TIMEOUT" default:"5s"`
    KeepAliveWithoutStream   bool          `json:"keep_alive_without_stream" yaml:"keep_alive_without_stream" env:"GRPC_CLIENT_KEEP_ALIVE_WITHOUT_STREAM" default:"true"`
    ConnectionTimeout        time.Duration `json:"connection_timeout" yaml:"connection_timeout" env:"GRPC_CLIENT_CONNECTION_TIMEOUT" default:"10s"`
    EnableTLS                bool          `json:"enable_tls" yaml:"enable_tls" env:"GRPC_CLIENT_ENABLE_TLS" default:"false"`
    InsecureTLS              bool          `json:"insecure_tls" yaml:"insecure_tls" env:"GRPC_CLIENT_INSECURE_TLS" default:"false"`
    ServerNameOverride       string        `json:"server_name_override" yaml:"server_name_override" env:"GRPC_CLIENT_SERVER_NAME_OVERRIDE"`
    CertFile                 string        `json:"cert_file" yaml:"cert_file" env:"GRPC_CLIENT_CERT_FILE"`
    Block                    bool          `json:"block" yaml:"block" env:"GRPC_CLIENT_BLOCK" default:"true"`
}

// DefaultClientConfig returns default gRPC client configuration
func DefaultClientConfig() ClientConfig {
    return ClientConfig{
        Host:                   "localhost",
        Port:                   9000,
        MaxRecvMsgSize:         4 * 1024 * 1024, // 4MB
        MaxSendMsgSize:         4 * 1024 * 1024, // 4MB
        KeepAliveTime:          30 * time.Second,
        KeepAliveTimeout:       5 * time.Second,
        KeepAliveWithoutStream: true,
        ConnectionTimeout:      10 * time.Second,
        EnableTLS:              false,
        InsecureTLS:            false,
        Block:                  true,
    }
}

// NewClient creates a new gRPC client with default configuration
func NewClient() (*Client, error) {
    return NewClientWithConfig(DefaultClientConfig())
}

// NewClientWithConfig creates a new gRPC client with custom configuration
func NewClientWithConfig(config ClientConfig) (*Client, error) {
    // Setup keep alive parameters
    keepAliveParams := keepalive.ClientParameters{
        Time:                config.KeepAliveTime,
        Timeout:             config.KeepAliveTimeout,
        PermitWithoutStream: config.KeepAliveWithoutStream,
    }

    // Setup dial options
    opts := []grpc.DialOption{
        grpc.WithKeepaliveParams(keepAliveParams),
        grpc.WithDefaultCallOptions(
            grpc.MaxCallRecvMsgSize(config.MaxRecvMsgSize),
            grpc.MaxCallSendMsgSize(config.MaxSendMsgSize),
        ),
    }

    // Setup credentials
    if config.EnableTLS {
        var creds credentials.TransportCredentials
        if config.InsecureTLS {
            creds = credentials.NewTLS(&tls.Config{
                InsecureSkipVerify: true,
                ServerName:         config.ServerNameOverride,
            })
        } else if config.CertFile != "" {
            var err error
            creds, err = credentials.NewClientTLSFromFile(config.CertFile, config.ServerNameOverride)
            if err != nil {
                return nil, fmt.Errorf("failed to load TLS credentials: %w", err)
            }
        } else {
            creds = credentials.NewTLS(&tls.Config{
                ServerName: config.ServerNameOverride,
            })
        }
        opts = append(opts, grpc.WithTransportCredentials(creds))
    } else {
        opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
    }

    // Setup blocking connection (fix: Block=true means use grpc.WithBlock())
    if config.Block {
        opts = append(opts, grpc.WithBlock())
    }

    // Create connection
    address := fmt.Sprintf("%s:%d", config.Host, config.Port)
    ctx, cancel := context.WithTimeout(context.Background(), config.ConnectionTimeout)
    defer cancel()

    conn, err := grpc.DialContext(ctx, address, opts...)
    if err != nil {
        return nil, fmt.Errorf("failed to connect to gRPC server at %s: %w", address, err)
    }

    return &Client{
        conn:   conn,
        config: config,
    }, nil
}

// GetConnection returns the underlying gRPC connection
func (c *Client) GetConnection() *grpc.ClientConn {
    return c.conn
}

// Close closes the gRPC client connection
func (c *Client) Close() error {
    if c.conn != nil {
        return c.conn.Close()
    }
    return nil
}

// IsConnected checks if the client is connected to the server
func (c *Client) IsConnected() bool {
    if c.conn == nil {
        return false
    }
    state := c.conn.GetState()
    return state == connectivity.Idle || state == connectivity.Ready
}

// GetAddress returns the server address
func (c *Client) GetAddress() string {
    return fmt.Sprintf("%s:%d", c.config.Host, c.config.Port)
}

// GetConfig returns the client configuration
func (c *Client) GetConfig() ClientConfig {
    return c.config
}

// WaitForConnection waits for the connection to be ready
func (c *Client) WaitForConnection(ctx context.Context) error {
    if c.conn == nil {
        return fmt.Errorf("connection is nil")
    }

    // Wait for connection to be ready
    for {
        state := c.conn.GetState()
        if state == connectivity.Ready {
            return nil
        }
        if state == connectivity.TransientFailure || state == connectivity.Shutdown {
            return fmt.Errorf("connection failed with state: %v", state)
        }

        // Wait for state change or context cancellation
        if !c.conn.WaitForStateChange(ctx, state) {
            return ctx.Err()
        }
    }
}

// HealthCheck provides a basic health check for the gRPC client connection
func (c *Client) HealthCheck() error {
    if c.conn == nil {
        return status.Error(codes.Unavailable, "gRPC connection is nil")
    }

    state := c.conn.GetState()
    switch state {
    case connectivity.Ready, connectivity.Idle:
        return nil
    case connectivity.Connecting:
        return status.Error(codes.Unavailable, "gRPC connection is still connecting")
    case connectivity.TransientFailure:
        return status.Error(codes.Unavailable, "gRPC connection has transient failure")
    case connectivity.Shutdown:
        return status.Error(codes.Unavailable, "gRPC connection is shutdown")
    default:
        return status.Error(codes.Unknown, fmt.Sprintf("unknown connection state: %v", state))
    }
}