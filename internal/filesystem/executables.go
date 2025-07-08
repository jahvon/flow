package filesystem

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/flowexec/tuikit/io"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"

	"github.com/jahvon/flow/types/executable"
	"github.com/jahvon/flow/types/workspace"
)

func EnsureExecutableDir(workspacePath, subPath string) error {
	if _, err := os.Stat(filepath.Join(workspacePath, subPath)); os.IsNotExist(err) {
		err = os.MkdirAll(filepath.Join(workspacePath, subPath), 0750)
		if err != nil {
			return errors.Wrap(err, "unable to create executable directory")
		}
	}
	return nil
}

func WriteFlowFile(cfgFile string, cfg *executable.FlowFile) error {
	file, err := os.OpenFile(filepath.Clean(cfgFile), os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		return errors.Wrap(err, "unable to open cfg file")
	}
	defer file.Close()

	if err := file.Truncate(0); err != nil {
		return errors.Wrap(err, "unable to truncate config file")
	}

	err = yaml.NewEncoder(file).Encode(cfg)
	if err != nil {
		return errors.Wrap(err, "unable to encode config file")
	}

	return nil
}

func LoadFlowFile(cfgFile string) (*executable.FlowFile, error) {
	file, err := os.Open(filepath.Clean(cfgFile))
	if err != nil {
		return nil, errors.Wrap(err, "unable to open config file")
	}
	defer file.Close()

	cfg := &executable.FlowFile{}
	err = yaml.NewDecoder(file).Decode(cfg)
	if err != nil {
		return nil, errors.Wrap(err, "unable to decode config file")
	}
	return cfg, nil
}

func LoadWorkspaceFlowFiles(
	logger io.Logger,
	workspaceCfg *workspace.Workspace,
) (executable.FlowFileList, error) {
	cfgFiles, err := findFlowFiles(logger, workspaceCfg)
	if err != nil {
		return nil, err
	}

	var cfgs executable.FlowFileList
	for _, cfgFile := range cfgFiles {
		cfg, err := LoadFlowFile(cfgFile)
		if err != nil {
			logger.Errorx("unable to load executable config file", "configFile", cfgFile, "err", err)
			continue
		}
		cfg.SetDefaults()
		cfg.SetContext(workspaceCfg.AssignedName(), workspaceCfg.Location(), cfgFile)
		cfgs = append(cfgs, cfg)
	}
	logger.Debugx(
		fmt.Sprintf("loaded %d config files", len(cfgs)),
		"workspace",
		workspaceCfg.AssignedName(),
	)

	return cfgs, nil
}

func findFlowFiles(logger io.Logger, workspaceCfg *workspace.Workspace) ([]string, error) {
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

	var cfgPaths []string
	walkDirFunc := func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			if errors.Is(err, fs.ErrNotExist) {
				logger.Debugx("cfg path does not exist", "path", path)
				return nil
			}
			return err
		}
		if isPathIncluded(logger, path, workspaceCfg.Location(), includePaths) {
			if isPathExcluded(logger, path, workspaceCfg.Location(), excludedPaths) {
				return filepath.SkipDir
			}

			if filepath.Ext(entry.Name()) == executable.FlowFileExt {
				cfgPaths = append(cfgPaths, path)
			}
		}
		return nil
	}

	if err := filepath.WalkDir(workspaceCfg.Location(), walkDirFunc); err != nil {
		return nil, err
	}
	return cfgPaths, nil
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
