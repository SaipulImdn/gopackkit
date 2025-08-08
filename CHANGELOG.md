# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- Complete CI/CD pipeline with GitHub Actions
- Docker support with multi-stage builds
- Comprehensive Makefile for development
- Security scanning with gosec and govulncheck
- Code quality checks with golangci-lint
- Contributing guidelines and documentation

## [v0.1.0] - 2025-08-09

### Added
- **Logger Module**: Multi-backend logging support (Simple, Logrus, Zap)
  - Configurable log levels and formats
  - JSON and text output formats
  - Thread-safe operations
  
- **HTTP Client Module**: Robust HTTP client with retry mechanism
  - Automatic retry with exponential backoff
  - Configurable timeout and retry attempts
  - Request/Response middleware support
  - JSON POST convenience methods
  
- **Config Module**: Environment variable loader with struct tag support
  - Automatic environment variable mapping
  - Default values support
  - Type conversion for various Go types
  - Struct tag based configuration
  
- **Validator Module**: Secure data validation focused on preventing vulnerabilities
  - Phone number validation (10-15 digits)
  - Safe email validation (basic @ check without regex injection)
  - Alphanumeric validation
  - No special characters validation
  - Required field validation
  
- **MinIO Module**: MinIO client for object storage operations
  - Presigned URL generation (GET, PUT, POST)
  - Object upload and download
  - Bucket operations
  - Configurable expiration times
  
- **JWT Module**: Complete JWT token management system
  - Token generation with custom claims
  - Token validation and parsing
  - Refresh token mechanism
  - Configurable expiration times
  - HMAC-SHA256 signing
  
- **Password Module**: Secure password hashing and validation
  - Bcrypt password hashing
  - Password strength validation
  - Configurable validation rules
  - Random password generation
  - Hash migration detection

### Security
- All validation functions designed to prevent injection attacks
- No regex patterns that could be vulnerable to ReDoS
- Secure password hashing with bcrypt
- JWT tokens signed with HMAC-SHA256
- Security-focused design throughout all modules

### Documentation
- Comprehensive README.md with examples for all modules
- Individual module documentation
- Complete web API example showing integration
- Security considerations documented
- Installation and usage instructions

### Dependencies
- golang-jwt/jwt v5.2.1 for JWT token management
- minio-go v7.0.70 for MinIO object storage
- logrus v1.9.3 for structured logging
- zap v1.27.0 for high-performance logging
- bcrypt from golang.org/x/crypto for password hashing

[Unreleased]: https://github.com/saipulimdn/gopackkit/compare/v0.1.0...HEAD
[v0.1.0]: https://github.com/saipulimdn/gopackkit/releases/tag/v0.1.0
