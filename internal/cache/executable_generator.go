package cache

import (
	"path/filepath"
	"strings"
	"time"

	"github.com/jahvon/tuikit/io"
	"github.com/pkg/errors"

	"github.com/jahvon/flow/config"
	"github.com/jahvon/flow/internal/fileparser"
	"github.com/jahvon/flow/internal/utils"
)

const generatedTag = "generated"

func generatedExecutables(
	logger io.Logger,
	wsName, wsPath, definitionNs, definitionPath string,
	files []string,
) (config.ExecutableList, error) {
	executables := make(config.ExecutableList, 0)
	for _, file := range files {
		expandedFile := utils.ExpandDirectory(logger, file, wsPath, definitionPath, nil)
		executable, err := executablesFromFile(logger, file, expandedFile)
		if err != nil {
			return nil, err
		}
		executable.SetContext(wsName, wsPath, definitionNs, definitionPath)
		executables = append(executables, executable)
	}

	return executables, nil
}

func executablesFromFile(logger io.Logger, fileBase, filePath string) (*config.Executable, error) {
	configMap, err := fileparser.ExecConfigMapFromFile(logger, filePath)
	if err != nil {
		return nil, err
	}

	executable := &config.Executable{
		Verb: config.Verb("exec"),
		Name: filepath.Base(fileBase),
		Type: &config.ExecutableTypeSpec{
			Exec: &config.ExecExecutableType{
				File: fileBase,
			},
		},
	}
	for key, value := range configMap {
		switch key {
		case fileparser.TimeoutConfigurationKey:
			dur, err := time.ParseDuration(value)
			if err != nil {
				return nil, errors.Wrapf(err, "unable to parse timeout duration %s", value)
			}
			executable.Timeout = dur
		case fileparser.VerbConfigurationKey:
			executable.Verb = config.Verb(value)
		case fileparser.NameConfigurationKey:
			executable.Name = value
		case fileparser.VisibilityConfigurationKey:
			v := config.Visibility(value)
			executable.Visibility = &v
		case fileparser.DescriptionConfigurationKey:
			executable.Description = value
		case fileparser.AliasConfigurationKey:
			values := make([]string, 0)
			for _, v := range strings.Split(value, fileparser.InternalListSeparator) {
				values = append(values, strings.TrimSpace(v))
			}
			executable.Aliases = values
		case fileparser.TagConfigurationKey:
			values := make([]string, 0)
			for _, v := range strings.Split(value, fileparser.InternalListSeparator) {
				values = append(values, strings.TrimSpace(v))
			}
			executable.Tags = values
		}
	}

	executable.Tags = append(executable.Tags, generatedTag)
	return executable, nil
}
