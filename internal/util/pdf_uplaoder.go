package util

import (
	"errors"
	"mime/multipart"
	"path/filepath"
	"strings"

	"github.com/kansaok/go-boilerplate/internal/config/storage"
)

func UploadPDF(file multipart.File, header *multipart.FileHeader) (string, error) {
	// Validate file type
	if !isValidPDF(header.Filename) {
		return "", errors.New("invalid PDF format")
	}

	// Process metadata, e.g., file size
	if header.Size > 10*1024*1024 { // 10 MB max
		return "", errors.New("file is too large")
	}

	// Save file to storage
	filePath, err := storage.SaveLocal(file, header.Filename)
	if err != nil {
		return "", err
	}

	return filePath, nil
}

func isValidPDF(filename string) bool {
	return strings.ToLower(filepath.Ext(filename)) == ".pdf"
}
