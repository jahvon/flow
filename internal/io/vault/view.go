package vault

import (
	"fmt"
	"path/filepath"

	"github.com/jahvon/tuikit"
	"github.com/jahvon/tuikit/types"
	"github.com/jahvon/tuikit/views"
	"github.com/jahvon/vault"
	"gopkg.in/yaml.v3"

	vault2 "github.com/jahvon/flow/internal/vault"
)

type vaultEntity struct {
	Name string `json:"name" yaml:"name"`
	Path string `json:"path" yaml:"path"`
	Type string `json:"type" yaml:"type"`

	Data map[string]interface{} `json:"data" yaml:"data"`
}

func (v *vaultEntity) YAML() (string, error) {
	data, err := yaml.Marshal(v)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func (v *vaultEntity) JSON() (string, error) {
	data, err := yaml.Marshal(v)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func (v *vaultEntity) Markdown() string {
	md := fmt.Sprintf(
		"# [Vault] %s\n\n**Path:** %s\n\n**Type:** %s\n\n",
		v.Name, v.Path, v.Type,
	)

	if v.Data != nil {
		md += "## Data\n\n"
		for key, value := range v.Data {
			md += fmt.Sprintf("**%s:** %v\n\n", key, value)
		}
	}

	return md
}

func NewVaultView(
	container *tuikit.Container,
	vaultName string,
) tuikit.View {
	v, err := vaultFromName(vaultName)
	if err != nil {
		container.HandleError(fmt.Errorf("failed to load vault: %w", err))
		return nil
	}
	return views.NewEntityView(container.RenderState(), v, types.EntityFormatDocument)
}

type vaultCollection struct {
	Vaults []*vaultEntity `json:"vaults" yaml:"vaults"`
}

func (vc *vaultCollection) Singular() string {
	return "vault"
}

func (vc *vaultCollection) Plural() string {
	return "vaults"
}

func (vc *vaultCollection) Items() []*types.EntityInfo {
	items := make([]*types.EntityInfo, len(vc.Vaults))
	for i, v := range vc.Vaults {
		items[i] = &types.EntityInfo{
			Header:    v.Name,
			SubHeader: v.Path,
		}
	}
	return items
}

func (vc *vaultCollection) YAML() (string, error) {
	data, err := yaml.Marshal(vc)
	if err != nil {
		return "", fmt.Errorf("failed to marshal vaults: %w", err)
	}
	return string(data), nil
}

func (vc *vaultCollection) JSON() (string, error) {
	data, err := yaml.Marshal(vc)
	if err != nil {
		return "", fmt.Errorf("failed to marshal vaults: %w", err)
	}
	return string(data), nil
}

func NewVaultListView(
	container *tuikit.Container,
	vaultNames []string,
) tuikit.View {
	vaults := &vaultCollection{Vaults: make([]*vaultEntity, 0, len(vaultNames))}
	for _, name := range vaultNames {
		v, err := vaultFromName(name)
		if err != nil {
			container.HandleError(fmt.Errorf("failed to load vault %s: %w", name, err))
			continue
		}
		vaults.Vaults = append(vaults.Vaults, v)
	}
	if len(vaults.Vaults) == 0 {
		container.HandleError(fmt.Errorf("no vaults found"))
		return nil
	}

	selectFunc := func(filterVal string) error {
		for _, v := range vaults.Vaults {
			if v.Name == filterVal {
				return container.SetView(NewVaultView(container, v.Name))
			}
		}
		return fmt.Errorf("vault not found")
	}

	return views.NewCollectionView(
		container.RenderState(),
		vaults,
		types.CollectionFormatList,
		selectFunc,
	)
}

func vaultFromName(vaultName string) (*vaultEntity, error) {
	cfg, err := vault.LoadConfigJSON(
		filepath.Join(vault2.CacheDirectory(fmt.Sprintf("configs/%s.json", vaultName))),
	)
	if err != nil {
		return nil, err
	}
	v := &vaultEntity{
		Name: cfg.ID,
		Type: string(cfg.Type),
	}

	data := make(map[string]interface{})
	switch cfg.Type {
	case vault.ProviderTypeAES256:
		v.Path = cfg.Aes.StoragePath
		data["sources"] = cfg.Aes.KeySource
		v.Data = data
	case vault.ProviderTypeAge:
		v.Path = cfg.Age.StoragePath
		data["sources"] = cfg.Age.IdentitySources
		data["recipients"] = cfg.Age.Recipients
		v.Data = data
	}

	return v, nil
}
