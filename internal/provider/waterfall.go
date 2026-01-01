package provider

import (
	"net/http"
	"time"

	"github.com/jackh54/mcinit/internal/cache"
)

// WaterfallProvider implements Provider for Waterfall (uses PaperMC API)
type WaterfallProvider struct {
	*PaperProvider
}

// NewWaterfallProvider creates a new WaterfallProvider
func NewWaterfallProvider() *WaterfallProvider {
	c, _ := cache.New()
	return &WaterfallProvider{
		PaperProvider: &PaperProvider{
			cache:       c,
			projectName: "waterfall",
			client: &http.Client{
				Timeout: 30 * time.Second,
			},
		},
	}
}

// GetName returns the provider name
func (p *WaterfallProvider) GetName() string {
	return "waterfall"
}
