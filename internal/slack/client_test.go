package slack

import (
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
