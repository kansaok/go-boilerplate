package storage

import (
	"crypto/rand"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

const uploadDir = "./uploads"

func generateSafeFilename(original string) string {
	ext := strings.ToLower(filepath.Ext(original))
	b := make([]byte, 16)
	rand.Read(b)
	return fmt.Sprintf("%x%s", b, ext)
}

func SaveLocal(file io.Reader, filename string) (string, error) {
	if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
		return "", err
	}

	safeName := generateSafeFilename(filename)

	filePath := filepath.Join(uploadDir, safeName)
	cleanPath := filepath.Clean(filePath)

	absUploadDir, err := filepath.Abs(uploadDir)
	if err != nil {
		return "", fmt.Errorf("failed to resolve upload directory: %v", err)
	}
	absFilePath, err := filepath.Abs(cleanPath)
	if err != nil {
		return "", fmt.Errorf("failed to resolve file path: %v", err)
	}

	if !strings.HasPrefix(absFilePath, absUploadDir+string(filepath.Separator)) {
		return "", fmt.Errorf("invalid file path: path traversal detected")
	}

	destFile, err := os.Create(cleanPath)
	if err != nil {
		return "", err
	}
	defer destFile.Close()

	if _, err := io.Copy(destFile, file); err != nil {
		return "", err
	}

	return cleanPath, nil
}
