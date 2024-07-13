package vault

import (
	"encoding/json"
	"fmt"

	"github.com/jahvon/tuikit/types"
	"gopkg.in/yaml.v3"
)

type SecretValue string

func (s SecretValue) ObfuscatedString() string {
	if s.Empty() {
		return ""
	}
	return "********"
}

func (s SecretValue) String() string {
	return s.ObfuscatedString()
}

func (s SecretValue) PlainTextString() string {
	return string(s)
}

func (s SecretValue) Empty() bool {
	return string(s) == ""
}

type Secret struct {
	Reference string `json:"reference" yaml:"reference"`
	Secret    string `json:"value"     yaml:"value"`
}

func NewSecret(reference string, secret string) *Secret {
	if err := ValidateReference(reference); err != nil {
		return nil
	}
	return &Secret{Reference: reference, Secret: secret}
}

type SecretList []Secret

type enrichedSecretList struct {
	Secrets SecretList `json:"secrets" yaml:"secrets"`
}

func (c *Secret) YAML() (string, error) {
	yamlBytes, err := yaml.Marshal(c)
	if err != nil {
		return "", fmt.Errorf("failed to marshal secret - %w", err)
	}
	return string(yamlBytes), nil
}

func (c *Secret) JSON() (string, error) {
	jsonBytes, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal secret - %w", err)
	}
	return string(jsonBytes), nil
}

func (c *Secret) Markdown() string {
	var mkdwn string
	mkdwn = fmt.Sprintf("# [Secret] %s\n", c.Reference)
	mkdwn += fmt.Sprintf("**Value**\n```\n%s\n```", c.Secret)
	return mkdwn
}

func (l SecretList) YAML() (string, error) {
	enriched := enrichedSecretList{Secrets: l}
	yamlBytes, err := yaml.Marshal(enriched)
	if err != nil {
		return "", fmt.Errorf("failed to marshal secret list - %w", err)
	}
	return string(yamlBytes), nil
}

func (l SecretList) JSON() (string, error) {
	enriched := enrichedSecretList{Secrets: l}
	jsonBytes, err := json.MarshalIndent(enriched, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal secret list - %w", err)
	}
	return string(jsonBytes), nil
}

func (l SecretList) FindByName(name string) *Secret {
	for _, secret := range l {
		if secret.Reference == name {
			return &secret
		}
	}
	return nil
}

func (l SecretList) Items() []*types.CollectionItem {
	items := make([]*types.CollectionItem, 0)
	for _, secret := range l {
		item := types.CollectionItem{
			Header: secret.Reference,
			ID:     secret.Reference,
		}
		items = append(items, &item)
	}
	return items
}

func (l SecretList) Singular() string {
	return "secret"
}

func (l SecretList) Plural() string {
	return "secrets"
}
