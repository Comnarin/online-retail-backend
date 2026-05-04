package storage

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/rs/zerolog/log"
)

// MinIOClient wraps the MinIO SDK for file storage operations.
type MinIOClient struct {
	client   *minio.Client
	bucket   string
	endpoint string
	useSSL   bool
}

// NewMinIO creates a new MinIO client, auto-creating the bucket if it doesn't exist.
func NewMinIO(endpoint, accessKey, secretKey, bucket string, useSSL bool) (*MinIOClient, error) {
	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MinIO: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	exists, err := client.BucketExists(ctx, bucket)
	if err != nil {
		return nil, fmt.Errorf("failed to check bucket existence: %w", err)
	}
	if !exists {
		if err := client.MakeBucket(ctx, bucket, minio.MakeBucketOptions{}); err != nil {
			return nil, fmt.Errorf("failed to create bucket: %w", err)
		}
		log.Info().Str("bucket", bucket).Msg("MinIO bucket created")

		// Set bucket policy to allow public read
		policy := fmt.Sprintf(`{
			"Version": "2012-10-17",
			"Statement": [{
				"Effect": "Allow",
				"Principal": {"AWS": ["*"]},
				"Action": ["s3:GetObject"],
				"Resource": ["arn:aws:s3:::%s/*"]
			}]
		}`, bucket)
		if err := client.SetBucketPolicy(ctx, bucket, policy); err != nil {
			log.Warn().Err(err).Msg("failed to set bucket policy, images may not be publicly accessible")
		}
	}

	return &MinIOClient{
		client:   client,
		bucket:   bucket,
		endpoint: endpoint,
		useSSL:   useSSL,
	}, nil
}

// UploadFile stores a file in MinIO and returns the public URL.
func (m *MinIOClient) UploadFile(ctx context.Context, objectName string, reader io.Reader, fileSize int64, contentType string) (string, error) {
	_, err := m.client.PutObject(ctx, m.bucket, objectName, reader, fileSize, minio.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload file: %w", err)
	}

	return m.GetPublicURL(objectName), nil
}

// GetPublicURL constructs a relative public URL for a stored object.
// This is served by the backend's storage proxy to ensure accessibility via tunnels/external IPs.
func (m *MinIOClient) GetPublicURL(objectName string) string {
	return "/" + objectName
}

// GetObject retrieves an object from MinIO for proxying.
func (m *MinIOClient) GetObject(ctx context.Context, objectName string) (io.ReadCloser, minio.ObjectInfo, error) {
	obj, err := m.client.GetObject(ctx, m.bucket, objectName, minio.GetObjectOptions{})
	if err != nil {
		return nil, minio.ObjectInfo{}, err
	}

	info, err := obj.Stat()
	if err != nil {
		obj.Close()
		return nil, minio.ObjectInfo{}, err
	}

	return obj, info, nil
}
