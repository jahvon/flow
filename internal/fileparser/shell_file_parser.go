package fileparser

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/flowexec/tuikit/io"
)

const (
	TimeoutConfigurationKey     = "timeout"
	VerbConfigurationKey        = "verb"
	NameConfigurationKey        = "name"
	AliasConfigurationKey       = "alias"
	DescriptionConfigurationKey = "description"
	VisibilityConfigurationKey  = "visibility"
	TagConfigurationKey         = "tag"

	InternalListSeparator = ","

	shellCommentPrefix   = "# "
	keyPrefix            = "f:"
	multiLineKeyPrefix   = "f|"
	descriptionSeparator = "\n"
	descriptionAlias     = "desc"
)

var multiLineDescriptionTag = fmt.Sprintf("<%s%s>", multiLineKeyPrefix, DescriptionConfigurationKey)

func ExecConfigMapFromFile(logger io.Logger, file string) (map[string]string, error) {
	if err := validateFile(file); err != nil {
		return nil, err
	}

	fileBytes, err := os.ReadFile(filepath.Clean(file))
	if err != nil {
		return nil, err
	}

	configMap := make(map[string]string)
	processingMultiLineDescription := false
	for _, line := range strings.Split(string(fileBytes), "\n") {
		isComment := strings.HasPrefix(line, strings.TrimSpace(shellCommentPrefix))
		if trimmedLine := strings.TrimSpace(line); !isComment && trimmedLine != "" {
			// If the line is not a comment or empty, break out of the loop.
			// All flow executable configuration should be at the top of the file.
			break
		}

		line = strings.TrimPrefix(line, shellCommentPrefix)
		if processingMultiLineDescription = processMultiLineDescription(
			line, configMap, processingMultiLineDescription,
		); processingMultiLineDescription {
			continue
		}

		for key, value := range parseConfigurations(logger, line) {
			switch key {
			case TimeoutConfigurationKey, VerbConfigurationKey, NameConfigurationKey, VisibilityConfigurationKey:
				configMap[key] = strings.TrimSpace(value)
			case DescriptionConfigurationKey, descriptionAlias:
				if existingValue, ok := configMap[DescriptionConfigurationKey]; ok {
					configMap[DescriptionConfigurationKey] =
						fmt.Sprintf("%s%s%s", existingValue, descriptionSeparator, value)
				} else {
					configMap[DescriptionConfigurationKey] = value
				}
			case AliasConfigurationKey, TagConfigurationKey:
				value = strings.TrimSpace(value)
				if existingValue, ok := configMap[key]; ok {
					configMap[key] = fmt.Sprintf("%s%s%s", existingValue, InternalListSeparator, value)
				} else {
					configMap[key] = value
				}
			}
		}
	}
	if len(configMap) == 0 {
		return nil, fmt.Errorf("no flow configurations found in file (%s)", file)
	}
	return configMap, nil
}

// This regex is used to extract all flow configurations from a shell script.
// The regex matches the following:
// - f:<key>=<value>
// - f:<key>="<value>"
// - f:<key>='<value>'
// - f:<key>="<value>",<key>='<value>',key=<value>
// - f:exampleKey='Example value\, with comma'
// - f:exampleKey="Example value\, with comma"
// - f:exampleKey="\'Example value with escaped quotes\'"
// - f:exampleKey='\"Example value with escaped quotes\"'
// and so on.
var flowConfigRegex = regexp.MustCompile(`f:\w+=(?:"(?:[^"\\]|\\.)*"|'(?:[^'\\]|\\.)*'|[^, ]+)`)

func validateFile(file string) error {
	info, err := os.Stat(file)
	if err != nil {
		return err
	} else if info.IsDir() {
		return fmt.Errorf("file (%s) is a directory", file)
	}
	ext := filepath.Ext(file)
	if ext != ".sh" {
		return fmt.Errorf("file (%s) is not a shell script", file)
	}
	return nil
}

func processMultiLineDescription(line string, configMap map[string]string, processing bool) bool {
	multiLineDescPrefixed := strings.HasPrefix(line, multiLineDescriptionTag)
	multiLineDescSuffixed := strings.HasSuffix(line, multiLineDescriptionTag)
	existingValue := configMap[DescriptionConfigurationKey]
	switch {
	case processing && multiLineDescPrefixed:
		processing = false
	case processing && multiLineDescSuffixed:
		processing = false
		line = strings.TrimSuffix(line, multiLineDescriptionTag)
		if existingValue != "" {
			configMap[DescriptionConfigurationKey] = fmt.Sprintf(
				"%s%s%s",
				existingValue,
				descriptionSeparator,
				line,
			)
		} else {
			configMap[DescriptionConfigurationKey] = line
		}
	case processing:
		if existingValue != "" {
			configMap[DescriptionConfigurationKey] = fmt.Sprintf(
				"%s%s%s",
				existingValue,
				descriptionSeparator,
				line,
			)
		} else {
			configMap[DescriptionConfigurationKey] = line
		}
	case multiLineDescPrefixed:
		processing = true
		line = strings.TrimPrefix(line, multiLineDescriptionTag)
		if trimmedLine := strings.TrimSpace(line); trimmedLine == "" {
			return processing
		} else if existingValue != "" {
			configMap[DescriptionConfigurationKey] = fmt.Sprintf(
				"%s%s%s",
				existingValue,
				descriptionSeparator,
				line,
			)
		} else {
			configMap[DescriptionConfigurationKey] = line
		}
	}
	return processing
}

func parseConfigurations(logger io.Logger, line string) map[string]string {
	configMap := make(map[string]string)
	matches := flowConfigRegex.FindAllString(line, -1)
	for _, match := range matches {
		split := strings.SplitN(match, "=", 2)
		key := strings.TrimSpace(strings.ToLower(strings.TrimPrefix(split[0], keyPrefix)))
		value := split[1]
		// Removing quotes if present
		if (value[0] == '"' && value[len(value)-1] == '"') || (value[0] == '\'' && value[len(value)-1] == '\'') {
			value = value[1 : len(value)-1]
		}
		// Replace escaped characters
		value = strings.ReplaceAll(value, `\"`, `"`)
		value = strings.ReplaceAll(value, `\'`, `'`)
		value = strings.ReplaceAll(value, `\\`, `\`)
		value = strings.ReplaceAll(value, `\,`, `,`)

		if !validateKey(key) {
			logger.Warnf("invalid key (%s) in configuration", key)
			continue
		}

		switch key {
		case AliasConfigurationKey, TagConfigurationKey:
			value = strings.TrimSpace(value)
			if existingValue, ok := configMap[key]; ok {
				configMap[key] = fmt.Sprintf("%s%s%s", existingValue, InternalListSeparator, value)
			} else {
				configMap[key] = value
			}
		case DescriptionConfigurationKey, descriptionAlias:
			if existingValue, ok := configMap[DescriptionConfigurationKey]; ok {
				configMap[DescriptionConfigurationKey] =
					fmt.Sprintf("%s%s%s", existingValue, descriptionSeparator, value)
			} else {
				configMap[DescriptionConfigurationKey] = value
			}
		default:
			configMap[key] = strings.TrimSpace(value)
		}
	}
	return configMap
}

func validateKey(key string) bool {
	switch key {
	case VerbConfigurationKey, NameConfigurationKey, DescriptionConfigurationKey, descriptionAlias,
		AliasConfigurationKey, VisibilityConfigurationKey, TagConfigurationKey,
		TimeoutConfigurationKey:
		return true
	default:
		return false
	}
}
