package cache

import (
	"fmt"

	"github.com/jahvon/tuikit/io"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"

	"github.com/jahvon/flow/internal/filesystem"
	"github.com/jahvon/flow/types/common"
	"github.com/jahvon/flow/types/executable"
)

const execCacheKey = "executables"

//go:generate mockgen -destination=mocks/mock_executable_cache.go -package=mocks github.com/jahvon/flow/internal/cache ExecutableCache
type ExecutableCache interface {
	Update(logger io.Logger) error
	GetExecutableByRef(logger io.Logger, ref executable.Ref) (*executable.Executable, error)
	GetExecutableList(logger io.Logger) (executable.ExecutableList, error)
}
type WorkspaceInfo struct {
	WorkspaceName string `json:"workspaceName" yaml:"workspaceName"`
	WorkspacePath string `json:"workspacePath" yaml:"workspacePath"`
}

type ExecutableCacheData struct {
	// Map of executable ref to config path
	ExecutableMap map[executable.Ref]string `json:"executableMap" yaml:"executableMap"`
	// Map of executable alias ref to primary executable ref
	AliasMap map[executable.Ref]executable.Ref `json:"aliasMap" yaml:"aliasMap"`
	// Map of config paths to their workspace / workspace path
	ConfigMap map[string]WorkspaceInfo `json:"configMap" yaml:"configMap"`

	loadedExecutables map[string]*executable.Executable
}

type ExecutableCacheImpl struct {
	Data           *ExecutableCacheData `json:",inline" yaml:",inline"`
	WorkspaceCache WorkspaceCache       `json:"-"       yaml:"-"`
}

func NewExecutableCache() ExecutableCache {
	if executableCache == nil {
		executableCache = &ExecutableCacheImpl{
			Data: &ExecutableCacheData{
				ExecutableMap: make(map[executable.Ref]string),
				AliasMap:      make(map[executable.Ref]executable.Ref),
				ConfigMap:     make(map[string]WorkspaceInfo),
			},
			WorkspaceCache: NewWorkspaceCache(),
		}
	}
	return executableCache
}

func (c *ExecutableCacheImpl) Update(logger io.Logger) error { //nolint:gocognit
	if c.Data == nil || c.WorkspaceCache == nil {
		logger.Debugf("Initializing executable cache data")
		ec, ok := NewExecutableCache().(*ExecutableCacheImpl)
		if !ok {
			return errors.New("unable to initialize executable cache")
		}
		c.Data = ec.Data
		c.WorkspaceCache = ec.WorkspaceCache
	} else {
		logger.Debugf("Updating executable cache data")
	}

	wsCacheData, err := c.WorkspaceCache.GetLatestData(logger)
	if err != nil {
		return fmt.Errorf("failed to get workspace cache data\n%w", err)
	}

	cacheData := c.Data
	for name, wsCfg := range wsCacheData.Workspaces {
		wsCfg.SetContext(name, wsCacheData.WorkspaceLocations[name])
		cfgs, err := filesystem.LoadWorkspaceFlowFiles(logger, wsCfg)
		if err != nil {
			logger.Errorx("failed to load workspace executable configs", "workspace", wsCfg.AssignedName(), "err", err)
			continue
		}
		for _, cfg := range cfgs {
			if len(cfg.FromFiles) > 0 {
				generated, err := generatedExecutables(
					logger,
					name,
					wsCfg.Location(),
					cfg.Namespace,
					cfg.ConfigPath(),
					cfg.FromFiles,
				)
				if err != nil {
					logger.Errorx(
						"failed to generate executables from files",
						"cfgPath", cfg.ConfigPath(),
						"err", err,
					)
				}
				cfg.Executables = append(cfg.Executables, generated...)
			}

			if cfg.Visibility == nil || common.Visibility(*cfg.Visibility).IsHidden() || len(cfg.Executables) == 0 {
				continue
			}
			for _, e := range cfg.Executables {
				if e == nil || (e.Visibility != nil && common.Visibility(*e.Visibility).IsHidden()) {
					continue
				}
				cacheData.ExecutableMap[e.Ref()] = cfg.ConfigPath()
				for _, ref := range enumerateExecutableAliasRefs(e) {
					cacheData.AliasMap[ref] = e.Ref()
				}
				cacheData.ConfigMap[cfg.ConfigPath()] = WorkspaceInfo{
					WorkspaceName: wsCfg.AssignedName(),
					WorkspacePath: wsCfg.Location(),
				}
			}
		}
	}

	data, err := yaml.Marshal(cacheData)
	if err != nil {
		return errors.Wrap(err, "unable to encode cache data")
	}

	err = filesystem.WriteLatestCachedData(execCacheKey, data)
	if err != nil {
		return errors.Wrap(err, "unable to write cache data")
	}

	logger.Debugx("Successfully updated executable cache data", "count", len(cacheData.ExecutableMap))
	return nil
}

func (c *ExecutableCacheImpl) GetExecutableByRef(logger io.Logger, ref executable.Ref) (*executable.Executable, error) {
	err := c.initExecutableCacheData(logger)
	if err != nil {
		return nil, err
	} else if c.Data == nil {
		return nil, errors.New("no cached executables found")
	}

	if c.Data.loadedExecutables == nil {
		c.Data.loadedExecutables = make(map[string]*executable.Executable)
	} else if exec, found := c.Data.loadedExecutables[ref.String()]; found {
		return exec, nil
	}

	cfgPath, found := c.Data.ExecutableMap[ref]
	if !found {
		if primaryRef, aliasFound := c.Data.AliasMap[ref]; aliasFound {
			cfgPath, found = c.Data.ExecutableMap[primaryRef]
			if !found {
				return nil, NewExecutableNotFoundError(ref.String())
			}
		} else {
			return nil, NewExecutableNotFoundError(ref.String())
		}
	}
	cfg, err := filesystem.LoadFlowFile(cfgPath)
	if err != nil {
		return nil, errors.Wrap(err, "unable to load executable config")
	}

	wsInfo, found := c.Data.ConfigMap[cfgPath]
	if !found {
		return nil, errors.Wrap(err, "unable to find workspace info for config")
	}

	cfg.SetDefaults()
	cfg.SetContext(wsInfo.WorkspaceName, wsInfo.WorkspacePath, cfgPath)

	generated, err := generatedExecutables(
		logger,
		wsInfo.WorkspaceName,
		wsInfo.WorkspacePath,
		cfg.Namespace,
		cfg.ConfigPath(),
		cfg.FromFiles,
	)
	if err != nil {
		logger.Warnx(
			"failed to generate executables from files",
			"cfgPath", cfgPath,
			"err", err,
		)
	}
	cfg.Executables = append(cfg.Executables, generated...)

	execs := executable.ExecutableList(cfg.Executables)
	exec, err := execs.FindByVerbAndID(ref.GetVerb(), ref.GetID())
	if err != nil {
		return nil, err
	} else if exec == nil {
		return nil, NewExecutableNotFoundError(ref.String())
	}

	c.Data.loadedExecutables[ref.String()] = exec

	return exec, nil
}

func (c *ExecutableCacheImpl) GetExecutableList(logger io.Logger) (executable.ExecutableList, error) {
	err := c.initExecutableCacheData(logger)
	if err != nil {
		return nil, err
	} else if c.Data == nil {
		return nil, errors.New("no cached executables found")
	}

	list := make(executable.ExecutableList, 0)
	for cfgPath := range c.Data.ConfigMap {
		cfg, err := filesystem.LoadFlowFile(cfgPath)
		if err != nil {
			logger.Errorx("unable to load executable config", "cfgPath", cfgPath, "err", err)
			continue
		}
		wsInfo, found := c.Data.ConfigMap[cfgPath]
		if !found {
			logger.Errorx("unable to find workspace info for config", "cfgPath", cfgPath)
			continue
		}
		cfg.SetDefaults()
		cfg.SetContext(wsInfo.WorkspaceName, wsInfo.WorkspacePath, cfgPath)

		generated, err := generatedExecutables(
			logger,
			wsInfo.WorkspaceName,
			wsInfo.WorkspacePath,
			cfg.Namespace,
			cfg.ConfigPath(),
			cfg.FromFiles,
		)
		if err != nil {
			logger.Warnx(
				"failed to generate executables from files",
				"cfgPath", cfgPath,
				"err", err,
			)
		}
		cfg.Executables = append(cfg.Executables, generated...)

		list = append(list, cfg.Executables...)
	}
	return list, nil
}

func (c *ExecutableCacheImpl) initExecutableCacheData(logger io.Logger) error {
	cacheData, err := filesystem.LoadLatestCachedData(execCacheKey)
	if err != nil {
		return errors.Wrap(err, "unable to load executable cache data")
	} else if cacheData == nil {
		if err := c.Update(logger); err != nil {
			return errors.Wrap(err, "unable to update executable cache data")
		}
	}

	c.Data = &ExecutableCacheData{}
	if err := yaml.Unmarshal(cacheData, c.Data); err != nil {
		return errors.Wrap(err, "unable to decode executable cache data")
	}
	return nil
}

func enumerateExecutableAliasRefs(exec *executable.Executable) executable.RefList {
	refs := make(executable.RefList, 0)

	for _, verb := range executable.RelatedVerbs(exec.Verb) {
		refs = append(refs, executable.NewRef(exec.ID(), verb))
		for _, id := range exec.AliasesIDs() {
			refs = append(refs, executable.NewRef(id, verb))
		}
	}

	return refs
}
