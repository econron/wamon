package models

import (
	"time"
)

// Category represents the type of entry
type Category string

const (
	Research           Category = "調べ物"
	Programming        Category = "プログラマ"
	ResearchAndProgram Category = "調べてプログラマ"
)

// Entry represents a single journal entry
type Entry struct {
	ID            string    `json:"id"`
	Category      Category  `json:"category"`
	ResearchTopic string    `json:"research_topic,omitempty"`
	ProgramTitle  string    `json:"program_title,omitempty"`
	Satisfaction  int       `json:"satisfaction"` // 1-5 scale
	CreatedAt     time.Time `json:"created_at"`
}

// NewEntry creates a new entry with a unique ID and current timestamp
func NewEntry(category Category, satisfaction int) *Entry {
	return &Entry{
		ID:           generateID(),
		Category:     category,
		Satisfaction: satisfaction,
		CreatedAt:    time.Now(),
	}
}

// generateID creates a simple timestamp-based ID
func generateID() string {
	return time.Now().Format("20060102150405")
}
