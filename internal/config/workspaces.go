package config

import (
	"fmt"
	"path/filepath"

	"github.com/jahvon/flow/internal/workspace"
)

func SetCurrentWorkspace(config *RootConfig, name string) error {
	if _, ok := config.Workspaces[name]; !ok {
		return fmt.Errorf("workspace %s does not exist", name)
	}

	config.CurrentWorkspace = name
	return writeConfigFile(config)
}

func SetWorkspace(config *RootConfig, name, location string) error {
	workspaceDir := filepath.Join(location, name)
	if err := workspace.SetDirectory(workspaceDir); err != nil {
		return err
	}

	config.Workspaces[name] = workspaceDir
	return writeConfigFile(config)
}

func DeleteWorkspace(config *RootConfig, name string) error {
	if _, found := config.Workspaces[name]; !found {
		return fmt.Errorf("workspace %s does not exist", name)
	}

	delete(config.Workspaces, name)
	return writeConfigFile(config)
}
