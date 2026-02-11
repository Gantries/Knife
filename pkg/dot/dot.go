// Package dot contains utilities can be used extreme frequently.
package dot

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/gantries/knife/pkg/fs"
)

// MaxLineLength define maximum line length for environment file.
const MaxLineLength = 2048

// Must treat return value as expected mostly.
func Must[T any](v T, err error) T {
	if err != nil {
		panic(err)
	}
	return v
}

// Env get environment from .env file in the nearest parent folder, should be used in unit test or similar scenes.
func Env(env string, maxLength int) (string, error) {
	parent, err := os.Getwd()
	if err != nil {
		return "", err
	}

	for last := ""; parent != last; last, parent = parent, filepath.Dir(parent) {
		dotenv := filepath.Join(parent, ".env")
		if _, err = os.Stat(dotenv); err == nil {
			if found, value, err := fs.WithScanner(dotenv, maxLength, func(line string) (bool, string) {
				if strings.HasPrefix(line, env+"=") {
					// stop on first match
					return true, line[len(env)+1:]
				}
				return false, ""
			}); found {
				return value, err
			}
		}
	}
	return "", fmt.Errorf("Environment file not found")
}
