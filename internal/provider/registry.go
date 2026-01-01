package provider

import (
	"fmt"
)

var providers = map[string]Provider{}

// Register registers a provider
func Register(name string, provider Provider) {
	providers[name] = provider
}

// Get returns a provider by name
func Get(name string) (Provider, error) {
	provider, ok := providers[name]
	if !ok {
		return nil, fmt.Errorf("provider not found: %s", name)
	}
	return provider, nil
}

// List returns all registered provider names
func List() []string {
	names := make([]string, 0, len(providers))
	for name := range providers {
		names = append(names, name)
	}
	return names
}

// Exists checks if a provider is registered
func Exists(name string) bool {
	_, ok := providers[name]
	return ok
}

func init() {
	// Register all providers
	Register("vanilla", NewVanillaProvider())
	Register("paper", NewPaperProvider())
	Register("purpur", NewPurpurProvider())
	Register("folia", NewFoliaProvider())
	Register("velocity", NewVelocityProvider())
	Register("waterfall", NewWaterfallProvider())
	Register("bungee", NewBungeeProvider())
}

