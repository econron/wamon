package cmd

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestExportDBError tests the export command with database error
func TestExportDBError(t *testing.T) {
	// Save original dbPath
	originalDBPath := dbPath
	defer func() { dbPath = originalDBPath }()

	// Set an invalid database path
	dbPath = "/path/to/nonexistent/db.sqlite"

	// Capture output
	output := captureOutput(func() {
		exportCmd.Run(exportCmd, []string{"test_export.json"})
	})

	// Verify error message
	assert.Contains(t, output, "データベースの初期化エラー")
}

// TestExportWriteError tests the export command with file write error
func TestExportWriteError(t *testing.T) {
	// Create test environment
	testDBPath, cleanup := setupTestEnvironment(t)
	defer cleanup()

	// Set database path
	originalDBPath := dbPath
	dbPath = testDBPath
	defer func() { dbPath = originalDBPath }()

	// Try to write to a location that should be unwritable
	unwritablePath := filepath.Join("/", "definitely_unwritable_system_path.json")

	// Capture output
	output := captureOutput(func() {
		exportCmd.Run(exportCmd, []string{unwritablePath})
	})

	// Verify error message
	// Note: This will only work if the test is run without root privileges
	// The test might be skipped on some systems where the path is writable
	if !strings.Contains(output, "エクスポートされました") {
		assert.Contains(t, output, "エクスポートエラー")
	}
}
