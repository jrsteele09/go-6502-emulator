package assembler

import (
	"bufio"
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

// Preprocessor handles file inclusion before tokenization
type Preprocessor struct {
	resolver      FileResolver
	maxDepth      int
	includedFiles map[string]bool // Track included files to prevent circular includes
}

// NewPreprocessor creates a new preprocessor with the given file resolver
func NewPreprocessor(resolver FileResolver) *Preprocessor {
	return &Preprocessor{
		resolver:      resolver,
		maxDepth:      10, // Reasonable default for include depth
		includedFiles: make(map[string]bool),
	}
}

// SetMaxDepth sets the maximum include depth to prevent infinite recursion
func (p *Preprocessor) SetMaxDepth(depth int) {
	p.maxDepth = depth
}

// Process processes the input, expanding all include directives
func (p *Preprocessor) Process(input io.Reader) (io.Reader, error) {
	// Reset included files for each processing session
	p.includedFiles = make(map[string]bool)

	result, err := p.processReader(input, 0)
	if err != nil {
		return nil, err
	}

	return strings.NewReader(result), nil
}

// processReader recursively processes a reader, expanding includes
func (p *Preprocessor) processReader(input io.Reader, depth int) (string, error) {
	if depth > p.maxDepth {
		return "", fmt.Errorf("maximum include depth (%d) exceeded", p.maxDepth)
	}

	var result strings.Builder
	scanner := bufio.NewScanner(input)
	lineNum := 0

	for scanner.Scan() {
		lineNum++
		line := scanner.Text()
		trimmedLine := strings.TrimSpace(line)

		// Check for include directives
		includePath := p.extractIncludePath(trimmedLine)
		if includePath != "" {
			// Prevent circular includes
			if p.includedFiles[includePath] {
				return "", fmt.Errorf("circular include detected: '%s' (line %d)", includePath, lineNum)
			}

			// Mark this file as included
			p.includedFiles[includePath] = true

			// Resolve and process the included file
			includeReader, err := p.resolver.Resolve(includePath)
			if err != nil {
				return "", fmt.Errorf("line %d: %w", lineNum, err)
			}

			includedContent, err := p.processReader(includeReader, depth+1)
			if err != nil {
				return "", fmt.Errorf("in file '%s': %w", includePath, err)
			}

			// Add the included content
			result.WriteString(includedContent)

			// Remove from included files set when done (allows including same file in different branches)
			delete(p.includedFiles, includePath)
		} else {
			// Regular line, just copy it
			result.WriteString(line)
			result.WriteString("\n")
		}
	}

	if err := scanner.Err(); err != nil {
		return "", fmt.Errorf("error reading input: %w", err)
	}

	return result.String(), nil
}

// extractIncludePath extracts the file path from include directives
// Supports both #include "file.asm" and .include "file.asm" formats
func (p *Preprocessor) extractIncludePath(line string) string {
	line = strings.TrimSpace(line)

	// Handle #include directive
	if strings.HasPrefix(line, "#include") {
		return p.extractQuotedPath(line[8:]) // Skip "#include"
	}

	// Handle .include directive
	if strings.HasPrefix(line, ".include") {
		return p.extractQuotedPath(line[8:]) // Skip ".include"
	}

	// Handle .INCLUDE directive (uppercase)
	if strings.HasPrefix(line, ".INCLUDE") {
		return p.extractQuotedPath(line[8:]) // Skip ".INCLUDE"
	}

	return ""
}

// extractQuotedPath extracts a quoted file path from the remaining part of an include line
func (p *Preprocessor) extractQuotedPath(remainder string) string {
	remainder = strings.TrimSpace(remainder)

	// Handle both single and double quotes
	if len(remainder) < 2 {
		return ""
	}

	// Check for double quotes
	if remainder[0] == '"' {
		if endPos := strings.Index(remainder[1:], `"`); endPos != -1 {
			return remainder[1 : endPos+1]
		}
	}

	// Check for single quotes
	if remainder[0] == '\'' {
		if endPos := strings.Index(remainder[1:], `'`); endPos != -1 {
			return remainder[1 : endPos+1]
		}
	}

	// Handle unquoted paths (space-delimited)
	parts := strings.Fields(remainder)
	if len(parts) > 0 {
		return parts[0]
	}

	return ""
}
