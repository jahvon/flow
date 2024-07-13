package cache

import (
	"path/filepath"
	"strings"
	"time"

	"github.com/jahvon/tuikit/io"
	"github.com/pkg/errors"

	"github.com/jahvon/flow/internal/fileparser"
	"github.com/jahvon/flow/internal/utils"
	"github.com/jahvon/flow/types/executable"
)

const generatedTag = "generated"

func generatedExecutables(
	logger io.Logger,
	wsName, wsPath, flowFileNs, flowFilePath string,
	files []string,
) (executable.ExecutableList, error) {
	executables := make(executable.ExecutableList, 0)
	for _, file := range files {
		expandedFile := utils.ExpandDirectory(logger, file, wsPath, flowFilePath, nil)
		exec, err := executablesFromFile(logger, file, expandedFile)
		if err != nil {
			return nil, err
		}
		exec.SetContext(wsName, wsPath, flowFileNs, flowFilePath)
		executables = append(executables, exec)
	}

	return executables, nil
}

func executablesFromFile(logger io.Logger, fileBase, filePath string) (*executable.Executable, error) {
	configMap, err := fileparser.ExecConfigMapFromFile(logger, filePath)
	if err != nil {
		return nil, err
	}

	exec := &executable.Executable{
		Verb: executable.Verb("exec"),
		Name: filepath.Base(fileBase),
		Exec: &executable.ExecExecutableType{
			File: fileBase,
		},
	}
	for key, value := range configMap {
		switch key {
		case fileparser.TimeoutConfigurationKey:
			dur, err := time.ParseDuration(value)
			if err != nil {
				return nil, errors.Wrapf(err, "unable to parse timeout duration %s", value)
			}
			exec.Timeout = dur
		case fileparser.VerbConfigurationKey:
			exec.Verb = executable.Verb(value)
		case fileparser.NameConfigurationKey:
			exec.Name = value
		case fileparser.VisibilityConfigurationKey:
			v := executable.ExecutableVisibility(value)
			exec.Visibility = &v
		case fileparser.DescriptionConfigurationKey:
			exec.Description = value
		case fileparser.AliasConfigurationKey:
			values := make([]string, 0)
			for _, v := range strings.Split(value, fileparser.InternalListSeparator) {
				values = append(values, strings.TrimSpace(v))
			}
			exec.Aliases = values
		case fileparser.TagConfigurationKey:
			values := make([]string, 0)
			for _, v := range strings.Split(value, fileparser.InternalListSeparator) {
				values = append(values, strings.TrimSpace(v))
			}
			exec.Tags = values
		}
	}

	exec.Tags = append(exec.Tags, generatedTag)
	return exec, nil
}
