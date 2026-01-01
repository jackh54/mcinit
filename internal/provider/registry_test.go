package provider

import (
	"testing"
)

func TestProviderRegistry(t *testing.T) {
	// Test that all expected providers are registered
	expectedProviders := []string{
		"vanilla",
		"paper",
		"purpur",
		"folia",
		"velocity",
		"waterfall",
		"bungee",
	}

	for _, name := range expectedProviders {
		t.Run(name, func(t *testing.T) {
			if !Exists(name) {
				t.Errorf("Provider %s not registered", name)
			}

			provider, err := Get(name)
			if err != nil {
				t.Errorf("Get(%s) error = %v", name, err)
			}

			if provider == nil {
				t.Errorf("Get(%s) returned nil provider", name)
			}

			if provider.GetName() != name {
				t.Errorf("Provider.GetName() = %s, want %s", provider.GetName(), name)
			}
		})
	}
}

func TestGetNonExistent(t *testing.T) {
	_, err := Get("nonexistent")
	if err == nil {
		t.Error("Get() should return error for non-existent provider")
	}
}

func TestList(t *testing.T) {
	providers := List()
	
	if len(providers) == 0 {
		t.Error("List() returned no providers")
	}

	// Check that it includes expected providers
	expectedCount := 7 // vanilla, paper, purpur, folia, velocity, waterfall, bungee
	if len(providers) != expectedCount {
		t.Errorf("List() returned %d providers, expected %d", len(providers), expectedCount)
	}
}

func TestProviderInterface(t *testing.T) {
	// Test that provider interface methods exist
	provider, _ := Get("paper")
	
	// These should not panic
	_ = provider.GetName()
	
	// Note: We don't test actual API calls here to avoid network dependencies
	// Those would be integration tests
}

