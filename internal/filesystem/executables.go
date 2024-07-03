package filesystem

import (
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

const ExecutableDefinitionExt = ".flow"

func EnsureExecutableDir(workspacePath, subPath string) error {
	if _, err := os.Stat(filepath.Join(workspacePath, subPath)); os.IsNotExist(err) {
		err = os.MkdirAll(filepath.Join(workspacePath, subPath), 0750)
		if err != nil {
			return errors.Wrap(err, "unable to create executable directory")
		}
	}
	return nil
}

func InitExecutables(
	template *config.ExecutableDefinitionTemplate,
	ws *config.WorkspaceConfig,
	name, subPath string,
) error {
	if err := EnsureExecutableDir(ws.Location(), subPath); err != nil {
		return errors.Wrap(err, "unable to ensure executable directory")
	}
	if err := RenderAndWriteExecutablesTemplate(template, ws, name, subPath); err != nil {
		return errors.Wrap(err, "unable to write executable definition template")
	}
	return nil
}

func WriteExecutableDefinition(definitionFile string, definition *config.ExecutableDefinition) error {
	file, err := os.OpenFile(filepath.Clean(definitionFile), os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		return errors.Wrap(err, "unable to open definition file")
	}
	defer file.Close()

	if err := file.Truncate(0); err != nil {
		return errors.Wrap(err, "unable to truncate definition file")
	}

	err = yaml.NewEncoder(file).Encode(definition)
	if err != nil {
		return errors.Wrap(err, "unable to encode definition file")
	}

	return nil
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
	logger io.Logger,
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

func findDefinitionFiles(logger io.Logger, workspaceCfg *config.WorkspaceConfig) ([]string, error) {
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
		if err != nil {
			if errors.Is(err, fs.ErrNotExist) {
				logger.Debugx("definition path does not exist", "path", path)
				return nil
			}
			return err
		}
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
func isPathIncluded(logger io.Logger, path, basePath string, includePaths []string) bool {
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
func isPathExcluded(logger io.Logger, path, basePath string, excludedPaths []string) bool {
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
