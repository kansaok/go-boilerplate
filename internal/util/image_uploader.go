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

var allowedImageMimeTypes = map[string]string{
	".jpg":  "image/jpeg",
	".jpeg": "image/jpeg",
	".png":  "image/png",
	".gif":  "image/gif",
	".webp": "image/webp",
}

func UploadImage(file multipart.File, header *multipart.FileHeader) (string, error) {
	ext := strings.ToLower(filepath.Ext(header.Filename))
	if allowedMime, ok := allowedImageMimeTypes[ext]; !ok {
		return "", errors.New("invalid image format")
	} else {
		if _, err := file.Seek(0, io.SeekStart); err != nil {
			return "", errors.New("failed to read file")
		}
		buf := make([]byte, 512)
		n, err := file.Read(buf)
		if err != nil {
			return "", errors.New("failed to read file for validation")
		}
		detectedMime := http.DetectContentType(buf[:n])
		if detectedMime != allowedMime {
			return "", errors.New("file content does not match declared extension")
		}
		if _, err := file.Seek(0, io.SeekStart); err != nil {
			return "", errors.New("failed to reset file pointer")
		}
	}

	if header.Size > 5*1024*1024 {
		return "", errors.New("file is too large")
	}

	filePath, err := storage.SaveLocal(file, header.Filename)
	if err != nil {
		return "", err
	}

	return filePath, nil
}
