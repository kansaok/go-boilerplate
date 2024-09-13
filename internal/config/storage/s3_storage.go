package storage

import (
	"bytes"
	"context"
	"fmt"
	"mime/multipart"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	appConfig "github.com/kansaok/go-boilerplate/internal/config"
)

var s3Client *s3.Client

func InitS3Client() error {
	// Load AWS SDK configuration
	cfg, err := awsConfig.LoadDefaultConfig(context.TODO())
	if err != nil {
		return fmt.Errorf("unable to load SDK config: %v", err)
	}
	s3Client = s3.NewFromConfig(cfg)
	return nil
}

// SaveToS3 save the file to AWS S3
func SaveToS3(file multipart.File, filename string) (string, error) {
	// Load configuration from config package
	cfg := appConfig.LoadConfig()

	if s3Client == nil {
		err := InitS3Client()
		if err != nil {
			return "", fmt.Errorf("failed to initialize S3 client: %v", err)
		}
	}

	// Convert file to byte buffer for S3 upload
	fileBuffer := &bytes.Buffer{}
	_, err := fileBuffer.ReadFrom(file)
	if err != nil {
		return "", fmt.Errorf("failed to read file: %v", err)
	}

	// Upload file to S3
	_, err = s3Client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(cfg.DatabaseConfig.BucketName),
		Key:    aws.String(filename),
		Body:   bytes.NewReader(fileBuffer.Bytes()),
	})

	if err != nil {
		return "", fmt.Errorf("failed to upload file to S3: %v", err)
	}

	return fmt.Sprintf("https://%s.s3.amazonaws.com/%s", cfg.DatabaseConfig.BucketName, filename), nil
}
