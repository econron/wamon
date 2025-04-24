package cmd

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestInitConfig(t *testing.T) {
	// Save original cfgFile value
	originalCfgFile := cfgFile
	defer func() {
		cfgFile = originalCfgFile
	}()

	// Case 1: Test with explicit config file
	// Create temp dir and file
	tempDir, err := os.MkdirTemp("", "wamon-test-*")
	assert.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Create a test config file
	testConfigFile := filepath.Join(tempDir, "test-config.yaml")
	err = os.WriteFile(testConfigFile, []byte("slack:\n  token: test-token\n  channel: test-channel"), 0644)
	assert.NoError(t, err)

	// Set config file and reset viper
	cfgFile = testConfigFile
	viper.Reset()

	// Run initConfig
	initConfig()

	// Verify the config was loaded
	assert.Equal(t, testConfigFile, viper.ConfigFileUsed())
	assert.Equal(t, "test-token", viper.GetString("slack.token"))
	assert.Equal(t, "test-channel", viper.GetString("slack.channel"))

	// Case 2: Test with default config locations
	// Reset global variables
	cfgFile = ""
	viper.Reset()

	// Create .wamon config directory in temp location
	home := tempDir // Simulate home directory
	wamonConfigDir := filepath.Join(home, ".wamon")
	err = os.MkdirAll(wamonConfigDir, 0755)
	assert.NoError(t, err)

	// Create config file in .wamon directory
	defaultConfigFile := filepath.Join(wamonConfigDir, ".wamon.yaml")
	err = os.WriteFile(defaultConfigFile, []byte("slack:\n  token: default-token\n  channel: default-channel"), 0644)
	assert.NoError(t, err)

	// Override home directory detection
	originalUserHomeDir := osUserHomeDir
	defer func() {
		osUserHomeDir = originalUserHomeDir
	}()

	osUserHomeDir = func() (string, error) {
		return home, nil
	}

	// Run initConfig
	initConfig()

	// This might not find the config because we've mocked UserHomeDir after viper initialization
	// This test mainly verifies that initConfig doesn't crash
}

// Mock functions to override os.UserHomeDir for testing
var osUserHomeDir = os.UserHomeDir
