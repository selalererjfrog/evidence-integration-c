package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWriteToFile(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "jira-helper-test")
	assert.NoError(t, err)
	defer os.RemoveAll(tempDir)

	tests := []struct {
		name        string
		filename    string
		data        []byte
		expectError bool
		setup       func()
		teardown    func()
	}{
		{
			name:        "Write to simple file",
			filename:    filepath.Join(tempDir, "test.txt"),
			data:        []byte("Hello, World!"),
			expectError: false,
		},
		{
			name:        "Write to file in non-existent directory",
			filename:    filepath.Join(tempDir, "subdir", "test.txt"),
			data:        []byte("Hello from subdir!"),
			expectError: false,
		},
		{
			name:        "Write to deeply nested directory",
			filename:    filepath.Join(tempDir, "a", "b", "c", "test.txt"),
			data:        []byte("Deep file"),
			expectError: false,
		},
		{
			name:        "Overwrite existing file",
			filename:    filepath.Join(tempDir, "existing.txt"),
			data:        []byte("New content"),
			expectError: false,
			setup: func() {
				// Create existing file with different content
				os.WriteFile(filepath.Join(tempDir, "existing.txt"), []byte("Old content"), 0644)
			},
		},
		{
			name:        "Write empty data",
			filename:    filepath.Join(tempDir, "empty.txt"),
			data:        []byte{},
			expectError: false,
		},
		{
			name:        "Write to current directory (. case)",
			filename:    "test-current-dir.txt",
			data:        []byte("Current directory test"),
			expectError: false,
			teardown: func() {
				os.Remove("test-current-dir.txt")
			},
		},
		{
			name:        "Invalid filename (directory with no write permission)",
			filename:    "/root/test.txt",
			data:        []byte("Should fail"),
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Run setup if provided
			if tt.setup != nil {
				tt.setup()
			}

			// Run test
			err := writeToFile(tt.filename, tt.data)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)

				// Verify file was created with correct content
				content, err := os.ReadFile(tt.filename)
				assert.NoError(t, err)
				assert.Equal(t, tt.data, content)

				// Verify file permissions
				info, err := os.Stat(tt.filename)
				assert.NoError(t, err)
				assert.Equal(t, os.FileMode(0644), info.Mode().Perm())
			}

			// Run teardown if provided
			if tt.teardown != nil {
				tt.teardown()
			}
		})
	}
}

func TestWriteToFile_DirectoryCreation(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "jira-helper-dirtest")
	assert.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Test that directories are created with correct permissions
	filename := filepath.Join(tempDir, "new", "path", "file.txt")
	data := []byte("test data")

	err = writeToFile(filename, data)
	assert.NoError(t, err)

	// Check that all directories were created
	dirPath := filepath.Join(tempDir, "new", "path")
	info, err := os.Stat(dirPath)
	assert.NoError(t, err)
	assert.True(t, info.IsDir())
	assert.Equal(t, os.FileMode(0755), info.Mode().Perm())

	// Check file content
	content, err := os.ReadFile(filename)
	assert.NoError(t, err)
	assert.Equal(t, data, content)
}
