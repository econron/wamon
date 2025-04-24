package models

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewEntry(t *testing.T) {
	// Test with Research category
	entry := NewEntry(Research, 5)
	assert.NotEmpty(t, entry.ID)
	assert.Equal(t, Research, entry.Category)
	assert.Equal(t, 5, entry.Satisfaction)
	assert.WithinDuration(t, time.Now(), entry.CreatedAt, 2*time.Second)

	// Test with Programming category
	entry = NewEntry(Programming, 3)
	assert.NotEmpty(t, entry.ID)
	assert.Equal(t, Programming, entry.Category)
	assert.Equal(t, 3, entry.Satisfaction)
	assert.WithinDuration(t, time.Now(), entry.CreatedAt, 2*time.Second)

	// Test with ResearchAndProgram category
	entry = NewEntry(ResearchAndProgram, 1)
	assert.NotEmpty(t, entry.ID)
	assert.Equal(t, ResearchAndProgram, entry.Category)
	assert.Equal(t, 1, entry.Satisfaction)
	assert.WithinDuration(t, time.Now(), entry.CreatedAt, 2*time.Second)
}
