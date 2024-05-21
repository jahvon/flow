package file

import (
	"os"
	"path/filepath"

	"github.com/pkg/errors"
)

const (
	flowDirName = "flow"

	FlowConfigDirEnvVar = "FLOW_CONFIG_DIR"
	FlowCacheDirEnvVar  = "FLOW_CACHE_DIR"
)

func UserConfigFilePath() string {
	return ConfigDirPath() + "/config.yaml"
}

func DefaultWorkspaceDir() string {
	return CachedDataDirPath() + "/default"
}

func LogsDir() string {
	return CachedDataDirPath() + "/logs"
}

func ConfigDirPath() string {
	if dir := os.Getenv(FlowConfigDirEnvVar); dir != "" {
		return dir
	}

	dirname, err := os.UserConfigDir()
	if err != nil {
		panic(errors.Wrap(err, "unable to get config directory"))
	}
	return filepath.Join(dirname, flowDirName)
}

func CachedDataDirPath() string {
	if dir := os.Getenv(FlowCacheDirEnvVar); dir != "" {
		return dir
	}

	dirname, err := os.UserCacheDir()
	if err != nil {
		panic(errors.Wrap(err, "unable to get cache directory"))
	}
	return filepath.Join(dirname, flowDirName)
}

func LatestCachedDataDir() string {
	return CachedDataDirPath() + "/latestcache"
}

func LatestCachedDataFilePath(cacheKey string) string {
	return filepath.Join(LatestCachedDataDir(), cacheKey)
}

func EnsureConfigDir() error {
	if _, err := os.Stat(ConfigDirPath()); os.IsNotExist(err) {
		err = os.MkdirAll(ConfigDirPath(), 0750)
		if err != nil {
			return errors.Wrap(err, "unable to create config directory")
		}
	} else if err != nil {
		return errors.Wrap(err, "unable to check for config directory")
	}
	return nil
}

func EnsureCachedDataDir() error {
	if _, err := os.Stat(LatestCachedDataDir()); os.IsNotExist(err) {
		err = os.MkdirAll(LatestCachedDataDir(), 0750)
		if err != nil {
			return errors.Wrap(err, "unable to create cache directory")
		}
	} else if err != nil {
		return errors.Wrap(err, "unable to check for cache directory")
	}

	return nil
}

func EnsureDefaultWorkspace() error {
	if _, err := os.Stat(DefaultWorkspaceDir()); os.IsNotExist(err) {
		err = os.MkdirAll(DefaultWorkspaceDir(), 0750)
		if err != nil {
			return errors.Wrap(err, "unable to create default workspace directory")
		}
	} else if err != nil {
		return errors.Wrap(err, "unable to check for default workspace directory")
	}
	return nil
}

func EnsureWorkspaceDir(workspacePath string) error {
	if _, err := os.Stat(workspacePath); os.IsNotExist(err) {
		err = os.MkdirAll(workspacePath, 0750)
		if err != nil {
			return errors.Wrap(err, "unable to create workspace directory")
		}
	} else if err != nil {
		return errors.Wrap(err, "unable to check for workspace directory")
	}
	return nil
}

func EnsureWorkspaceConfig(workspaceName, workspacePath string) error {
	if _, err := os.Stat(filepath.Join(workspacePath, WorkspaceConfigFileName)); os.IsNotExist(err) {
		return InitWorkspaceConfig(workspaceName, workspacePath)
	} else if err != nil {
		return errors.Wrapf(err, "unable to check for workspace %s config file", workspaceName)
	}
	return nil
}

func EnsureExecutableDir(workspacePath, subPath string) error {
	if _, err := os.Stat(filepath.Join(workspacePath, subPath)); os.IsNotExist(err) {
		err = os.MkdirAll(filepath.Join(workspacePath, subPath), 0750)
		if err != nil {
			return errors.Wrap(err, "unable to create executable directory")
		}
	}
	return nil
}

func EnsureLogsDir() error {
	if _, err := os.Stat(LogsDir()); os.IsNotExist(err) {
		err = os.MkdirAll(LogsDir(), 0750)
		if err != nil {
			return errors.Wrap(err, "unable to create logs directory")
		}
	} else if err != nil {
		return errors.Wrap(err, "unable to check for logs directory")
	}
	return nil
}
