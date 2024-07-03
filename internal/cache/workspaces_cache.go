package cache

import (
	"github.com/jahvon/tuikit/io"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"

	"github.com/jahvon/flow/config"
	"github.com/jahvon/flow/internal/filesystem"
)

const wsCacheKey = "workspace"

//go:generate mockgen -destination=mocks/mock_workspace_cache.go -package=mocks github.com/jahvon/flow/internal/cache WorkspaceCache
type WorkspaceCache interface {
	Update(logger io.Logger) error
	GetData() *WorkspaceCacheData
	GetLatestData(logger io.Logger) (*WorkspaceCacheData, error)
	GetWorkspaceConfigList(logger io.Logger) (config.WorkspaceConfigList, error)
}
type WorkspaceCacheData struct {
	// Map of workspace name to workspace config
	Workspaces map[string]*config.WorkspaceConfig `yaml:"workspaces"`
	// Map of workspace name to workspace path
	WorkspaceLocations map[string]string `yaml:"workspaceLocations"`
}

type WorkspaceCacheImpl struct {
	Data *WorkspaceCacheData
}

func NewWorkspaceCache() WorkspaceCache {
	if workspaceCache == nil {
		workspaceCache = &WorkspaceCacheImpl{
			Data: &WorkspaceCacheData{
				Workspaces:         make(map[string]*config.WorkspaceConfig),
				WorkspaceLocations: make(map[string]string),
			},
		}
	}
	return workspaceCache
}

func (c *WorkspaceCacheImpl) Update(logger io.Logger) error {
	if c.Data == nil {
		logger.Debugf("Initializing workspace cache data")
		wc, ok := NewWorkspaceCache().(*WorkspaceCacheImpl)
		if !ok {
			return errors.New("unable to initialize workspace cache")
		}
		c.Data = wc.Data
	} else {
		logger.Debugf("Updating workspace cache data")
	}

	userCfg, err := filesystem.LoadUserConfig()
	if err != nil {
		return err
	}

	cacheData := c.Data
	for name, path := range userCfg.Workspaces {
		wsCfg, err := filesystem.LoadWorkspaceConfig(name, path)
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

	err = filesystem.WriteLatestCachedData(wsCacheKey, data)
	if err != nil {
		return errors.Wrap(err, "unable to write cache data")
	}

	logger.Debugx("Successfully updated workspace cache data", "count", len(cacheData.Workspaces))
	return nil
}

func (c *WorkspaceCacheImpl) GetData() *WorkspaceCacheData {
	return c.Data
}

func (c *WorkspaceCacheImpl) GetLatestData(logger io.Logger) (*WorkspaceCacheData, error) {
	cacheData, err := filesystem.LoadLatestCachedData(wsCacheKey)
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

func (c *WorkspaceCacheImpl) GetWorkspaceConfigList(logger io.Logger) (config.WorkspaceConfigList, error) {
	cache, err := c.GetLatestData(logger)
	if err != nil {
		return nil, err
	}

	wsCfgs := make(config.WorkspaceConfigList, 0, len(c.Data.Workspaces))
	for wsName, wsCfg := range cache.Workspaces {
		wsCfg.SetContext(wsName, cache.WorkspaceLocations[wsName])
		wsCfgs = append(wsCfgs, *wsCfg)
	}
	return wsCfgs, nil
}
