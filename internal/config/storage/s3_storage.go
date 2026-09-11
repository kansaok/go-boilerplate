package storage

import (
	"bytes"
	"context"
	"fmt"
	"mime/multipart"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	appConfig "github.com/kansaok/go-boilerplate/internal/config"
)

var s3Client *s3.Client

func InitS3Client() error {
	cfg, err := awsConfig.LoadDefaultConfig(context.TODO())
	if err != nil {
		return fmt.Errorf("unable to load SDK config: %v", err)
	}
	s3Client = s3.NewFromConfig(cfg)
	return nil
}

func SaveToS3(file multipart.File, filename string) (string, error) {
	cfg := appConfig.LoadConfig()

	if s3Client == nil {
		err := InitS3Client()
		if err != nil {
			return "", fmt.Errorf("failed to initialize S3 client: %v", err)
		}
	}

	fileBuffer := &bytes.Buffer{}
	_, err := fileBuffer.ReadFrom(file)
	if err != nil {
		return "", fmt.Errorf("failed to read file: %v", err)
	}

	contentType := detectContentType(fileBuffer.Bytes())

	_, err = s3Client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket:               aws.String(cfg.DatabaseConfig.BucketName),
		Key:                  aws.String(filename),
		Body:                 bytes.NewReader(fileBuffer.Bytes()),
		ServerSideEncryption: types.ServerSideEncryptionAes256,
		ContentType:          aws.String(contentType),
	})

	if err != nil {
		return "", fmt.Errorf("failed to upload file to S3: %v", err)
	}

	return fmt.Sprintf("https://%s.s3.amazonaws.com/%s", cfg.DatabaseConfig.BucketName, filename), nil
}

func detectContentType(data []byte) string {
	if len(data) < 4 {
		return "application/octet-stream"
	}
	if data[0] == 0xFF && data[1] == 0xD8 {
		return "image/jpeg"
	}
	if data[0] == 0x89 && data[1] == 0x50 && data[2] == 0x4E && data[3] == 0x47 {
		return "image/png"
	}
	if data[0] == 0x25 && data[1] == 0x50 && data[2] == 0x44 && data[3] == 0x46 {
		return "application/pdf"
	}
	return "application/octet-stream"
}
