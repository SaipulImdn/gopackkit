package password

import (
	"crypto/rand"
	"errors"
	"fmt"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// Manager handles password operations
type Manager struct {
	config Config
}

// Config holds password configuration
type Config struct {
	MinLength      int  `json:"min_length" yaml:"min_length" env:"PASSWORD_MIN_LENGTH" default:"8"`
	MaxLength      int  `json:"max_length" yaml:"max_length" env:"PASSWORD_MAX_LENGTH" default:"128"`
	RequireUpper   bool `json:"require_upper" yaml:"require_upper" env:"PASSWORD_REQUIRE_UPPER" default:"true"`
	RequireLower   bool `json:"require_lower" yaml:"require_lower" env:"PASSWORD_REQUIRE_LOWER" default:"true"`
	RequireDigit   bool `json:"require_digit" yaml:"require_digit" env:"PASSWORD_REQUIRE_DIGIT" default:"true"`
	RequireSpecial bool `json:"require_special" yaml:"require_special" env:"PASSWORD_REQUIRE_SPECIAL" default:"false"`
	BcryptCost     int  `json:"bcrypt_cost" yaml:"bcrypt_cost" env:"PASSWORD_BCRYPT_COST" default:"12"`
}

// PasswordStrength represents password strength level
type PasswordStrength int

const (
	StrengthWeak PasswordStrength = iota
	StrengthFair
	StrengthGood
	StrengthStrong
	StrengthVeryStrong
)

func (ps PasswordStrength) String() string {
	switch ps {
	case StrengthWeak:
		return "Weak"
	case StrengthFair:
		return "Fair"
	case StrengthGood:
		return "Good"
	case StrengthStrong:
		return "Strong"
	case StrengthVeryStrong:
		return "Very Strong"
	default:
		return "Unknown"
	}
}

// PasswordValidation represents password validation result
type PasswordValidation struct {
	Valid       bool             `json:"valid"`
	Strength    PasswordStrength `json:"strength"`
	Score       int              `json:"score"`
	Errors      []string         `json:"errors,omitempty"`
	Suggestions []string         `json:"suggestions,omitempty"`
}

// HashedPassword represents a hashed password with metadata
type HashedPassword struct {
	Hash      string    `json:"hash"`
	Algorithm string    `json:"algorithm"`
	Cost      int       `json:"cost"`
	CreatedAt time.Time `json:"created_at"`
}

var (
	ErrPasswordTooShort   = errors.New("password is too short")
	ErrPasswordTooLong    = errors.New("password is too long")
	ErrPasswordTooWeak    = errors.New("password is too weak")
	ErrInvalidHash        = errors.New("invalid password hash")
	ErrHashingFailed      = errors.New("password hashing failed")
	ErrVerificationFailed = errors.New("password verification failed")
)

// New creates a new password manager with default configuration
func New() *Manager {
	return NewWithConfig(Config{
		MinLength:      8,
		MaxLength:      128,
		RequireUpper:   true,
		RequireLower:   true,
		RequireDigit:   true,
		RequireSpecial: false,
		BcryptCost:     12,
	})
}

// NewWithConfig creates a new password manager with custom configuration
func NewWithConfig(config Config) *Manager {
	// Set defaults if not provided
	if config.MinLength == 0 {
		config.MinLength = 8
	}
	if config.MaxLength == 0 {
		config.MaxLength = 128
	}
	if config.BcryptCost == 0 {
		config.BcryptCost = 12
	}

	// Validate bcrypt cost (4-31)
	if config.BcryptCost < 4 {
		config.BcryptCost = 4
	}
	if config.BcryptCost > 31 {
		config.BcryptCost = 31
	}

	return &Manager{config: config}
}

// Hash creates a bcrypt hash of the password
func (pm *Manager) Hash(password string) (*HashedPassword, error) {
	// Validate password first
	validation := pm.Validate(password)
	if !validation.Valid {
		return nil, fmt.Errorf("%w: %s", ErrPasswordTooWeak, strings.Join(validation.Errors, ", "))
	}

	// Generate bcrypt hash
	hashBytes, err := bcrypt.GenerateFromPassword([]byte(password), pm.config.BcryptCost)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrHashingFailed, err)
	}

	return &HashedPassword{
		Hash:      string(hashBytes),
		Algorithm: "bcrypt",
		Cost:      pm.config.BcryptCost,
		CreatedAt: time.Now(),
	}, nil
}

// HashString creates a bcrypt hash and returns only the hash string
func (pm *Manager) HashString(password string) (string, error) {
	hashedPassword, err := pm.Hash(password)
	if err != nil {
		return "", err
	}
	return hashedPassword.Hash, nil
}

// Verify verifies a password against its hash
func (pm *Manager) Verify(password, hash string) error {
	if hash == "" {
		return ErrInvalidHash
	}

	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return ErrVerificationFailed
		}
		return fmt.Errorf("%w: %v", ErrVerificationFailed, err)
	}

	return nil
}

// VerifyAndCheck verifies password and returns additional information
func (pm *Manager) VerifyAndCheck(password, hash string) (bool, error) {
	err := pm.Verify(password, hash)
	if err != nil {
		if errors.Is(err, ErrVerificationFailed) {
			return false, nil // Wrong password, but no error
		}
		return false, err // Actual error occurred
	}
	return true, nil
}

// Validate validates password according to configured rules
func (pm *Manager) Validate(password string) PasswordValidation {
	validation := PasswordValidation{
		Valid:       true,
		Errors:      []string{},
		Suggestions: []string{},
	}

	// Check length
	if len(password) < pm.config.MinLength {
		validation.Valid = false
		validation.Errors = append(validation.Errors, fmt.Sprintf("Password must be at least %d characters long", pm.config.MinLength))
		validation.Suggestions = append(validation.Suggestions, "Use a longer password")
	}

	if len(password) > pm.config.MaxLength {
		validation.Valid = false
		validation.Errors = append(validation.Errors, fmt.Sprintf("Password must not exceed %d characters", pm.config.MaxLength))
	}

	// Check character requirements
	if pm.config.RequireUpper && !containsUpper(password) {
		validation.Valid = false
		validation.Errors = append(validation.Errors, "Password must contain at least one uppercase letter")
		validation.Suggestions = append(validation.Suggestions, "Add uppercase letters (A-Z)")
	}

	if pm.config.RequireLower && !containsLower(password) {
		validation.Valid = false
		validation.Errors = append(validation.Errors, "Password must contain at least one lowercase letter")
		validation.Suggestions = append(validation.Suggestions, "Add lowercase letters (a-z)")
	}

	if pm.config.RequireDigit && !containsDigit(password) {
		validation.Valid = false
		validation.Errors = append(validation.Errors, "Password must contain at least one digit")
		validation.Suggestions = append(validation.Suggestions, "Add numbers (0-9)")
	}

	if pm.config.RequireSpecial && !containsSpecial(password) {
		validation.Valid = false
		validation.Errors = append(validation.Errors, "Password must contain at least one special character")
		validation.Suggestions = append(validation.Suggestions, "Add special characters (!@#$%^&*)")
	}

	// Calculate strength and score
	validation.Strength, validation.Score = pm.calculateStrength(password)

	return validation
}

// calculateStrength calculates password strength and score
func (pm *Manager) calculateStrength(password string) (PasswordStrength, int) {
	score := 0

	// Length scoring
	length := len(password)
	if length >= 8 {
		score += 1
	}
	if length >= 12 {
		score += 1
	}
	if length >= 16 {
		score += 1
	}

	// Character variety scoring
	if containsLower(password) {
		score += 1
	}
	if containsUpper(password) {
		score += 1
	}
	if containsDigit(password) {
		score += 1
	}
	if containsSpecial(password) {
		score += 2
	}

	// Additional patterns
	if !hasRepeatingChars(password) {
		score += 1
	}
	if !hasSequentialChars(password) {
		score += 1
	}

	// Determine strength based on score
	switch {
	case score >= 9:
		return StrengthVeryStrong, score
	case score >= 7:
		return StrengthStrong, score
	case score >= 5:
		return StrengthGood, score
	case score >= 3:
		return StrengthFair, score
	default:
		return StrengthWeak, score
	}
}

// GenerateRandomPassword generates a random password
func (pm *Manager) GenerateRandomPassword(length int) (string, error) {
	if length < pm.config.MinLength {
		length = pm.config.MinLength
	}
	if length > pm.config.MaxLength {
		length = pm.config.MaxLength
	}

	// Character sets
	lowercase := "abcdefghijklmnopqrstuvwxyz"
	uppercase := "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	digits := "0123456789"
	special := "!@#$%^&*()_+-=[]{}|;:,.<>?"

	var charset string
	var required []byte

	// Build charset and required characters based on config
	if pm.config.RequireLower {
		charset += lowercase
		required = append(required, lowercase[randomInt(len(lowercase))])
	}
	if pm.config.RequireUpper {
		charset += uppercase
		required = append(required, uppercase[randomInt(len(uppercase))])
	}
	if pm.config.RequireDigit {
		charset += digits
		required = append(required, digits[randomInt(len(digits))])
	}
	if pm.config.RequireSpecial {
		charset += special
		required = append(required, special[randomInt(len(special))])
	}

	// If no requirements, use all characters
	if charset == "" {
		charset = lowercase + uppercase + digits
	}

	// Generate password
	password := make([]byte, length)

	// Place required characters first
	for i, char := range required {
		if i < length {
			password[i] = char
		}
	}

	// Fill remaining positions
	for i := len(required); i < length; i++ {
		password[i] = charset[randomInt(len(charset))]
	}

	// Shuffle the password
	for i := range password {
		j := randomInt(len(password))
		password[i], password[j] = password[j], password[i]
	}

	return string(password), nil
}

// NeedsRehash checks if password hash needs to be updated (cost changed)
func (pm *Manager) NeedsRehash(hash string) bool {
	cost, err := bcrypt.Cost([]byte(hash))
	if err != nil {
		return true
	}
	return cost != pm.config.BcryptCost
}

// GetConfig returns current configuration
func (pm *Manager) GetConfig() Config {
	return pm.config
}

// UpdateConfig updates password manager configuration
func (pm *Manager) UpdateConfig(config Config) {
	pm.config = config
}

// Helper functions

func containsUpper(s string) bool {
	for _, char := range s {
		if char >= 'A' && char <= 'Z' {
			return true
		}
	}
	return false
}

func containsLower(s string) bool {
	for _, char := range s {
		if char >= 'a' && char <= 'z' {
			return true
		}
	}
	return false
}

func containsDigit(s string) bool {
	for _, char := range s {
		if char >= '0' && char <= '9' {
			return true
		}
	}
	return false
}

func containsSpecial(s string) bool {
	specialChars := "!@#$%^&*()_+-=[]{}|;:,.<>?"
	for _, char := range s {
		for _, special := range specialChars {
			if char == special {
				return true
			}
		}
	}
	return false
}

func hasRepeatingChars(s string) bool {
	if len(s) < 3 {
		return false
	}

	for i := 0; i < len(s)-2; i++ {
		if s[i] == s[i+1] && s[i+1] == s[i+2] {
			return true
		}
	}
	return false
}

func hasSequentialChars(s string) bool {
	if len(s) < 3 {
		return false
	}

	sequences := []string{
		"abcdefghijklmnopqrstuvwxyz",
		"ABCDEFGHIJKLMNOPQRSTUVWXYZ",
		"0123456789",
		"qwertyuiopasdfghjklzxcvbnm", // keyboard layout
	}

	for _, seq := range sequences {
		if containsSequence(s, seq, 3) || containsSequence(s, reverse(seq), 3) {
			return true
		}
	}
	return false
}

func containsSequence(s, sequence string, length int) bool {
	for i := 0; i <= len(sequence)-length; i++ {
		subseq := sequence[i : i+length]
		if strings.Contains(strings.ToLower(s), subseq) {
			return true
		}
	}
	return false
}

func reverse(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

func randomInt(max int) int {
	if max <= 0 {
		return 0
	}

	bytes := make([]byte, 4)
	_, err := rand.Read(bytes)
	if err != nil {
		// Fallback to time-based randomness
		return int(time.Now().UnixNano()) % max
	}

	// Convert bytes to int
	n := int(bytes[0])<<24 | int(bytes[1])<<16 | int(bytes[2])<<8 | int(bytes[3])
	if n < 0 {
		n = -n
	}
	return n % max
}
