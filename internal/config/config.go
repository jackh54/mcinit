package config

import (
	"time"
)

// Config represents the mcinit.json configuration
type Config struct {
	Version      string        `json:"version"`
	Server       ServerConfig  `json:"server"`
	Java         JavaConfig    `json:"java"`
	JVM          JVMConfig     `json:"jvm"`
	ServerConfig ServerProps   `json:"serverConfig"`
	Plugins      PluginsConfig `json:"plugins"`
	EULA         EULAConfig    `json:"eula"`
	Paths        PathsConfig   `json:"paths"`
	CreatedAt    time.Time     `json:"createdAt"`
	UpdatedAt    time.Time     `json:"updatedAt"`
}

// ServerConfig represents server-specific configuration
type ServerConfig struct {
	Type             string `json:"type"`
	MinecraftVersion string `json:"minecraftVersion"`
	Build            string `json:"build,omitempty"`
	JarPath          string `json:"jarPath"`
	Name             string `json:"name"`
	DownloadURL      string `json:"downloadUrl,omitempty"`
	SHA256           string `json:"sha256,omitempty"`
	SHA1             string `json:"sha1,omitempty"`
}

// JavaConfig represents Java installation configuration
type JavaConfig struct {
	Version         string `json:"version"`
	Path            string `json:"path"`
	DetectedVersion string `json:"detectedVersion,omitempty"`
}

// JVMConfig represents JVM flags configuration
type JVMConfig struct {
	Xms         string   `json:"xms"`
	Xmx         string   `json:"xmx"`
	Flags       string   `json:"flags"`
	CustomFlags []string `json:"customFlags,omitempty"`
}

// ServerProps represents server.properties configuration
type ServerProps struct {
	Port        int    `json:"port"`
	NoGUI       bool   `json:"nogui"`
	MaxPlayers  int    `json:"maxPlayers"`
	OnlineMode  bool   `json:"onlineMode"`
	Difficulty  string `json:"difficulty"`
}

// PluginsConfig represents plugin linking configuration
type PluginsConfig struct {
	Links []PluginLink `json:"links,omitempty"`
}

// PluginLink represents a plugin link configuration
type PluginLink struct {
	Source      string `json:"source"`
	Mode        string `json:"mode"` // "copy" or "symlink"
	AutoRestart bool   `json:"autoRestart"`
}

// EULAConfig represents EULA acceptance configuration
type EULAConfig struct {
	Accepted   bool      `json:"accepted"`
	AcceptedAt time.Time `json:"acceptedAt,omitempty"`
}

// PathsConfig represents path configuration
type PathsConfig struct {
	ServerDir string `json:"serverDir"`
	CacheDir  string `json:"cacheDir"`
}

// Validate validates the configuration
func (c *Config) Validate() error {
	// Basic validation
	if c.Server.Type == "" {
		return &ValidationError{Field: "server.type", Message: "server type is required"}
	}

	if c.Server.MinecraftVersion == "" {
		return &ValidationError{Field: "server.minecraftVersion", Message: "Minecraft version is required"}
	}

	if c.Server.JarPath == "" {
		return &ValidationError{Field: "server.jarPath", Message: "jar path is required"}
	}

	if c.JVM.Xmx == "" {
		return &ValidationError{Field: "jvm.xmx", Message: "maximum heap size (Xmx) is required"}
	}

	if c.ServerConfig.Port < 1 || c.ServerConfig.Port > 65535 {
		return &ValidationError{Field: "serverConfig.port", Message: "port must be between 1 and 65535"}
	}

	return nil
}

// ValidationError represents a configuration validation error
type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	return e.Field + ": " + e.Message
}

