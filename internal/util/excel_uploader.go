package util

import (
	"errors"
	"mime/multipart"
	"path/filepath"
	"strings"

	"github.com/kansaok/go-boilerplate/internal/config/storage"
)

var allowedExcelTypes = []string{".xls", ".xlsx"}

func UploadExcel(file multipart.File, header *multipart.FileHeader) (string, error) {
	// Validate file type
	if !isValidExcel(header.Filename) {
		return "", errors.New("invalid Excel format")
	}

	// Process metadata, e.g., file size
	if header.Size > 15*1024*1024 { // 15 MB max
		return "", errors.New("file is too large")
	}

	// Save file to storage
	filePath, err := storage.SaveLocal(file, header.Filename)
	if err != nil {
		return "", err
	}

	return filePath, nil
}

func isValidExcel(filename string) bool {
	ext := strings.ToLower(filepath.Ext(filename))
	for _, allowedType := range allowedExcelTypes {
		if ext == allowedType {
			return true
		}
	}
	return false
}
