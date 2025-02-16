package vault

import "github.com/jahvon/flow/internal/services/vault/providers"

func init() {
	Providers.Register(providers.DefaultProviderName, &providers.DefaultProvider{})
	Providers.Register(providers.CLIProviderName, &providers.CLIProvider{})
}

type Provider interface {
	New(config map[string]interface{}) (Adapter, error)
}

var Providers = &Registry{
	providers: make(map[string]Provider),
}
