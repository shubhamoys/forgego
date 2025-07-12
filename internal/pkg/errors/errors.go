package errors

import "errors"

var (
	ErrPromptFailed       = errors.New("failed to get user input")
	ErrInvalidProjectName = errors.New("invalid project name: must not be empty or contain spaces or slashes")
	ErrInvalidPackageName = errors.New("invalid package name: must start with 'github.com/'")
	ErrFileOperation      = errors.New("file operation failed")
	ErrGitOperation       = errors.New("git operation failed")
)
