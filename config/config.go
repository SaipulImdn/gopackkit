// Package config provides functionality for loading configuration from environment variables
// with support for struct tags and default values.
package config

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
)

// Load loads configuration from environment variables and applies defaults
func Load(config interface{}) error {
	return LoadFromEnv(config)
}

// LoadFromEnv loads configuration from environment variables
func LoadFromEnv(config interface{}) error {
	v := reflect.ValueOf(config)
	if v.Kind() != reflect.Ptr || v.Elem().Kind() != reflect.Struct {
		return fmt.Errorf("config must be a pointer to a struct")
	}

	return loadStruct(v.Elem(), reflect.TypeOf(config).Elem())
}

// LoadFromJSON loads configuration from a JSON file
func LoadFromJSON(filename string, config interface{}) error {
	data, err := os.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("failed to read config file: %w", err)
	}

	if err := json.Unmarshal(data, config); err != nil {
		return fmt.Errorf("failed to unmarshal JSON config: %w", err)
	}

	// Apply environment variables and defaults after loading JSON
	return LoadFromEnv(config)
}

func loadStruct(v reflect.Value, t reflect.Type) error {
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		fieldType := t.Field(i)

		// Skip unexported fields
		if !field.CanSet() {
			continue
		}

		// Handle nested structs
		if field.Kind() == reflect.Struct {
			if err := loadStruct(field, fieldType.Type); err != nil {
				return err
			}
			continue
		}

		// Get environment variable name
		envName := fieldType.Tag.Get("env")
		if envName == "" {
			// Use field name in uppercase as default
			envName = strings.ToUpper(fieldType.Name)
		}

		// Get environment variable value
		envValue := os.Getenv(envName)

		// Use default value if env var is not set
		if envValue == "" {
			defaultValue := fieldType.Tag.Get("default")
			if defaultValue != "" {
				envValue = defaultValue
			}
		}

		// Set field value
		if envValue != "" {
			if err := setFieldValue(field, envValue); err != nil {
				return fmt.Errorf("failed to set field %s: %w", fieldType.Name, err)
			}
		}
	}

	return nil
}

func setFieldValue(field reflect.Value, value string) error {
	switch field.Kind() {
	case reflect.String:
		field.SetString(value)

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		intValue, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return fmt.Errorf("invalid integer value: %s", value)
		}
		field.SetInt(intValue)

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		uintValue, err := strconv.ParseUint(value, 10, 64)
		if err != nil {
			return fmt.Errorf("invalid unsigned integer value: %s", value)
		}
		field.SetUint(uintValue)

	case reflect.Float32, reflect.Float64:
		floatValue, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return fmt.Errorf("invalid float value: %s", value)
		}
		field.SetFloat(floatValue)

	case reflect.Bool:
		boolValue, err := strconv.ParseBool(value)
		if err != nil {
			return fmt.Errorf("invalid boolean value: %s", value)
		}
		field.SetBool(boolValue)

	case reflect.Slice:
		// Handle string slices (comma-separated values)
		if field.Type().Elem().Kind() == reflect.String {
			values := strings.Split(value, ",")
			slice := reflect.MakeSlice(field.Type(), len(values), len(values))
			for i, v := range values {
				slice.Index(i).SetString(strings.TrimSpace(v))
			}
			field.Set(slice)
		} else {
			return fmt.Errorf("unsupported slice type: %s", field.Type())
		}

	default:
		return fmt.Errorf("unsupported field type: %s", field.Kind())
	}

	return nil
}

// MustLoad loads configuration and panics on error
func MustLoad(config interface{}) {
	if err := Load(config); err != nil {
		panic(fmt.Sprintf("failed to load configuration: %v", err))
	}
}

// MustLoadFromJSON loads configuration from JSON and panics on error
func MustLoadFromJSON(filename string, config interface{}) {
	if err := LoadFromJSON(filename, config); err != nil {
		panic(fmt.Sprintf("failed to load configuration from JSON: %v", err))
	}
}
