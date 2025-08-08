package validator

import (
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
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
	case "url":
		return validateURL(fieldName, value)
	case "alpha":
		return validateAlpha(fieldName, value)
	case "numeric":
		return validateNumeric(fieldName, value)
	case "len":
		return validateLength(fieldName, value, ruleValue)
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
