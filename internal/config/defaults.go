package config

import (
	"time"
)

// DefaultConfig returns a new Config with default values
func DefaultConfig() *Config {
	now := time.Now().UTC()

	return &Config{
		Version: "1.0.0",
		Server: ServerConfig{
			Type:             "paper",
			MinecraftVersion: "",
			Build:            "latest",
			JarPath:          "server.jar",
			Name:             "dev-server",
		},
		Java: JavaConfig{
			Version: "auto",
			Path:    "auto",
		},
		JVM: JVMConfig{
			Xms:         "2G",
			Xmx:         "4G",
			Flags:       "aikar",
			CustomFlags: []string{},
		},
		ServerConfig: ServerProps{
			Port:        25565,
			NoGUI:       true,
			MaxPlayers:  20,
			OnlineMode:  false,
			Difficulty:  "easy",
		},
		Plugins: PluginsConfig{
			Links: []PluginLink{},
		},
		EULA: EULAConfig{
			Accepted: false,
		},
		Paths: PathsConfig{
			ServerDir: ".",
			CacheDir:  "", // Will be populated from cache system
		},
		CreatedAt: now,
		UpdatedAt: now,
	}
}

// ApplyDefaults fills in missing values with defaults
func (c *Config) ApplyDefaults() {
	if c.Version == "" {
		c.Version = "1.0.0"
	}

	if c.Server.Type == "" {
		c.Server.Type = "paper"
	}

	if c.Server.Build == "" {
		c.Server.Build = "latest"
	}

	if c.Server.JarPath == "" {
		c.Server.JarPath = "server.jar"
	}

	if c.Server.Name == "" {
		c.Server.Name = "dev-server"
	}

	if c.Java.Version == "" {
		c.Java.Version = "auto"
	}

	if c.Java.Path == "" {
		c.Java.Path = "auto"
	}

	if c.JVM.Xms == "" {
		c.JVM.Xms = "2G"
	}

	if c.JVM.Xmx == "" {
		c.JVM.Xmx = "4G"
	}

	if c.JVM.Flags == "" {
		c.JVM.Flags = "aikar"
	}

	if c.JVM.CustomFlags == nil {
		c.JVM.CustomFlags = []string{}
	}

	if c.ServerConfig.Port == 0 {
		c.ServerConfig.Port = 25565
	}

	if c.ServerConfig.MaxPlayers == 0 {
		c.ServerConfig.MaxPlayers = 20
	}

	if c.ServerConfig.Difficulty == "" {
		c.ServerConfig.Difficulty = "easy"
	}

	if c.Plugins.Links == nil {
		c.Plugins.Links = []PluginLink{}
	}

	if c.Paths.ServerDir == "" {
		c.Paths.ServerDir = "."
	}

	if c.CreatedAt.IsZero() {
		c.CreatedAt = time.Now().UTC()
	}

	if c.UpdatedAt.IsZero() {
		c.UpdatedAt = time.Now().UTC()
	}
}

