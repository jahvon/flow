package cache

import (
	"fmt"

	"github.com/jahvon/tuikit/io"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"

	"github.com/jahvon/flow/config"
	"github.com/jahvon/flow/internal/filesystem"
)

const execCacheKey = "executables"

//go:generate mockgen -destination=mocks/mock_executable_cache.go -package=mocks github.com/jahvon/flow/internal/cache ExecutableCache
type ExecutableCache interface {
	Update(logger io.Logger) error
	GetExecutableByRef(logger io.Logger, ref config.Ref) (*config.Executable, error)
	GetExecutableList(logger io.Logger) (config.ExecutableList, error)
}
type WorkspaceInfo struct {
	WorkspaceName string `json:"workspaceName" yaml:"workspaceName"`
	WorkspacePath string `json:"workspacePath" yaml:"workspacePath"`
}

type ExecutableCacheData struct {
	// Map of executable ref to definition path
	ExecutableMap map[config.Ref]string `json:"executableMap" yaml:"executableMap"`
	// Map of executable alias ref to primary executable ref
	AliasMap map[config.Ref]config.Ref `json:"aliasMap" yaml:"aliasMap"`
	// Map of definition paths to their workspace / workspace path
	DefinitionMap map[string]WorkspaceInfo `json:"definitionMap" yaml:"definitionMap"`

	loadedExecutables map[string]*config.Executable
}

type ExecutableCacheImpl struct {
	Data           *ExecutableCacheData `json:",inline" yaml:",inline"`
	WorkspaceCache WorkspaceCache       `json:"-"       yaml:"-"`
}

func NewExecutableCache() ExecutableCache {
	if executableCache == nil {
		executableCache = &ExecutableCacheImpl{
			Data: &ExecutableCacheData{
				ExecutableMap: make(map[config.Ref]string),
				AliasMap:      make(map[config.Ref]config.Ref),
				DefinitionMap: make(map[string]WorkspaceInfo),
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
		definitions, err := filesystem.LoadWorkspaceExecutableDefinitions(logger, wsCfg)
		if err != nil {
			logger.Errorx("failed to load workspace executable definitions", "workspace", wsCfg.AssignedName(), "err", err)
			continue
		}
		for _, def := range definitions {
			if len(def.FromFiles) > 0 {
				generated, err := generatedExecutables(
					logger,
					name,
					wsCfg.Location(),
					def.Namespace,
					def.DefinitionPath(),
					def.FromFiles,
				)
				if err != nil {
					logger.Errorx(
						"failed to generate executables from files",
						"definitionPath", def.DefinitionPath(),
						"err", err,
					)
				}
				def.Executables = append(def.Executables, generated...)
			}

			if def == nil || def.Visibility == config.VisibilityHidden || len(def.Executables) == 0 {
				continue
			}
			for _, e := range def.Executables {
				if e == nil || (e.Visibility != nil && *e.Visibility == config.VisibilityHidden) {
					continue
				}
				cacheData.ExecutableMap[e.Ref()] = def.DefinitionPath()
				for _, ref := range enumerateExecutableAliasRefs(e) {
					cacheData.AliasMap[ref] = e.Ref()
				}
				cacheData.DefinitionMap[def.DefinitionPath()] = WorkspaceInfo{
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

func (c *ExecutableCacheImpl) GetExecutableByRef(logger io.Logger, ref config.Ref) (*config.Executable, error) {
	err := c.initExecutableCacheData(logger)
	if err != nil {
		return nil, err
	} else if c.Data == nil {
		return nil, errors.New("no cached executables found")
	}

	if c.Data.loadedExecutables == nil {
		c.Data.loadedExecutables = make(map[string]*config.Executable)
	} else if executable, found := c.Data.loadedExecutables[ref.String()]; found {
		return executable, nil
	}

	definitionPath, found := c.Data.ExecutableMap[ref]
	if !found {
		if primaryRef, aliasFound := c.Data.AliasMap[ref]; aliasFound {
			definitionPath, found = c.Data.ExecutableMap[primaryRef]
			if !found {
				return nil, NewExecutableNotFoundError(ref.String())
			}
		} else {
			return nil, NewExecutableNotFoundError(ref.String())
		}
	}
	definition, err := filesystem.LoadExecutableDefinition(definitionPath)
	if err != nil {
		return nil, errors.Wrap(err, "unable to load executable definition")
	}

	wsInfo, found := c.Data.DefinitionMap[definitionPath]
	if !found {
		return nil, errors.Wrap(err, "unable to find workspace info for definition")
	}

	definition.SetDefaults()
	definition.SetContext(wsInfo.WorkspaceName, wsInfo.WorkspacePath, definitionPath)

	generated, err := generatedExecutables(
		logger,
		wsInfo.WorkspaceName,
		wsInfo.WorkspacePath,
		definition.Namespace,
		definition.DefinitionPath(),
		definition.FromFiles,
	)
	if err != nil {
		logger.Warnx(
			"failed to generate executables from files",
			"definitionPath", definitionPath,
			"err", err,
		)
	}
	definition.Executables = append(definition.Executables, generated...)

	executable, err := definition.Executables.FindByVerbAndID(ref.GetVerb(), ref.GetID())
	if err != nil {
		return nil, err
	} else if executable == nil {
		return nil, NewExecutableNotFoundError(ref.String())
	}

	c.Data.loadedExecutables[ref.String()] = executable

	return executable, nil
}

func (c *ExecutableCacheImpl) GetExecutableList(logger io.Logger) (config.ExecutableList, error) {
	err := c.initExecutableCacheData(logger)
	if err != nil {
		return nil, err
	} else if c.Data == nil {
		return nil, errors.New("no cached executables found")
	}

	list := make(config.ExecutableList, 0)
	for definitionPath := range c.Data.DefinitionMap {
		definition, err := filesystem.LoadExecutableDefinition(definitionPath)
		if err != nil {
			logger.Errorx("unable to load executable definition", "definitionPath", definitionPath, "err", err)
			continue
		}
		wsInfo, found := c.Data.DefinitionMap[definitionPath]
		if !found {
			logger.Errorx("unable to find workspace info for definition", "definitionPath", definitionPath)
			continue
		}
		definition.SetDefaults()
		definition.SetContext(wsInfo.WorkspaceName, wsInfo.WorkspacePath, definitionPath)

		generated, err := generatedExecutables(
			logger,
			wsInfo.WorkspaceName,
			wsInfo.WorkspacePath,
			definition.Namespace,
			definition.DefinitionPath(),
			definition.FromFiles,
		)
		if err != nil {
			logger.Warnx(
				"failed to generate executables from files",
				"definitionPath", definitionPath,
				"err", err,
			)
		}
		definition.Executables = append(definition.Executables, generated...)

		list = append(list, definition.Executables...)
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

func enumerateExecutableAliasRefs(executable *config.Executable) []config.Ref {
	refs := make([]config.Ref, 0)

	for _, verb := range config.RelatedVerbs(executable.Verb) {
		refs = append(refs, config.NewRef(executable.ID(), verb))
		for _, id := range executable.AliasesIDs() {
			refs = append(refs, config.NewRef(id, verb))
		}
	}

	return refs
}
