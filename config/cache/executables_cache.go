package cache

import (
	"fmt"

	"gopkg.in/yaml.v3"

	"github.com/jahvon/flow/config"
	"github.com/jahvon/flow/config/file"
)

const execCacheKey = "executables"

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

type ExecutableCache struct {
	Data *ExecutableCacheData `json:",inline" yaml:",inline"`
}

func NewExecutableCache() *ExecutableCache {
	if executableCache == nil {
		executableCache = &ExecutableCache{}
	}
	return executableCache
}

func (c *ExecutableCache) Update() error { //nolint:gocognit
	if c.Data == nil {
		log.Debug().Msg("Initializing executable cache data")
		c.Data = &ExecutableCacheData{
			ExecutableMap: make(map[config.Ref]string),
		}
	} else {
		log.Debug().Msg("Updating executable cache data")
	}

	if c.Data.ExecutableMap == nil {
		c.Data.ExecutableMap = make(map[config.Ref]string)
	}
	if c.Data.DefinitionMap == nil {
		c.Data.DefinitionMap = make(map[string]WorkspaceInfo)
	}

	wsCache := NewWorkspaceCache()
	wsCacheData, err := wsCache.Get()
	if err != nil {
		return fmt.Errorf("failed to get workspace cache data - %w", err)
	}

	cacheData := &ExecutableCacheData{
		ExecutableMap: make(map[config.Ref]string),
		AliasMap:      make(map[config.Ref]config.Ref),
		DefinitionMap: make(map[string]WorkspaceInfo),
	}
	for _, wsCfg := range wsCacheData.Workspaces {
		definitions, err := file.LoadWorkspaceExecutableDefinitions(wsCfg)
		if err != nil {
			return fmt.Errorf("failed to load workspace executable definitions - %w", err)
		}
		for _, def := range definitions {
			if def == nil || def.Visibility == config.VisibilityHidden || len(def.Executables) == 0 {
				continue
			}
			for _, e := range def.Executables {
				if e == nil || e.Visibility == config.VisibilityHidden {
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
	log.Trace().Int("executables", len(cacheData.ExecutableMap)).Msg("Successfully loaded executable definitions")

	data, err := yaml.Marshal(cacheData)
	if err != nil {
		return fmt.Errorf("unable to encode cache data - %w", err)
	}

	err = file.WriteLatestCachedData(execCacheKey, data)
	if err != nil {
		return fmt.Errorf("unable to write cache data - %w", err)
	}

	c.Data = cacheData
	log.Trace().Msg("Successfully updated executable cache data")

	return nil
}

func (c *ExecutableCache) GetExecutableByRef(ref config.Ref) (*config.Executable, error) {
	err := c.initExecutableCacheData()
	if err != nil {
		return nil, err
	} else if c.Data == nil {
		return nil, fmt.Errorf("no cached executables found")
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
				return nil, fmt.Errorf("unable to find executable with reference %s", ref)
			}
		} else {
			return nil, fmt.Errorf("unable to find executable with reference %s", ref)
		}
	}
	definition, err := file.LoadExecutableDefinition(definitionPath)
	if err != nil {
		return nil, fmt.Errorf("unable to load executable definition - %w", err)
	}

	wsInfo, found := c.Data.DefinitionMap[definitionPath]
	if !found {
		return nil, fmt.Errorf("unable to find workspace info for definition %s", definitionPath)
	}

	definition.SetDefaults()
	definition.SetContext(wsInfo.WorkspaceName, wsInfo.WorkspacePath, definitionPath)

	executable, err := definition.Executables.FindByVerbAndID(ref.GetVerb(), ref.GetID())
	if err != nil {
		return nil, err
	} else if executable == nil {
		return nil, fmt.Errorf("unable to find executable with reference %s", ref)
	}

	c.Data.loadedExecutables[ref.String()] = executable

	return executable, nil
}

func (c *ExecutableCache) GetExecutableList() (config.ExecutableList, error) {
	err := c.initExecutableCacheData()
	if err != nil {
		return nil, err
	} else if c.Data == nil {
		return nil, fmt.Errorf("no cached executables found")
	}

	list := make(config.ExecutableList, 0)
	for definitionPath := range c.Data.DefinitionMap {
		definition, err := file.LoadExecutableDefinition(definitionPath)
		if err != nil {
			return nil, fmt.Errorf("unable to load executable definition - %w", err)
		}
		wsInfo, found := c.Data.DefinitionMap[definitionPath]
		if !found {
			return nil, fmt.Errorf("unable to find workspace info for definition %s", definitionPath)
		}
		definition.SetDefaults()
		definition.SetContext(wsInfo.WorkspaceName, wsInfo.WorkspacePath, definitionPath)
		list = append(list, definition.Executables...)
	}
	return list, nil
}

func (c *ExecutableCache) initExecutableCacheData() error {
	if c.Data != nil {
		return nil
	}

	cacheData, err := file.LoadLatestCachedData(execCacheKey)
	if err != nil {
		return fmt.Errorf("unable to load executable cache data - %w", err)
	} else if cacheData == nil {
		if err := c.Update(); err != nil {
			return fmt.Errorf("unable to get updated workspace cache data - %w", err)
		}
	}

	c.Data = &ExecutableCacheData{}
	if err := yaml.Unmarshal(cacheData, c.Data); err != nil {
		return fmt.Errorf("unable to decode executable cache data - %w", err)
	}
	log.Trace().Int("executables", len(c.Data.ExecutableMap)).Msg("Fetched executable cache data")

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
