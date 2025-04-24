package cmd

import (
	"bufio"
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/econron/wamon/internal/db"
	"github.com/econron/wamon/internal/models"
	"github.com/stretchr/testify/assert"
)

// outputCaptor captures stdout output
type outputCaptor struct {
	output string
	wg     sync.WaitGroup
}

// captureOutput reads from the provided reader until EOF
func (c *outputCaptor) captureOutput(r io.Reader) {
	c.wg.Add(1)
	scanner := bufio.NewScanner(r)
	var captured string
	for scanner.Scan() {
		captured += scanner.Text() + "\n"
	}
	c.output = captured
	c.wg.Done()
}

// wait waits for the capture goroutine to complete
func (c *outputCaptor) wait() {
	c.wg.Wait()
}

type ExportedEntry struct {
	ID   string `json:"id"`
	TS   string `json:"ts"`
	Cat  string `json:"cat"`
	Body string `json:"body"`
}

func TestExportCmd(t *testing.T) {
	// テストディレクトリを作成
	tempDir, err := os.MkdirTemp("", "wamon-export-test-*")
	assert.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// テスト用DBファイルのパスを設定
	testDBPath := filepath.Join(tempDir, "test.db")
	originalDBPath := dbPath
	dbPath = testDBPath
	defer func() {
		dbPath = originalDBPath
	}()

	// テスト用のデータベースを作成
	database, err := db.NewDB(testDBPath)
	assert.NoError(t, err)
	defer database.Close()

	// テスト用のエントリを作成
	testEntries := createExportTestEntries(t)

	// エントリをデータベースに保存
	for _, entry := range testEntries {
		err := database.SaveEntry(entry)
		assert.NoError(t, err)
	}

	// テスト用のエクスポートファイルのパスを設定
	testExportPath := filepath.Join(tempDir, "export.json")

	// エクスポートコマンドを実行
	exportCmd.Run(exportCmd, []string{testExportPath})

	// エクスポートされたファイルが存在することを確認
	_, err = os.Stat(testExportPath)
	assert.NoError(t, err)

	// エクスポートされたデータを読み込み
	data, err := os.ReadFile(testExportPath)
	assert.NoError(t, err)

	// 各行が有効なJSONであることを確認
	lines := strings.Split(strings.TrimSpace(string(data)), "\n")
	assert.Equal(t, len(testEntries), len(lines))

	// 各行のJSONをパースして内容を検証
	for i, line := range lines {
		var exportedEntry ExportedEntry
		err := json.Unmarshal([]byte(line), &exportedEntry)
		assert.NoError(t, err)

		// 新しいエントリが先に来るので、逆順で比較
		entry := testEntries[len(testEntries)-i-1]

		// IDが一致することを確認
		assert.Equal(t, entry.ID, exportedEntry.ID)

		// タイムスタンプがISO8601形式であることを確認
		_, err = time.Parse(time.RFC3339, exportedEntry.TS)
		assert.NoError(t, err)

		// カテゴリが一致することを確認
		assert.Equal(t, string(entry.Category), exportedEntry.Cat)

		// 内容がエントリのカテゴリに応じて正しく設定されていることを確認
		switch entry.Category {
		case models.Research:
			assert.Equal(t, entry.ResearchTopic, exportedEntry.Body)
		case models.Programming:
			assert.Equal(t, entry.ProgramTitle, exportedEntry.Body)
		case models.ResearchAndProgram:
			expected := entry.ResearchTopic + " - " + entry.ProgramTitle
			assert.Equal(t, expected, exportedEntry.Body)
		}
	}
}

// TestExportCmdDefaultFilename tests the export command with the default filename
func TestExportCmdDefaultFilename(t *testing.T) {
	// 現在のディレクトリを記録
	currentDir, err := os.Getwd()
	assert.NoError(t, err)

	// テストディレクトリを作成
	tempDir, err := os.MkdirTemp("", "wamon-export-default-*")
	assert.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// テストディレクトリに移動
	err = os.Chdir(tempDir)
	assert.NoError(t, err)
	defer os.Chdir(currentDir)

	// テスト用DBファイルのパスを設定
	testDBPath := filepath.Join(tempDir, "test.db")
	originalDBPath := dbPath
	dbPath = testDBPath
	defer func() {
		dbPath = originalDBPath
	}()

	// テスト用のデータベースを作成
	database, err := db.NewDB(testDBPath)
	assert.NoError(t, err)
	defer database.Close()

	// テスト用のエントリを作成
	testEntries := createExportTestEntries(t)

	// エントリをデータベースに保存
	for _, entry := range testEntries {
		err := database.SaveEntry(entry)
		assert.NoError(t, err)
	}

	// エクスポートコマンドを引数なしで実行
	exportCmd.Run(exportCmd, []string{})

	// デフォルトのエクスポートファイルが存在することを確認
	defaultPath := "wamon_export.json"
	_, err = os.Stat(defaultPath)
	assert.NoError(t, err)

	// エクスポートされたデータを読み込み
	data, err := os.ReadFile(defaultPath)
	assert.NoError(t, err)

	// 各行が有効なJSONであることを確認
	lines := strings.Split(strings.TrimSpace(string(data)), "\n")
	assert.Equal(t, len(testEntries), len(lines))
}

func TestExportEmptyDatabase(t *testing.T) {
	// テストディレクトリを作成
	tempDir, err := os.MkdirTemp("", "wamon-export-empty-*")
	assert.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// テスト用DBファイルのパスを設定
	testDBPath := filepath.Join(tempDir, "empty.db")
	originalDBPath := dbPath
	dbPath = testDBPath
	defer func() {
		dbPath = originalDBPath
	}()

	// テスト用のデータベースを作成（エントリなし）
	database, err := db.NewDB(testDBPath)
	assert.NoError(t, err)
	defer database.Close()

	// テスト用のエクスポートファイルのパスを設定
	testExportPath := filepath.Join(tempDir, "empty_export.json")

	// エクスポートコマンドを実行
	exportCmd.Run(exportCmd, []string{testExportPath})

	// エクスポートされたファイルが存在することを確認
	_, err = os.Stat(testExportPath)
	assert.NoError(t, err)

	// エクスポートされたデータを読み込み
	data, err := os.ReadFile(testExportPath)
	assert.NoError(t, err)

	// 空のデータベースなので、ファイルも空であることを確認
	assert.Empty(t, strings.TrimSpace(string(data)))
}

// TestExportCmdWithSince tests the export command with the --since flag
func TestExportCmdWithSince(t *testing.T) {
	// テストディレクトリを作成
	tempDir, err := os.MkdirTemp("", "wamon-export-since-*")
	assert.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// テスト用DBファイルのパスを設定
	testDBPath := filepath.Join(tempDir, "test.db")
	originalDBPath := dbPath
	dbPath = testDBPath
	defer func() {
		dbPath = originalDBPath
	}()

	// テスト用のデータベースを作成
	database, err := db.NewDB(testDBPath)
	assert.NoError(t, err)
	defer database.Close()

	// 現在時刻を基準にしたエントリを作成
	now := time.Now()
	testEntries := []*models.Entry{
		{
			ID:            "entry1",
			Category:      models.Research,
			ResearchTopic: "Old entry",
			ProgramTitle:  "",
			Satisfaction:  4,
			CreatedAt:     now.Add(-48 * time.Hour), // 48時間前
		},
		{
			ID:            "entry2",
			Category:      models.Programming,
			ResearchTopic: "",
			ProgramTitle:  "Recent entry",
			Satisfaction:  5,
			CreatedAt:     now.Add(-12 * time.Hour), // 12時間前
		},
		{
			ID:            "entry3",
			Category:      models.ResearchAndProgram,
			ResearchTopic: "Very recent entry",
			ProgramTitle:  "Testing",
			Satisfaction:  3,
			CreatedAt:     now.Add(-1 * time.Hour), // 1時間前
		},
	}

	// エントリをデータベースに保存
	for _, entry := range testEntries {
		err := database.SaveEntry(entry)
		assert.NoError(t, err)
	}

	// テスト用のエクスポートファイルのパスを設定
	testExportPath := filepath.Join(tempDir, "export_since.json")

	// --since フラグを設定してエクスポートコマンドを実行 (24時間以内のエントリのみ)
	exportCmd.Flags().Set("since", "24h")
	exportCmd.Run(exportCmd, []string{testExportPath})
	exportCmd.Flags().Set("since", "") // 後のテストに影響しないようにリセット

	// エクスポートされたファイルが存在することを確認
	_, err = os.Stat(testExportPath)
	assert.NoError(t, err)

	// エクスポートされたデータを読み込み
	data, err := os.ReadFile(testExportPath)
	assert.NoError(t, err)

	// 各行が有効なJSONであることを確認
	lines := strings.Split(strings.TrimSpace(string(data)), "\n")

	// 過去24時間のエントリのみエクスポートされていることを確認 (entry2と3のみ)
	assert.Equal(t, 2, len(lines))

	// 各エントリを検証
	exportedIDs := make([]string, 0)
	for _, line := range lines {
		var exportedEntry ExportedEntry
		err := json.Unmarshal([]byte(line), &exportedEntry)
		assert.NoError(t, err)
		exportedIDs = append(exportedIDs, exportedEntry.ID)
	}

	// 期待されるIDが含まれていることを確認
	assert.Contains(t, exportedIDs, "entry2")
	assert.Contains(t, exportedIDs, "entry3")
	assert.NotContains(t, exportedIDs, "entry1") // 48時間前のエントリは含まれないはず
}

// createExportTestEntries テスト用のエントリを作成するヘルパー関数
func createExportTestEntries(t *testing.T) []*models.Entry {
	return []*models.Entry{
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
}
