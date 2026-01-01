package provider

import (
	"net/http"
	"time"

	"github.com/jackh54/mcinit/internal/cache"
)

// FoliaProvider implements Provider for Folia (uses PaperMC API)
type FoliaProvider struct {
	*PaperProvider
}

// NewFoliaProvider creates a new FoliaProvider
func NewFoliaProvider() *FoliaProvider {
	c, _ := cache.New()
	return &FoliaProvider{
		PaperProvider: &PaperProvider{
			cache:       c,
			projectName: "folia",
			client: &http.Client{
				Timeout: 30 * time.Second,
			},
		},
	}
}

// GetName returns the provider name
func (p *FoliaProvider) GetName() string {
	return "folia"
}
