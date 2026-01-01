package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()

	if cfg.Version != "1.0.0" {
		t.Errorf("Expected version 1.0.0, got %s", cfg.Version)
	}

	if cfg.Server.Type != "paper" {
		t.Errorf("Expected default server type 'paper', got %s", cfg.Server.Type)
	}

	if cfg.JVM.Xms != "2G" {
		t.Errorf("Expected default Xms '2G', got %s", cfg.JVM.Xms)
	}

	if cfg.JVM.Xmx != "4G" {
		t.Errorf("Expected default Xmx '4G', got %s", cfg.JVM.Xmx)
	}

	if cfg.ServerConfig.Port != 25565 {
		t.Errorf("Expected default port 25565, got %d", cfg.ServerConfig.Port)
	}
}

func TestConfigValidation(t *testing.T) {
	tests := []struct {
		name    string
		cfg     *Config
		wantErr bool
	}{
		{
			name: "valid config",
			cfg: func() *Config {
				c := DefaultConfig()
				c.Server.MinecraftVersion = "1.21.4"
				return c
			}(),
			wantErr: false,
		},
		{
			name: "missing server type",
			cfg: &Config{
				Server:       ServerConfig{MinecraftVersion: "1.21.4", JarPath: "server.jar"},
				JVM:          JVMConfig{Xmx: "4G"},
				ServerConfig: ServerProps{Port: 25565},
			},
			wantErr: true,
		},
		{
			name: "invalid port",
			cfg: &Config{
				Server:       ServerConfig{Type: "paper", MinecraftVersion: "1.21.4", JarPath: "server.jar"},
				JVM:          JVMConfig{Xmx: "4G"},
				ServerConfig: ServerProps{Port: 99999},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.cfg.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestConfigSaveAndLoad(t *testing.T) {
	// Create temp directory
	tempDir, err := os.MkdirTemp("", "mcinit-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer func() { _ = os.RemoveAll(tempDir) }()

	cfgPath := filepath.Join(tempDir, "mcinit.json")

	// Create and save config
	cfg := DefaultConfig()
	cfg.Server.Name = "test-server"
	cfg.Server.MinecraftVersion = "1.21.4"

	if err := Save(cfg, cfgPath); err != nil {
		t.Fatalf("Failed to save config: %v", err)
	}

	// Check file exists
	if !Exists(cfgPath) {
		t.Error("Config file should exist after Save()")
	}

	// Load config
	loaded, err := Load(cfgPath)
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	// Verify loaded data
	if loaded.Server.Name != "test-server" {
		t.Errorf("Expected name 'test-server', got %s", loaded.Server.Name)
	}

	if loaded.Server.MinecraftVersion != "1.21.4" {
		t.Errorf("Expected version '1.21.4', got %s", loaded.Server.MinecraftVersion)
	}
}

func TestApplyDefaults(t *testing.T) {
	cfg := &Config{}
	cfg.ApplyDefaults()

	if cfg.Version == "" {
		t.Error("Version should not be empty after ApplyDefaults()")
	}

	if cfg.Server.Type == "" {
		t.Error("Server type should not be empty after ApplyDefaults()")
	}

	if cfg.JVM.Xms == "" {
		t.Error("Xms should not be empty after ApplyDefaults()")
	}
}

func TestConfigJSON(t *testing.T) {
	cfg := DefaultConfig()
	cfg.Server.MinecraftVersion = "1.21.4"

	// Marshal to JSON
	data, err := json.Marshal(cfg)
	if err != nil {
		t.Fatalf("Failed to marshal config: %v", err)
	}

	// Unmarshal back
	var loaded Config
	if err := json.Unmarshal(data, &loaded); err != nil {
		t.Fatalf("Failed to unmarshal config: %v", err)
	}

	if loaded.Server.MinecraftVersion != "1.21.4" {
		t.Errorf("Expected version '1.21.4', got %s", loaded.Server.MinecraftVersion)
	}
}

func TestLoadNonExistent(t *testing.T) {
	_, err := Load("/nonexistent/path/mcinit.json")
	if err == nil {
		t.Error("Load() should return error for non-existent file")
	}
}

