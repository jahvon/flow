package providers

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
	return "", nil
}

func (a *CLIAdapter) Set(key string, value string) error {
	return nil
}

func (a *CLIAdapter) Delete(key string) error {
	return nil
}

func (a *CLIAdapter) List() (map[string]string, error) {
	return nil, nil
}
