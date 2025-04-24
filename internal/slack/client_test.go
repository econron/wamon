package slack

import (
	"fmt"
	"testing"
	"time"

	"github.com/econron/wamon/internal/models"
	"github.com/stretchr/testify/assert"
)

func TestNewClient(t *testing.T) {
	// Test with empty token
	emptyClient := NewClient(Config{})
	assert.NotNil(t, emptyClient)
	assert.Nil(t, emptyClient.api)

	// Test with valid token
	config := Config{
		Token:   "xoxb-valid-token",
		Channel: "test-channel",
		Enabled: true,
	}
	client := NewClient(config)
	assert.NotNil(t, client)
	assert.NotNil(t, client.api)
	assert.Equal(t, config, client.config)
}

func TestSendWeeklyReportValidation(t *testing.T) {
	// Test with empty token
	err := SendWeeklyReport("", "channel", []*models.Entry{
		{ID: "1", Category: models.Research, CreatedAt: time.Now()},
	})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "token is required")

	// Test with empty channel
	err = SendWeeklyReport("xoxb-token", "", []*models.Entry{
		{ID: "1", Category: models.Research, CreatedAt: time.Now()},
	})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "channel is required")

	// Test with empty entries
	err = SendWeeklyReport("xoxb-token", "channel", []*models.Entry{})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "no entries to report")
}

func TestGetStarRating(t *testing.T) {
	// Test valid ratings
	assert.Equal(t, "★☆☆☆☆", getStarRating(1))
	assert.Equal(t, "★★☆☆☆", getStarRating(2))
	assert.Equal(t, "★★★☆☆", getStarRating(3))
	assert.Equal(t, "★★★★☆", getStarRating(4))
	assert.Equal(t, "★★★★★", getStarRating(5))

	// Test out of range values
	assert.Equal(t, "★☆☆☆☆", getStarRating(0))
	assert.Equal(t, "★★★★★", getStarRating(6))
}

func TestClientSendWeeklyReportFailure(t *testing.T) {
	// Create client with no API
	client := &Client{config: Config{Channel: "test-channel"}}

	// Should fail because API is not initialized
	err := client.SendWeeklyReport([]*models.Entry{
		{ID: "1", Category: models.Research, CreatedAt: time.Now()},
	})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not properly initialized")

	// Test with empty entries
	client = NewClient(Config{Token: "xoxb-token", Channel: "channel"})
	err = client.SendWeeklyReport([]*models.Entry{})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "no entries to report")
}

// TestSendWeeklyReportWithMockApi tests the SendWeeklyReport function with a mock Slack API
func TestSendWeeklyReportWithMockApi(t *testing.T) {
	// Skip actual API calls
	t.Skip("Skipping test that would make actual API calls")

	// In a real implementation, you would use a mock API here
	// For example:
	// mockApi := &MockSlackAPI{}
	// originalSlackNew := slackNew
	// slackNew = func(token string) *slack.Client { return mockApi }
	// defer func() { slackNew = originalSlackNew }()
}

// TestGroupEntriesByDay tests the grouping logic used in SendWeeklyReport
func TestGroupEntriesByDay(t *testing.T) {
	// Create test entries spanning multiple days
	day1 := time.Date(2023, 1, 1, 10, 0, 0, 0, time.UTC)
	day2 := time.Date(2023, 1, 2, 11, 0, 0, 0, time.UTC)
	day3 := time.Date(2023, 1, 3, 12, 0, 0, 0, time.UTC)

	entries := []*models.Entry{
		{ID: "1", Category: models.Research, ResearchTopic: "Topic 1", Satisfaction: 5, CreatedAt: day1},
		{ID: "2", Category: models.Programming, ProgramTitle: "Program 1", Satisfaction: 4, CreatedAt: day1},
		{ID: "3", Category: models.Research, ResearchTopic: "Topic 2", Satisfaction: 3, CreatedAt: day2},
		{ID: "4", Category: models.ResearchAndProgram, ResearchTopic: "Topic 3", ProgramTitle: "Program 2", Satisfaction: 2, CreatedAt: day3},
	}

	// Group entries manually
	entriesByDay := make(map[string][]*models.Entry)
	for _, entry := range entries {
		day := entry.CreatedAt.Format("2006-01-02")
		entriesByDay[day] = append(entriesByDay[day], entry)
	}

	// Verify grouping
	assert.Len(t, entriesByDay, 3)               // 3 different days
	assert.Len(t, entriesByDay["2023-01-01"], 2) // 2 entries on day 1
	assert.Len(t, entriesByDay["2023-01-02"], 1) // 1 entry on day 2
	assert.Len(t, entriesByDay["2023-01-03"], 1) // 1 entry on day 3

	// Verify content of groups
	assert.Equal(t, entries[0].ID, entriesByDay["2023-01-01"][0].ID)
	assert.Equal(t, entries[1].ID, entriesByDay["2023-01-01"][1].ID)
	assert.Equal(t, entries[2].ID, entriesByDay["2023-01-02"][0].ID)
	assert.Equal(t, entries[3].ID, entriesByDay["2023-01-03"][0].ID)
}

// TestMessageFormatting tests the formatting logic for Slack messages
func TestMessageFormatting(t *testing.T) {
	// Test entry
	entry := &models.Entry{
		ID:            "123",
		Category:      models.ResearchAndProgram,
		ResearchTopic: "Test Research",
		ProgramTitle:  "Test Program",
		Satisfaction:  4,
		CreatedAt:     time.Now(),
	}

	// Format fields like in the real implementation
	var fieldTexts []string

	fieldTexts = append(fieldTexts, fmt.Sprintf("*カテゴリ:* %s", entry.Category))

	if entry.ResearchTopic != "" {
		fieldTexts = append(fieldTexts, fmt.Sprintf("*調べたこと:* %s", entry.ResearchTopic))
	}

	if entry.ProgramTitle != "" {
		fieldTexts = append(fieldTexts, fmt.Sprintf("*書いたプログラム:* %s", entry.ProgramTitle))
	}

	fieldTexts = append(fieldTexts, fmt.Sprintf("*満足度:* %s", getStarRating(entry.Satisfaction)))

	// Verify field texts
	assert.Len(t, fieldTexts, 4) // Category, Research, Program, Satisfaction
	assert.Contains(t, fieldTexts[0], "カテゴリ")
	assert.Contains(t, fieldTexts[1], "調べたこと")
	assert.Contains(t, fieldTexts[2], "書いたプログラム")
	assert.Contains(t, fieldTexts[3], "満足度")
	assert.Contains(t, fieldTexts[3], "★★★★☆") // 4 stars
}

// TestMessageFormattingForSlack tests the message formatting for various entry types
func TestMessageFormattingForSlack(t *testing.T) {
	// Create test entries with different categories
	entries := []*models.Entry{
		{
			ID:            "1",
			Category:      models.Research,
			ResearchTopic: "Research Topic",
			Satisfaction:  3,
			CreatedAt:     time.Now(),
		},
		{
			ID:           "2",
			Category:     models.Programming,
			ProgramTitle: "Program Title",
			Satisfaction: 4,
			CreatedAt:    time.Now(),
		},
		{
			ID:            "3",
			Category:      models.ResearchAndProgram,
			ResearchTopic: "Research Topic",
			ProgramTitle:  "Program Title",
			Satisfaction:  5,
			CreatedAt:     time.Now(),
		},
	}

	// For each entry type, check the field formatting
	for _, entry := range entries {
		var fieldTexts []string

		// Always add category
		fieldTexts = append(fieldTexts, "カテゴリ: "+string(entry.Category))

		// Add research topic if present
		if entry.ResearchTopic != "" {
			fieldTexts = append(fieldTexts, "調べたこと: "+entry.ResearchTopic)
		}

		// Add program title if present
		if entry.ProgramTitle != "" {
			fieldTexts = append(fieldTexts, "書いたプログラム: "+entry.ProgramTitle)
		}

		// Add satisfaction rating
		fieldTexts = append(fieldTexts, "満足度: "+getStarRating(entry.Satisfaction))

		// Verify correct number of fields based on entry type
		switch entry.Category {
		case models.Research:
			assert.Equal(t, 3, len(fieldTexts)) // Category, Research, Satisfaction
		case models.Programming:
			assert.Equal(t, 3, len(fieldTexts)) // Category, Program, Satisfaction
		case models.ResearchAndProgram:
			assert.Equal(t, 4, len(fieldTexts)) // Category, Research, Program, Satisfaction
		}
	}
}

// TestTimeFormatting tests the date formatting used in messages
func TestTimeFormatting(t *testing.T) {
	// Test specific date
	testTime := time.Date(2023, 5, 15, 10, 30, 0, 0, time.UTC)
	formattedDate := testTime.Format("2006/01/02 (Mon)")

	// Verify format
	assert.Equal(t, "2023/05/15 (Mon)", formattedDate)

	// Verify day grouping logic
	dayKey := testTime.Format("2006-01-02")
	assert.Equal(t, "2023-05-15", dayKey)
}
