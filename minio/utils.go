package minio

import (
	"fmt"
	"time"
)

// DefaultPresignedURLOptions returns default options for presigned URLs
func DefaultPresignedURLOptions() *PresignedURLOptions {
	return &PresignedURLOptions{
		Expiry: 24 * time.Hour,
	}
}

// ShortPresignedURLOptions returns options for short-lived presigned URLs (1 hour)
func ShortPresignedURLOptions() *PresignedURLOptions {
	return &PresignedURLOptions{
		Expiry: 1 * time.Hour,
	}
}

// LongPresignedURLOptions returns options for long-lived presigned URLs (7 days)
func LongPresignedURLOptions() *PresignedURLOptions {
	return &PresignedURLOptions{
		Expiry: 7 * 24 * time.Hour,
	}
}

// CustomPresignedURLOptions creates options with custom expiry
func CustomPresignedURLOptions(expiry time.Duration) *PresignedURLOptions {
	return &PresignedURLOptions{
		Expiry: expiry,
	}
}

// WithContentType adds content type header to presigned URL options
func (opts *PresignedURLOptions) WithContentType(contentType string) *PresignedURLOptions {
	if opts.ExtraHeaders == nil {
		opts.ExtraHeaders = make(map[string]string)
	}
	opts.ExtraHeaders["Content-Type"] = contentType
	return opts
}

// ValidateBucketName validates bucket name according to AWS S3 rules
func ValidateBucketName(bucketName string) error {
	if len(bucketName) < 3 || len(bucketName) > 63 {
		return fmt.Errorf("bucket name must be between 3 and 63 characters long")
	}
	
	// Basic validation - you can extend this with more AWS S3 rules
	if bucketName[0] == '-' || bucketName[len(bucketName)-1] == '-' {
		return fmt.Errorf("bucket name cannot start or end with a hyphen")
	}
	
	return nil
}

// ValidateObjectName validates object name
func ValidateObjectName(objectName string) error {
	if objectName == "" {
		return fmt.Errorf("object name cannot be empty")
	}
	
	if len(objectName) > 1024 {
		return fmt.Errorf("object name cannot exceed 1024 characters")
	}
	
	return nil
}
