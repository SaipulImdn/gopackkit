# MinIO Module

The MinIO module provides a comprehensive client for MinIO object storage operations with support for presigned URLs, object management, and secure file operations.

## Features

- **Presigned URL Generation**: GET, PUT, and POST presigned URLs
- **Object Operations**: Upload, download, delete, and existence checking
- **Bucket Management**: Create, list, and manage buckets
- **Secure Configuration**: TLS support and credential management
- **Error Handling**: Comprehensive error handling and validation
- **Flexible Expiry**: Configurable URL expiration times

## Installation

```bash
go get github.com/saipulimdn/gopackkit/minio
```

## Quick Start

### Basic Configuration

```go
package main

import (
    "log"
    "time"
    
    "github.com/saipulimdn/gopackkit/minio"
)

func main() {
    // Basic configuration
    config := minio.Config{
        Endpoint:        "localhost:9000",
        AccessKeyID:     "minioadmin",
        SecretAccessKey: "minioadmin",
        UseSSL:          false,
        Region:          "us-east-1",
    }
    
    client, err := minio.New(config)
    if err != nil {
        log.Fatal("Failed to create MinIO client:", err)
    }
    
    // Generate presigned URL for file upload
    uploadURL, err := client.GeneratePresignedPutURL(
        "my-bucket", 
        "uploads/document.pdf", 
        1*time.Hour,
    )
    if err != nil {
        log.Fatal("Failed to generate upload URL:", err)
    }
    
    log.Printf("Upload URL: %s", uploadURL)
}
```

### Production Configuration

```go
package main

import (
    "github.com/saipulimdn/gopackkit/minio"
)

func main() {
    // Production configuration with TLS
    config := minio.Config{
        Endpoint:        "s3.amazonaws.com",
        AccessKeyID:     "YOUR_ACCESS_KEY",
        SecretAccessKey: "YOUR_SECRET_KEY",
        UseSSL:          true,
        Region:          "us-west-2",
    }
    
    client, err := minio.New(config)
    if err != nil {
        log.Fatal("Failed to create MinIO client:", err)
    }
    
    // Your application logic here
}
```

## Configuration

```go
type Config struct {
    Endpoint        string // MinIO server endpoint
    AccessKeyID     string // Access key ID
    SecretAccessKey string // Secret access key
    UseSSL          bool   // Use HTTPS
    Region          string // AWS region
}
```

### Environment Variables

Configure MinIO using environment variables:

```bash
export MINIO_ENDPOINT=localhost:9000
export MINIO_ACCESS_KEY_ID=minioadmin
export MINIO_SECRET_ACCESS_KEY=minioadmin
export MINIO_USE_SSL=false
export MINIO_REGION=us-east-1
```

## Presigned URLs

### Generate Upload URL (PUT)

```go
package main

import (
    "time"
    "github.com/saipulimdn/gopackkit/minio"
)

func generateUploadURL(client *minio.Client, bucket, object string) (string, error) {
    // Generate presigned PUT URL for file upload
    // URL valid for 1 hour
    uploadURL, err := client.GeneratePresignedPutURL(bucket, object, 1*time.Hour)
    if err != nil {
        return "", err
    }
    
    return uploadURL, nil
}

func main() {
    client, err := minio.New(minio.Config{
        Endpoint:        "localhost:9000",
        AccessKeyID:     "minioadmin",
        SecretAccessKey: "minioadmin",
        UseSSL:          false,
        Region:          "us-east-1",
    })
    if err != nil {
        log.Fatal(err)
    }
    
    // Generate upload URL
    uploadURL, err := generateUploadURL(client, "documents", "uploads/file.pdf")
    if err != nil {
        log.Fatal("Failed to generate upload URL:", err)
    }
    
    log.Printf("Upload your file to: %s", uploadURL)
    
    // Client can now use this URL to upload files directly to MinIO
    // Example: curl -X PUT -T "local-file.pdf" "$uploadURL"
}
```

### Generate Download URL (GET)

```go
package main

import (
    "time"
    "github.com/saipulimdn/gopackkit/minio"
)

func generateDownloadURL(client *minio.Client, bucket, object string) (string, error) {
    // Generate presigned GET URL for file download
    // URL valid for 24 hours
    downloadURL, err := client.GeneratePresignedGetURL(bucket, object, 24*time.Hour)
    if err != nil {
        return "", err
    }
    
    return downloadURL, nil
}

func main() {
    client, err := minio.New(minio.Config{
        Endpoint:        "localhost:9000",
        AccessKeyID:     "minioadmin",
        SecretAccessKey: "minioadmin",
        UseSSL:          false,
    })
    if err != nil {
        log.Fatal(err)
    }
    
    // Check if file exists before generating download URL
    exists, err := client.ObjectExists("documents", "uploads/file.pdf")
    if err != nil {
        log.Fatal("Failed to check object existence:", err)
    }
    
    if !exists {
        log.Fatal("File does not exist")
    }
    
    // Generate download URL
    downloadURL, err := generateDownloadURL(client, "documents", "uploads/file.pdf")
    if err != nil {
        log.Fatal("Failed to generate download URL:", err)
    }
    
    log.Printf("Download your file from: %s", downloadURL)
}
```

### Generate POST Form URL

```go
package main

import (
    "encoding/json"
    "time"
    "github.com/saipulimdn/gopackkit/minio"
)

func generatePostFormURL(client *minio.Client, bucket, object string) (map[string]string, error) {
    // Generate presigned POST form data
    // URL and form fields valid for 1 hour
    formData, err := client.GeneratePresignedPostURL(bucket, object, 1*time.Hour)
    if err != nil {
        return nil, err
    }
    
    return formData, nil
}

func main() {
    client, err := minio.New(minio.Config{
        Endpoint:        "localhost:9000",
        AccessKeyID:     "minioadmin",
        SecretAccessKey: "minioadmin",
        UseSSL:          false,
    })
    if err != nil {
        log.Fatal(err)
    }
    
    // Generate POST form data
    formData, err := generatePostFormURL(client, "uploads", "images/photo.jpg")
    if err != nil {
        log.Fatal("Failed to generate POST form:", err)
    }
    
    // Print form data for frontend usage
    jsonData, _ := json.MarshalIndent(formData, "", "  ")
    log.Printf("POST form data:\n%s", jsonData)
    
    // Frontend can use this data to create a multipart form
    // for direct browser uploads to MinIO
}
```

## Object Operations

### Upload Object

```go
package main

import (
    "github.com/saipulimdn/gopackkit/minio"
)

func uploadFile(client *minio.Client, bucket, object, filePath string) error {
    // Upload file from local filesystem to MinIO
    err := client.UploadObject(bucket, object, filePath)
    if err != nil {
        return err
    }
    
    log.Printf("Successfully uploaded %s to %s/%s", filePath, bucket, object)
    return nil
}

func main() {
    client, err := minio.New(minio.Config{
        Endpoint:        "localhost:9000",
        AccessKeyID:     "minioadmin",
        SecretAccessKey: "minioadmin",
        UseSSL:          false,
    })
    if err != nil {
        log.Fatal(err)
    }
    
    // Upload a file
    err = uploadFile(client, "documents", "reports/monthly-report.pdf", "./local-report.pdf")
    if err != nil {
        log.Fatal("Upload failed:", err)
    }
}
```

### Download Object

```go
package main

import (
    "github.com/saipulimdn/gopackkit/minio"
)

func downloadFile(client *minio.Client, bucket, object, filePath string) error {
    // Download object from MinIO to local filesystem
    err := client.DownloadObject(bucket, object, filePath)
    if err != nil {
        return err
    }
    
    log.Printf("Successfully downloaded %s/%s to %s", bucket, object, filePath)
    return nil
}

func main() {
    client, err := minio.New(minio.Config{
        Endpoint:        "localhost:9000",
        AccessKeyID:     "minioadmin",
        SecretAccessKey: "minioadmin",
        UseSSL:          false,
    })
    if err != nil {
        log.Fatal(err)
    }
    
    // Download a file
    err = downloadFile(client, "documents", "reports/monthly-report.pdf", "./downloaded-report.pdf")
    if err != nil {
        log.Fatal("Download failed:", err)
    }
}
```

### Delete Object

```go
package main

import (
    "github.com/saipulimdn/gopackkit/minio"
)

func deleteFile(client *minio.Client, bucket, object string) error {
    // Check if object exists before deleting
    exists, err := client.ObjectExists(bucket, object)
    if err != nil {
        return err
    }
    
    if !exists {
        log.Printf("Object %s/%s does not exist", bucket, object)
        return nil
    }
    
    // Delete the object
    err = client.DeleteObject(bucket, object)
    if err != nil {
        return err
    }
    
    log.Printf("Successfully deleted %s/%s", bucket, object)
    return nil
}

func main() {
    client, err := minio.New(minio.Config{
        Endpoint:        "localhost:9000",
        AccessKeyID:     "minioadmin",
        SecretAccessKey: "minioadmin",
        UseSSL:          false,
    })
    if err != nil {
        log.Fatal(err)
    }
    
    // Delete a file
    err = deleteFile(client, "temp-uploads", "old-file.txt")
    if err != nil {
        log.Fatal("Delete failed:", err)
    }
}
```

### Check Object Existence

```go
package main

import (
    "github.com/saipulimdn/gopackkit/minio"
)

func checkFileExists(client *minio.Client, bucket, object string) {
    exists, err := client.ObjectExists(bucket, object)
    if err != nil {
        log.Printf("Error checking object existence: %v", err)
        return
    }
    
    if exists {
        log.Printf("Object %s/%s exists", bucket, object)
    } else {
        log.Printf("Object %s/%s does not exist", bucket, object)
    }
}

func main() {
    client, err := minio.New(minio.Config{
        Endpoint:        "localhost:9000",
        AccessKeyID:     "minioadmin",
        SecretAccessKey: "minioadmin",
        UseSSL:          false,
    })
    if err != nil {
        log.Fatal(err)
    }
    
    // Check if various files exist
    files := []struct {
        bucket string
        object string
    }{
        {"documents", "reports/monthly-report.pdf"},
        {"images", "photos/vacation.jpg"},
        {"backups", "database-backup.sql"},
    }
    
    for _, file := range files {
        checkFileExists(client, file.bucket, file.object)
    }
}
```

## Advanced Usage

### File Upload Service

```go
package main

import (
    "encoding/json"
    "fmt"
    "net/http"
    "path/filepath"
    "strings"
    "time"
    
    "github.com/saipulimdn/gopackkit/minio"
)

type FileUploadService struct {
    client *minio.Client
    bucket string
}

type UploadResponse struct {
    UploadURL string            `json:"upload_url"`
    FormData  map[string]string `json:"form_data,omitempty"`
    ObjectKey string            `json:"object_key"`
    ExpiresAt time.Time         `json:"expires_at"`
}

func NewFileUploadService(config minio.Config, bucket string) (*FileUploadService, error) {
    client, err := minio.New(config)
    if err != nil {
        return nil, err
    }
    
    return &FileUploadService{
        client: client,
        bucket: bucket,
    }, nil
}

func (s *FileUploadService) GenerateUploadURL(filename string, uploadType string) (*UploadResponse, error) {
    // Generate unique object key
    objectKey := fmt.Sprintf("uploads/%d/%s", time.Now().Unix(), filename)
    
    // Validate file extension
    ext := strings.ToLower(filepath.Ext(filename))
    allowedExtensions := []string{".jpg", ".jpeg", ".png", ".pdf", ".doc", ".docx", ".txt"}
    
    isAllowed := false
    for _, allowed := range allowedExtensions {
        if ext == allowed {
            isAllowed = true
            break
        }
    }
    
    if !isAllowed {
        return nil, fmt.Errorf("file type %s not allowed", ext)
    }
    
    expiresAt := time.Now().Add(1 * time.Hour)
    
    switch uploadType {
    case "direct":
        // Generate PUT URL for direct upload
        uploadURL, err := s.client.GeneratePresignedPutURL(s.bucket, objectKey, 1*time.Hour)
        if err != nil {
            return nil, err
        }
        
        return &UploadResponse{
            UploadURL: uploadURL,
            ObjectKey: objectKey,
            ExpiresAt: expiresAt,
        }, nil
        
    case "form":
        // Generate POST form data for browser upload
        formData, err := s.client.GeneratePresignedPostURL(s.bucket, objectKey, 1*time.Hour)
        if err != nil {
            return nil, err
        }
        
        return &UploadResponse{
            FormData:  formData,
            ObjectKey: objectKey,
            ExpiresAt: expiresAt,
        }, nil
        
    default:
        return nil, fmt.Errorf("invalid upload type: %s", uploadType)
    }
}

func (s *FileUploadService) GenerateDownloadURL(objectKey string) (string, error) {
    // Check if file exists
    exists, err := s.client.ObjectExists(s.bucket, objectKey)
    if err != nil {
        return "", err
    }
    
    if !exists {
        return "", fmt.Errorf("file not found: %s", objectKey)
    }
    
    // Generate download URL valid for 24 hours
    downloadURL, err := s.client.GeneratePresignedGetURL(s.bucket, objectKey, 24*time.Hour)
    if err != nil {
        return "", err
    }
    
    return downloadURL, nil
}

func (s *FileUploadService) DeleteFile(objectKey string) error {
    return s.client.DeleteObject(s.bucket, objectKey)
}

// HTTP Handlers
func (s *FileUploadService) uploadHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }
    
    filename := r.URL.Query().Get("filename")
    uploadType := r.URL.Query().Get("type")
    
    if filename == "" {
        http.Error(w, "filename parameter required", http.StatusBadRequest)
        return
    }
    
    if uploadType == "" {
        uploadType = "direct"
    }
    
    response, err := s.GenerateUploadURL(filename, uploadType)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(response)
}

func (s *FileUploadService) downloadHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodGet {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }
    
    objectKey := r.URL.Query().Get("key")
    if objectKey == "" {
        http.Error(w, "key parameter required", http.StatusBadRequest)
        return
    }
    
    downloadURL, err := s.GenerateDownloadURL(objectKey)
    if err != nil {
        http.Error(w, err.Error(), http.StatusNotFound)
        return
    }
    
    // Redirect to presigned URL
    http.Redirect(w, r, downloadURL, http.StatusTemporaryRedirect)
}

func main() {
    config := minio.Config{
        Endpoint:        "localhost:9000",
        AccessKeyID:     "minioadmin",
        SecretAccessKey: "minioadmin",
        UseSSL:          false,
        Region:          "us-east-1",
    }
    
    service, err := NewFileUploadService(config, "file-uploads")
    if err != nil {
        log.Fatal("Failed to create upload service:", err)
    }
    
    // Setup HTTP handlers
    http.HandleFunc("/upload", service.uploadHandler)
    http.HandleFunc("/download", service.downloadHandler)
    
    log.Println("File upload service starting on :8080")
    log.Fatal(http.ListenAndServe(":8080", nil))
}
```

### Image Processing Pipeline

```go
package main

import (
    "fmt"
    "image"
    "image/jpeg"
    "image/png"
    "os"
    "path/filepath"
    "strings"
    
    "github.com/saipulimdn/gopackkit/minio"
    "github.com/nfnt/resize" // Image resizing library
)

type ImageProcessor struct {
    client *minio.Client
    bucket string
}

func NewImageProcessor(config minio.Config, bucket string) (*ImageProcessor, error) {
    client, err := minio.New(config)
    if err != nil {
        return nil, err
    }
    
    return &ImageProcessor{
        client: client,
        bucket: bucket,
    }, nil
}

func (p *ImageProcessor) ProcessImage(originalKey string) error {
    // Download original image
    tempFile := filepath.Join("/tmp", "original_"+filepath.Base(originalKey))
    err := p.client.DownloadObject(p.bucket, originalKey, tempFile)
    if err != nil {
        return fmt.Errorf("failed to download original image: %w", err)
    }
    defer os.Remove(tempFile)
    
    // Open and decode image
    file, err := os.Open(tempFile)
    if err != nil {
        return err
    }
    defer file.Close()
    
    img, format, err := image.Decode(file)
    if err != nil {
        return fmt.Errorf("failed to decode image: %w", err)
    }
    
    // Generate thumbnails in different sizes
    sizes := []struct {
        name   string
        width  uint
        height uint
    }{
        {"thumbnail", 150, 150},
        {"small", 300, 300},
        {"medium", 600, 600},
        {"large", 1200, 1200},
    }
    
    for _, size := range sizes {
        // Resize image
        resized := resize.Resize(size.width, size.height, img, resize.Lanczos3)
        
        // Create thumbnail file
        thumbnailFile := filepath.Join("/tmp", fmt.Sprintf("%s_%s", size.name, filepath.Base(originalKey)))
        thumb, err := os.Create(thumbnailFile)
        if err != nil {
            continue
        }
        
        // Encode and save thumbnail
        switch format {
        case "jpeg":
            jpeg.Encode(thumb, resized, &jpeg.Options{Quality: 85})
        case "png":
            png.Encode(thumb, resized)
        }
        thumb.Close()
        
        // Upload thumbnail to MinIO
        thumbnailKey := p.generateThumbnailKey(originalKey, size.name)
        err = p.client.UploadObject(p.bucket, thumbnailKey, thumbnailFile)
        if err != nil {
            log.Printf("Failed to upload thumbnail %s: %v", size.name, err)
        } else {
            log.Printf("Successfully created %s thumbnail: %s", size.name, thumbnailKey)
        }
        
        // Clean up thumbnail file
        os.Remove(thumbnailFile)
    }
    
    return nil
}

func (p *ImageProcessor) generateThumbnailKey(originalKey, size string) string {
    ext := filepath.Ext(originalKey)
    nameWithoutExt := strings.TrimSuffix(originalKey, ext)
    return fmt.Sprintf("%s_%s%s", nameWithoutExt, size, ext)
}

func (p *ImageProcessor) GetImageURLs(originalKey string) (map[string]string, error) {
    urls := make(map[string]string)
    
    // Original image URL
    originalURL, err := p.client.GeneratePresignedGetURL(p.bucket, originalKey, 24*time.Hour)
    if err != nil {
        return nil, err
    }
    urls["original"] = originalURL
    
    // Thumbnail URLs
    sizes := []string{"thumbnail", "small", "medium", "large"}
    for _, size := range sizes {
        thumbnailKey := p.generateThumbnailKey(originalKey, size)
        
        // Check if thumbnail exists
        exists, err := p.client.ObjectExists(p.bucket, thumbnailKey)
        if err != nil || !exists {
            continue
        }
        
        // Generate URL
        url, err := p.client.GeneratePresignedGetURL(p.bucket, thumbnailKey, 24*time.Hour)
        if err != nil {
            continue
        }
        
        urls[size] = url
    }
    
    return urls, nil
}
```

## Best Practices

### Security

1. **Use HTTPS in Production**: Always enable SSL for production environments
2. **Secure Credentials**: Store access keys securely, preferably in environment variables
3. **Limited Permissions**: Use IAM policies to limit MinIO access permissions
4. **Validate File Types**: Always validate uploaded file types and sizes
5. **Short-lived URLs**: Use reasonable expiration times for presigned URLs

### Performance

1. **Connection Reuse**: Reuse MinIO client instances across requests
2. **Concurrent Operations**: Use goroutines for batch operations
3. **Appropriate Timeouts**: Set reasonable timeouts for operations
4. **Regional Proximity**: Use MinIO regions close to your application

### Error Handling

```go
package main

import (
    "errors"
    "github.com/saipulimdn/gopackkit/minio"
    "github.com/minio/minio-go/v7"
)

func handleMinIOError(err error) {
    if err == nil {
        return
    }
    
    // Check for specific MinIO errors
    var minioErr minio.ErrorResponse
    if errors.As(err, &minioErr) {
        switch minioErr.Code {
        case "NoSuchBucket":
            log.Println("Bucket does not exist")
        case "NoSuchKey":
            log.Println("Object does not exist")
        case "AccessDenied":
            log.Println("Access denied - check credentials and permissions")
        case "InvalidBucketName":
            log.Println("Invalid bucket name")
        default:
            log.Printf("MinIO error: %s - %s", minioErr.Code, minioErr.Message)
        }
    } else {
        log.Printf("General error: %v", err)
    }
}
```

## Testing

### Unit Tests

```go
package main

import (
    "testing"
    "time"
    
    "github.com/saipulimdn/gopackkit/minio"
)

func TestMinIOClient(t *testing.T) {
    // Test client creation
    config := minio.Config{
        Endpoint:        "localhost:9000",
        AccessKeyID:     "minioadmin",
        SecretAccessKey: "minioadmin",
        UseSSL:          false,
        Region:          "us-east-1",
    }
    
    client, err := minio.New(config)
    if err != nil {
        t.Fatalf("Failed to create MinIO client: %v", err)
    }
    
    // Test presigned URL generation
    url, err := client.GeneratePresignedGetURL("test-bucket", "test-object", 1*time.Hour)
    if err != nil {
        t.Errorf("Failed to generate presigned URL: %v", err)
    }
    
    if url == "" {
        t.Error("Generated URL should not be empty")
    }
}

func TestConfigValidation(t *testing.T) {
    // Test invalid configuration
    invalidConfig := minio.Config{
        Endpoint: "", // Empty endpoint
    }
    
    _, err := minio.New(invalidConfig)
    if err == nil {
        t.Error("Expected error for invalid configuration")
    }
}
```

## Troubleshooting

### Common Issues

1. **Connection Refused**: Check MinIO server is running and endpoint is correct
2. **Access Denied**: Verify access key ID and secret access key
3. **SSL Errors**: Check SSL configuration and certificates
4. **Invalid Bucket Name**: Ensure bucket names follow naming conventions
5. **Object Not Found**: Verify object key and bucket exist

### Debug Mode

```go
// Enable debug logging for MinIO operations
import "github.com/minio/minio-go/v7/pkg/set"

func enableDebug() {
    // This would depend on the specific MinIO client implementation
    // You might need to implement custom logging in your wrapper
}
```
