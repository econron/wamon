package cmd

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestExportCommandDirectly tests the export command functionality directly
func TestExportCommandDirectly(t *testing.T) {
	// Skip in normal test runs to avoid creating files
	if testing.Short() {
		t.Skip("Skipping in short mode")
	}

	// Create temporary directory for test
	tempDir, err := os.MkdirTemp("", "wamon-export-direct-test-*")
	assert.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Set export file path
	exportFilePath := filepath.Join(tempDir, "export.json")

	// Set database path to an in-memory database
	originalDBPath := dbPath
	dbPath = ":memory:"
	defer func() {
		dbPath = originalDBPath
	}()

	// Run the export command
	exportCmd.Run(exportCmd, []string{exportFilePath})

	// Check if the file was created (even if empty)
	_, err = os.Stat(exportFilePath)
	assert.NoError(t, err)
}
