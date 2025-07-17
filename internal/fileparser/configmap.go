package fileparser

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/pkg/errors"

	"github.com/flowexec/flow/types/executable"
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

	commentPrefix        = "# "
	keyPrefix            = "f:"
	multiLineKeyPrefix   = "f|"
	descriptionSeparator = "\n"
	descriptionAlias     = "desc"
)

var multiLineDescriptionTag = fmt.Sprintf("<%s%s>", multiLineKeyPrefix, DescriptionConfigurationKey)

func ExtractExecConfigMap(data, prefix string) (map[string]string, error) {
	configMap := make(map[string]string)
	processingMultiLineDescription := false
	for _, line := range strings.Split(data, "\n") {
		isComment := strings.HasPrefix(line, strings.TrimSpace(prefix))
		if trimmedLine := strings.TrimSpace(line); !isComment && trimmedLine != "" {
			// If the line is not a comment or empty, break out of the loop.
			// All flow executable configuration should be at the top of the file.
			break
		}

		line = strings.TrimPrefix(line, commentPrefix)
		if processingMultiLineDescription = processMultiLineDescription(
			line, configMap, processingMultiLineDescription,
		); processingMultiLineDescription {
			continue
		}

		cfg, err := parseConfigurations(line)
		if err != nil {
			return nil, fmt.Errorf("unable to extract executable configurations: %w", err)
		}

		for key, value := range cfg {
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
	return configMap, nil
}

func applyConfig(exec *executable.Executable, key, value string) error {
	switch key {
	case TimeoutConfigurationKey:
		dur, err := time.ParseDuration(value)
		if err != nil {
			return errors.Wrapf(err, "unable to parse timeout duration %s", value)
		}
		exec.Timeout = &dur
	case VerbConfigurationKey:
		exec.Verb = executable.Verb(value)
	case NameConfigurationKey:
		exec.Name = value
	case VisibilityConfigurationKey:
		v := executable.ExecutableVisibility(value)
		exec.Visibility = &v
	case DescriptionConfigurationKey:
		exec.Description = value
	case AliasConfigurationKey:
		values := make([]string, 0)
		for _, v := range strings.Split(value, InternalListSeparator) {
			values = append(values, strings.TrimSpace(v))
		}
		exec.Aliases = values
	case TagConfigurationKey:
		values := make([]string, 0)
		for _, v := range strings.Split(value, InternalListSeparator) {
			values = append(values, strings.TrimSpace(v))
		}
		exec.Tags = values
	}
	return nil
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

func parseConfigurations(line string) (map[string]string, error) {
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
			return nil, fmt.Errorf("invalid key (%s)", key)
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
	return configMap, nil
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
