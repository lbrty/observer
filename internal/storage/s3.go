package storage

import (
	"context"
	"fmt"
	"io"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"

	"github.com/lbrty/observer/internal/config"
)

// S3Storage stores files in an S3-compatible bucket.
type S3Storage struct {
	client *s3.Client
	bucket string
}

// NewS3Storage creates an S3Storage. Returns an error if Bucket is empty.
// When AccessKey and SecretKey are both empty, the SDK default credential
// provider chain is used (env vars → shared credentials file → IAM role).
func NewS3Storage(cfg config.StorageConfig) (*S3Storage, error) {
	if cfg.S3Bucket == "" {
		return nil, fmt.Errorf("S3_BUCKET is required")
	}

	opts := []func(*awsconfig.LoadOptions) error{
		awsconfig.WithRegion(cfg.S3Region),
	}
	if cfg.S3AccessKey != "" && cfg.S3SecretKey != "" {
		opts = append(opts, awsconfig.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(cfg.S3AccessKey, cfg.S3SecretKey, ""),
		))
	}

	awsCfg, err := awsconfig.LoadDefaultConfig(context.Background(), opts...)
	if err != nil {
		return nil, fmt.Errorf("load AWS config: %w", err)
	}

	var clientOpts []func(*s3.Options)
	if cfg.S3Endpoint != "" {
		clientOpts = append(clientOpts, func(o *s3.Options) {
			o.BaseEndpoint = aws.String(cfg.S3Endpoint)
			o.UsePathStyle = true // required for MinIO and most S3-compatible endpoints
		})
	}

	return &S3Storage{
		client: s3.NewFromConfig(awsCfg, clientOpts...),
		bucket: cfg.S3Bucket,
	}, nil
}

func (s *S3Storage) Save(ctx context.Context, path string, r io.Reader) error {
	_, err := s.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(path),
		Body:   r,
	})
	if err != nil {
		return fmt.Errorf("s3 put object: %w", err)
	}
	return nil
}

func (s *S3Storage) Open(ctx context.Context, path string) (io.ReadCloser, error) {
	out, err := s.client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(path),
	})
	if err != nil {
		return nil, fmt.Errorf("s3 get object: %w", err)
	}
	return out.Body, nil
}

func (s *S3Storage) Delete(ctx context.Context, path string) error {
	_, err := s.client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(path),
	})
	if err != nil {
		return fmt.Errorf("s3 delete object: %w", err)
	}
	return nil
}
