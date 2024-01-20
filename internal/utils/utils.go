package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
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
		targetDir = filepath.Dir(execPath)
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

func PathFromWd(path string) string {
	wd, err := os.Getwd()
	if err != nil {
		log.Warn().Err(err).Msg("unable to get working directory for relative path")
		return path
	}
	relPath, err := filepath.Rel(wd, path)
	if err != nil {
		log.Warn().Err(err).Msg("unable to get relative path")
		return path
	}
	return relPath
}

func ValidateOneOf(fieldName string, vals ...interface{}) error {
	var count int
	for _, val := range vals {
		if val != nil && (reflect.ValueOf(val).Kind() == reflect.Ptr && !reflect.ValueOf(val).IsNil()) {
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

func IsMultiLine(s string) bool {
	return NumLines(s) > 1
}

func NumLines(s string) int {
	return strings.Count(s, "\n") + 1
}

// WrapLines Replace every n space with a newline character, leaving at most maxWords words per line.
func WrapLines(text string, maxWords int) string {
	trimmed := strings.TrimSpace(text)
	words := strings.Split(trimmed, " ")
	var lines []string
	var line string
	for i, word := range words {
		if i%maxWords == 0 && i != 0 {
			lines = append(lines, line)
			line = ""
		}
		line += word + " "
	}
	lines = append(lines, line)
	return strings.Join(lines, "\n")
}

// ShortenString shortens a string to maxLen characters, appending "..." if the string is longer than maxLen.
func ShortenString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}
