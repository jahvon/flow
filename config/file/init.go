package file

import (
	"os"
	"path/filepath"

	"github.com/pkg/errors"

	"github.com/jahvon/flow/config"
)

func InitUserConfig() error {
	if err := EnsureDefaultWorkspace(); err != nil {
		return errors.Wrap(err, "unable to ensure default workspace")
	}
	defaultWsName := "default"
	if _, err := os.Stat(filepath.Join(DefaultWorkspacePath, WorkspaceConfigFileName)); os.IsNotExist(err) {
		if err := InitWorkspaceConfig(defaultWsName, DefaultWorkspacePath); err != nil {
			return errors.Wrap(err, "unable to initialize default workspace config")
		}
	} else if err != nil {
		return errors.Wrap(err, "unable to check for default workspace config")
	}

	defaultCfg := &config.UserConfig{
		Workspaces:       map[string]string{defaultWsName: DefaultWorkspacePath},
		CurrentWorkspace: defaultWsName,
		WorkspaceMode:    config.WorkspaceModeDynamic,
		Interactive: &config.InteractiveConfig{
			Enabled: true,
		},
		DefaultLogMode: "logfmt",
	}

	_, err := os.Create(UserConfigPath)
	if err != nil {
		return errors.Wrap(err, "unable to create config file")
	}
	err = WriteUserConfig(defaultCfg)
	if err != nil {
		return errors.Wrap(err, "unable to write default config")
	}
	return nil
}

func InitWorkspaceConfig(name, path string) error {
	wsCfg := config.DefaultWorkspaceConfig(name)

	if err := EnsureWorkspaceDir(path); err != nil {
		return errors.Wrap(err, "unable to ensure workspace directory")
	}

	if err := WriteWorkspaceConfig(path, wsCfg); err != nil {
		return errors.Wrap(err, "unable to write workspace config")
	}
	return nil
}

func InitExecutables(
	template *config.ExecutableDefinitionTemplate,
	ws *config.WorkspaceConfig,
	name, subPath string,
) error {
	if err := EnsureExecutableDir(ws.Location(), subPath); err != nil {
		return errors.Wrap(err, "unable to ensure executable directory")
	}
	if err := RenderAndWriteExecutablesTemplate(template, ws, name, subPath); err != nil {
		return errors.Wrap(err, "unable to write executable definition template")
	}
	return nil
}
