package config

import "time"

const DefaultTimeout = 30 * time.Minute

func DefaultWorkspaceConfig(name string) *WorkspaceConfig {
	return &WorkspaceConfig{
		DisplayName: name,
	}
}
