package util

import (
	"errors"
	"mime/multipart"
	"path/filepath"
	"strings"

	"github.com/kansaok/go-boilerplate/internal/config/storage"
)

func UploadTXT(file multipart.File, header *multipart.FileHeader) (string, error) {
	// Validate file type
	if !isValidTXT(header.Filename) {
		return "", errors.New("invalid TXT format")
	}

	// Process metadata, e.g., file size
	if header.Size > 2*1024*1024 { // 2 MB max
		return "", errors.New("file is too large")
	}

	// Save file to storage
	filePath, err := storage.SaveLocal(file, header.Filename)
	if err != nil {
		return "", err
	}

	return filePath, nil
}

func isValidTXT(filename string) bool {
	return strings.ToLower(filepath.Ext(filename)) == ".txt"
}
