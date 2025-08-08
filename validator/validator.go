package validator

import (
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

// Constants for common error messages
const (
	ErrMsgValueMustBeString = "value must be a string"
	ErrMsgFieldRequired     = "field is required"
	ErrMsgInvalidMinValue   = "invalid min value"
	ErrMsgInvalidMaxValue   = "invalid max value"
	ErrMsgInvalidLenValue   = "invalid length value"
)

// ValidationError represents a validation error
type ValidationError struct {
	Field   string
	Value   interface{}
	Tag     string
	Message string
}

func (e ValidationError) Error() string {
	return fmt.Sprintf("validation failed for field '%s': %s", e.Field, e.Message)
}

// ValidationErrors represents multiple validation errors
type ValidationErrors []ValidationError

func (ve ValidationErrors) Error() string {
	if len(ve) == 0 {
		return ""
	}
	
	var messages []string
	for _, err := range ve {
		messages = append(messages, err.Error())
	}
	
	return strings.Join(messages, "; ")
}

// Validate validates a struct based on struct tags
func Validate(v interface{}) error {
	return ValidateStruct(v)
}

// ValidateStruct validates a struct and returns validation errors
func ValidateStruct(s interface{}) error {
	val := reflect.ValueOf(s)
	typ := reflect.TypeOf(s)
	
	// Handle pointers
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
		typ = typ.Elem()
	}
	
	if val.Kind() != reflect.Struct {
		return fmt.Errorf("expected struct, got %s", val.Kind())
	}
	
	var errors ValidationErrors
	
	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldType := typ.Field(i)
		
		// Skip unexported fields
		if !field.CanInterface() {
			continue
		}
		
		// Get validation tag
		validateTag := fieldType.Tag.Get("validate")
		if validateTag == "" {
			continue
		}
		
		// Validate field
		if err := validateField(fieldType.Name, field.Interface(), validateTag); err != nil {
			if ve, ok := err.(ValidationErrors); ok {
				errors = append(errors, ve...)
			} else if e, ok := err.(ValidationError); ok {
				errors = append(errors, e)
			}
		}
	}
	
	if len(errors) > 0 {
		return errors
	}
	
	return nil
}

// validateField validates a single field
func validateField(fieldName string, value interface{}, tag string) error {
	var errors ValidationErrors
	
	// Split validation rules
	rules := strings.Split(tag, ",")
	
	for _, rule := range rules {
		rule = strings.TrimSpace(rule)
		if rule == "" {
			continue
		}
		
		// Parse rule
		parts := strings.SplitN(rule, "=", 2)
		ruleName := parts[0]
		var ruleValue string
		if len(parts) > 1 {
			ruleValue = parts[1]
		}
		
		// Apply validation rule
		if err := applyValidationRule(fieldName, value, ruleName, ruleValue); err != nil {
			errors = append(errors, *err)
		}
	}
	
	if len(errors) > 0 {
		return errors
	}
	
	return nil
}

// applyValidationRule applies a specific validation rule
func applyValidationRule(fieldName string, value interface{}, ruleName, ruleValue string) *ValidationError {
	switch ruleName {
	case "required":
		return validateRequired(fieldName, value)
	case "min":
		return validateMin(fieldName, value, ruleValue)
	case "max":
		return validateMax(fieldName, value, ruleValue)
	case "email":
		return validateEmail(fieldName, value)
	case "email_safe":
		return validateEmailSafe(fieldName, value)
	case "phone":
		return validatePhone(fieldName, value)
	case "phone_id":
		return validatePhoneIndonesia(fieldName, value)
	case "url":
		return validateURL(fieldName, value)
	case "alpha":
		return validateAlpha(fieldName, value)
	case "numeric":
		return validateNumeric(fieldName, value)
	case "len":
		return validateLength(fieldName, value, ruleValue)
	case "alphanumeric":
		return validateAlphanumeric(fieldName, value)
	case "no_special_chars":
		return validateNoSpecialChars(fieldName, value)
	default:
		return &ValidationError{
			Field:   fieldName,
			Value:   value,
			Tag:     ruleName,
			Message: fmt.Sprintf("unknown validation rule: %s", ruleName),
		}
	}
}

// validateRequired checks if a value is not empty
func validateRequired(fieldName string, value interface{}) *ValidationError {
	if isEmpty(value) {
		return &ValidationError{
			Field:   fieldName,
			Value:   value,
			Tag:     "required",
			Message: "field is required",
		}
	}
	return nil
}

// validateMin validates minimum value/length
func validateMin(fieldName string, value interface{}, minStr string) *ValidationError {
	min, err := strconv.Atoi(minStr)
	if err != nil {
		return &ValidationError{
			Field:   fieldName,
			Value:   value,
			Tag:     "min",
			Message: "invalid min value",
		}
	}
	
	switch v := value.(type) {
	case string:
		if len(v) < min {
			return &ValidationError{
				Field:   fieldName,
				Value:   value,
				Tag:     "min",
				Message: fmt.Sprintf("length must be at least %d", min),
			}
		}
	case int, int8, int16, int32, int64:
		val := reflect.ValueOf(v).Int()
		if int(val) < min {
			return &ValidationError{
				Field:   fieldName,
				Value:   value,
				Tag:     "min",
				Message: fmt.Sprintf("value must be at least %d", min),
			}
		}
	}
	
	return nil
}

// validateMax validates maximum value/length
func validateMax(fieldName string, value interface{}, maxStr string) *ValidationError {
	max, err := strconv.Atoi(maxStr)
	if err != nil {
		return &ValidationError{
			Field:   fieldName,
			Value:   value,
			Tag:     "max",
			Message: "invalid max value",
		}
	}
	
	switch v := value.(type) {
	case string:
		if len(v) > max {
			return &ValidationError{
				Field:   fieldName,
				Value:   value,
				Tag:     "max",
				Message: fmt.Sprintf("length must not exceed %d", max),
			}
		}
	case int, int8, int16, int32, int64:
		val := reflect.ValueOf(v).Int()
		if int(val) > max {
			return &ValidationError{
				Field:   fieldName,
				Value:   value,
				Tag:     "max",
				Message: fmt.Sprintf("value must not exceed %d", max),
			}
		}
	}
	
	return nil
}

// validateEmail validates email format
func validateEmail(fieldName string, value interface{}) *ValidationError {
	str, ok := value.(string)
	if !ok {
		return &ValidationError{
			Field:   fieldName,
			Value:   value,
			Tag:     "email",
			Message: "value must be a string",
		}
	}
	
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(str) {
		return &ValidationError{
			Field:   fieldName,
			Value:   value,
			Tag:     "email",
			Message: "invalid email format",
		}
	}
	
	return nil
}

// validateURL validates URL format
func validateURL(fieldName string, value interface{}) *ValidationError {
	str, ok := value.(string)
	if !ok {
		return &ValidationError{
			Field:   fieldName,
			Value:   value,
			Tag:     "url",
			Message: "value must be a string",
		}
	}
	
	urlRegex := regexp.MustCompile(`^https?://[^\s/$.?#].[^\s]*$`)
	if !urlRegex.MatchString(str) {
		return &ValidationError{
			Field:   fieldName,
			Value:   value,
			Tag:     "url",
			Message: "invalid URL format",
		}
	}
	
	return nil
}

// validateAlpha validates alphabetic characters only
func validateAlpha(fieldName string, value interface{}) *ValidationError {
	str, ok := value.(string)
	if !ok {
		return &ValidationError{
			Field:   fieldName,
			Value:   value,
			Tag:     "alpha",
			Message: "value must be a string",
		}
	}
	
	alphaRegex := regexp.MustCompile(`^[a-zA-Z]+$`)
	if !alphaRegex.MatchString(str) {
		return &ValidationError{
			Field:   fieldName,
			Value:   value,
			Tag:     "alpha",
			Message: "value must contain only alphabetic characters",
		}
	}
	
	return nil
}

// validateNumeric validates numeric characters only
func validateNumeric(fieldName string, value interface{}) *ValidationError {
	str, ok := value.(string)
	if !ok {
		return &ValidationError{
			Field:   fieldName,
			Value:   value,
			Tag:     "numeric",
			Message: "value must be a string",
		}
	}
	
	numericRegex := regexp.MustCompile(`^[0-9]+$`)
	if !numericRegex.MatchString(str) {
		return &ValidationError{
			Field:   fieldName,
			Value:   value,
			Tag:     "numeric",
			Message: "value must contain only numeric characters",
		}
	}
	
	return nil
}

// validateLength validates exact length
func validateLength(fieldName string, value interface{}, lenStr string) *ValidationError {
	expectedLen, err := strconv.Atoi(lenStr)
	if err != nil {
		return &ValidationError{
			Field:   fieldName,
			Value:   value,
			Tag:     "len",
			Message: "invalid length value",
		}
	}
	
	str, ok := value.(string)
	if !ok {
		return &ValidationError{
			Field:   fieldName,
			Value:   value,
			Tag:     "len",
			Message: "value must be a string",
		}
	}
	
	if len(str) != expectedLen {
		return &ValidationError{
			Field:   fieldName,
			Value:   value,
			Tag:     "len",
			Message: fmt.Sprintf("length must be exactly %d", expectedLen),
		}
	}
	
	return nil
}

// isEmpty checks if a value is considered empty
func isEmpty(value interface{}) bool {
	if value == nil {
		return true
	}
	
	v := reflect.ValueOf(value)
	switch v.Kind() {
	case reflect.String:
		return v.String() == ""
	case reflect.Slice, reflect.Map, reflect.Array:
		return v.Len() == 0
	case reflect.Ptr, reflect.Interface:
		return v.IsNil()
	default:
		return false
	}
}

// validateEmailSafe validates email format with secure pattern (no regex injection risk)
func validateEmailSafe(fieldName string, value interface{}) *ValidationError {
	str, ok := value.(string)
	if !ok {
		return &ValidationError{
			Field:   fieldName,
			Value:   value,
			Tag:     "email_safe",
			Message: "value must be a string",
		}
	}
	
	// Simple but secure email validation
	// Check basic format: something@something.domain
	if !strings.Contains(str, "@") {
		return &ValidationError{
			Field:   fieldName,
			Value:   value,
			Tag:     "email_safe",
			Message: "invalid email format: missing @ symbol",
		}
	}
	
	parts := strings.Split(str, "@")
	if len(parts) != 2 {
		return &ValidationError{
			Field:   fieldName,
			Value:   value,
			Tag:     "email_safe",
			Message: "invalid email format: multiple @ symbols",
		}
	}
	
	localPart := parts[0]
	domainPart := parts[1]
	
	// Validate local part (before @)
	if len(localPart) == 0 || len(localPart) > 64 {
		return &ValidationError{
			Field:   fieldName,
			Value:   value,
			Tag:     "email_safe",
			Message: "invalid email format: local part length must be 1-64 characters",
		}
	}
	
	// Validate domain part (after @)
	if len(domainPart) == 0 || len(domainPart) > 255 {
		return &ValidationError{
			Field:   fieldName,
			Value:   value,
			Tag:     "email_safe",
			Message: "invalid email format: domain part length must be 1-255 characters",
		}
	}
	
	// Domain must contain at least one dot
	if !strings.Contains(domainPart, ".") {
		return &ValidationError{
			Field:   fieldName,
			Value:   value,
			Tag:     "email_safe",
			Message: "invalid email format: domain must contain at least one dot",
		}
	}
	
	// Check for invalid characters (basic security check)
	invalidChars := []string{" ", "\t", "\n", "\r", "<", ">", "[", "]", "\\", ",", ";", ":"}
	for _, char := range invalidChars {
		if strings.Contains(str, char) {
			return &ValidationError{
				Field:   fieldName,
				Value:   value,
				Tag:     "email_safe",
				Message: "invalid email format: contains invalid characters",
			}
		}
	}
	
	return nil
}

// validatePhone validates phone number (10-15 digits)
func validatePhone(fieldName string, value interface{}) *ValidationError {
	str, ok := value.(string)
	if !ok {
		return &ValidationError{
			Field:   fieldName,
			Value:   value,
			Tag:     "phone",
			Message: "value must be a string",
		}
	}
	
	// Remove common phone number separators
	cleaned := strings.ReplaceAll(str, " ", "")
	cleaned = strings.ReplaceAll(cleaned, "-", "")
	cleaned = strings.ReplaceAll(cleaned, "(", "")
	cleaned = strings.ReplaceAll(cleaned, ")", "")
	cleaned = strings.ReplaceAll(cleaned, "+", "")
	
	// Check if it's all digits
	if !isAllDigits(cleaned) {
		return &ValidationError{
			Field:   fieldName,
			Value:   value,
			Tag:     "phone",
			Message: "phone number must contain only digits and valid separators",
		}
	}
	
	// Check length (10-15 digits as per international standards)
	if len(cleaned) < 10 || len(cleaned) > 15 {
		return &ValidationError{
			Field:   fieldName,
			Value:   value,
			Tag:     "phone",
			Message: "phone number must be between 10 and 15 digits",
		}
	}
	
	return nil
}

// validatePhoneIndonesia validates Indonesian phone number format
func validatePhoneIndonesia(fieldName string, value interface{}) *ValidationError {
	str, ok := value.(string)
	if !ok {
		return &ValidationError{
			Field:   fieldName,
			Value:   value,
			Tag:     "phone_id",
			Message: "value must be a string",
		}
	}
	
	// Remove common separators
	cleaned := strings.ReplaceAll(str, " ", "")
	cleaned = strings.ReplaceAll(cleaned, "-", "")
	cleaned = strings.ReplaceAll(cleaned, "(", "")
	cleaned = strings.ReplaceAll(cleaned, ")", "")
	
	// Handle Indonesian country code
	if strings.HasPrefix(cleaned, "+62") {
		cleaned = "0" + cleaned[3:] // Convert +62xxx to 0xxx
	} else if strings.HasPrefix(cleaned, "62") && len(cleaned) > 10 {
		cleaned = "0" + cleaned[2:] // Convert 62xxx to 0xxx
	}
	
	// Check if it's all digits
	if !isAllDigits(cleaned) {
		return &ValidationError{
			Field:   fieldName,
			Value:   value,
			Tag:     "phone_id",
			Message: "Indonesian phone number must contain only digits",
		}
	}
	
	// Indonesian mobile numbers: 08xx-xxxx-xxxx (10-13 digits)
	// Indonesian landline: 0xx-xxxx-xxxx (10-11 digits)
	if !strings.HasPrefix(cleaned, "0") {
		return &ValidationError{
			Field:   fieldName,
			Value:   value,
			Tag:     "phone_id",
			Message: "Indonesian phone number must start with 0",
		}
	}
	
	if len(cleaned) < 10 || len(cleaned) > 13 {
		return &ValidationError{
			Field:   fieldName,
			Value:   value,
			Tag:     "phone_id",
			Message: "Indonesian phone number must be between 10 and 13 digits",
		}
	}
	
	return nil
}

// validateAlphanumeric validates that string contains only letters and numbers
func validateAlphanumeric(fieldName string, value interface{}) *ValidationError {
	str, ok := value.(string)
	if !ok {
		return &ValidationError{
			Field:   fieldName,
			Value:   value,
			Tag:     "alphanumeric",
			Message: "value must be a string",
		}
	}
	
	for _, char := range str {
		if !((char >= 'a' && char <= 'z') || (char >= 'A' && char <= 'Z') || (char >= '0' && char <= '9')) {
			return &ValidationError{
				Field:   fieldName,
				Value:   value,
				Tag:     "alphanumeric",
				Message: "value must contain only letters and numbers",
			}
		}
	}
	
	return nil
}

// validateNoSpecialChars validates that string doesn't contain potentially dangerous characters
func validateNoSpecialChars(fieldName string, value interface{}) *ValidationError {
	str, ok := value.(string)
	if !ok {
		return &ValidationError{
			Field:   fieldName,
			Value:   value,
			Tag:     "no_special_chars",
			Message: "value must be a string",
		}
	}
	
	// List of potentially dangerous characters for security
	dangerousChars := []string{
		"<", ">", "&", "\"", "'", "/", "\\", ";", ":", "|", 
		"*", "?", "[", "]", "{", "}", "$", "`", "!", "@",
		"#", "%", "^", "(", ")", "=", "+", "~",
	}
	
	for _, char := range dangerousChars {
		if strings.Contains(str, char) {
			return &ValidationError{
				Field:   fieldName,
				Value:   value,
				Tag:     "no_special_chars",
				Message: fmt.Sprintf("value contains forbidden character: %s", char),
			}
		}
	}
	
	return nil
}

// isAllDigits checks if string contains only digits
func isAllDigits(s string) bool {
	for _, char := range s {
		if char < '0' || char > '9' {
			return false
		}
	}
	return true
}
