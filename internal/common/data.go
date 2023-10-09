package common

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/jahvon/flow/internal/io"
)

const (
	flowDirName = "flow"
)

var log = io.Log()

func ConfigDirPath() string {
	dirname, err := os.UserConfigDir()
	if err != nil {
		log.Fatal().Err(err).Msg("unable to get config directory")
	}
	return filepath.Join(dirname, flowDirName)
}

func EnsureConfigDir() error {
	if _, err := os.Stat(ConfigDirPath()); os.IsNotExist(err) {
		err = os.MkdirAll(ConfigDirPath(), 0755)
		if err != nil {
			return fmt.Errorf("unable to create config directory - %v", err)
		}
	} else if err != nil {
		return fmt.Errorf("unable to check for config directory - %v", err)
	}
	return nil
}

func CachedDataDirPath() string {
	dirname, err := os.UserCacheDir()
	if err != nil {
		log.Fatal().Err(err).Msg("unable to get cache directory")
	}
	return filepath.Join(dirname, flowDirName)
}

func EnsureCachedDataDir() error {
	if _, err := os.Stat(CachedDataDirPath()); os.IsNotExist(err) {
		err = os.MkdirAll(CachedDataDirPath(), 0755)
		if err != nil {
			return fmt.Errorf("unable to create cache directory - %v", err)
		}
	} else if err != nil {
		return fmt.Errorf("unable to check for cache directory - %v", err)
	}
	return nil
}
