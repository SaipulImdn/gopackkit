package grpc

import (
	"context"
	"fmt"
	"net"
	"sync"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
)

// Server represents a gRPC server instance
type Server struct {
	server   *grpc.Server
	listener net.Listener
	config   Config
	mu       sync.RWMutex
	running  bool
}

// Config holds gRPC server configuration
type Config struct {
	Host                        string        `json:"host" yaml:"host" env:"GRPC_HOST" default:"0.0.0.0"`
	Port                        int           `json:"port" yaml:"port" env:"GRPC_PORT" default:"9000"`
	MaxRecvMsgSize              int           `json:"max_recv_msg_size" yaml:"max_recv_msg_size" env:"GRPC_MAX_RECV_MSG_SIZE" default:"4194304"`     // 4MB
	MaxSendMsgSize              int           `json:"max_send_msg_size" yaml:"max_send_msg_size" env:"GRPC_MAX_SEND_MSG_SIZE" default:"4194304"`     // 4MB
	MaxConnectionIdle           time.Duration `json:"max_connection_idle" yaml:"max_connection_idle" env:"GRPC_MAX_CONNECTION_IDLE" default:"60s"`
	MaxConnectionAge            time.Duration `json:"max_connection_age" yaml:"max_connection_age" env:"GRPC_MAX_CONNECTION_AGE" default:"300s"`
	MaxConnectionAgeGrace       time.Duration `json:"max_connection_age_grace" yaml:"max_connection_age_grace" env:"GRPC_MAX_CONNECTION_AGE_GRACE" default:"10s"`
	KeepAliveTime               time.Duration `json:"keep_alive_time" yaml:"keep_alive_time" env:"GRPC_KEEP_ALIVE_TIME" default:"60s"`
	KeepAliveTimeout            time.Duration `json:"keep_alive_timeout" yaml:"keep_alive_timeout" env:"GRPC_KEEP_ALIVE_TIMEOUT" default:"5s"`
	KeepAliveEnforcementPolicy  bool          `json:"keep_alive_enforcement_policy" yaml:"keep_alive_enforcement_policy" env:"GRPC_KEEP_ALIVE_ENFORCEMENT_POLICY" default:"true"`
	EnableReflection            bool          `json:"enable_reflection" yaml:"enable_reflection" env:"GRPC_ENABLE_REFLECTION" default:"false"`
	EnableTLS                   bool          `json:"enable_tls" yaml:"enable_tls" env:"GRPC_ENABLE_TLS" default:"false"`
	CertFile                    string        `json:"cert_file" yaml:"cert_file" env:"GRPC_CERT_FILE"`
	KeyFile                     string        `json:"key_file" yaml:"key_file" env:"GRPC_KEY_FILE"`
}

// DefaultConfig returns default gRPC server configuration
func DefaultConfig() Config {
	return Config{
		Host:                        "0.0.0.0",
		Port:                        9000,
		MaxRecvMsgSize:             4 * 1024 * 1024, // 4MB
		MaxSendMsgSize:             4 * 1024 * 1024, // 4MB
		MaxConnectionIdle:          60 * time.Second,
		MaxConnectionAge:           300 * time.Second,
		MaxConnectionAgeGrace:      10 * time.Second,
		KeepAliveTime:              60 * time.Second,
		KeepAliveTimeout:           5 * time.Second,
		KeepAliveEnforcementPolicy: true,
		EnableReflection:           false,
		EnableTLS:                  false,
	}
}

// NewServer creates a new gRPC server with default configuration
func NewServer() *Server {
	return NewServerWithConfig(DefaultConfig())
}

// NewServerWithConfig creates a new gRPC server with custom configuration
func NewServerWithConfig(config Config) *Server {
	// Setup keep alive parameters
	keepAliveParams := keepalive.ServerParameters{
		MaxConnectionIdle:     config.MaxConnectionIdle,
		MaxConnectionAge:      config.MaxConnectionAge,
		MaxConnectionAgeGrace: config.MaxConnectionAgeGrace,
		Time:                  config.KeepAliveTime,
		Timeout:               config.KeepAliveTimeout,
	}

	// Setup enforcement policy
	keepAliveEnforcementPolicy := keepalive.EnforcementPolicy{
		MinTime:             30 * time.Second,
		PermitWithoutStream: config.KeepAliveEnforcementPolicy,
	}

	// Setup server options
	opts := []grpc.ServerOption{
		grpc.MaxRecvMsgSize(config.MaxRecvMsgSize),
		grpc.MaxSendMsgSize(config.MaxSendMsgSize),
		grpc.KeepaliveParams(keepAliveParams),
		grpc.KeepaliveEnforcementPolicy(keepAliveEnforcementPolicy),
	}

	// Create gRPC server
	server := grpc.NewServer(opts...)

	// Enable reflection if configured
	if config.EnableReflection {
		reflection.Register(server)
	}

	return &Server{
		server: server,
		config: config,
	}
}

// RegisterService registers a gRPC service
func (s *Server) RegisterService(desc *grpc.ServiceDesc, impl interface{}) {
	s.server.RegisterService(desc, impl)
}

// GetServer returns the underlying gRPC server
func (s *Server) GetServer() *grpc.Server {
	return s.server
}

// Start starts the gRPC server
func (s *Server) Start() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.running {
		return fmt.Errorf("gRPC server is already running")
	}

	// Create listener
	address := fmt.Sprintf("%s:%d", s.config.Host, s.config.Port)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		return fmt.Errorf("failed to listen on %s: %w", address, err)
	}

	s.listener = listener
	s.running = true

	// Start serving
	go func() {
		if err := s.server.Serve(listener); err != nil {
			s.mu.Lock()
			s.running = false
			s.mu.Unlock()
		}
	}()

	return nil
}

// Stop gracefully stops the gRPC server
func (s *Server) Stop() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.running {
		return fmt.Errorf("gRPC server is not running")
	}

	s.server.GracefulStop()
	s.running = false

	return nil
}

// ForceStop forcefully stops the gRPC server
func (s *Server) ForceStop() {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.server.Stop()
	s.running = false
}

// IsRunning returns true if the server is running
func (s *Server) IsRunning() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.running
}

// GetAddress returns the server address
func (s *Server) GetAddress() string {
	if s.listener == nil {
		return fmt.Sprintf("%s:%d", s.config.Host, s.config.Port)
	}
	return s.listener.Addr().String()
}

// GetConfig returns the server configuration
func (s *Server) GetConfig() Config {
	return s.config
}

// Serve starts the server and blocks until it's stopped
func (s *Server) Serve() error {
	if err := s.Start(); err != nil {
		return err
	}

	// Wait for server to stop
	for s.IsRunning() {
		time.Sleep(100 * time.Millisecond)
	}

	return nil
}

// ServeWithContext starts the server and blocks until context is cancelled
func (s *Server) ServeWithContext(ctx context.Context) error {
	if err := s.Start(); err != nil {
		return err
	}

	// Wait for context cancellation or server stop
	for {
		select {
		case <-ctx.Done():
			return s.Stop()
		default:
			if !s.IsRunning() {
				return nil
			}
			time.Sleep(100 * time.Millisecond)
		}
	}
}

// HealthCheck provides a basic health check for the gRPC server
func (s *Server) HealthCheck() error {
	if !s.IsRunning() {
		return status.Error(codes.Unavailable, "gRPC server is not running")
	}
	return nil
}
