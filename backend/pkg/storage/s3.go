package storage

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type S3Client struct {
	client  *s3.Client
	bucket  string
	region  string
	baseURL string
}

// NewS3Client creates an S3-compatible client
// For Cloudflare R2, set endpoint to: https://<account_id>.r2.cloudflarestorage.com
func NewS3Client(region, accessKeyID, secretAccessKey, bucket, baseURL, endpoint string) (*S3Client, error) {
	log.Printf("[S3] Initializing S3 client - Region: %s, Bucket: %s, Endpoint: %s, BaseURL: %s", region, bucket, endpoint, baseURL)

	if accessKeyID == "" || secretAccessKey == "" {
		log.Println("[S3] ERROR: AWS_ACCESS_KEY_ID or AWS_SECRET_ACCESS_KEY not set")
		return nil, fmt.Errorf("S3 credentials not configured")
	}

	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(region),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
			accessKeyID,
			secretAccessKey,
			"",
		)),
	)
	if err != nil {
		log.Printf("[S3] ERROR loading AWS config: %v", err)
		return nil, err
	}

	// Create S3 client with optional custom endpoint (for R2, MinIO, etc.)
	var client *s3.Client
	if endpoint != "" {
		log.Printf("[S3] Using custom endpoint: %s", endpoint)
		client = s3.NewFromConfig(cfg, func(o *s3.Options) {
			o.BaseEndpoint = aws.String(endpoint)
			o.UsePathStyle = true // Required for R2
		})
	} else {
		log.Println("[S3] Using default AWS S3 endpoint")
		client = s3.NewFromConfig(cfg)
	}

	// Default to S3 URL if no custom base URL provided
	if baseURL == "" {
		baseURL = fmt.Sprintf("https://%s.s3.%s.amazonaws.com", bucket, region)
	}

	log.Printf("[S3] S3 client initialized successfully - Public URL base: %s", baseURL)

	return &S3Client{
		client:  client,
		bucket:  bucket,
		region:  region,
		baseURL: baseURL,
	}, nil
}

func (s *S3Client) GetPresignedUploadURL(key, contentType string) (string, error) {
	log.Printf("[S3] Creating presigned URL - Bucket: %s, Key: %s, ContentType: %s", s.bucket, key, contentType)

	presignClient := s3.NewPresignClient(s.client)

	req, err := presignClient.PresignPutObject(context.TODO(), &s3.PutObjectInput{
		Bucket:      aws.String(s.bucket),
		Key:         aws.String(key),
		ContentType: aws.String(contentType),
	}, func(opts *s3.PresignOptions) {
		opts.Expires = 15 * time.Minute
	})
	if err != nil {
		log.Printf("[S3] ERROR creating presigned URL: %v", err)
		return "", err
	}

	log.Printf("[S3] Presigned URL created successfully")
	return req.URL, nil
}

func (s *S3Client) GetPublicURL(key string) string {
	return fmt.Sprintf("%s/%s", s.baseURL, key)
}

func (s *S3Client) DeleteObject(key string) error {
	_, err := s.client.DeleteObject(context.TODO(), &s3.DeleteObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
	})
	return err
}

