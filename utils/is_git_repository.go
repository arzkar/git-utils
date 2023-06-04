package utils

import (
	"os"
	"path/filepath"
)

func IsGitRepository(path string) bool {
	_, err := os.Stat(filepath.Join(path, ".git"))
	return err == nil
}