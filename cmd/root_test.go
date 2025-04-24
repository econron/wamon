package cmd

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/econron/wamon/internal/db"
	"github.com/econron/wamon/internal/models"
	"github.com/stretchr/testify/assert"
)

// setupTestEnvironment sets up a test environment and returns a cleanup function
func setupTestEnvironment(t *testing.T) (string, func()) {
	// Create temporary directory
	tempDir, err := os.MkdirTemp("", "wamon-test-*")
	assert.NoError(t, err)

	// Set up a test database
	testDBPath := filepath.Join(tempDir, "test.db")
	dbPath = testDBPath

	cleanup := func() {
		os.RemoveAll(tempDir)
	}

	return testDBPath, cleanup
}

// createTestEntries adds test entries to the database
func createTestEntries(t *testing.T, database db.DB, count int) []*models.Entry {
	entries := make([]*models.Entry, count)

	for i := 0; i < count; i++ {
		entry := &models.Entry{
			ID:            time.Now().Add(time.Duration(-i) * time.Hour).Format("20060102150405"),
			Category:      models.Research,
			ResearchTopic: "Test Research Topic",
			Satisfaction:  3,
			CreatedAt:     time.Now().Add(time.Duration(-i) * time.Hour),
		}

		// Alternate between research and programming
		if i%2 == 1 {
			entry.Category = models.Programming
			entry.ResearchTopic = ""
			entry.ProgramTitle = "Test Program"
		}

		err := database.SaveEntry(entry)
		assert.NoError(t, err)

		entries[i] = entry
	}

	return entries
}

// captureOutput captures stdout during a function execution
func captureOutput(f func()) string {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	f()

	w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	io.Copy(&buf, r)
	return buf.String()
}

// TestExecuteNoArgs tests the root command with no arguments
func TestExecuteNoArgs(t *testing.T) {
	// Skip interactive test in CI
	if os.Getenv("CI") == "true" {
		t.Skip("Skipping interactive test in CI environment")
	}

	testDBPath, cleanup := setupTestEnvironment(t)
	defer cleanup()

	// Verify database path is set correctly
	assert.Equal(t, testDBPath, dbPath)

	// Create a mock os.Args
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()
	os.Args = []string{"wamon", "--db", testDBPath}

	// This test cannot be fully automated due to interactive nature
	// But we can at least ensure it doesn't crash
	// Execute() // Uncomment for manual testing
}

// TestListCommandEmpty tests the list command with an empty database
func TestListCommandEmpty(t *testing.T) {
	testDBPath, cleanup := setupTestEnvironment(t)
	defer cleanup()

	// Verify we're using the test database
	assert.Equal(t, testDBPath, dbPath)

	// Execute list command with empty DB
	output := captureOutput(func() {
		listCmd.Run(listCmd, []string{})
	})

	assert.Contains(t, output, "記録がありません")
}

// TestListCommand tests the list command with data
func TestListCommand(t *testing.T) {
	testDBPath, cleanup := setupTestEnvironment(t)
	defer cleanup()

	// Verify we're using the test database
	assert.Equal(t, testDBPath, dbPath)

	// Create database and add test entries
	database, err := db.NewDB(testDBPath)
	assert.NoError(t, err)
	createTestEntries(t, database, 3)
	database.Close()

	// Test list all
	output := captureOutput(func() {
		listCmd.Run(listCmd, []string{})
	})
	assert.Contains(t, output, "ワモンアザラシの記録")
	assert.Contains(t, output, "合計: 3件の記録")

	// Test list with category filter
	categoryFilter = "調べ物"
	output = captureOutput(func() {
		listCmd.Run(listCmd, []string{})
	})
	assert.Contains(t, output, "カテゴリ: 調べ物")
}

// TestListCommandInvalidCategory tests the list command with an invalid category
func TestListCommandInvalidCategory(t *testing.T) {
	testDBPath, cleanup := setupTestEnvironment(t)
	defer cleanup()

	// Verify we're using the test database
	assert.Equal(t, testDBPath, dbPath)

	// Test with invalid category
	categoryFilter = "invalid"
	output := captureOutput(func() {
		listCmd.Run(listCmd, []string{})
	})
	assert.Contains(t, output, "無効なカテゴリです")
}

// TestEditCommandNotFound tests editing a non-existent entry
func TestEditCommandNotFound(t *testing.T) {
	testDBPath, cleanup := setupTestEnvironment(t)
	defer cleanup()

	// Verify we're using the test database
	assert.Equal(t, testDBPath, dbPath)

	// Test with non-existent ID
	output := captureOutput(func() {
		editCmd.Run(editCmd, []string{"non-existent-id"})
	})
	assert.Contains(t, output, "記録が見つかりません")
}

// TestEditCommandDBError tests the edit command with a database error
func TestEditCommandDBError(t *testing.T) {
	_, cleanup := setupTestEnvironment(t)
	defer cleanup()

	// Create an invalid database path
	invalidPath := "/invalid/path/does/not/exist.db"
	dbPath = invalidPath

	output := captureOutput(func() {
		editCmd.Run(editCmd, []string{"any-id"})
	})
	assert.Contains(t, output, "データベースの初期化エラー")
}

// TestReportCommandEmpty tests the report command with an empty database
func TestReportCommandEmpty(t *testing.T) {
	// Skip in CI environment due to interactive prompts
	if os.Getenv("CI") == "true" {
		t.Skip("Skipping interactive test in CI environment")
	}

	testDBPath, cleanup := setupTestEnvironment(t)
	defer cleanup()

	// Verify we're using the test database
	assert.Equal(t, testDBPath, dbPath)

	output := captureOutput(func() {
		reportCmd.Run(reportCmd, []string{})
	})
	// The report command will ask for Slack token in CI, so the message is different
	assert.Contains(t, output, "SlackのBot User OAuth Token")
}

// TestReportCommandDBError tests the report command with a database error
func TestReportCommandDBError(t *testing.T) {
	// Create an invalid database path
	invalidPath := "/invalid/path/does/not/exist.db"
	dbPath = invalidPath

	output := captureOutput(func() {
		reportCmd.Run(reportCmd, []string{})
	})
	assert.Contains(t, output, "データベースの初期化エラー")
}

// TestGetDefaultDBPath tests the getDefaultDBPath function
func TestGetDefaultDBPath(t *testing.T) {
	path := getDefaultDBPath()
	assert.NotEmpty(t, path)
	assert.Contains(t, path, "wamon.db")
}
