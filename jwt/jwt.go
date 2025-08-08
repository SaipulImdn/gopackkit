package jwt

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// TokenManager handles JWT token operations
type TokenManager struct {
	config Config
}

// Config holds JWT configuration
type Config struct {
	SecretKey       string        `json:"secret_key" yaml:"secret_key" env:"JWT_SECRET_KEY"`
	AccessTokenTTL  time.Duration `json:"access_token_ttl" yaml:"access_token_ttl" env:"JWT_ACCESS_TOKEN_TTL" default:"15m"`
	RefreshTokenTTL time.Duration `json:"refresh_token_ttl" yaml:"refresh_token_ttl" env:"JWT_REFRESH_TOKEN_TTL" default:"7d"`
	Issuer          string        `json:"issuer" yaml:"issuer" env:"JWT_ISSUER" default:"gopackkit"`
	Algorithm       string        `json:"algorithm" yaml:"algorithm" env:"JWT_ALGORITHM" default:"HS256"`
}

// Claims represents JWT claims structure
type Claims struct {
	UserID   string                 `json:"user_id"`
	Username string                 `json:"username,omitempty"`
	Email    string                 `json:"email,omitempty"`
	Roles    []string               `json:"roles,omitempty"`
	Custom   map[string]interface{} `json:"custom,omitempty"`
	jwt.RegisteredClaims
}

// TokenPair represents access and refresh token pair
type TokenPair struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	ExpiresAt    time.Time `json:"expires_at"`
	TokenType    string    `json:"token_type"`
}

// TokenInfo represents decoded token information
type TokenInfo struct {
	Claims    *Claims   `json:"claims"`
	Valid     bool      `json:"valid"`
	Expired   bool      `json:"expired"`
	IssuedAt  time.Time `json:"issued_at"`
	ExpiresAt time.Time `json:"expires_at"`
}

var (
	ErrInvalidToken     = errors.New("invalid token")
	ErrExpiredToken     = errors.New("token has expired")
	ErrInvalidSignature = errors.New("invalid token signature")
	ErrMissingSecretKey = errors.New("secret key is required")
)

// New creates a new JWT token manager
func New(config Config) (*TokenManager, error) {
	if config.SecretKey == "" {
		return nil, ErrMissingSecretKey
	}

	// Set default values
	if config.AccessTokenTTL == 0 {
		config.AccessTokenTTL = 15 * time.Minute
	}
	if config.RefreshTokenTTL == 0 {
		config.RefreshTokenTTL = 7 * 24 * time.Hour
	}
	if config.Issuer == "" {
		config.Issuer = "gopackkit"
	}
	if config.Algorithm == "" {
		config.Algorithm = "HS256"
	}

	return &TokenManager{
		config: config,
	}, nil
}

// GenerateTokenPair generates both access and refresh tokens
func (tm *TokenManager) GenerateTokenPair(userID, username, email string, roles []string, customClaims map[string]interface{}) (*TokenPair, error) {
	now := time.Now()
	accessExpiry := now.Add(tm.config.AccessTokenTTL)
	refreshExpiry := now.Add(tm.config.RefreshTokenTTL)

	// Generate access token
	accessClaims := &Claims{
		UserID:   userID,
		Username: username,
		Email:    email,
		Roles:    roles,
		Custom:   customClaims,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(accessExpiry),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			Issuer:    tm.config.Issuer,
			Subject:   userID,
			ID:        fmt.Sprintf("access_%d", now.Unix()),
		},
	}

	accessToken, err := tm.generateToken(accessClaims)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	// Generate refresh token
	refreshClaims := &Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(refreshExpiry),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			Issuer:    tm.config.Issuer,
			Subject:   userID,
			ID:        fmt.Sprintf("refresh_%d", now.Unix()),
		},
	}

	refreshToken, err := tm.generateToken(refreshClaims)
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresAt:    accessExpiry,
		TokenType:    "Bearer",
	}, nil
}

// GenerateAccessToken generates only access token
func (tm *TokenManager) GenerateAccessToken(userID, username, email string, roles []string, customClaims map[string]interface{}) (string, error) {
	now := time.Now()
	expiry := now.Add(tm.config.AccessTokenTTL)

	claims := &Claims{
		UserID:   userID,
		Username: username,
		Email:    email,
		Roles:    roles,
		Custom:   customClaims,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiry),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			Issuer:    tm.config.Issuer,
			Subject:   userID,
			ID:        fmt.Sprintf("access_%d", now.Unix()),
		},
	}

	return tm.generateToken(claims)
}

// ValidateToken validates and parses a JWT token
func (tm *TokenManager) ValidateToken(tokenString string) (*TokenInfo, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(tm.config.SecretKey), nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return &TokenInfo{
				Valid:   false,
				Expired: true,
			}, ErrExpiredToken
		}
		return &TokenInfo{
			Valid:   false,
			Expired: false,
		}, ErrInvalidToken
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return &TokenInfo{
			Valid:   false,
			Expired: false,
		}, ErrInvalidToken
	}

	return &TokenInfo{
		Claims:    claims,
		Valid:     true,
		Expired:   false,
		IssuedAt:  claims.IssuedAt.Time,
		ExpiresAt: claims.ExpiresAt.Time,
	}, nil
}

// RefreshToken generates a new access token from a valid refresh token
func (tm *TokenManager) RefreshToken(refreshTokenString string) (*TokenPair, error) {
	tokenInfo, err := tm.ValidateToken(refreshTokenString)
	if err != nil {
		return nil, fmt.Errorf("invalid refresh token: %w", err)
	}

	if !tokenInfo.Valid {
		return nil, ErrInvalidToken
	}

	// Generate new token pair
	return tm.GenerateTokenPair(
		tokenInfo.Claims.UserID,
		tokenInfo.Claims.Username,
		tokenInfo.Claims.Email,
		tokenInfo.Claims.Roles,
		tokenInfo.Claims.Custom,
	)
}

// ExtractUserID extracts user ID from token without full validation
func (tm *TokenManager) ExtractUserID(tokenString string) (string, error) {
	token, _, err := new(jwt.Parser).ParseUnverified(tokenString, &Claims{})
	if err != nil {
		return "", fmt.Errorf("failed to parse token: %w", err)
	}

	claims, ok := token.Claims.(*Claims)
	if !ok {
		return "", ErrInvalidToken
	}

	return claims.UserID, nil
}

// GetTokenExpiry returns token expiration time
func (tm *TokenManager) GetTokenExpiry(tokenString string) (time.Time, error) {
	tokenInfo, err := tm.ValidateToken(tokenString)
	if err != nil {
		return time.Time{}, err
	}

	return tokenInfo.ExpiresAt, nil
}

// IsTokenExpired checks if token is expired
func (tm *TokenManager) IsTokenExpired(tokenString string) bool {
	tokenInfo, err := tm.ValidateToken(tokenString)
	if err != nil {
		return true
	}

	return tokenInfo.Expired || time.Now().After(tokenInfo.ExpiresAt)
}

// generateToken creates a signed JWT token
func (tm *TokenManager) generateToken(claims *Claims) (string, error) {
	var signingMethod jwt.SigningMethod

	switch tm.config.Algorithm {
	case "HS256":
		signingMethod = jwt.SigningMethodHS256
	case "HS384":
		signingMethod = jwt.SigningMethodHS384
	case "HS512":
		signingMethod = jwt.SigningMethodHS512
	default:
		signingMethod = jwt.SigningMethodHS256
	}

	token := jwt.NewWithClaims(signingMethod, claims)
	return token.SignedString([]byte(tm.config.SecretKey))
}

// GetConfig returns current configuration
func (tm *TokenManager) GetConfig() Config {
	return tm.config
}

// UpdateConfig updates token manager configuration
func (tm *TokenManager) UpdateConfig(config Config) error {
	if config.SecretKey == "" {
		return ErrMissingSecretKey
	}

	tm.config = config
	return nil
}
