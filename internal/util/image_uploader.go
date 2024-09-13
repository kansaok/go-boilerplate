package util

import (
	"errors"
	"mime/multipart"
	"path/filepath"
	"strings"

	"github.com/kansaok/go-boilerplate/internal/config/storage"
)

var allowedImageTypes = []string{"image/jpeg", "image/png", "image/gif"}

func UploadImage(file multipart.File, header *multipart.FileHeader) (string, error) {
	// Validate file type
	if !isValidImage(header.Filename) {
		return "", errors.New("invalid image format")
	}

	// Process metadata, e.g., file size
	if header.Size > 5*1024*1024 { // 5 MB max
		return "", errors.New("file is too large")
	}

	// Save file to storage
	filePath, err := storage.SaveLocal(file, header.Filename)
	if err != nil {
		return "", err
	}

	return filePath, nil
}

func isValidImage(filename string) bool {
	ext := strings.ToLower(filepath.Ext(filename))
	for _, allowedType := range allowedImageTypes {
		if strings.HasSuffix(allowedType, ext) {
			return true
		}
	}
	return false
}
