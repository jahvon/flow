package providers

const DefaultProviderName = "default"

type DefaultProvider struct{}

func (p *DefaultProvider) New(config map[string]interface{}) (Adapter, error) {
	return &DefaultAdapter{}, nil
}

type DefaultAdapter struct{}

func (a *DefaultAdapter) Get(key string) (string, error) {
	return "", nil
}

func (a *DefaultAdapter) Set(key string, value string) error {
	return nil
}

func (a *DefaultAdapter) Delete(key string) error {
	return nil
}

func (a *DefaultAdapter) List() (map[string]string, error) {
	return nil, nil
}
