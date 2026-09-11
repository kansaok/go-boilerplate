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

func UploadTXT(file multipart.File, header *multipart.FileHeader) (string, error) {
	ext := strings.ToLower(filepath.Ext(header.Filename))
	if ext != ".txt" {
		return "", errors.New("invalid TXT format")
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
	if !strings.HasPrefix(detectedMime, "text/") {
		return "", errors.New("file content does not match TXT format")
	}
	if _, err := file.Seek(0, io.SeekStart); err != nil {
		return "", errors.New("failed to reset file pointer")
	}

	if header.Size > 2*1024*1024 {
		return "", errors.New("file is too large")
	}

	filePath, err := storage.SaveLocal(file, header.Filename)
	if err != nil {
		return "", err
	}

	return filePath, nil
}
