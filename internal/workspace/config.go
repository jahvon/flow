package workspace

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/samber/lo"
	"gopkg.in/yaml.v3"

	"github.com/jahvon/flow/internal/executable"
)

type Config struct {
	DisplayName string             `json:"displayName"           yaml:"displayName"`
	Description string             `json:"description,omitempty" yaml:"description,omitempty"`
	Tags        []string           `json:"tags,omitempty"        yaml:"tags,omitempty"`
	Git         *GitConfig         `json:"git,omitempty"         yaml:"git,omitempty"`
	Executables *ExecutablesConfig `json:"executables,omitempty" yaml:"executables,omitempty"`

	assignedName string
	location     string
}

type GitConfig struct {
	Enabled    bool `json:"enabled"              yaml:"enabled"`
	PullOnSync bool `json:"pullOnSync,omitempty" yaml:"pullOnSync,omitempty"`
}

type ExecutablesConfig struct {
	Included    []string                         `json:"included,omitempty"    yaml:"included,omitempty"`
	Excluded    []string                         `json:"excluded,omitempty"    yaml:"excluded,omitempty"`
	Preferences map[string]executable.Preference `json:"preferences,omitempty" yaml:"preferences,omitempty"`
}

func (c *Config) Validate() error {
	return nil
}

func (c *Config) AssignedName() string {
	return c.assignedName
}

func (c *Config) Location() string {
	return c.location
}

func (c *Config) SetContext(name, location string) {
	c.assignedName = name
	c.location = location
}

func (c *Config) HasAnyTags(tags []string) bool {
	if len(tags) == 0 {
		return true
	}

	return lo.Some(c.Tags, tags)
}

func defaultConfig(name string) *Config {
	return &Config{
		DisplayName: name,
		Git:         &GitConfig{Enabled: false},
	}
}

func NewWorkspace(name, location string) (*Config, error) {
	if info, err := os.Stat(location); os.IsNotExist(err) {
		err = os.MkdirAll(location, 0750)
		if err != nil {
			return nil, fmt.Errorf("unable to create workspace directory - %w", err)
		}
	} else if err != nil {
		return nil, fmt.Errorf("unable to check for workspace directory - %w", err)
	} else if !info.IsDir() {
		return nil, fmt.Errorf("workspace path (%s) exists but is not a directory", location)
	}

	config := defaultConfig(name)
	config.assignedName = name
	config.location = location
	configInfo, err := os.Stat(location + "/" + ConfigFileName)
	switch {
	case os.IsNotExist(err):
		err = writeConfigFile(location, config)
		if err != nil {
			return nil, fmt.Errorf("unable to write workspace config file - %w", err)
		}
	case err != nil:
		return nil, fmt.Errorf("unable to check for workspace config file - %w", err)
	case configInfo.IsDir():
		return nil, fmt.Errorf("workspace config file (%s) exists but is a directory", location+"/"+ConfigFileName)
	default:
		config, err = LoadConfig(name, location)
		if err != nil {
			return nil, fmt.Errorf("unable to load workspace config file - %w", err)
		}
	}

	return config, nil
}

func LoadConfig(workspaceName, workspacePath string) (*Config, error) {
	config := &Config{}
	file, err := os.Open(filepath.Clean(workspacePath + "/" + ConfigFileName))
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf(
				"workspace config file does not exist. please recreate the %s workspace",
				workspaceName,
			)
		}
		return nil, fmt.Errorf("unable to open workspace config file - %w", err)
	}
	defer file.Close()

	err = yaml.NewDecoder(file).Decode(config)
	if err != nil {
		return nil, fmt.Errorf("unable to decode workspace config file - %w", err)
	}

	config.assignedName = workspaceName
	config.location = workspacePath
	return config, nil
}
