package storage

import (
	"bytes"
	"context"
	"crypto/rand"
	"fmt"
	"io"
	"mime/multipart"
	"net/url"
	"path/filepath"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	appConfig "github.com/kansaok/go-boilerplate/internal/config"
)

// maxS3UploadBytes batas ukuran default file yang diunggah ke S3 (20MB).
const maxS3UploadBytes = 20 << 20

var s3Client *s3.Client

func InitS3Client() error {
	cfg, err := awsConfig.LoadDefaultConfig(context.TODO())
	if err != nil {
		return fmt.Errorf("unable to load SDK config: %v", err)
	}
	s3Client = s3.NewFromConfig(cfg)
	return nil
}

// SaveToS3 mengunggah file ke S3 dengan batas ukuran, object key acak yang
// aman (mencegah injeksi path pada key/URL), dan enkripsi server-side.
// maxSize <= 0 menggunakan batas default (20MB).
func SaveToS3(file multipart.File, filename string, maxSize int64) (string, error) {
	cfg := appConfig.LoadConfig()

	if s3Client == nil {
		err := InitS3Client()
		if err != nil {
			return "", fmt.Errorf("failed to initialize S3 client: %v", err)
		}
	}

	if maxSize <= 0 {
		maxSize = maxS3UploadBytes
	}

	limited := io.LimitReader(file, maxSize+1)
	fileBuffer := &bytes.Buffer{}
	written, err := fileBuffer.ReadFrom(limited)
	if err != nil {
		return "", fmt.Errorf("failed to read file: %v", err)
	}
	if written > maxSize {
		return "", fmt.Errorf("file exceeds the maximum allowed size of %d bytes", maxSize)
	}

	contentType := detectContentType(fileBuffer.Bytes())

	key := generateSafeObjectKey(filename)

	_, err = s3Client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket:               aws.String(cfg.DatabaseConfig.BucketName),
		Key:                  aws.String(key),
		Body:                 bytes.NewReader(fileBuffer.Bytes()),
		ServerSideEncryption: types.ServerSideEncryptionAes256,
		ContentType:          aws.String(contentType),
	})

	if err != nil {
		return "", fmt.Errorf("failed to upload file to S3: %v", err)
	}

	u := url.URL{
		Scheme: "https",
		Host:   fmt.Sprintf("%s.s3.amazonaws.com", cfg.DatabaseConfig.BucketName),
		Path:   key,
	}
	return u.String(), nil
}

// generateSafeObjectKey membuat object key acak berpola `hex`+`ext`
// sehingga nama file asli tidak bocor ke URL dan key tidak bisa
// disisipkan sebagai bagian dari path lain.
func generateSafeObjectKey(filename string) string {
	ext := strings.ToLower(filepath.Ext(filename))
	b := make([]byte, 16)
	rand.Read(b)
	return fmt.Sprintf("%x%s", b, ext)
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