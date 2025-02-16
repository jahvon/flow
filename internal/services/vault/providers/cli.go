package providers

import "github.com/jahvon/flow/internal/services/run"

const CLIProviderName = "cli"

type CLIProvider struct{}

func (p *CLIProvider) New(config map[string]interface{}) (Adapter, error) {
	return &CLIAdapter{}, nil
}

type CLIAdapter struct {
	getCommand    string
	setCommand    string
	deleteCommand string
	listCommand   string
}

func (a *CLIAdapter) Get(key string) (string, error) {
	run.RunCmd()
}
