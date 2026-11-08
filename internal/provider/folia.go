package provider

// FoliaProvider implements Provider for Folia (uses the PaperMC Fill API)
type FoliaProvider struct {
	*PaperProvider
}

// NewFoliaProvider creates a new FoliaProvider
func NewFoliaProvider() *FoliaProvider {
	return &FoliaProvider{PaperProvider: newFillProvider("folia")}
}

// GetName returns the provider name
func (p *FoliaProvider) GetName() string {
	return "folia"
}
