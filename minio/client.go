package minio

import (
	"context"
	"fmt"
	"net/url"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

// Client wraps MinIO client with additional functionality
type Client struct {
	client *minio.Client
	config Config
}

// Config holds MinIO client configuration
type Config struct {
	Endpoint        string `json:"endpoint" yaml:"endpoint" env:"MINIO_ENDPOINT"`
	AccessKeyID     string `json:"access_key_id" yaml:"access_key_id" env:"MINIO_ACCESS_KEY_ID"`
	SecretAccessKey string `json:"secret_access_key" yaml:"secret_access_key" env:"MINIO_SECRET_ACCESS_KEY"`
	UseSSL          bool   `json:"use_ssl" yaml:"use_ssl" env:"MINIO_USE_SSL" default:"false"`
	Region          string `json:"region" yaml:"region" env:"MINIO_REGION" default:"us-east-1"`
}

// PresignedURLOptions holds options for presigned URL generation
type PresignedURLOptions struct {
	Expiry      time.Duration     `json:"expiry"`
	ReqParams   url.Values        `json:"req_params"`
	ExtraHeaders map[string]string `json:"extra_headers"`
}

// DeleteObjectOptions holds options for object deletion
type DeleteObjectOptions struct {
	VersionID string `json:"version_id"`
}

// ObjectInfo holds object information
type ObjectInfo struct {
	Key          string            `json:"key"`
	LastModified time.Time         `json:"last_modified"`
	Size         int64             `json:"size"`
	ContentType  string            `json:"content_type"`
	Metadata     map[string]string `json:"metadata"`
	VersionID    string            `json:"version_id"`
	ETag         string            `json:"etag"`
}

// New creates a new MinIO client with default configuration
func New(config Config) (*Client, error) {
	return NewWithConfig(config)
}

// NewWithConfig creates a new MinIO client with custom configuration
func NewWithConfig(config Config) (*Client, error) {
	if config.Endpoint == "" {
		return nil, fmt.Errorf("endpoint is required")
	}
	if config.AccessKeyID == "" {
		return nil, fmt.Errorf("access key ID is required")
	}
	if config.SecretAccessKey == "" {
		return nil, fmt.Errorf("secret access key is required")
	}

	// Initialize MinIO client
	minioClient, err := minio.New(config.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(config.AccessKeyID, config.SecretAccessKey, ""),
		Secure: config.UseSSL,
		Region: config.Region,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create MinIO client: %w", err)
	}

	return &Client{
		client: minioClient,
		config: config,
	}, nil
}

// GetPresignedURL generates a presigned URL for GET operations
func (c *Client) GetPresignedURL(ctx context.Context, bucketName, objectName string, opts *PresignedURLOptions) (*url.URL, error) {
	if opts == nil {
		opts = &PresignedURLOptions{
			Expiry: 24 * time.Hour, // Default 24 hours
		}
	}

	if opts.Expiry == 0 {
		opts.Expiry = 24 * time.Hour
	}

	presignedURL, err := c.client.PresignedGetObject(ctx, bucketName, objectName, opts.Expiry, opts.ReqParams)
	if err != nil {
		return nil, fmt.Errorf("failed to generate presigned GET URL: %w", err)
	}

	return presignedURL, nil
}

// PutPresignedURL generates a presigned URL for PUT operations
func (c *Client) PutPresignedURL(ctx context.Context, bucketName, objectName string, opts *PresignedURLOptions) (*url.URL, error) {
	if opts == nil {
		opts = &PresignedURLOptions{
			Expiry: 24 * time.Hour, // Default 24 hours
		}
	}

	if opts.Expiry == 0 {
		opts.Expiry = 24 * time.Hour
	}

	presignedURL, err := c.client.PresignedPutObject(ctx, bucketName, objectName, opts.Expiry)
	if err != nil {
		return nil, fmt.Errorf("failed to generate presigned PUT URL: %w", err)
	}

	return presignedURL, nil
}

// PostPresignedURL generates a presigned URL for POST operations (form upload)
func (c *Client) PostPresignedURL(ctx context.Context, bucketName, objectName string, opts *PresignedURLOptions) (*url.URL, map[string]string, error) {
	if opts == nil {
		opts = &PresignedURLOptions{
			Expiry: 24 * time.Hour, // Default 24 hours
		}
	}

	if opts.Expiry == 0 {
		opts.Expiry = 24 * time.Hour
	}

	policy := minio.NewPostPolicy()
	policy.SetBucket(bucketName)
	policy.SetKey(objectName)
	policy.SetExpires(time.Now().UTC().Add(opts.Expiry))

	// Add extra headers as conditions if provided
	if opts.ExtraHeaders != nil {
		for key, value := range opts.ExtraHeaders {
			policy.SetContentType(value)
			if key == "Content-Type" {
				break
			}
		}
	}

	presignedURL, formData, err := c.client.PresignedPostPolicy(ctx, policy)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to generate presigned POST URL: %w", err)
	}

	return presignedURL, formData, nil
}

// DeleteObject deletes an object from the bucket
func (c *Client) DeleteObject(ctx context.Context, bucketName, objectName string, opts *DeleteObjectOptions) error {
	var minioOpts minio.RemoveObjectOptions
	
	if opts != nil && opts.VersionID != "" {
		minioOpts.VersionID = opts.VersionID
	}

	err := c.client.RemoveObject(ctx, bucketName, objectName, minioOpts)
	if err != nil {
		return fmt.Errorf("failed to delete object: %w", err)
	}

	return nil
}

// DeleteObjects deletes multiple objects from the bucket
func (c *Client) DeleteObjects(ctx context.Context, bucketName string, objectNames []string) error {
	objectsCh := make(chan minio.ObjectInfo)

	// Send object names for deletion
	go func() {
		defer close(objectsCh)
		for _, objectName := range objectNames {
			objectsCh <- minio.ObjectInfo{
				Key: objectName,
			}
		}
	}()

	// Delete objects
	opts := minio.RemoveObjectsOptions{
		GovernanceBypass: true,
	}

	for rErr := range c.client.RemoveObjects(ctx, bucketName, objectsCh, opts) {
		if rErr.Err != nil {
			return fmt.Errorf("failed to delete object %s: %w", rErr.ObjectName, rErr.Err)
		}
	}

	return nil
}

// GetObjectInfo retrieves object information
func (c *Client) GetObjectInfo(ctx context.Context, bucketName, objectName string) (*ObjectInfo, error) {
	objInfo, err := c.client.StatObject(ctx, bucketName, objectName, minio.StatObjectOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get object info: %w", err)
	}

	return &ObjectInfo{
		Key:          objInfo.Key,
		LastModified: objInfo.LastModified,
		Size:         objInfo.Size,
		ContentType:  objInfo.ContentType,
		Metadata:     objInfo.UserMetadata,
		VersionID:    objInfo.VersionID,
		ETag:         objInfo.ETag,
	}, nil
}

// ListObjects lists objects in a bucket with optional prefix
func (c *Client) ListObjects(ctx context.Context, bucketName, prefix string, recursive bool) ([]ObjectInfo, error) {
	opts := minio.ListObjectsOptions{
		Prefix:    prefix,
		Recursive: recursive,
	}

	var objects []ObjectInfo
	for object := range c.client.ListObjects(ctx, bucketName, opts) {
		if object.Err != nil {
			return nil, fmt.Errorf("failed to list objects: %w", object.Err)
		}

		objects = append(objects, ObjectInfo{
			Key:          object.Key,
			LastModified: object.LastModified,
			Size:         object.Size,
			ContentType:  object.ContentType,
			ETag:         object.ETag,
		})
	}

	return objects, nil
}

// ObjectExists checks if an object exists in the bucket
func (c *Client) ObjectExists(ctx context.Context, bucketName, objectName string) (bool, error) {
	_, err := c.client.StatObject(ctx, bucketName, objectName, minio.StatObjectOptions{})
	if err != nil {
		// Check if error is "object not found"
		if minio.ToErrorResponse(err).Code == "NoSuchKey" {
			return false, nil
		}
		return false, fmt.Errorf("failed to check object existence: %w", err)
	}
	return true, nil
}

// BucketExists checks if a bucket exists
func (c *Client) BucketExists(ctx context.Context, bucketName string) (bool, error) {
	exists, err := c.client.BucketExists(ctx, bucketName)
	if err != nil {
		return false, fmt.Errorf("failed to check bucket existence: %w", err)
	}
	return exists, nil
}

// CreateBucket creates a new bucket
func (c *Client) CreateBucket(ctx context.Context, bucketName string, region string) error {
	if region == "" {
		region = c.config.Region
	}

	err := c.client.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{
		Region: region,
	})
	if err != nil {
		return fmt.Errorf("failed to create bucket: %w", err)
	}

	return nil
}

// DeleteBucket deletes an empty bucket
func (c *Client) DeleteBucket(ctx context.Context, bucketName string) error {
	err := c.client.RemoveBucket(ctx, bucketName)
	if err != nil {
		return fmt.Errorf("failed to delete bucket: %w", err)
	}

	return nil
}
