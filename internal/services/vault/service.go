package vault

import (
	"fmt"

	"github.com/jahvon/flow/internal/services/vault/providers"
	"github.com/jahvon/flow/types/config"
)

//go:generate mockgen -destination=mocks/mock_provider.go -package=mocks github.com/jahvon/flow/internal/services/vault Provider
type Provider interface {
	New(config map[string]interface{}) (providers.Adapter, error)
}

func NewVault(cfg *config.Vault) (providers.Adapter, error) {
	if cfg == nil || cfg.Current == "" {
		provider, _ := Providers.Get(providers.DefaultProviderName)
		return provider.New(nil)
	}

	if len(cfg.Vaults) == 0 {
		return nil, fmt.Errorf("no vaults configured")
	}

	v, exists := cfg.Vaults[cfg.Current]
	if !exists {
		return nil, fmt.Errorf("unknown vault: %s", cfg.Current)
	}
	vault, ok := v.(config.Entry)
	if !ok {
		return nil, fmt.Errorf("invalid vault configuration")
	}

	provider, exists := Providers.Get(string(vault.Provider))
	if !exists {
		return nil, fmt.Errorf("unknown vault provider: %s", vault.Provider)
	}

	return provider.New(vault.Config)
}
