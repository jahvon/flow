package file

import (
	"fmt"
	"os"
	"path/filepath"
)

const (
	flowDirName = "flow"
)

var (
	UserConfigPath       = ConfigDirPath() + "/config.yaml"
	DefaultWorkspacePath = CachedDataDirPath() + "/default"
	LatestCacheDataPath  = CachedDataDirPath() + "/latestcache"
)

func ConfigDirPath() string {
	dirname, err := os.UserConfigDir()
	if err != nil {
		log.Panic().Err(err).Msg("unable to get config directory")
	}
	return filepath.Join(dirname, flowDirName)
}

func CachedDataDirPath() string {
	dirname, err := os.UserCacheDir()
	if err != nil {
		log.Panic().Err(err).Msg("unable to get cache directory")
	}
	return filepath.Join(dirname, flowDirName)
}

func LatestCachedDataFilePath(cacheKey string) string {
	return filepath.Join(LatestCacheDataPath, cacheKey)
}

func EnsureConfigDir() error {
	if _, err := os.Stat(ConfigDirPath()); os.IsNotExist(err) {
		err = os.MkdirAll(ConfigDirPath(), 0750)
		if err != nil {
			return fmt.Errorf("unable to create config directory - %w", err)
		}
		log.Info().Msgf("created config directory at %s", ConfigDirPath())
	} else if err != nil {
		return fmt.Errorf("unable to check for config directory - %w", err)
	}
	return nil
}

func EnsureCachedDataDir() error {
	if _, err := os.Stat(LatestCacheDataPath); os.IsNotExist(err) {
		err = os.MkdirAll(LatestCacheDataPath, 0750)
		if err != nil {
			return fmt.Errorf("unable to create cache directory - %w", err)
		}
		log.Info().Msgf("created cache directory at %s", LatestCacheDataPath)
	} else if err != nil {
		return fmt.Errorf("unable to check for cache directory - %w", err)
	}

	return nil
}

func EnsureDefaultWorkspace() error {
	if _, err := os.Stat(DefaultWorkspacePath); os.IsNotExist(err) {
		err = os.MkdirAll(DefaultWorkspacePath, 0750)
		if err != nil {
			return fmt.Errorf("unable to create default workspace directory - %w", err)
		}
		log.Info().Msgf("created default workspace directory at %s", DefaultWorkspacePath)
	} else if err != nil {
		return fmt.Errorf("unable to check for default workspace directory - %w", err)
	}
	return nil
}

func EnsureWorkspaceDir(workspacePath string) error {
	if _, err := os.Stat(workspacePath); os.IsNotExist(err) {
		err = os.MkdirAll(workspacePath, 0750)
		if err != nil {
			return fmt.Errorf("unable to create workspace directory - %w", err)
		}
		log.Info().Msgf("created workspace directory at %s", workspacePath)
	} else if err != nil {
		return fmt.Errorf("unable to check for workspace directory - %w", err)
	}
	return nil
}

func EnsureWorkspaceConfig(workspaceName, workspacePath string) error {
	if _, err := os.Stat(filepath.Join(workspacePath, WorkspaceConfigFileName)); os.IsNotExist(err) {
		return InitWorkspaceConfig(workspaceName, workspacePath)
	} else if err != nil {
		return fmt.Errorf("unable to check for workspace %s  config directory - %w", workspaceName, err)
	}
	return nil
}
