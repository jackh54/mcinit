package provider

// WaterfallProvider implements Provider for Waterfall (uses the PaperMC Fill API)
type WaterfallProvider struct {
	*PaperProvider
}

// NewWaterfallProvider creates a new WaterfallProvider
func NewWaterfallProvider() *WaterfallProvider {
	return &WaterfallProvider{PaperProvider: newFillProvider("waterfall")}
}

// GetName returns the provider name
func (p *WaterfallProvider) GetName() string {
	return "waterfall"
}
