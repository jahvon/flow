package config

import (
	"encoding/json"
	"fmt"

	"github.com/jahvon/tuikit/types"
	"gopkg.in/yaml.v3"
)

type SecretConfig struct {
	Name   string `json:"name"  yaml:"name"`
	Secret string `json:"value" yaml:"value"`
}

type SecretConfigList []SecretConfig

type enrichedSecretConfigList struct {
	Secrets SecretConfigList `json:"secrets" yaml:"secrets"`
}

func (c *SecretConfig) AssignedName() string {
	return c.Name
}

func (c *SecretConfig) SecretValue() string {
	return c.Secret
}

func (c *SecretConfig) YAML() (string, error) {
	yamlBytes, err := yaml.Marshal(c)
	if err != nil {
		return "", fmt.Errorf("failed to marshal secret config - %w", err)
	}
	return string(yamlBytes), nil
}

func (c *SecretConfig) JSON() (string, error) {
	jsonBytes, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal secret config - %w", err)
	}
	return string(jsonBytes), nil
}

func (c *SecretConfig) Markdown() string {
	var mkdwn string
	mkdwn = fmt.Sprintf("# [Secret] %s\n", c.AssignedName())
	mkdwn += fmt.Sprintf("## Value\n%s\n", c.SecretValue())
	return mkdwn
}

func (l SecretConfigList) YAML() (string, error) {
	enriched := enrichedSecretConfigList{Secrets: l}
	yamlBytes, err := yaml.Marshal(enriched)
	if err != nil {
		return "", fmt.Errorf("failed to marshal secret config list - %w", err)
	}
	return string(yamlBytes), nil
}

func (l SecretConfigList) JSON() (string, error) {
	enriched := enrichedSecretConfigList{Secrets: l}
	jsonBytes, err := json.MarshalIndent(enriched, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal secret config list - %w", err)
	}
	return string(jsonBytes), nil
}

func (l SecretConfigList) FindByName(name string) *SecretConfig {
	for _, secret := range l {
		if secret.AssignedName() == name {
			return &secret
		}
	}
	return nil
}

func (l SecretConfigList) Items() []*types.CollectionItem {
	items := make([]*types.CollectionItem, 0)
	for _, secret := range l {
		item := types.CollectionItem{
			Header:    secret.AssignedName(),
			SubHeader: secret.SecretValue(),
			ID:        secret.AssignedName(),
		}
		items = append(items, &item)
	}
	return items
}

func (l SecretConfigList) Singular() string {
	return "secret"
}

func (l SecretConfigList) Plural() string {
	return "secrets"
}
