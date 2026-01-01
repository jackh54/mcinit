package provider

import (
	"net/http"
	"time"

	"github.com/jackh54/mcinit/internal/cache"
)

// VelocityProvider implements Provider for Velocity (uses PaperMC API)
type VelocityProvider struct {
	*PaperProvider
}

// NewVelocityProvider creates a new VelocityProvider
func NewVelocityProvider() *VelocityProvider {
	c, _ := cache.New()
	return &VelocityProvider{
		PaperProvider: &PaperProvider{
			cache:       c,
			projectName: "velocity",
			client: &http.Client{
				Timeout: 30 * time.Second,
			},
		},
	}
}

// GetName returns the provider name
func (p *VelocityProvider) GetName() string {
	return "velocity"
}
