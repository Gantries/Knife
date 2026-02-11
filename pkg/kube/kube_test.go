package kube

import (
	"os"
	"path/filepath"
	"testing"
)

func TestNamespace(t *testing.T) {
	// Since FileNamespace is a constant pointing to Kubernetes paths,
	// and we're not running in Kubernetes, this should return an error
	_, err := Namespace()
	if err == nil {
		t.Error("Namespace() expected error when not in Kubernetes")
	}
}

func TestToken(t *testing.T) {
	// Since FileToken is a constant pointing to Kubernetes paths,
	// and we're not running in Kubernetes, this should return an error
	_, err := Token()
	if err == nil {
		t.Error("Token() expected error when not in Kubernetes")
	}
}

func TestCertificate(t *testing.T) {
	// Since FileCa is a constant pointing to Kubernetes paths,
	// and we're not running in Kubernetes, this should return an error
	_, err := Certificate()
	if err == nil {
		t.Error("Certificate() expected error when not in Kubernetes")
	}
}

func TestReadfile(t *testing.T) {
	t.Run("non-existent file", func(t *testing.T) {
		_, err := readfile("/non/existent/file")
		if err == nil {
			t.Error("readfile() expected error for non-existent file")
		}
	})

	t.Run("valid file", func(t *testing.T) {
		tmpDir := t.TempDir()
		testFile := filepath.Join(tmpDir, "test.txt")
		content := "test content"

		err := os.WriteFile(testFile, []byte(content), 0644)
		if err != nil {
			t.Fatalf("failed to create test file: %v", err)
		}

		result, err := readfile(testFile)
		if err != nil {
			t.Errorf("readfile() unexpected error: %v", err)
		}
		if result != content {
			t.Errorf("readfile() = %v, want %v", result, content)
		}
	})
}
