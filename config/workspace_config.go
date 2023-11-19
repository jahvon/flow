package config

type WorkspaceConfig struct {
	DisplayName string             `json:"displayName"           yaml:"displayName"`
	Description string             `json:"description,omitempty" yaml:"description,omitempty"`
	Tags        Tags               `json:"tags,omitempty"        yaml:"tags,omitempty"`
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
	Included []string `json:"included,omitempty" yaml:"included,omitempty"`
	Excluded []string `json:"excluded,omitempty" yaml:"excluded,omitempty"`
}

func (c *WorkspaceConfig) AssignedName() string {
	return c.assignedName
}

func (c *WorkspaceConfig) Location() string {
	return c.location
}

func (c *WorkspaceConfig) SetContext(name, location string) {
	c.assignedName = name
	c.location = location
}
