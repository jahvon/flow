package cache

import (
	"fmt"

	"github.com/jahvon/flow/config"
	"github.com/jahvon/flow/config/file"
)

type WorkspaceInfo struct {
	WorkspaceName string `json:"workspaceName" yaml:"workspaceName"`
	WorkspacePath string `json:"workspacePath" yaml:"workspacePath"`
}

type ExecutableCacheData struct {
	// Map of executable ID to definition path
	ExecutableMap map[config.Ref]string `json:"executableMap" yaml:"executableMap"`
	// Map of definition paths to their workspace / workspace path
	DefinitionMap map[string]WorkspaceInfo `json:"definitionMap" yaml:"definitionMap"`

	loadedExecutables map[string]*config.Executable
}

type ExecutableCache struct {
	Data *ExecutableCacheData `json:",inline" yaml:",inline"`
}

func NewExecutableCache() *ExecutableCache {
	return &ExecutableCache{}
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
				c.Data.ExecutableMap[e.Ref()] = def.DefinitionPath()
				c.Data.DefinitionMap[def.DefinitionPath()] = WorkspaceInfo{
					WorkspaceName: wsCfg.AssignedName(),
					WorkspacePath: wsCfg.Location(),
				}
			}
		}
	}

	log.Trace().Int("executables", len(c.Data.ExecutableMap)).Msg("Successfully loaded executable definitions")
	return nil
}

func (c *ExecutableCache) GetExecutableByRef(ref config.Ref) (*config.Executable, error) {
	if c.Data == nil || c.Data.ExecutableMap == nil || c.Data.DefinitionMap == nil {
		if err := c.Update(); err != nil {
			return nil, fmt.Errorf("unable to get updated executable cache data - %w", err)
		}
	}

	if c.Data.loadedExecutables == nil {
		c.Data.loadedExecutables = make(map[string]*config.Executable)
	} else if executable, found := c.Data.loadedExecutables[ref.String()]; found {
		return executable, nil
	}

	definitionPath, found := c.Data.ExecutableMap[ref]
	if !found {
		return nil, fmt.Errorf("unable to find executable with reference %s", ref)
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
		return nil, fmt.Errorf("unable to find executable with reference %s - %w", ref, err)
	} else if executable == nil {
		return nil, fmt.Errorf("unable to find executable with reference %s", ref)
	}

	c.Data.loadedExecutables[ref.String()] = executable

	return executable, nil
}

func (c *ExecutableCache) GetExecutableList() (config.ExecutableList, error) {
	if c.Data == nil || c.Data.ExecutableMap == nil || c.Data.DefinitionMap == nil {
		if err := c.Update(); err != nil {
			return nil, fmt.Errorf("unable to get updated executable cache data - %w", err)
		}
	}

	list := make(config.ExecutableList, 0)
	for _, definitionPath := range c.Data.ExecutableMap {
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
