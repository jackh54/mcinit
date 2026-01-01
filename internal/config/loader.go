package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/jackh54/mcinit/internal/utils"
)

// Load loads a configuration from a file
func Load(path string) (*Config, error) {
	absPath, err := utils.AbsolutePath(path)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve config path: %w", err)
	}

	data, err := os.ReadFile(absPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("config file not found: %s", absPath)
		}
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	// Apply defaults for any missing values
	cfg.ApplyDefaults()

	// Validate configuration
	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	return &cfg, nil
}

// Save saves a configuration to a file
func Save(cfg *Config, path string) error {
	absPath, err := utils.AbsolutePath(path)
	if err != nil {
		return fmt.Errorf("failed to resolve config path: %w", err)
	}

	// Update timestamp
	cfg.UpdatedAt = time.Now().UTC()

	// Ensure parent directory exists
	dir := filepath.Dir(absPath)
	if err := utils.EnsureDir(dir); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	// Marshal config to JSON with indentation
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	// Write to file
	if err := os.WriteFile(absPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

// Exists checks if a configuration file exists
func Exists(path string) bool {
	absPath, err := utils.AbsolutePath(path)
	if err != nil {
		return false
	}
	return utils.PathExists(absPath)
}

// LoadOrCreate loads an existing config or creates a default one
func LoadOrCreate(path string) (*Config, error) {
	if Exists(path) {
		return Load(path)
	}

	cfg := DefaultConfig()
	return cfg, nil
}
