package utils

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/shubhamoys/forgego/internal/pkg/errors"
)

// CreateDir creates a directory at the specified path if it doesn't exist.
func CreateDir(path string) error {
	if err := os.MkdirAll(path, 0755); err != nil {
		return fmt.Errorf("%w: failed to create directory %s: %v", errors.ErrFileOperation, path, err)
	}
	return nil
}

// WriteFile writes content to a file at the specified path, creating parent directories if needed.
func WriteFile(path, content string) error {
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return fmt.Errorf("%w: failed to create parent directories for %s: %v", errors.ErrFileOperation, path, err)
	}
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		return fmt.Errorf("%w: failed to write file %s: %v", errors.ErrFileOperation, path, err)
	}
	return nil
}

// InitGit initializes a Git repository in the specified directory.
func InitGit(path string) error {
	cmd := exec.Command("git", "init")
	cmd.Dir = path
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("%w: failed to initialize Git repository in %s: %v", errors.ErrGitOperation, path, err)
	}
	return nil
}
