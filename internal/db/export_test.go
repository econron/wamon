package db

import (
	"encoding/json"
	"os"
	"testing"
	"time"

	"github.com/econron/wamon/internal/models"
	"github.com/stretchr/testify/assert"
)

func TestExportEntries(t *testing.T) {
	// Setup test database with some entries
	db := setupTestDB(t)

	// Create test entries with fixed timestamps for easier assertions
	entries := []*models.Entry{
		{
			ID:            "20220101120000",
			Category:      models.Research,
			ResearchTopic: "How to write unit tests",
			ProgramTitle:  "",
			Satisfaction:  4,
			CreatedAt:     time.Date(2022, 1, 1, 12, 0, 0, 0, time.UTC),
		},
		{
			ID:            "20220102130000",
			Category:      models.Programming,
			ResearchTopic: "",
			ProgramTitle:  "Refactor connection pool",
			Satisfaction:  5,
			CreatedAt:     time.Date(2022, 1, 2, 13, 0, 0, 0, time.UTC),
		},
		{
			ID:            "20220103140000",
			Category:      models.ResearchAndProgram,
			ResearchTopic: "SQL optimization",
			ProgramTitle:  "Implement query caching",
			Satisfaction:  3,
			CreatedAt:     time.Date(2022, 1, 3, 14, 0, 0, 0, time.UTC),
		},
	}

	// Save entries to database
	for _, entry := range entries {
		err := db.SaveEntry(entry)
		assert.NoError(t, err)
	}

	// Create temporary file for export
	tempFile, err := os.CreateTemp("", "wamon-export-*.json")
	assert.NoError(t, err)
	tempFile.Close()
	defer os.Remove(tempFile.Name())

	// Test ExportEntries
	err = db.(*SQLiteDB).ExportEntries(tempFile.Name())
	assert.NoError(t, err)

	// Read exported file and verify contents
	exportedData, err := os.ReadFile(tempFile.Name())
	assert.NoError(t, err)

	// Split by newlines to get individual JSON objects
	lines := splitLines(string(exportedData))
	assert.Equal(t, len(entries), len(lines))

	// Parse each line and verify it matches our entries
	for i, line := range lines {
		var exportedEntry map[string]interface{}
		err := json.Unmarshal([]byte(line), &exportedEntry)
		assert.NoError(t, err)

		// Verify basic fields
		assert.Equal(t, entries[len(entries)-i-1].ID, exportedEntry["id"])

		// Format timestamp to match expected format (ISO 8601)
		expectedTime := entries[len(entries)-i-1].CreatedAt.Format(time.RFC3339)
		assert.Equal(t, expectedTime, exportedEntry["ts"])

		// Check category field
		assert.Equal(t, string(entries[len(entries)-i-1].Category), exportedEntry["cat"])

		// Verify body contains appropriate content based on entry type
		switch entries[len(entries)-i-1].Category {
		case models.Research:
			assert.Equal(t, entries[len(entries)-i-1].ResearchTopic, exportedEntry["body"])
		case models.Programming:
			assert.Equal(t, entries[len(entries)-i-1].ProgramTitle, exportedEntry["body"])
		case models.ResearchAndProgram:
			combinedBody := entries[len(entries)-i-1].ResearchTopic + " - " + entries[len(entries)-i-1].ProgramTitle
			assert.Equal(t, combinedBody, exportedEntry["body"])
		}
	}

	// Test export when there are no entries
	emptyDB := setupTestDB(t)
	tempEmptyFile, err := os.CreateTemp("", "wamon-export-empty-*.json")
	assert.NoError(t, err)
	tempEmptyFile.Close()
	defer os.Remove(tempEmptyFile.Name())

	err = emptyDB.(*SQLiteDB).ExportEntries(tempEmptyFile.Name())
	assert.NoError(t, err)

	// Empty file should exist but be empty
	emptyExportedData, err := os.ReadFile(tempEmptyFile.Name())
	assert.NoError(t, err)
	assert.Empty(t, string(emptyExportedData))
}

// Helper function to split the file into lines
func splitLines(data string) []string {
	var result []string

	// Remove any trailing newline
	if len(data) > 0 && data[len(data)-1] == '\n' {
		data = data[:len(data)-1]
	}

	// If there's no data, return empty slice
	if data == "" {
		return result
	}

	// Split data by newlines and filter out empty lines
	lines := []byte(data)
	start := 0
	for i, b := range lines {
		if b == '\n' {
			line := string(lines[start:i])
			if line != "" {
				result = append(result, line)
			}
			start = i + 1
		}
	}

	// Add the last line if it's not empty
	if start < len(lines) {
		line := string(lines[start:])
		if line != "" {
			result = append(result, line)
		}
	}

	return result
}

func TestExportEntriesErrorHandling(t *testing.T) {
	// Setup test database
	db := setupTestDB(t)

	// Test with invalid path (directory that doesn't exist)
	err := db.(*SQLiteDB).ExportEntries("/invalid/path/that/doesnt/exist.json")
	assert.Error(t, err)

	// Create a directory with the same name as our intended file to force an error
	tempDir, err := os.MkdirTemp("", "wamon-export-dir-*")
	assert.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Try to export to a directory instead of a file
	err = db.(*SQLiteDB).ExportEntries(tempDir)
	assert.Error(t, err)
}

func TestExportEntriesSince(t *testing.T) {
	// Setup test database with some entries
	db := setupTestDB(t)

	// Create test entries with different timestamps
	entries := []*models.Entry{
		{
			ID:            "20220101120000",
			Category:      models.Research,
			ResearchTopic: "Old entry",
			ProgramTitle:  "",
			Satisfaction:  4,
			CreatedAt:     time.Date(2022, 1, 1, 12, 0, 0, 0, time.UTC),
		},
		{
			ID:            "20220102130000",
			Category:      models.Programming,
			ResearchTopic: "",
			ProgramTitle:  "Medium entry",
			Satisfaction:  5,
			CreatedAt:     time.Date(2022, 1, 2, 13, 0, 0, 0, time.UTC),
		},
		{
			ID:            "20220103140000",
			Category:      models.ResearchAndProgram,
			ResearchTopic: "Recent entry",
			ProgramTitle:  "Testing",
			Satisfaction:  3,
			CreatedAt:     time.Date(2022, 1, 3, 14, 0, 0, 0, time.UTC),
		},
	}

	// Save entries to database
	for _, entry := range entries {
		err := db.SaveEntry(entry)
		assert.NoError(t, err)
	}

	// Create temporary file for export
	tempFile, err := os.CreateTemp("", "wamon-export-since-*.json")
	assert.NoError(t, err)
	tempFile.Close()
	defer os.Remove(tempFile.Name())

	// Set a cutoff date to filter by (include only entries from January 2, 2022 and later)
	cutoffDate := time.Date(2022, 1, 2, 0, 0, 0, 0, time.UTC)

	// Test ExportEntriesSince
	err = db.(*SQLiteDB).ExportEntriesSince(tempFile.Name(), cutoffDate)
	assert.NoError(t, err)

	// Read exported file and verify contents
	exportedData, err := os.ReadFile(tempFile.Name())
	assert.NoError(t, err)

	// Split by newlines to get individual JSON objects
	lines := splitLines(string(exportedData))
	assert.Equal(t, 2, len(lines)) // Should only include the two most recent entries

	// Parse each line and verify the entries are the correct ones
	var exportedIDs []string
	for _, line := range lines {
		var exportedEntry map[string]interface{}
		err := json.Unmarshal([]byte(line), &exportedEntry)
		assert.NoError(t, err)
		exportedIDs = append(exportedIDs, exportedEntry["id"].(string))
	}

	// Check that we only have the two most recent entries
	assert.Contains(t, exportedIDs, "20220102130000")
	assert.Contains(t, exportedIDs, "20220103140000")
	assert.NotContains(t, exportedIDs, "20220101120000") // This one should be filtered out

	// Test with a cutoff that excludes all entries
	futureCutoff := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	tempEmptyFile, err := os.CreateTemp("", "wamon-export-empty-since-*.json")
	assert.NoError(t, err)
	tempEmptyFile.Close()
	defer os.Remove(tempEmptyFile.Name())

	err = db.(*SQLiteDB).ExportEntriesSince(tempEmptyFile.Name(), futureCutoff)
	assert.NoError(t, err)

	// Empty file should exist but have no entries
	emptyExportedData, err := os.ReadFile(tempEmptyFile.Name())
	assert.NoError(t, err)
	assert.Empty(t, string(emptyExportedData))
}
