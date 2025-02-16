package vault

import (
	"fmt"

	"github.com/jahvon/flow/internal/services/vault/providers"
)

func NewVault(config *Config) (Adapter, error) {
	if config.Vault == nil {
		provider, _ := Providers.Get(providers.DefaultProviderName)
		return provider.New(nil)
	}

	provider, exists := Providers.Get(config.Vault.Provider)
	if !exists {
		return nil, fmt.Errorf("unknown vault provider: %s", config.Vault.Provider)
	}

	return provider.New(config.Vault.Config)
}
