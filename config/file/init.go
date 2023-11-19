package file

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/jahvon/flow/config"
)

func InitUserConfig() error {
	if err := EnsureDefaultWorkspace(); err != nil {
		return fmt.Errorf("unable to ensure default workspace - %w", err)
	}
	defaultWsName := "default"
	if _, err := os.Stat(filepath.Join(DefaultWorkspacePath, WorkspaceConfigFileName)); os.IsNotExist(err) {
		if err := InitWorkspaceConfig(defaultWsName, DefaultWorkspacePath); err != nil {
			return fmt.Errorf("unable to initialize default workspace config - %w", err)
		}
	} else if err != nil {
		return fmt.Errorf("unable to check for default workspace config - %w", err)
	}
	log.Trace().Msg("default workspace initialized")

	defaultCfg := &config.UserConfig{
		Workspaces:       map[string]string{defaultWsName: DefaultWorkspacePath},
		CurrentWorkspace: defaultWsName,
	}

	_, err := os.Create(UserConfigPath)
	if err != nil {
		return fmt.Errorf("unable to create config file - %w", err)
	}
	err = WriteUserConfig(defaultCfg)
	if err != nil {
		return fmt.Errorf("unable to write config file - %w", err)
	}
	log.Trace().Msg("user config initialized")
	return nil
}

func InitWorkspaceConfig(name, path string) error {
	wsCfg := config.DefaultWorkspaceConfig(name)

	if err := EnsureWorkspaceDir(path); err != nil {
		return fmt.Errorf("unable to ensure workspace directory - %w", err)
	}

	if err := WriteWorkspaceConfig(path, wsCfg); err != nil {
		return fmt.Errorf("unable to write workspace config - %w", err)
	}
	log.Trace().Msgf("workspace config initialized for %s", name)
	return nil
}
