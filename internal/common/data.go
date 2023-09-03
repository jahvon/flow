package common

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/jahvon/tbox/internal/io"
)

const dataDirName = ".tbox"

var log = io.Log()

func DataDirPath() string {
	dirname, err := os.UserHomeDir()
	if err != nil {
		log.Fatal().Err(err).Msg("unable to get home directory")
	}
	return filepath.Join(dirname, dataDirName)
}

func EnsureDataDir() error {
	if _, err := os.Stat(DataDirPath()); os.IsNotExist(err) {
		err = os.MkdirAll(DataDirPath(), 0755)
		if err != nil {
			return fmt.Errorf("unable to create data directory - %v", err)
		}
	} else if err != nil {
		return fmt.Errorf("unable to check for data directory - %v", err)
	}
	return nil
}
