package db

import (
	"database/sql"
	"os"
	"testing"
	"time"

	"github.com/econron/wamon/internal/models"
	"github.com/stretchr/testify/assert"
)

// TestMain manages the test setup and teardown for all database tests
func TestMain(m *testing.M) {
	// テスト実行
	code := m.Run()
	os.Exit(code)
}

// setupTestDB はテスト用のインメモリデータベースを設定します
func setupTestDB(t *testing.T) DB {
	// テスト用のインメモリデータベースを設定（各テストごとに新しいインスタンス）
	db, err := NewDB(":memory:")
	assert.NoError(t, err)

	// テスト終了時にクリーンアップ
	t.Cleanup(func() {
		db.Close()
	})

	return db
}

// createTestEntry はテスト用のエントリを作成します
func createTestEntry() *models.Entry {
	return &models.Entry{
		ID:            "20220101120000",
		Category:      models.Research,
		ResearchTopic: "テストトピック",
		ProgramTitle:  "",
		Satisfaction:  5,
		CreatedAt:     time.Now(),
	}
}

func TestNewDB(t *testing.T) {
	// テスト用の一時ファイルパス
	tempDBPath := "test_db.sqlite"
	defer os.Remove(tempDBPath) // テスト終了後にファイルを削除

	// 初期化をテスト
	database, err := NewDB(tempDBPath)
	assert.NoError(t, err)
	assert.NotNil(t, database)
	database.Close()

	// DBファイルが作成されたことを確認
	_, err = os.Stat(tempDBPath)
	assert.NoError(t, err)

	// メモリDBでのテスト
	memoryDB, err := NewDB(":memory:")
	assert.NoError(t, err)
	assert.NotNil(t, memoryDB)
	memoryDB.Close()
}

func TestNewDBFailure(t *testing.T) {
	// 書き込み権限のない場所にDBを作成しようとする
	if os.Getuid() == 0 {
		// rootユーザーの場合はスキップ
		t.Skip("Skipping test when running as root")
	}

	// '/root'などの特権ディレクトリにアクセスを試みる
	invalidPath := "/root/nonexistent_dir/test.db"
	db, err := NewDB(invalidPath)
	assert.Error(t, err)
	assert.Nil(t, db)
}

func TestGetDB(t *testing.T) {
	// テスト用の一時ファイルパス
	tempDBPath := "test_get_db.sqlite"
	defer os.Remove(tempDBPath) // テスト終了後にファイルを削除

	// シングルトンインスタンスの取得をテスト
	db1, err := GetDB(tempDBPath)
	assert.NoError(t, err)
	assert.NotNil(t, db1)

	// 2回目の呼び出しで同じインスタンスが返されることを確認
	db2, err := GetDB(tempDBPath)
	assert.NoError(t, err)
	assert.Equal(t, db1, db2)
}

func TestGetDBFailure(t *testing.T) {
	// 無効なパスでGetDBを呼び出す
	if os.Getuid() == 0 {
		t.Skip("Skipping test when running as root")
	}

	invalidPath := "/root/nonexistent_dir/test.db"
	db, err := GetDB(invalidPath)
	assert.Error(t, err)
	assert.Nil(t, db)
}

func TestCreateDirIfNotExists(t *testing.T) {
	// テスト用の一時ディレクトリ
	tempDir := "test_dir"
	defer os.RemoveAll(tempDir) // テスト終了後にディレクトリを削除

	// ディレクトリが存在しない場合
	err := createDirIfNotExists(tempDir)
	assert.NoError(t, err)

	// ディレクトリが作成されたことを確認
	_, err = os.Stat(tempDir)
	assert.NoError(t, err)

	// すでに存在する場合
	err = createDirIfNotExists(tempDir)
	assert.NoError(t, err)
}

func TestCreateDirIfNotExistsFailure(t *testing.T) {
	// 書き込み権限のない場所にディレクトリを作成しようとする
	if os.Getuid() == 0 {
		t.Skip("Skipping test when running as root")
	}

	err := createDirIfNotExists("/root/nonexistent_dir")
	assert.Error(t, err)
}

func TestSaveEntry(t *testing.T) {
	db := setupTestDB(t)

	// テスト用のエントリ作成
	entry := createTestEntry()

	// エントリを保存
	err := db.SaveEntry(entry)
	assert.NoError(t, err)

	// 保存したエントリを取得
	savedEntry, err := db.GetEntryByID(entry.ID)
	assert.NoError(t, err)
	assert.Equal(t, entry.ID, savedEntry.ID)
	assert.Equal(t, entry.Category, savedEntry.Category)
	assert.Equal(t, entry.ResearchTopic, savedEntry.ResearchTopic)
	assert.Equal(t, entry.Satisfaction, savedEntry.Satisfaction)
}

func TestSaveEntryFailure(t *testing.T) {
	db := setupTestDB(t)

	// NULLにできない必須フィールドをNULLにしてエラーを発生させる
	invalidEntry := &models.Entry{
		ID:        "",          // 空のID（主キー）
		Category:  "",          // 空のカテゴリ（NOT NULL制約）
		CreatedAt: time.Time{}, // ゼロ時間
	}

	err := db.SaveEntry(invalidEntry)
	assert.Error(t, err)
}

func TestUpdateEntry(t *testing.T) {
	db := setupTestDB(t)

	// テスト用のエントリ作成と保存
	entry := createTestEntry()
	err := db.SaveEntry(entry)
	assert.NoError(t, err)

	// エントリを更新
	entry.ResearchTopic = "更新されたトピック"
	entry.Satisfaction = 4
	err = db.UpdateEntry(entry)
	assert.NoError(t, err)

	// 更新したエントリを取得して確認
	updatedEntry, err := db.GetEntryByID(entry.ID)
	assert.NoError(t, err)
	assert.Equal(t, "更新されたトピック", updatedEntry.ResearchTopic)
	assert.Equal(t, 4, updatedEntry.Satisfaction)
}

func TestUpdateEntryFailure(t *testing.T) {
	db := setupTestDB(t)

	// 存在しないIDでの更新
	nonExistentEntry := createTestEntry()
	nonExistentEntry.ID = "nonexistent"
	err := db.UpdateEntry(nonExistentEntry)
	assert.Error(t, err)
	assert.Equal(t, sql.ErrNoRows, err)

	// 無効なデータで更新を試みる
	entry := createTestEntry()
	err = db.SaveEntry(entry)
	assert.NoError(t, err)

	// カテゴリをNULLにしてみる（NOT NULL制約に違反）
	entry.Category = ""
	err = db.UpdateEntry(entry)
	assert.Error(t, err)
}

func TestGetAllEntries(t *testing.T) {
	db := setupTestDB(t)

	// テスト用のエントリを複数作成
	entry1 := createTestEntry()
	entry2 := &models.Entry{
		ID:           "20220102120000",
		Category:     models.Programming,
		ProgramTitle: "テストプログラム",
		Satisfaction: 4,
		CreatedAt:    time.Now().Add(-24 * time.Hour), // 1日前
	}

	// エントリを保存
	err := db.SaveEntry(entry1)
	assert.NoError(t, err)
	err = db.SaveEntry(entry2)
	assert.NoError(t, err)

	// 全エントリを取得
	entries, err := db.GetAllEntries()
	assert.NoError(t, err)
	assert.Len(t, entries, 2)
	// 新しい順に並んでいることを確認（entry1が先）
	assert.Equal(t, entry1.ID, entries[0].ID)
}

func TestGetAllEntriesEmpty(t *testing.T) {
	db := setupTestDB(t)

	// エントリが1つもない状態でGetAllEntriesを呼び出す
	entries, err := db.GetAllEntries()
	assert.NoError(t, err)
	assert.Empty(t, entries)
}

func TestGetEntriesByCategory(t *testing.T) {
	db := setupTestDB(t)

	// 異なるカテゴリのエントリを作成
	researchEntry := createTestEntry() // Research カテゴリ
	programEntry := &models.Entry{
		ID:           "20220102120000",
		Category:     models.Programming,
		ProgramTitle: "テストプログラム",
		Satisfaction: 4,
		CreatedAt:    time.Now(),
	}

	// エントリを保存
	err := db.SaveEntry(researchEntry)
	assert.NoError(t, err)
	err = db.SaveEntry(programEntry)
	assert.NoError(t, err)

	// Research カテゴリのエントリを取得
	researchEntries, err := db.GetEntriesByCategory(models.Research)
	assert.NoError(t, err)
	assert.Len(t, researchEntries, 1)
	assert.Equal(t, models.Research, researchEntries[0].Category)

	// Programming カテゴリのエントリを取得
	programEntries, err := db.GetEntriesByCategory(models.Programming)
	assert.NoError(t, err)
	assert.Len(t, programEntries, 1)
	assert.Equal(t, models.Programming, programEntries[0].Category)
}

func TestGetEntriesByCategoryEmpty(t *testing.T) {
	db := setupTestDB(t)

	// 存在しないカテゴリでフィルタリング
	entries, err := db.GetEntriesByCategory(models.Category("unknown"))
	assert.NoError(t, err)
	assert.Empty(t, entries)
}

func TestGetEntryByID(t *testing.T) {
	db := setupTestDB(t)

	// テスト用のエントリ作成と保存
	entry := createTestEntry()
	err := db.SaveEntry(entry)
	assert.NoError(t, err)

	// IDでエントリを取得
	foundEntry, err := db.GetEntryByID(entry.ID)
	assert.NoError(t, err)
	assert.Equal(t, entry.ID, foundEntry.ID)
	assert.Equal(t, entry.Category, foundEntry.Category)
}

func TestGetEntryByIDNotFound(t *testing.T) {
	db := setupTestDB(t)

	// 存在しないIDの場合
	_, err := db.GetEntryByID("nonexistent")
	assert.Error(t, err)
	assert.Equal(t, sql.ErrNoRows, err)
}

func TestGetEntryCount(t *testing.T) {
	db := setupTestDB(t)

	// 最初はエントリなし
	count, err := db.GetEntryCount()
	assert.NoError(t, err)
	assert.Equal(t, 0, count)

	// エントリを追加
	err = db.SaveEntry(createTestEntry())
	assert.NoError(t, err)

	// カウントを再取得
	count, err = db.GetEntryCount()
	assert.NoError(t, err)
	assert.Equal(t, 1, count)
}

func TestGetEntriesFromLastWeek(t *testing.T) {
	db := setupTestDB(t)

	// 1週間以内のエントリ
	recentEntry := createTestEntry()
	err := db.SaveEntry(recentEntry)
	assert.NoError(t, err)

	// 8日前のエントリ（1週間外）
	oldEntry := &models.Entry{
		ID:           "20220102120000",
		Category:     models.Programming,
		ProgramTitle: "古いエントリ",
		Satisfaction: 3,
		CreatedAt:    time.Now().AddDate(0, 0, -8),
	}
	err = db.SaveEntry(oldEntry)
	assert.NoError(t, err)

	// 1週間以内のエントリのみを取得
	entries, err := db.GetEntriesFromLastWeek()
	assert.NoError(t, err)
	assert.Len(t, entries, 1)
	assert.Equal(t, recentEntry.ID, entries[0].ID)
}

func TestGetEntriesFromLastWeekEmpty(t *testing.T) {
	db := setupTestDB(t)

	// 1週間以内のエントリがない場合
	entries, err := db.GetEntriesFromLastWeek()
	assert.NoError(t, err)
	assert.Empty(t, entries)
}

func TestClose(t *testing.T) {
	// DBを作成して閉じる
	db, err := NewDB(":memory:")
	assert.NoError(t, err)

	err = db.Close()
	assert.NoError(t, err)

	// 閉じた後のDBで操作を試みる（別のテストとして追加）
	sqlDB := db.(*SQLiteDB)
	_, err = sqlDB.db.Exec("SELECT 1")
	assert.Error(t, err)
}

// 後方互換性テスト
func TestBackwardCompatibility(t *testing.T) {
	// テスト用の一時ファイルパス
	tempDBPath := "test_compat.sqlite"
	defer os.Remove(tempDBPath) // テスト終了後にファイルを削除

	// 初期化をテスト
	sqlDB, err := InitDB(tempDBPath)
	assert.NoError(t, err)
	assert.NotNil(t, sqlDB)
	defer sqlDB.Close()
}
