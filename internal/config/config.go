package config

import (
	"os"
	"path/filepath"

	"github.com/econron/wamon/internal/slack"
	"github.com/spf13/viper"
)

// AppConfig holds all application configuration
type AppConfig struct {
	Slack slack.Config
}

// LoadConfig loads the application configuration from viper and environment variables
func LoadConfig() (*AppConfig, error) {
	// First try environment variables
	slackToken := os.Getenv("WAMON_SLACK_TOKEN")
	slackChannel := os.Getenv("WAMON_SLACK_CHANNEL")

	// If not set in env vars, use viper config
	if slackToken == "" {
		slackToken = viper.GetString("slack.token")
	}
	if slackChannel == "" {
		slackChannel = viper.GetString("slack.channel")
	}

	// Set enabled if we have a token
	enabled := slackToken != ""

	config := &AppConfig{
		Slack: slack.Config{
			Token:   slackToken,
			Channel: slackChannel,
			Enabled: enabled,
		},
	}

	return config, nil
}

// SaveSlackConfig saves Slack configuration permanently
func SaveSlackConfig(token, channel string) error {
	// Set the values
	viper.Set("slack.token", token)
	viper.Set("slack.channel", channel)
	viper.Set("slack.enabled", true)

	// Ensure config directory exists
	configDir := filepath.Join(os.Getenv("HOME"), ".wamon")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return err
	}

	// Save to file
	configFile := filepath.Join(configDir, ".wamon.yaml")
	if err := viper.WriteConfigAs(configFile); err != nil {
		return err
	}

	// 設定を再読み込み
	viper.SetConfigFile(configFile)
	return viper.ReadInConfig()
}

// SetDefaults sets default values for the configuration
func SetDefaults() {
	// Slack defaults
	viper.SetDefault("slack.enabled", false)
	viper.SetDefault("slack.channel", "general")
}
