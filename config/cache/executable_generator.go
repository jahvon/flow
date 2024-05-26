package cache

import (
	"path/filepath"
	"strings"
	"time"

	"github.com/jahvon/tuikit/io"
	"github.com/pkg/errors"

	"github.com/jahvon/flow/config"
	"github.com/jahvon/flow/internal/fileparser"
)

const generatedTag = "generated"

func generatedExecutables(logger io.Logger, definitionPath string, files []string) (config.ExecutableList, error) {
	executables := make(config.ExecutableList, 0)
	for _, file := range files {
		shFile := filepath.Join(filepath.Dir(definitionPath), file)
		executable, err := executablesFromFile(logger, shFile)
		if err != nil {
			return nil, err
		}
		executables = append(executables, executable)
	}

	return executables, nil
}

func executablesFromFile(logger io.Logger, shFile string) (*config.Executable, error) {
	configMap, err := fileparser.ExecConfigMapFromFile(logger, shFile)
	if err != nil {
		return nil, err
	}

	filename := filepath.Base(shFile)
	executable := &config.Executable{
		Verb: config.Verb("exec"),
		Name: filepath.Base(filename),
		Type: &config.ExecutableTypeSpec{
			Exec: &config.ExecExecutableType{
				File: filename,
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
