package db

import (
	"database/sql"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/econron/wamon/internal/models"

	_ "github.com/mattn/go-sqlite3"
)

var (
	db   *sql.DB
	once sync.Once
)

// InitDB initializes the database and creates tables if they don't exist
func InitDB(dbPath string) (*sql.DB, error) {
	var err error
	once.Do(func() {
		// Ensure directory exists
		if dbPath != ":memory:" {
			dirPath := filepath.Dir(dbPath)
			if dirPath != "." && dirPath != "" {
				// Ensure the directory exists
				err = createDirIfNotExists(dirPath)
				if err != nil {
					return
				}
			}
		}

		// Open the database
		db, err = sql.Open("sqlite3", dbPath)
		if err != nil {
			return
		}

		// Create tables if they don't exist
		err = createTables()
	})

	return db, err
}

// createDirIfNotExists creates a directory if it doesn't exist
func createDirIfNotExists(dirPath string) error {
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		return os.MkdirAll(dirPath, 0755)
	}
	return nil
}

// createTables creates the necessary tables
func createTables() error {
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

// SaveEntry saves an entry to the database
func SaveEntry(entry *models.Entry) error {
	_, err := db.Exec(
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
func UpdateEntry(entry *models.Entry) error {
	result, err := db.Exec(
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
func GetAllEntries() ([]*models.Entry, error) {
	rows, err := db.Query(`
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
func GetEntriesByCategory(category models.Category) ([]*models.Entry, error) {
	rows, err := db.Query(`
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
func GetEntryByID(id string) (*models.Entry, error) {
	entry := &models.Entry{}
	var category string
	err := db.QueryRow(`
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
func GetEntryCount() (int, error) {
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM entries").Scan(&count)
	return count, err
}

// GetEntriesFromLastWeek retrieves entries from the past 7 days
func GetEntriesFromLastWeek() ([]*models.Entry, error) {
	// Calculate the timestamp for 7 days ago
	oneWeekAgo := time.Now().AddDate(0, 0, -7)

	rows, err := db.Query(`
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
