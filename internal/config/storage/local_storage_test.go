package storage

import (
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestSaveLocal_RejectsPathTraversal(t *testing.T) {
	traversalPayloads := []string{
		"../../etc/crontab",
		"..%2f..%2fetc%2fpasswd",
		"../../../../etc/cron.d/malicious",
		"/etc/passwd",
		"sub/../../../../root/.bashrc",
	}

	for _, payload := range traversalPayloads {
		path, err := SaveLocal(strings.NewReader("payload"), payload)
		if err != nil {
			t.Fatalf("SaveLocal returned unexpected error for %q: %v", payload, err)
		}

		absUploadDir, _ := filepath.Abs(uploadDir)
		absPath, _ := filepath.Abs(path)

		// The returned path must always stay inside the upload directory
		if !strings.HasPrefix(absPath, absUploadDir+string(filepath.Separator)) {
			t.Errorf("path traversal escaped upload dir for payload %q: %s", payload, path)
		}
	}
}

func TestSaveLocal_GeneratesUniqueSanitizedFilenames(t *testing.T) {
	first, err := SaveLocal(strings.NewReader("hello"), "photo.jpg")
	if err != nil {
		t.Fatalf("SaveLocal failed: %v", err)
	}
	second, err := SaveLocal(strings.NewReader("hello"), "photo.jpg")
	if err != nil {
		t.Fatalf("SaveLocal failed: %v", err)
	}

	if first == second {
		t.Fatal("duplicate uploads should yield different randomized filenames")
	}

	if filepath.Base(first) == "photo.jpg" {
		t.Fatal("original filename must not be used directly")
	}

	if filepath.Ext(first) != ".jpg" {
		t.Errorf("expected .jpg extension preserved, got %q", filepath.Ext(first))
	}
}

func TestSaveLocal_PreservesFileContents(t *testing.T) {
	body := "file contents for verification"
	path, err := SaveLocal(strings.NewReader(body), "data.txt")
	if err != nil {
		t.Fatalf("SaveLocal failed: %v", err)
	}

	f, err := os.Open(path)
	if err != nil {
		t.Fatalf("failed to open stored file: %v", err)
	}
	defer f.Close()

	content, err := io.ReadAll(f)
	if err != nil {
		t.Fatalf("failed to read stored file: %v", err)
	}
	if string(content) != body {
		t.Errorf("stored content mismatch: got %q want %q", string(content), body)
	}
}