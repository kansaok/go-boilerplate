package util

import (
	"errors"
	"io"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/kansaok/go-boilerplate/internal/config/storage"
)

var allowedExcelExtensions = map[string]bool{
	".xls":  true,
	".xlsx": true,
}

var allowedExcelMimeTypes = map[string]bool{
	"application/vnd.ms-excel":            true,
	"application/vnd.openxmlformats-officedocument.spreadsheetml.sheet": true,
	"application/zip":                     true,
}

func UploadExcel(file multipart.File, header *multipart.FileHeader) (string, error) {
	ext := strings.ToLower(filepath.Ext(header.Filename))
	if !allowedExcelExtensions[ext] {
		return "", errors.New("invalid Excel format")
	}

	if _, err := file.Seek(0, io.SeekStart); err != nil {
		return "", errors.New("failed to read file")
	}
	buf := make([]byte, 512)
	n, err := file.Read(buf)
	if err != nil {
		return "", errors.New("failed to read file for validation")
	}
	detectedMime := http.DetectContentType(buf[:n])
	if !allowedExcelMimeTypes[detectedMime] {
		return "", errors.New("file content does not match Excel format")
	}
	if _, err := file.Seek(0, io.SeekStart); err != nil {
		return "", errors.New("failed to reset file pointer")
	}

	if header.Size > 15*1024*1024 {
		return "", errors.New("file is too large")
	}

	filePath, err := storage.SaveLocal(file, header.Filename)
	if err != nil {
		return "", err
	}

	return filePath, nil
}
