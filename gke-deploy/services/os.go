package services

import (
	"context"
	"io/ioutil"
	"os"
)

// OS implements the OSService interface.
type OS struct{}

// NewOS returns a new OS object.
func NewOS(ctx context.Context) (*OS, error) {
	return &OS{}, nil
}

// Stat gets a file description for a file filename.
func (o *OS) Stat(ctx context.Context, filename string) (os.FileInfo, error) {
	if filename == "-" {
		return os.Stdin.Stat()
	}
	return os.Stat(filename)
}

// ReadDir gets file descriptions for all files contained in a directory dirname.
func (o *OS) ReadDir(ctx context.Context, dirname string) ([]os.FileInfo, error) {
	return ioutil.ReadDir(dirname)
}

// ReadFile gets the entire contents of a file filename as bytes.
func (o *OS) ReadFile(ctx context.Context, filename string) ([]byte, error) {
	if filename == "-" {
		return ioutil.ReadAll(os.Stdin)
	}
	return ioutil.ReadFile(filename)
}

// WriteFile writes data to a file.
func (o *OS) WriteFile(ctx context.Context, filename string, data []byte, perm os.FileMode) error {
	return ioutil.WriteFile(filename, data, perm)
}

// MkdirAll creates a directory dirname, including all parent directories if they do not exist.
func (o *OS) MkdirAll(ctx context.Context, dirname string, perm os.FileMode) error {
	return os.MkdirAll(dirname, perm)
}

// RemoveAll removes path and any children it contains.
func (o *OS) RemoveAll(ctx context.Context, path string) error {
	return os.RemoveAll(path)
}

// TempDir creates a new temporary directory in the directory dir.
func (o *OS) TempDir(ctx context.Context, dir, pattern string) (string, error) {
	return ioutil.TempDir(dir, pattern)
}
