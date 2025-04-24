package db

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/econron/wamon/internal/models"
	"github.com/stretchr/testify/assert"
)

func TestImportEntries(t *testing.T) {
	// Create a temporary database
	tempDir, err := os.MkdirTemp("", "wamon-import-test-*")
	assert.NoError(t, err)
	defer os.RemoveAll(tempDir)

	dbPath := filepath.Join(tempDir, "test.db")
	db, err := NewDB(dbPath)
	assert.NoError(t, err)
	defer db.Close()

	// Create a test export file with sample data
	exportFilePath := filepath.Join(tempDir, "test_export.json")
	file, err := os.Create(exportFilePath)
	assert.NoError(t, err)

	// Write test entries in JSON Lines format
	testEntries := []string{
		`{"id":"20220101120000","ts":"2022-01-01T12:00:00+09:00","cat":"research","body":"量子コンピューティングの基礎"}`,
		`{"id":"20220102130000","ts":"2022-01-02T13:00:00+09:00","cat":"programming","body":"Goのコンカレンシーパターン"}`,
		`{"id":"20220103140000","ts":"2022-01-03T14:00:00+09:00","cat":"research_and_programming","body":"機械学習モデルの実装 - TensorFlowによる実装"}`,
	}

	for _, entry := range testEntries {
		_, err = file.WriteString(entry + "\n")
		assert.NoError(t, err)
	}
	file.Close()

	// Test ImportEntries
	count, err := db.(*SQLiteDB).ImportEntries(exportFilePath)
	assert.NoError(t, err)
	assert.Equal(t, 3, count, "インポートしたエントリ数が一致しません")

	// Verify imported entries
	entries, err := db.GetAllEntries()
	assert.NoError(t, err)
	assert.Equal(t, 3, len(entries), "データベース内のエントリ数が一致しません")

	// Check specific entries
	// エントリはcreated_at DESCでソートされるので、逆順で取得される
	assert.Equal(t, "20220103140000", entries[0].ID)
	assert.Equal(t, models.ResearchAndProgram, entries[0].Category)
	assert.Equal(t, "機械学習モデルの実装", entries[0].ResearchTopic)
	assert.Equal(t, "TensorFlowによる実装", entries[0].ProgramTitle)
	assert.Equal(t, 3, entries[0].Satisfaction) // デフォルト値

	assert.Equal(t, "20220102130000", entries[1].ID)
	assert.Equal(t, models.Programming, entries[1].Category)
	assert.Equal(t, "Goのコンカレンシーパターン", entries[1].ProgramTitle)

	assert.Equal(t, "20220101120000", entries[2].ID)
	assert.Equal(t, models.Research, entries[2].Category)
	assert.Equal(t, "量子コンピューティングの基礎", entries[2].ResearchTopic)

	// Test duplicate handling by importing the same file again
	count, err = db.(*SQLiteDB).ImportEntries(exportFilePath)
	assert.NoError(t, err)
	assert.Equal(t, 0, count, "重複エントリがスキップされていません")

	// Check entry count remains the same
	entries, err = db.GetAllEntries()
	assert.NoError(t, err)
	assert.Equal(t, 3, len(entries), "重複インポート後のエントリ数が変わっています")
}

func TestImportEntriesErrorHandling(t *testing.T) {
	// Create a temporary database
	tempDir, err := os.MkdirTemp("", "wamon-import-error-test-*")
	assert.NoError(t, err)
	defer os.RemoveAll(tempDir)

	dbPath := filepath.Join(tempDir, "test.db")
	db, err := NewDB(dbPath)
	assert.NoError(t, err)
	defer db.Close()

	// Test with non-existent file
	_, err = db.(*SQLiteDB).ImportEntries("/path/that/does/not/exist.json")
	assert.Error(t, err, "存在しないファイルのインポートがエラーを返していません")
	assert.Contains(t, err.Error(), "no such file")

	// Test with invalid JSON
	invalidJSONPath := filepath.Join(tempDir, "invalid.json")
	file, err := os.Create(invalidJSONPath)
	assert.NoError(t, err)
	_, err = file.WriteString("This is not a valid JSON\n")
	assert.NoError(t, err)
	err = file.Sync() // 確実に書き込み
	assert.NoError(t, err)
	file.Close()

	_, err = db.(*SQLiteDB).ImportEntries(invalidJSONPath)
	assert.Error(t, err, "不正なJSONのインポートがエラーを返していません")
	assert.Contains(t, err.Error(), "JSON解析エラー")

	// Test with invalid category
	invalidCategoryPath := filepath.Join(tempDir, "invalid_category.json")
	file, err = os.Create(invalidCategoryPath)
	assert.NoError(t, err)
	_, err = file.WriteString(`{"id":"12345","ts":"2022-01-01T12:00:00+09:00","cat":"invalid_category","body":"test"}`)
	assert.NoError(t, err)
	err = file.Sync() // 確実に書き込み
	assert.NoError(t, err)
	file.Close()

	_, err = db.(*SQLiteDB).ImportEntries(invalidCategoryPath)
	assert.Error(t, err, "不正なカテゴリのインポートがエラーを返していません")
	assert.Contains(t, err.Error(), "不明なカテゴリ")

	// Test with missing required fields
	missingFieldsPath := filepath.Join(tempDir, "missing_fields.json")
	file, err = os.Create(missingFieldsPath)
	assert.NoError(t, err)
	_, err = file.WriteString(`{"id":"12345"}`) // Missing ts, cat, body
	assert.NoError(t, err)
	err = file.Sync() // 確実に書き込み
	assert.NoError(t, err)
	file.Close()

	_, err = db.(*SQLiteDB).ImportEntries(missingFieldsPath)
	assert.Error(t, err, "必須フィールドが欠けているJSONのインポートがエラーを返していません")
	assert.Contains(t, err.Error(), "不正な日時形式")
}

func TestImportWithMixedValidAndInvalidEntries(t *testing.T) {
	// Create a temporary database
	tempDir, err := os.MkdirTemp("", "wamon-import-mixed-test-*")
	assert.NoError(t, err)
	defer os.RemoveAll(tempDir)

	dbPath := filepath.Join(tempDir, "test.db")
	db, err := NewDB(dbPath)
	assert.NoError(t, err)
	defer db.Close()

	// Create a test file with valid and invalid entries
	mixedFilePath := filepath.Join(tempDir, "mixed.json")
	file, err := os.Create(mixedFilePath)
	assert.NoError(t, err)

	// First entry is valid
	_, err = file.WriteString(`{"id":"20220101120000","ts":"2022-01-01T12:00:00+09:00","cat":"research","body":"Test Research"}`)
	assert.NoError(t, err)
	_, err = file.WriteString("\n") // 改行
	assert.NoError(t, err)
	// Second entry is invalid (missing body)
	_, err = file.WriteString(`{"id":"20220102130000","ts":"2022-01-02T13:00:00+09:00","cat":"programming"}`)
	assert.NoError(t, err)
	err = file.Sync() // 確実に書き込み
	assert.NoError(t, err)
	file.Close()

	// Import should fail at the second entry
	_, err = db.(*SQLiteDB).ImportEntries(mixedFilePath)
	assert.Error(t, err, "不正なエントリを含むファイルのインポートがエラーを返していません")
	assert.Contains(t, err.Error(), "不正な本文形式")

	// Check that no entries were imported (transaction should be rolled back)
	count, err := db.GetEntryCount()
	assert.NoError(t, err)
	assert.Equal(t, 0, count, "エラー時にもエントリがインポートされています")
}
