package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/jahvon/tuikit/io"
	"github.com/pkg/errors"
)

// ExpandDirectory expands the directory field of an executable to an absolute path.
// The following transformations are applied:
// - empty dir -> execPath
// - // -> wsPath + dir path
// - . -> current working directory + dir path
// - relative path -> execPath + dir path
// - ${envVar} -> expanded to the value from env map.
func ExpandDirectory(logger io.Logger, dir, wsPath, execPath string, env map[string]string) string {
	execDir := filepath.Dir(execPath)
	var targetDir string
	switch {
	case dir == "":
		targetDir = execDir
	case strings.HasPrefix(dir, "//"):
		targetDir = strings.Replace(dir, "//", wsPath+"/", 1)
	case dir == "." || strings.HasPrefix(dir, "./"):
		wd, err := os.Getwd()
		if err != nil {
			logger.Warnx("unable to get working directory for relative path expansion", "err", err)
			targetDir = filepath.Join(execDir, dir)
		} else {
			targetDir = filepath.Join(wd, dir[1:])
		}
	case strings.HasPrefix(dir, "~/"):
		homeDir, err := os.UserHomeDir()
		if err != nil {
			logger.Warnx("unable to get user home directory for relative path expansion", "err", err)
			targetDir = filepath.Join(execDir, dir)
		} else {
			targetDir = filepath.Join(homeDir, dir[2:])
		}
	case strings.HasPrefix(dir, "/"):
		targetDir = dir
	default:
		targetDir = execDir + "/" + dir
	}

	targetDir = os.Expand(targetDir, func(key string) string {
		val, found := env[key]
		if !found {
			logger.Warnx("unable to find env key in directory for expansion", "key", key)
		}
		return val
	})
	return filepath.Clean(targetDir)
}

func PathFromWd(path string) (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return path, errors.Wrap(err, "unable to get working directory")
	}
	relPath, err := filepath.Rel(wd, path)
	if err != nil {
		return path, errors.Wrap(err, "unable to get relative path")
	}
	return relPath, nil
}

func ValidateOneOf(fieldName string, vals ...interface{}) error {
	var count int
	for _, val := range vals {
		if val == nil {
			continue
		}

		isPtr := reflect.ValueOf(val).Kind() == reflect.Ptr && !reflect.ValueOf(val).IsNil()
		isVal := reflect.ValueOf(val).Kind() != reflect.Ptr && !reflect.ValueOf(val).IsZero()
		if isPtr || isVal {
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

func IsZeroValue(v interface{}) bool {
	if v == nil {
		return true
	}
	return reflect.DeepEqual(v, reflect.Zero(reflect.TypeOf(v)).Interface())
}
