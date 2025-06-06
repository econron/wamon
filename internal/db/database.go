package db

import (
	"bufio"
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/econron/wamon/internal/models"

	_ "github.com/mattn/go-sqlite3"
)

// DB is the interface for database operations
type DB interface {
	SaveEntry(entry *models.Entry) error
	UpdateEntry(entry *models.Entry) error
	GetAllEntries() ([]*models.Entry, error)
	GetEntriesByCategory(category models.Category) ([]*models.Entry, error)
	GetEntryByID(id string) (*models.Entry, error)
	GetEntryCount() (int, error)
	GetEntriesFromLastWeek() ([]*models.Entry, error)
	GetEntriesSince(since time.Time) ([]*models.Entry, error)
	ExportEntries(filePath string) error
	ExportEntriesSince(filePath string, since time.Time) error
	ImportEntries(filePath string) (int, error)
	Close() error
}

// SQLiteDB implements the DB interface with SQLite
type SQLiteDB struct {
	db *sql.DB
}

var (
	instance *SQLiteDB
	once     sync.Once
)

// NewDB creates a new DB instance with dependency injection
func NewDB(dbPath string) (DB, error) {
	var err error
	db, err := initDB(dbPath)
	if err != nil {
		return nil, err
	}
	return &SQLiteDB{db: db}, nil
}

// GetDB returns a global singleton instance of the database
// This is maintained for backward compatibility but should be avoided
// in favor of dependency injection with NewDB
func GetDB(dbPath string) (DB, error) {
	var err error
	once.Do(func() {
		var db *sql.DB
		db, err = initDB(dbPath)
		if err != nil {
			return
		}
		instance = &SQLiteDB{db: db}
	})
	if err != nil {
		return nil, err
	}
	return instance, nil
}

// initDB initializes the database and creates tables if they don't exist
func initDB(dbPath string) (*sql.DB, error) {
	// Ensure directory exists
	if dbPath != ":memory:" {
		dirPath := filepath.Dir(dbPath)
		if dirPath != "." && dirPath != "" {
			// Ensure the directory exists
			err := createDirIfNotExists(dirPath)
			if err != nil {
				return nil, err
			}
		}
	}

	// Open the database
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}

	// Create tables if they don't exist
	err = createTables(db)
	if err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}

// createDirIfNotExists creates a directory if it doesn't exist
func createDirIfNotExists(dirPath string) error {
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		return os.MkdirAll(dirPath, 0755)
	}
	return nil
}

// createTables creates the necessary tables
func createTables(db *sql.DB) error {
	// Create entries table
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS entries (
			id TEXT PRIMARY KEY,
			category TEXT NOT NULL,
			research_topic TEXT,
			program_title TEXT,
			satisfaction INTEGER NOT NULL,
			created_at TIMESTAMP NOT NULL
		)
	`)
	return err
}

// Close closes the database connection
func (s *SQLiteDB) Close() error {
	if s.db != nil {
		return s.db.Close()
	}
	return nil
}

// SaveEntry saves an entry to the database
func (s *SQLiteDB) SaveEntry(entry *models.Entry) error {
	_, err := s.db.Exec(
		`INSERT INTO entries (id, category, research_topic, program_title, satisfaction, created_at)
		 VALUES (?, ?, ?, ?, ?, ?)`,
		entry.ID,
		entry.Category,
		entry.ResearchTopic,
		entry.ProgramTitle,
		entry.Satisfaction,
		entry.CreatedAt,
	)
	return err
}

// UpdateEntry updates an existing entry in the database
func (s *SQLiteDB) UpdateEntry(entry *models.Entry) error {
	result, err := s.db.Exec(
		`UPDATE entries 
		 SET category = ?, research_topic = ?, program_title = ?, satisfaction = ?
		 WHERE id = ?`,
		entry.Category,
		entry.ResearchTopic,
		entry.ProgramTitle,
		entry.Satisfaction,
		entry.ID,
	)
	if err != nil {
		return err
	}

	// Check if the entry was actually updated
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

// GetAllEntries retrieves all entries from the database
func (s *SQLiteDB) GetAllEntries() ([]*models.Entry, error) {
	rows, err := s.db.Query(`
		SELECT id, category, research_topic, program_title, satisfaction, created_at
		FROM entries
		ORDER BY created_at DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var entries []*models.Entry
	for rows.Next() {
		entry := &models.Entry{}
		var category string
		err := rows.Scan(
			&entry.ID,
			&category,
			&entry.ResearchTopic,
			&entry.ProgramTitle,
			&entry.Satisfaction,
			&entry.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		entry.Category = models.Category(category)
		entries = append(entries, entry)
	}

	return entries, nil
}

// GetEntriesByCategory retrieves entries by category
func (s *SQLiteDB) GetEntriesByCategory(category models.Category) ([]*models.Entry, error) {
	rows, err := s.db.Query(`
		SELECT id, category, research_topic, program_title, satisfaction, created_at
		FROM entries
		WHERE category = ?
		ORDER BY created_at DESC
	`, category)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var entries []*models.Entry
	for rows.Next() {
		entry := &models.Entry{}
		var category string
		err := rows.Scan(
			&entry.ID,
			&category,
			&entry.ResearchTopic,
			&entry.ProgramTitle,
			&entry.Satisfaction,
			&entry.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		entry.Category = models.Category(category)
		entries = append(entries, entry)
	}

	return entries, nil
}

// GetEntryByID retrieves an entry by ID
func (s *SQLiteDB) GetEntryByID(id string) (*models.Entry, error) {
	entry := &models.Entry{}
	var category string
	err := s.db.QueryRow(`
		SELECT id, category, research_topic, program_title, satisfaction, created_at
		FROM entries
		WHERE id = ?
	`, id).Scan(
		&entry.ID,
		&category,
		&entry.ResearchTopic,
		&entry.ProgramTitle,
		&entry.Satisfaction,
		&entry.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	entry.Category = models.Category(category)
	return entry, nil
}

// GetEntryCount returns the total number of entries
func (s *SQLiteDB) GetEntryCount() (int, error) {
	var count int
	err := s.db.QueryRow("SELECT COUNT(*) FROM entries").Scan(&count)
	return count, err
}

// GetEntriesFromLastWeek retrieves entries from the past 7 days
func (s *SQLiteDB) GetEntriesFromLastWeek() ([]*models.Entry, error) {
	// Calculate the timestamp for 7 days ago
	oneWeekAgo := time.Now().AddDate(0, 0, -7)

	rows, err := s.db.Query(`
		SELECT id, category, research_topic, program_title, satisfaction, created_at
		FROM entries
		WHERE created_at >= ?
		ORDER BY created_at DESC
	`, oneWeekAgo)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var entries []*models.Entry
	for rows.Next() {
		entry := &models.Entry{}
		var category string
		err := rows.Scan(
			&entry.ID,
			&category,
			&entry.ResearchTopic,
			&entry.ProgramTitle,
			&entry.Satisfaction,
			&entry.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		entry.Category = models.Category(category)
		entries = append(entries, entry)
	}

	return entries, nil
}

// GetEntriesSince retrieves entries created after the specified time
func (s *SQLiteDB) GetEntriesSince(since time.Time) ([]*models.Entry, error) {
	rows, err := s.db.Query(`
		SELECT id, category, research_topic, program_title, satisfaction, created_at
		FROM entries
		WHERE created_at >= ?
		ORDER BY created_at DESC
	`, since)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var entries []*models.Entry
	for rows.Next() {
		entry := &models.Entry{}
		var category string
		err := rows.Scan(
			&entry.ID,
			&category,
			&entry.ResearchTopic,
			&entry.ProgramTitle,
			&entry.Satisfaction,
			&entry.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		entry.Category = models.Category(category)
		entries = append(entries, entry)
	}

	return entries, nil
}

// ExportEntries exports all entries from the database to a JSON file
// Each entry is written as a separate JSON object on its own line (JSON Lines format)
func (s *SQLiteDB) ExportEntries(filePath string) error {
	// Get all entries ordered by creation time (newest first)
	entries, err := s.GetAllEntries()
	if err != nil {
		return err
	}

	// If there are no entries, create an empty file
	if len(entries) == 0 {
		emptyFile, err := os.Create(filePath)
		if err != nil {
			return err
		}
		defer emptyFile.Close()
		return nil
	}

	// Create the directory if it doesn't exist
	dir := filepath.Dir(filePath)
	if dir != "." {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return err
		}
	}

	// Create or truncate the output file
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	// Write each entry as a JSON object on its own line
	for _, entry := range entries {
		// Use the original category string directly
		catStr := string(entry.Category)

		// Create simplified export format object
		exportEntry := map[string]interface{}{
			"id":  entry.ID,
			"ts":  entry.CreatedAt.Format(time.RFC3339),
			"cat": catStr,
		}

		// Set body based on category
		switch entry.Category {
		case models.Research:
			exportEntry["body"] = entry.ResearchTopic
		case models.Programming:
			exportEntry["body"] = entry.ProgramTitle
		case models.ResearchAndProgram:
			exportEntry["body"] = entry.ResearchTopic + " - " + entry.ProgramTitle
		}

		// Convert to JSON
		jsonData, err := json.Marshal(exportEntry)
		if err != nil {
			return err
		}

		// Write the JSON line
		_, err = file.Write(jsonData)
		if err != nil {
			return err
		}

		// Add newline
		_, err = file.Write([]byte("\n"))
		if err != nil {
			return err
		}
	}

	return nil
}

// ExportEntriesSince exports entries from the database created after the specified time to a JSON file
// Each entry is written as a separate JSON object on its own line (JSON Lines format)
func (s *SQLiteDB) ExportEntriesSince(filePath string, since time.Time) error {
	// Get all entries since the specified time, ordered by creation time (newest first)
	entries, err := s.GetEntriesSince(since)
	if err != nil {
		return err
	}

	// If there are no entries, create an empty file
	if len(entries) == 0 {
		emptyFile, err := os.Create(filePath)
		if err != nil {
			return err
		}
		defer emptyFile.Close()
		return nil
	}

	// Create the directory if it doesn't exist
	dir := filepath.Dir(filePath)
	if dir != "." {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return err
		}
	}

	// Create or truncate the output file
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	// Write each entry as a JSON object on its own line
	for _, entry := range entries {
		// Use the original category string directly
		catStr := string(entry.Category)

		// Create simplified export format object
		exportEntry := map[string]interface{}{
			"id":  entry.ID,
			"ts":  entry.CreatedAt.Format(time.RFC3339),
			"cat": catStr,
		}

		// Set body based on category
		switch entry.Category {
		case models.Research:
			exportEntry["body"] = entry.ResearchTopic
		case models.Programming:
			exportEntry["body"] = entry.ProgramTitle
		case models.ResearchAndProgram:
			exportEntry["body"] = entry.ResearchTopic + " - " + entry.ProgramTitle
		}

		// Convert to JSON
		jsonData, err := json.Marshal(exportEntry)
		if err != nil {
			return err
		}

		// Write the JSON line
		_, err = file.Write(jsonData)
		if err != nil {
			return err
		}

		// Add newline
		_, err = file.Write([]byte("\n"))
		if err != nil {
			return err
		}
	}

	return nil
}

// ImportEntries imports entries from a JSON file into the database
func (s *SQLiteDB) ImportEntries(filePath string) (int, error) {
	// Open the file
	file, err := os.Open(filePath)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	// Begin transaction
	tx, err := s.db.Begin()
	if err != nil {
		return 0, fmt.Errorf("トランザクション開始エラー: %v", err)
	}

	// Read the file line by line
	scanner := bufio.NewScanner(file)
	importedCount := 0

	// Prepare function for rollback in case of error
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}

		// Parse the JSON line
		var data map[string]interface{}
		if err := json.Unmarshal([]byte(line), &data); err != nil {
			return importedCount, fmt.Errorf("JSON解析エラー: %v", err)
		}

		// Extract fields
		id, ok := data["id"].(string)
		if !ok {
			return importedCount, fmt.Errorf("不正なID形式: %v", data["id"])
		}

		// Check if entry with this ID already exists
		var count int
		err := tx.QueryRow("SELECT COUNT(*) FROM entries WHERE id = ?", id).Scan(&count)
		if err != nil {
			return importedCount, fmt.Errorf("データベースクエリエラー: %v", err)
		}

		if count > 0 {
			// Entry already exists, skip
			continue
		}

		// Parse timestamp
		tsStr, ok := data["ts"].(string)
		if !ok {
			return importedCount, fmt.Errorf("不正な日時形式: %v", data["ts"])
		}
		createdAt, err := time.Parse(time.RFC3339, tsStr)
		if err != nil {
			return importedCount, fmt.Errorf("日時の解析エラー: %v", err)
		}

		// Parse category
		catStr, ok := data["cat"].(string)
		if !ok {
			return importedCount, fmt.Errorf("不正なカテゴリ形式: %v", data["cat"])
		}

		// Map the export category names back to internal model categories
		var category models.Category
		switch catStr {
		case "research", "調べ物":
			category = models.Research
		case "programming", "プログラマ":
			category = models.Programming
		case "research_and_programming", "調べてプログラマ":
			category = models.ResearchAndProgram
		default:
			return importedCount, fmt.Errorf("不明なカテゴリ: %v", catStr)
		}

		// Parse body
		body, ok := data["body"].(string)
		if !ok {
			return importedCount, fmt.Errorf("不正な本文形式: %v", data["body"])
		}

		// Create entry
		entry := &models.Entry{
			ID:           id,
			Category:     category,
			CreatedAt:    createdAt,
			Satisfaction: 3, // Default satisfaction if not specified
		}

		// Set content based on category
		switch category {
		case models.Research:
			entry.ResearchTopic = body
		case models.Programming:
			entry.ProgramTitle = body
		case models.ResearchAndProgram:
			// Try to split combined body
			parts := strings.Split(body, " - ")
			if len(parts) == 2 {
				entry.ResearchTopic = parts[0]
				entry.ProgramTitle = parts[1]
			} else {
				// Fallback to using the entire body as research topic
				entry.ResearchTopic = body
			}
		}

		// Save the entry within transaction
		_, err = tx.Exec(
			`INSERT INTO entries (id, category, research_topic, program_title, satisfaction, created_at)
			VALUES (?, ?, ?, ?, ?, ?)`,
			entry.ID,
			entry.Category,
			entry.ResearchTopic,
			entry.ProgramTitle,
			entry.Satisfaction,
			entry.CreatedAt,
		)
		if err != nil {
			return importedCount, fmt.Errorf("エントリの保存エラー: %v", err)
		}

		importedCount++
	}

	if err := scanner.Err(); err != nil {
		return importedCount, fmt.Errorf("ファイル読み込みエラー: %v", err)
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return importedCount, fmt.Errorf("トランザクションコミットエラー: %v", err)
	}

	return importedCount, nil
}

// Backward compatibility functions that use the global db instance
// These should be avoided in new code in favor of using the DB interface

// InitDB initializes the database and creates tables if they don't exist
// Deprecated: Use NewDB instead
func InitDB(dbPath string) (*sql.DB, error) {
	return initDB(dbPath)
}
