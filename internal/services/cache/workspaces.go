//nolint:cyclop
package cache

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"

	"github.com/jahvon/flow/internal/common"
	"github.com/jahvon/flow/internal/config"
	"github.com/jahvon/flow/internal/io"
	"github.com/jahvon/flow/internal/workspace"
)

type WorkspaceCacheData struct {
	Workspaces         map[string]*workspace.Config `yaml:"workspaces"`
	WorkspaceLocations map[string]string            `yaml:"workspaceLocations"`
}

var (
	log           = io.Log().With().Str("service", "cache").Logger()
	cacheFilePath = common.CachedDataDirPath() + "/latest_cache"
)

func Update() (*WorkspaceCacheData, error) {
	log.Info().Msg("Starting sync")
	rootCfg := config.LoadConfig()
	if rootCfg == nil {
		return nil, fmt.Errorf("failed to load config")
	}

	if info, err := os.Stat(cacheFilePath); err != nil && os.IsNotExist(err) {
		log.Debug().Msg("Cache data file does not exist, creating")
		if _, err := os.Create(cacheFilePath); err != nil {
			return nil, fmt.Errorf("unable to create cache data file - %w", err)
		}
	} else if err != nil {
		return nil, fmt.Errorf("unable to stat cache data file - %w", err)
	} else if info.IsDir() {
		return nil, fmt.Errorf("cache data file is a directory")
	}

	cacheData := &WorkspaceCacheData{
		Workspaces:         make(map[string]*workspace.Config),
		WorkspaceLocations: make(map[string]string),
	}
	for name, path := range rootCfg.Workspaces {
		wsCfg, err := workspace.LoadConfig(name, path)
		if err != nil {
			return nil, fmt.Errorf("failed loading workspace config: %w", err)
		} else if wsCfg == nil {
			log.Err(err).Msgf("config not found for workspace %s", name)
			continue
		}
		cacheData.Workspaces[name] = wsCfg
		cacheData.WorkspaceLocations[name] = path
	}
	log.Debug().Int("workspaces", len(cacheData.Workspaces)).Msg("Successfully loaded workspace configs")

	data, err := yaml.Marshal(cacheData)
	if err != nil {
		return nil, fmt.Errorf("unable to encode cache data - %w", err)
	}

	cacheFile, err := os.OpenFile(cacheFilePath, os.O_RDWR, 0644)
	if err != nil {
		return nil, fmt.Errorf("unable to open cache data file - %w", err)
	}

	if err := cacheFile.Truncate(0); err != nil {
		return nil, fmt.Errorf("unable to truncate cache data file - %w", err)
	}

	if _, err := cacheFile.Write(data); err != nil {
		return nil, fmt.Errorf("unable to write cache data file - %w", err)
	}

	log.Debug().Msg("Completed cache sync")
	return cacheData, nil
}

func Get() (*WorkspaceCacheData, error) {
	if info, err := os.Stat(cacheFilePath); err != nil {
		if os.IsNotExist(err) {
			log.Debug().Msg("Cache data file does not exist, running sync first")
			return Update()
		}
		return nil, fmt.Errorf("unable to stat cache data file - %w", err)
	} else if info.IsDir() {
		return nil, fmt.Errorf("cache data file is a directory")
	}

	cacheFile, err := os.Open(cacheFilePath)
	if err != nil {
		return nil, fmt.Errorf("unable to open cache data file - %w", err)
	}
	defer cacheFile.Close()

	cacheData := WorkspaceCacheData{}
	if err := yaml.NewDecoder(cacheFile).Decode(&cacheData); err != nil {
		return nil, fmt.Errorf("unable to decode cache data file - %w", err)
	}
	for ws, wsCfg := range cacheData.Workspaces {
		if location, ok := cacheData.WorkspaceLocations[ws]; ok {
			wsCfg.SetContext(ws, location)
			continue
		} else {
			wsCfg.SetContext(ws, "unknown")
		}
	}
	log.Debug().Int("workspaces", len(cacheData.Workspaces)).Msg("Loaded cache data")
	return &cacheData, nil
}
