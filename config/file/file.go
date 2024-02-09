package file

import (
	"bufio"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/jahvon/tuikit/io"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"

	"github.com/jahvon/flow/config"
)

const (
	WorkspaceConfigFileName = "flow.yaml"
	ExecutableDefinitionExt = ".flow"
)

func WriteUserConfig(config *config.UserConfig) error {
	file, err := os.OpenFile(filepath.Clean(UserConfigPath), os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		return errors.Wrap(err, "unable to open config file")
	}
	defer file.Close()

	if err := file.Truncate(0); err != nil {
		return errors.Wrap(err, "unable to truncate config file")
	}

	err = yaml.NewEncoder(file).Encode(config)
	if err != nil {
		return errors.Wrap(err, "unable to encode config file")
	}

	return nil
}

func LoadUserConfig() (*config.UserConfig, error) {
	if err := EnsureConfigDir(); err != nil {
		return nil, errors.Wrap(err, "unable to ensure existence of config directory")
	}

	if _, err := os.Stat(UserConfigPath); os.IsNotExist(err) {
		if err := InitUserConfig(); err != nil {
			return nil, errors.Wrap(err, "unable to initialize config file")
		}
	} else if err != nil {
		return nil, errors.Wrap(err, "unable to stat config file")
	}

	file, err := os.Open(UserConfigPath)
	if err != nil {
		return nil, errors.Wrap(err, "unable to open config file")
	}
	defer file.Close()

	userCfg := &config.UserConfig{}
	err = yaml.NewDecoder(file).Decode(userCfg)
	if err != nil {
		return nil, errors.Wrap(err, "unable to decode config file")
	}

	if err := userCfg.Validate(); err != nil {
		return nil, errors.Wrap(err, "encountered validation error")
	}

	return userCfg, nil
}

func WriteWorkspaceConfig(workspacePath string, config *config.WorkspaceConfig) error {
	wsFile := filepath.Join(workspacePath, WorkspaceConfigFileName)
	file, err := os.OpenFile(filepath.Clean(wsFile), os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		return errors.Wrap(err, "unable to open workspace config file")
	}
	defer file.Close()

	if err := file.Truncate(0); err != nil {
		return errors.Wrap(err, "unable to truncate workspace config file")
	}

	err = yaml.NewEncoder(file).Encode(config)
	if err != nil {
		return errors.Wrap(err, "unable to encode workspace config file")
	}

	return nil
}

func LoadWorkspaceConfig(workspaceName, workspacePath string) (*config.WorkspaceConfig, error) {
	if err := EnsureWorkspaceDir(workspacePath); err != nil {
		return nil, errors.Wrap(err, "unable to ensure workspace directory")
	} else if err := EnsureWorkspaceConfig(workspaceName, workspacePath); err != nil {
		return nil, errors.Wrap(err, "unable to ensure workspace config file")
	}

	wsCfg := &config.WorkspaceConfig{}
	wsFile := filepath.Join(workspacePath, WorkspaceConfigFileName)
	file, err := os.Open(filepath.Clean(wsFile))
	if err != nil {
		return nil, errors.Wrap(err, "unable to open workspace config file")
	}
	defer file.Close()

	err = yaml.NewDecoder(file).Decode(wsCfg)
	if err != nil {
		return nil, errors.Wrap(err, "unable to decode workspace config file")
	}

	wsCfg.SetContext(workspaceName, workspacePath)
	return wsCfg, nil
}

func LoadExecutableDefinition(definitionFile string) (*config.ExecutableDefinition, error) {
	file, err := os.Open(filepath.Clean(definitionFile))
	if err != nil {
		return nil, errors.Wrap(err, "unable to open definition file")
	}
	defer file.Close()

	definition := &config.ExecutableDefinition{}
	err = yaml.NewDecoder(file).Decode(definition)
	if err != nil {
		return nil, errors.Wrap(err, "unable to decode definition file")
	}
	return definition, nil
}

func LoadWorkspaceExecutableDefinitions(
	logger *io.Logger,
	workspaceCfg *config.WorkspaceConfig,
) (config.ExecutableDefinitionList, error) {
	definitionFiles, err := findDefinitionFiles(logger, workspaceCfg)
	if err != nil {
		return nil, err
	}

	var definitions config.ExecutableDefinitionList
	for _, definitionFile := range definitionFiles {
		definition, err := LoadExecutableDefinition(definitionFile)
		if err != nil {
			logger.Errorx("unable to load executable definition file", "definitionFile", definitionFile, "err", err)
			continue
		}
		definition.SetContext(workspaceCfg.AssignedName(), workspaceCfg.Location(), definitionFile)
		definitions = append(definitions, definition)
	}
	logger.Debugx(
		fmt.Sprintf("loaded %d definitions", len(definitions)),
		"workspace",
		workspaceCfg.AssignedName(),
	)

	return definitions, nil
}

func WriteLatestCachedData(cacheKey string, data []byte) error {
	if err := EnsureCachedDataDir(); err != nil {
		return errors.Wrap(err, "unable to ensure existence of cache directory")
	}

	file, err := os.OpenFile(filepath.Clean(LatestCachedDataFilePath(cacheKey)), os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		return errors.Wrap(err, "unable to open cache data file")
	}
	defer file.Close()

	if err := file.Truncate(0); err != nil {
		return errors.Wrap(err, "unable to truncate cache data file")
	}

	if _, err := file.Write(data); err != nil {
		return errors.Wrap(err, "unable to write cache data file")
	}

	return nil
}

func LoadLatestCachedData(cacheKey string) ([]byte, error) {
	if err := EnsureCachedDataDir(); err != nil {
		return nil, errors.Wrap(err, "unable to ensure existence of cache directory")
	}

	if _, err := os.Stat(LatestCachedDataFilePath(cacheKey)); os.IsNotExist(err) {
		return nil, nil
	} else if err != nil {
		return nil, errors.Wrap(err, "unable to stat cache data file")
	}

	file, err := os.Open(LatestCachedDataFilePath(cacheKey))
	if err != nil {
		return nil, errors.Wrap(err, "unable to open cache data file")
	}
	defer file.Close()

	data := make([]byte, 0)
	buf := bufio.NewReader(file)
	for {
		var line []byte
		line, err = buf.ReadBytes('\n')
		if err != nil {
			break
		}
		data = append(data, line...)
	}
	if err.Error() != "EOF" {
		return nil, errors.Wrap(err, "unable to read cache data file")
	}

	return data, nil
}

func findDefinitionFiles(logger *io.Logger, workspaceCfg *config.WorkspaceConfig) ([]string, error) {
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
		if isPathIncluded(logger, path, workspaceCfg.Location(), includePaths) {
			if isPathExcluded(logger, path, workspaceCfg.Location(), excludedPaths) {
				return filepath.SkipDir
			}

			if filepath.Ext(entry.Name()) == ExecutableDefinitionExt {
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
func isPathIncluded(logger *io.Logger, path, basePath string, includePaths []string) bool {
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
			logger.Errorx(
				"unable to regex match path against include path",
				"path", path,
				"includePath", includePath,
				"err", err,
			)
			continue
		}
		return isMatch
	}
	return false
}

// IsPathExcluded returns true if the path is in any of the excluded paths.
func isPathExcluded(logger *io.Logger, path, basePath string, excludedPaths []string) bool {
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
			logger.Errorx(
				"unable to regex match path against excluded path",
				"path", path,
				"excludedPath", excludedPath,
				"err", err,
			)
			continue
		}
		return isMatch
	}
	return false
}
