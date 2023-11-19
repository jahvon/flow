package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/rs/zerolog/log"
)

// ExpandDirectory expands the directory field of an executable to an absolute path.
// The following transformations are applied:
// - empty dir -> execPath
// - // -> wsPath + dir path
// - . -> current working directory + dir path
// - relative path -> execPath + dir path
// - ${envVar} -> expanded to the value from env map.
func ExpandDirectory(dir, wsPath, execPath string, env map[string]string) string {
	var targetDir string
	switch {
	case dir == "":
		targetDir = execPath
	case strings.HasPrefix(dir, "//"):
		targetDir = strings.Replace(dir, "//", wsPath+"/", 1)
	case dir == "." || strings.HasPrefix(dir, "./"):
		wd, err := os.Getwd()
		if err != nil {
			log.Warn().Err(err).Msg("unable to get working directory for relative path expansion")
			targetDir = filepath.Join(execPath, dir)
		} else {
			targetDir = filepath.Join(wd, dir[1:])
		}
	case strings.HasPrefix(targetDir, "~/"):
		homeDir, err := os.UserHomeDir()
		if err != nil {
			log.Warn().Err(err).Msg("unable to get user home directory for relative path expansion")
			targetDir = filepath.Join(execPath, dir)
		} else {
			targetDir = filepath.Join(homeDir, dir[2:])
		}
	default:
		targetDir = execPath + "/" + targetDir
	}

	targetDir = os.Expand(targetDir, func(key string) string {
		val, found := env[key]
		if !found {
			log.Warn().Str("key", key).Msg("unable to find env key in directory for expansion")
		}
		return val
	})
	return filepath.Clean(targetDir)
}

func ValidateOneOf(fieldName string, vals ...interface{}) error {
	var count int
	for _, val := range vals {
		if val != nil {
			count++
		}
	}
	if count == 0 {
		return fmt.Errorf("must define at least one %s", fieldName)
	} else if count > 1 {
		return fmt.Errorf("must define only one %s", fieldName)
	}
	return nil
}
