package storage

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/UDL-TF/UnitedAPI/internal/config"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

// MinIOClient wraps the MinIO client with additional functionality
type MinIOClient struct {
	client     *minio.Client
	bucketName string
}

// NewMinIOClient creates a new MinIO client instance
func NewMinIOClient(cfg *config.MinIOConfig) (*MinIOClient, error) {
	// Initialize MinIO client object.
	minioClient, err := minio.New(cfg.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.AccessKeyID, cfg.SecretAccessKey, ""),
		Secure: cfg.UseSSL,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to initialize MinIO client: %v", err)
	}

	client := &MinIOClient{
		client:     minioClient,
		bucketName: cfg.BucketName,
	}

	// Ensure bucket exists
	if err := client.ensureBucketExists(); err != nil {
		return nil, fmt.Errorf("failed to ensure bucket exists: %v", err)
	}

	return client, nil
}

// ensureBucketExists checks if bucket exists and creates it if it doesn't
func (mc *MinIOClient) ensureBucketExists() error {
	ctx := context.Background()

	exists, err := mc.client.BucketExists(ctx, mc.bucketName)
	if err != nil {
		return fmt.Errorf("error checking if bucket exists: %v", err)
	}

	if !exists {
		err = mc.client.MakeBucket(ctx, mc.bucketName, minio.MakeBucketOptions{})
		if err != nil {
			return fmt.Errorf("error creating bucket: %v", err)
		}
		fmt.Printf("Successfully created bucket: %s\n", mc.bucketName)
	}

	return nil
}

// UploadFile uploads a file to MinIO
func (mc *MinIOClient) UploadFile(objectName string, reader io.Reader, objectSize int64, contentType string) error {
	ctx := context.Background()

	_, err := mc.client.PutObject(ctx, mc.bucketName, objectName, reader, objectSize, minio.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		return fmt.Errorf("failed to upload file: %v", err)
	}

	return nil
}

// DownloadFile downloads a file from MinIO
func (mc *MinIOClient) DownloadFile(objectName string) (*minio.Object, error) {
	ctx := context.Background()

	object, err := mc.client.GetObject(ctx, mc.bucketName, objectName, minio.GetObjectOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to download file: %v", err)
	}

	return object, nil
}

// DeleteFile deletes a file from MinIO
func (mc *MinIOClient) DeleteFile(objectName string) error {
	ctx := context.Background()

	err := mc.client.RemoveObject(ctx, mc.bucketName, objectName, minio.RemoveObjectOptions{})
	if err != nil {
		return fmt.Errorf("failed to delete file: %v", err)
	}

	return nil
}

// ListFiles lists all files in the bucket
func (mc *MinIOClient) ListFiles(prefix string) ([]minio.ObjectInfo, error) {
	ctx := context.Background()
	var objects []minio.ObjectInfo

	objectCh := mc.client.ListObjects(ctx, mc.bucketName, minio.ListObjectsOptions{
		Prefix:    prefix,
		Recursive: true,
	})

	for object := range objectCh {
		if object.Err != nil {
			return nil, fmt.Errorf("error listing files: %v", object.Err)
		}
		objects = append(objects, object)
	}

	return objects, nil
}

// GetPresignedURL generates a presigned URL for direct access
func (mc *MinIOClient) GetPresignedURL(objectName string, method string, expiry time.Duration) (string, error) {
	ctx := context.Background()

	switch method {
	case "GET":
		url, err := mc.client.PresignedGetObject(ctx, mc.bucketName, objectName, expiry, nil)
		if err != nil {
			return "", fmt.Errorf("failed to generate presigned URL: %v", err)
		}
		return url.String(), nil
	case "PUT":
		url, err := mc.client.PresignedPutObject(ctx, mc.bucketName, objectName, expiry)
		if err != nil {
			return "", fmt.Errorf("failed to generate presigned URL: %v", err)
		}
		return url.String(), nil
	default:
		return "", fmt.Errorf("unsupported method: %s", method)
	}
}

// GetFileInfo gets metadata information about a file
func (mc *MinIOClient) GetFileInfo(objectName string) (minio.ObjectInfo, error) {
	ctx := context.Background()

	objInfo, err := mc.client.StatObject(ctx, mc.bucketName, objectName, minio.StatObjectOptions{})
	if err != nil {
		return minio.ObjectInfo{}, fmt.Errorf("failed to get file info: %v", err)
	}

	return objInfo, nil
}
