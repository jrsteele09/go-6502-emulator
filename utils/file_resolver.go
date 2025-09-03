package utils

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// FileResolver is an interface for resolving file paths to readers
type FileResolver interface {
	Resolve(path string) (io.Reader, error)
}

// OSFileResolver resolves file paths using the OS filesystem
type OSFileResolver struct {
	baseDir string
}

// NewOSFileResolver creates a new OS-based file resolver
func NewOSFileResolver(baseDir string) *OSFileResolver {
	return &OSFileResolver{baseDir: baseDir}
}

// Resolve opens a file from the filesystem
func (r *OSFileResolver) Resolve(path string) (io.Reader, error) {
	// If path is relative, make it relative to baseDir
	if !filepath.IsAbs(path) {
		path = filepath.Join(r.baseDir, path)
	}

	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open include file '%s': %w", path, err)
	}

	// Read the entire file into memory so we can close it immediately
	content, err := io.ReadAll(file)
	file.Close()
	if err != nil {
		return nil, fmt.Errorf("failed to read include file '%s': %w", path, err)
	}

	return bytes.NewReader(content), nil
}

// MemoryFileResolver resolves file paths from an in-memory map (useful for testing)
type MemoryFileResolver struct {
	files map[string]string
}

// NewMemoryFileResolver creates a new memory-based file resolver
func NewMemoryFileResolver(files map[string]string) *MemoryFileResolver {
	return &MemoryFileResolver{files: files}
}

// Resolve returns a reader for the specified file from memory
func (r *MemoryFileResolver) Resolve(path string) (io.Reader, error) {
	content, exists := r.files[path]
	if !exists {
		return nil, fmt.Errorf("include file '%s' not found in memory resolver", path)
	}
	return strings.NewReader(content), nil
}
