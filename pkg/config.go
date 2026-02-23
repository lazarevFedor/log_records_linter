package pkg

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"sync"
)

type Config struct {
	// EnableLowercaseStart checks if log messages start with a lowercase letter.
	EnableLowercaseStart bool `json:"enable_lowercase_start"`

	// EnableEnglishOnly checks if log messages contain only English letters, digits, spaces, and allowed punctuation.
	EnableEnglishOnly bool `json:"enable_english_only"`

	// EnableNoSpecialChars checks if log messages do not contain special characters or emoji.
	EnableNoSpecialChars bool `json:"enable_no_special_chars"`

	// EnableSensitivePatterns checks if log messages do not contain sensitive information
	EnableSensitivePatterns bool `json:"enable_sensitive_patterns"`
}

// currentConfig holds the global configuration for the analyzer, protected by configMu for concurrent access.
var (
	configMu      sync.RWMutex
	currentConfig *Config
)

// DefaultConfig returns a Config with default settings, which can be used when no custom configuration is provided.
func DefaultConfig() *Config {
	return &Config{
		EnableLowercaseStart:    true,
		EnableEnglishOnly:       true,
		EnableNoSpecialChars:    true,
		EnableSensitivePatterns: true,
	}
}

// GetConfig returns the current configuration if it was set.
func GetConfig() *Config {
	configMu.RLock()
	defer configMu.RUnlock()
	return currentConfig
}

// ResolveConfig returns a config from SetConfig or loads it from the given path.
func ResolveConfig(path string) (*Config, error) {
	if cfg := GetConfig(); cfg != nil {
		return cfg, nil
	}
	cfg, err := LoadConfig(path)
	if err != nil {
		return DefaultConfig(), nil
	}
	return cfg, nil
}

// LoadConfig reads and parses config from a JSON file.
func LoadConfig(path string) (*Config, error) {
	if strings.TrimSpace(path) == "" {
		return DefaultConfig(), nil
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read config: %w", err)
	}

	cfg := DefaultConfig()
	if err := json.Unmarshal(data, cfg); err != nil {
		return nil, fmt.Errorf("parse config: %w", err)
	}

	return cfg, nil
}
