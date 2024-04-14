package cache

import (
	"github.com/jahvon/tuikit/io"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"

	"github.com/jahvon/flow/config"
	"github.com/jahvon/flow/config/file"
)

const wsCacheKey = "workspace"

type WorkspaceCacheData struct {
	// Map of workspace name to workspace config
	Workspaces map[string]*config.WorkspaceConfig `yaml:"workspaces"`
	// Map of workspace name to workspace path
	WorkspaceLocations map[string]string `yaml:"workspaceLocations"`
}

type WorkspaceCache struct {
	Data *WorkspaceCacheData
}

func NewWorkspaceCache() *WorkspaceCache {
	if workspaceCache == nil {
		workspaceCache = &WorkspaceCache{}
	}
	return workspaceCache
}

func (c *WorkspaceCache) Update(logger io.Logger) error {
	if c.Data == nil {
		logger.Debugf("Initializing workspace cache data")
	} else {
		logger.Debugf("Updating workspace cache data")
	}

	userCfg, err := file.LoadUserConfig()
	if err != nil {
		return err
	}

	cacheData := &WorkspaceCacheData{
		Workspaces:         make(map[string]*config.WorkspaceConfig),
		WorkspaceLocations: make(map[string]string),
	}
	for name, path := range userCfg.Workspaces {
		wsCfg, err := file.LoadWorkspaceConfig(name, path)
		if err != nil {
			return errors.Wrap(err, "failed loading workspace config")
		} else if wsCfg == nil {
			logger.Errorx("config not found for workspace", "name", name, "path", path)
			continue
		}
		cacheData.Workspaces[name] = wsCfg
		cacheData.WorkspaceLocations[name] = path
	}
	data, err := yaml.Marshal(cacheData)
	if err != nil {
		return errors.Wrap(err, "unable to encode cache data")
	}

	err = file.WriteLatestCachedData(wsCacheKey, data)
	if err != nil {
		return errors.Wrap(err, "unable to write cache data")
	}

	c.Data = cacheData
	logger.Debugx("Successfully updated workspace cache data", "count", len(cacheData.Workspaces))
	return nil
}

func (c *WorkspaceCache) Get(logger io.Logger) (*WorkspaceCacheData, error) {
	if c.Data != nil {
		return c.Data, nil
	}

	cacheData, err := file.LoadLatestCachedData(wsCacheKey)
	if err != nil {
		return nil, errors.Wrap(err, "unable to load workspace cache data")
	} else if cacheData == nil {
		if err := c.Update(logger); err != nil {
			return nil, errors.Wrap(err, "unable to get updated workspace cache data")
		}
	}

	c.Data = &WorkspaceCacheData{}
	if err := yaml.Unmarshal(cacheData, c.Data); err != nil {
		return nil, errors.Wrap(err, "unable to decode workspace cache data")
	}
	return c.Data, nil
}
