package file

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"gopkg.in/yaml.v3"

	"github.com/jahvon/flow/config"
	"github.com/jahvon/flow/internal/io"
)

var log = io.Log().With().Str("scope", "config/file").Logger()

const (
	WorkspaceConfigFileName = "flow.yaml"
	ExecutableDefinitionExt = ".flow"
)

func WriteUserConfig(config *config.UserConfig) error {
	file, err := os.OpenFile(filepath.Clean(UserConfigPath), os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		return fmt.Errorf("unable to open config file - %w", err)
	}
	defer file.Close()

	if err := file.Truncate(0); err != nil {
		return fmt.Errorf("unable to truncate config file - %w", err)
	}

	err = yaml.NewEncoder(file).Encode(config)
	if err != nil {
		return fmt.Errorf("unable to encode config file - %w", err)
	}

	return nil
}

func LoadUserConfig() *config.UserConfig {
	if err := EnsureConfigDir(); err != nil {
		log.Panic().Err(err).Msg("encountered issue in ensuring existence of config directory")
	}

	if _, err := os.Stat(UserConfigPath); os.IsNotExist(err) {
		if err := InitUserConfig(); err != nil {
			log.Panic().Err(err).Msg("unable to initialize config file")
		}
	} else if err != nil {
		log.Panic().Err(err).Msg("unable to stat config file")
	}

	file, err := os.Open(UserConfigPath)
	if err != nil {
		log.Panic().Err(err).Msg("unable to open config file")
	}
	defer file.Close()

	userCfg := &config.UserConfig{}
	err = yaml.NewDecoder(file).Decode(userCfg)
	if err != nil {
		log.Panic().Err(err).Msg("unable to decode config file")
	}

	if err := userCfg.Validate(); err != nil {
		log.Panic().Err(err).Msg("encountered validation error")
	}

	return userCfg
}

func WriteWorkspaceConfig(workspacePath string, config *config.WorkspaceConfig) error {
	file, err := os.OpenFile(filepath.Join(workspacePath, WorkspaceConfigFileName), os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		return fmt.Errorf("unable to open workspace config file - %w", err)
	}
	defer file.Close()

	if err := file.Truncate(0); err != nil {
		return fmt.Errorf("unable to truncate workspace config file - %w", err)
	}

	err = yaml.NewEncoder(file).Encode(config)
	if err != nil {
		return fmt.Errorf("unable to encode workspace config file - %w", err)
	}

	return nil
}

func LoadWorkspaceConfig(workspaceName, workspacePath string) (*config.WorkspaceConfig, error) {
	if err := EnsureWorkspaceDir(workspacePath); err != nil {
		return nil, fmt.Errorf("unable to ensure workspace directory - %w", err)
	} else if err := EnsureWorkspaceConfig(workspaceName, workspacePath); err != nil {
		return nil, fmt.Errorf("unable to ensure workspace config file - %w", err)
	}

	wsCfg := &config.WorkspaceConfig{}
	file, err := os.Open(filepath.Join(workspacePath, WorkspaceConfigFileName))
	if err != nil {
		return nil, fmt.Errorf("unable to open workspace config file - %w", err)
	}
	defer file.Close()

	err = yaml.NewDecoder(file).Decode(wsCfg)
	if err != nil {
		return nil, fmt.Errorf("unable to decode workspace config file - %w", err)
	}

	wsCfg.SetContext(workspaceName, workspacePath)
	return wsCfg, nil
}

func LoadExecutableDefinition(definitionFile string) (*config.ExecutableDefinition, error) {
	file, err := os.Open(filepath.Clean(definitionFile))
	if err != nil {
		return nil, fmt.Errorf("unable to open definition file - %w", err)
	}
	defer file.Close()

	definition := &config.ExecutableDefinition{}
	err = yaml.NewDecoder(file).Decode(definition)
	if err != nil {
		return nil, fmt.Errorf("unable to decode definition file - %w", err)
	}
	return definition, nil
}

func LoadWorkspaceExecutableDefinitions(workspaceCfg *config.WorkspaceConfig) (config.ExecutableDefinitionList, error) {
	definitionFiles, err := findDefinitionFiles(workspaceCfg)
	if err != nil {
		return nil, err
	}

	var definitions config.ExecutableDefinitionList
	for _, definitionFile := range definitionFiles {
		definition, err := LoadExecutableDefinition(definitionFile)
		if err != nil {
			return nil, err
		}
		definition.SetContext(workspaceCfg.AssignedName(), workspaceCfg.Location(), definitionFile)
		definitions = append(definitions, definition)
	}
	log.Trace().Msgf("loaded %d definitions", len(definitions))

	return definitions, nil
}

func WriteLatestCachedData(cacheKey string, data []byte) error {
	if err := EnsureCachedDataDir(); err != nil {
		return fmt.Errorf("unable to ensure existence of cache directory - %w", err)
	}

	file, err := os.OpenFile(filepath.Clean(LatestCachedDataFilePath(cacheKey)), os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		return fmt.Errorf("unable to open cache data file - %w", err)
	}
	defer file.Close()

	if err := file.Truncate(0); err != nil {
		return fmt.Errorf("unable to truncate cache data file - %w", err)
	}

	if _, err := file.Write(data); err != nil {
		return fmt.Errorf("unable to write cache data file - %w", err)
	}

	return nil
}

func findDefinitionFiles(workspaceCfg *config.WorkspaceConfig) ([]string, error) {
	var includePaths, excludedPaths []string
	if workspaceCfg.Executables != nil {
		includePaths = workspaceCfg.Executables.Included
		if len(includePaths) == 0 {
			includePaths = []string{workspaceCfg.Location()}
		}

		excludedPaths = workspaceCfg.Executables.Excluded
	} else {
		includePaths = []string{workspaceCfg.Location()}
	}

	var definitionPaths []string
	walkDirFunc := func(path string, entry fs.DirEntry, err error) error {
		if isPathIncluded(path, workspaceCfg.Location(), includePaths) {
			if isPathExcluded(path, workspaceCfg.Location(), excludedPaths) {
				return filepath.SkipDir
			}

			if filepath.Ext(entry.Name()) == ExecutableDefinitionExt {
				log.Trace().Msgf("found definition file %s", path)
				definitionPaths = append(definitionPaths, path)
			}
		}
		return nil
	}

	if err := filepath.WalkDir(workspaceCfg.Location(), walkDirFunc); err != nil {
		return nil, err
	}
	return definitionPaths, nil
}

// IsPathIn returns true if the path is in any of the include paths.
func isPathIncluded(path, basePath string, includePaths []string) bool {
	if includePaths == nil {
		return true
	}

	for _, p := range includePaths {
		includePath := p
		if strings.HasPrefix(includePath, "//") {
			includePath = strings.Replace(includePath, "//", basePath+"/", 1)
		}

		if path == includePath || strings.HasPrefix(path, includePath) {
			return true
		}

		isMatch, err := regexp.MatchString(includePath, path)
		if err != nil {
			log.Err(err).Msgf("unable to regex match path %s against include path %s", path, includePath)
			continue
		}
		return isMatch
	}
	return false
}

// IsPathExcluded returns true if the path is in any of the excluded paths.
func isPathExcluded(path, basePath string, excludedPaths []string) bool {
	if excludedPaths == nil {
		return false
	}

	for _, p := range excludedPaths {
		excludedPath := p
		if strings.HasPrefix(excludedPath, "//") {
			excludedPath = strings.Replace(excludedPath, "//", basePath+"/", 1)
		}

		if path == excludedPath || strings.HasPrefix(path, excludedPath) {
			return true
		}

		isMatch, err := regexp.MatchString(excludedPath, path)
		if err != nil {
			log.Err(err).Msgf("unable to regex match path %s against excluded path %s", path, excludedPath)
			continue
		}
		return isMatch
	}
	return false
}
