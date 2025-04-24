package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/econron/wamon/internal/slack"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestSetDefaults(t *testing.T) {
	// Reset viper
	viper.Reset()

	// Set defaults
	SetDefaults()

	// Verify defaults are set
	assert.False(t, viper.GetBool("slack.enabled"))
	assert.Equal(t, "general", viper.GetString("slack.channel"))
}

func TestLoadConfigFromEnvironment(t *testing.T) {
	// Save original environment variables
	originalToken := os.Getenv("WAMON_SLACK_TOKEN")
	originalChannel := os.Getenv("WAMON_SLACK_CHANNEL")

	// Restore environment variables after test
	defer func() {
		os.Setenv("WAMON_SLACK_TOKEN", originalToken)
		os.Setenv("WAMON_SLACK_CHANNEL", originalChannel)
	}()

	// Set test environment variables
	os.Setenv("WAMON_SLACK_TOKEN", "test-token")
	os.Setenv("WAMON_SLACK_CHANNEL", "test-channel")

	// Reset viper and set defaults
	viper.Reset()
	SetDefaults()

	// Load config
	config, err := LoadConfig()
	assert.NoError(t, err)
	assert.NotNil(t, config)

	// Verify config values from environment
	assert.Equal(t, "test-token", config.Slack.Token)
	assert.Equal(t, "test-channel", config.Slack.Channel)
	assert.True(t, config.Slack.Enabled)
}

func TestLoadConfigFromViper(t *testing.T) {
	// Save original environment variables
	originalToken := os.Getenv("WAMON_SLACK_TOKEN")
	originalChannel := os.Getenv("WAMON_SLACK_CHANNEL")

	// Restore environment variables after test
	defer func() {
		os.Setenv("WAMON_SLACK_TOKEN", originalToken)
		os.Setenv("WAMON_SLACK_CHANNEL", originalChannel)
	}()

	// Clear environment variables
	os.Unsetenv("WAMON_SLACK_TOKEN")
	os.Unsetenv("WAMON_SLACK_CHANNEL")

	// Reset viper and set test values
	viper.Reset()
	viper.Set("slack.token", "viper-token")
	viper.Set("slack.channel", "viper-channel")

	// Load config
	config, err := LoadConfig()
	assert.NoError(t, err)
	assert.NotNil(t, config)

	// Verify config values from viper
	assert.Equal(t, "viper-token", config.Slack.Token)
	assert.Equal(t, "viper-channel", config.Slack.Channel)
	assert.True(t, config.Slack.Enabled)
}

func TestLoadConfigWithoutToken(t *testing.T) {
	// Save original environment variables
	originalToken := os.Getenv("WAMON_SLACK_TOKEN")
	originalChannel := os.Getenv("WAMON_SLACK_CHANNEL")

	// Restore environment variables after test
	defer func() {
		os.Setenv("WAMON_SLACK_TOKEN", originalToken)
		os.Setenv("WAMON_SLACK_CHANNEL", originalChannel)
	}()

	// Clear environment variables
	os.Unsetenv("WAMON_SLACK_TOKEN")
	os.Unsetenv("WAMON_SLACK_CHANNEL")

	// Reset viper and set defaults
	viper.Reset()
	SetDefaults()

	// Load config
	config, err := LoadConfig()
	assert.NoError(t, err)
	assert.NotNil(t, config)

	// Verify Slack is disabled when token is empty
	assert.Equal(t, "", config.Slack.Token)
	assert.Equal(t, "general", config.Slack.Channel)
	assert.False(t, config.Slack.Enabled)
}

func TestSaveSlackConfig(t *testing.T) {
	// Create temporary directory for config
	tempDir, err := os.MkdirTemp("", "wamon-config-test-*")
	assert.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Save original HOME environment variable
	originalHome := os.Getenv("HOME")

	// Set HOME to temp directory
	os.Setenv("HOME", tempDir)

	// Restore original HOME after test
	defer os.Setenv("HOME", originalHome)

	// Reset viper
	viper.Reset()

	// Save Slack config
	err = SaveSlackConfig("new-token", "new-channel")
	assert.NoError(t, err)

	// Check if config file was created
	configFile := filepath.Join(tempDir, ".wamon", ".wamon.yaml")
	_, err = os.Stat(configFile)
	assert.NoError(t, err)

	// Verify config values were saved
	assert.Equal(t, "new-token", viper.GetString("slack.token"))
	assert.Equal(t, "new-channel", viper.GetString("slack.channel"))
	assert.True(t, viper.GetBool("slack.enabled"))

	// Load config to verify
	config, err := LoadConfig()
	assert.NoError(t, err)
	assert.NotNil(t, config)

	// Check that the loaded config has the correct values
	assert.Equal(t, slack.Config{
		Token:   "new-token",
		Channel: "new-channel",
		Enabled: true,
	}, config.Slack)
}
