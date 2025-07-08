package cache

import (
	"encoding/json"

	"github.com/flowexec/tuikit"
	"github.com/flowexec/tuikit/types"
	"github.com/flowexec/tuikit/views"
	"gopkg.in/yaml.v3"
)

type cacheData struct {
	Cache map[string]string `json:"cache" yaml:"cache"`
}

func (d *cacheData) Items() []*types.EntityInfo {
	items := make([]*types.EntityInfo, 0, len(d.Cache))
	for key, value := range d.Cache {
		items = append(items, &types.EntityInfo{
			Header:    key,
			SubHeader: value,
			ID:        key,
		})
	}
	return items
}

func (d *cacheData) Singular() string {
	return "Entry"
}

func (d *cacheData) Plural() string {
	return "Entries"
}

func (d *cacheData) YAML() (string, error) {
	data, err := yaml.Marshal(d)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func (d *cacheData) JSON() (string, error) {
	data, err := json.MarshalIndent(d, "", "  ")
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func NewCacheListView(
	container *tuikit.Container,
	cache map[string]string,
) tuikit.View {
	data := &cacheData{Cache: cache}
	return views.NewCollectionView(container.RenderState(), data, types.CollectionFormatList, nil)
}
