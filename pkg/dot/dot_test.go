package dot

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

func Test_Must(t *testing.T) {
	if Must(func() (string, error) { return "", nil }()) != "" {
		t.Fatalf("expected empty string but not")
	}

	if Must(func() (int, error) { return 1, nil }()) != 1 {
		t.Fatalf("expected 1 but not")
	}

	defer func() {
		if err := recover(); err == nil {
			t.Fatalf("expected panic but not")
		}
	}()
	Must(func() (string, error) { return "", fmt.Errorf("fail") }())
}

func Test_Env(t *testing.T) {
	// Skip if no .env file exists in current or parent directories
	cwd, err := os.Getwd()
	if err != nil {
		t.Skip("Cannot get working directory")
		return
	}

	hasEnvFile := false
	for last := ""; cwd != last; last = cwd {
		cwd = filepath.Dir(cwd)
		dotenv := filepath.Join(cwd, ".env")
		if _, err := os.Stat(dotenv); err == nil {
			hasEnvFile = true
			break
		}
	}

	if !hasEnvFile {
		t.Skip("No .env file found in current or parent directories")
		return
	}

	key, err := Env("TEST_KEY", MaxLineLength)
	if err != nil {
		t.Fatalf("error to get env: %s", "TEST_KEY")
	}
	if len(key) <= 0 {
		t.Fatalf("expected a valid string but empty")
	}
	_, err = Env("ANYTHING_NOT_EXIST", MaxLineLength)
	if err == nil {
		t.Fatalf("expected error but not")
	}
}
