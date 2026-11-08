package provider

// VelocityProvider implements Provider for Velocity (uses the PaperMC Fill API)
type VelocityProvider struct {
	*PaperProvider
}

// NewVelocityProvider creates a new VelocityProvider
func NewVelocityProvider() *VelocityProvider {
	return &VelocityProvider{PaperProvider: newFillProvider("velocity")}
}

// GetName returns the provider name
func (p *VelocityProvider) GetName() string {
	return "velocity"
}
