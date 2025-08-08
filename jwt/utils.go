package jwt

import (
	"encoding/base64"
	"encoding/json"
	"strings"
	"time"
)

// DecodeTokenPayload decodes JWT payload without verification (for debugging)
func DecodeTokenPayload(tokenString string) (map[string]interface{}, error) {
	parts := strings.Split(tokenString, ".")
	if len(parts) != 3 {
		return nil, ErrInvalidToken
	}

	// Decode payload (second part)
	payload, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return nil, err
	}

	var claims map[string]interface{}
	if err := json.Unmarshal(payload, &claims); err != nil {
		return nil, err
	}

	return claims, nil
}

// GetTokenHeader returns JWT header information
func GetTokenHeader(tokenString string) (map[string]interface{}, error) {
	parts := strings.Split(tokenString, ".")
	if len(parts) != 3 {
		return nil, ErrInvalidToken
	}

	// Decode header (first part)
	header, err := base64.RawURLEncoding.DecodeString(parts[0])
	if err != nil {
		return nil, err
	}

	var headerMap map[string]interface{}
	if err := json.Unmarshal(header, &headerMap); err != nil {
		return nil, err
	}

	return headerMap, nil
}

// ValidateTokenFormat checks if token has valid JWT format
func ValidateTokenFormat(tokenString string) bool {
	parts := strings.Split(tokenString, ".")
	return len(parts) == 3
}

// ExtractBearerToken extracts token from "Bearer <token>" format
func ExtractBearerToken(authHeader string) string {
	if authHeader == "" {
		return ""
	}

	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		return ""
	}

	return parts[1]
}

// IsTokenNearExpiry checks if token will expire within specified duration
func (tm *TokenManager) IsTokenNearExpiry(tokenString string, threshold time.Duration) bool {
	expiry, err := tm.GetTokenExpiry(tokenString)
	if err != nil {
		return true
	}

	return time.Until(expiry) <= threshold
}

// GetTokenRemainingTime returns remaining time until token expires
func (tm *TokenManager) GetTokenRemainingTime(tokenString string) (time.Duration, error) {
	expiry, err := tm.GetTokenExpiry(tokenString)
	if err != nil {
		return 0, err
	}

	remaining := time.Until(expiry)
	if remaining < 0 {
		return 0, nil
	}

	return remaining, nil
}

// CreateDefaultConfig creates a default JWT configuration
func CreateDefaultConfig(secretKey string) Config {
	return Config{
		SecretKey:       secretKey,
		AccessTokenTTL:  15 * time.Minute,
		RefreshTokenTTL: 7 * 24 * time.Hour,
		Issuer:          "gopackkit",
		Algorithm:       "HS256",
	}
}

// CreateCustomConfig creates a custom JWT configuration
func CreateCustomConfig(secretKey string, accessTTL, refreshTTL time.Duration, issuer string) Config {
	return Config{
		SecretKey:       secretKey,
		AccessTokenTTL:  accessTTL,
		RefreshTokenTTL: refreshTTL,
		Issuer:          issuer,
		Algorithm:       "HS256",
	}
}
