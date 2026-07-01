package storage

import (
	"io"
	"os"
	"path/filepath"
)

const uploadDir = "./uploads" // The folder where the file is stored

// SaveLocal stores the file to the local filesystem.
func SaveLocal(file io.Reader, filename string) (string, error) {
	// Create the folder if it doesn't exist.
	if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
		return "", err
	}

	// Determine the complete path of the file.
	filePath := filepath.Join(uploadDir, filename)

	// Create the file at the specified location.
	destFile, err := os.Create(filePath)
	if err != nil {
		return "", err
	}
	defer destFile.Close()

	// Copy the file contents to the destination.
	if _, err := io.Copy(destFile, file); err != nil {
		return "", err
	}

	return filePath, nil
}
