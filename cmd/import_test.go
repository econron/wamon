package cmd

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestImportCommand(t *testing.T) {
	// Create a temporary directory
	tempDir, err := os.MkdirTemp("", "wamon-import-cmd-test-*")
	assert.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Create test database
	dbPath := filepath.Join(tempDir, "test.db")
	os.Args = []string{"wamon", "--db", dbPath}

	// Create test export file
	exportFilePath := filepath.Join(tempDir, "import_test.json")
	file, err := os.Create(exportFilePath)
	assert.NoError(t, err)

	// Write test data
	testData := `{"id":"20220101120000","ts":"2022-01-01T12:00:00+09:00","cat":"research","body":"テストデータ"}`
	_, err = file.WriteString(testData)
	assert.NoError(t, err)
	err = file.Sync() // 確実に書き込み
	assert.NoError(t, err)
	file.Close()

	// Capture stdout
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Run the import command
	os.Args = []string{"wamon", "--db", dbPath, "import", exportFilePath}
	err = importCmd.Execute()
	assert.NoError(t, err)

	// Restore stdout
	w.Close()
	os.Stdout = oldStdout

	// Read output
	var buf bytes.Buffer
	io.Copy(&buf, r)
	output := buf.String()

	// Check output contains expected messages
	assert.Contains(t, output, "エントリを正常にインポートしました")
	assert.Contains(t, output, "現在のデータベースには合計")
}

func TestImportCommandWithInvalidFile(t *testing.T) {
	// Create a temporary directory
	tempDir, err := os.MkdirTemp("", "wamon-import-cmd-invalid-test-*")
	assert.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Create test database
	dbPath := filepath.Join(tempDir, "test.db")
	os.Args = []string{"wamon", "--db", dbPath}

	// Non-existent file
	nonExistentFile := filepath.Join(tempDir, "does_not_exist.json")

	// Capture stdout
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Run the import command with non-existent file
	os.Args = []string{"wamon", "--db", dbPath, "import", nonExistentFile}
	err = importCmd.Execute()
	assert.NoError(t, err) // Command itself should not error, but output should show error

	// Restore stdout
	w.Close()
	os.Stdout = oldStdout

	// Read output
	var buf bytes.Buffer
	io.Copy(&buf, r)
	output := buf.String()

	// Check output contains error message
	assert.Contains(t, output, "インポートエラー")
}
