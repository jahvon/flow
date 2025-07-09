//go:build darwin
// +build darwin

package open

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func open(uri string) *exec.Cmd {
	args := []string{uri}
	if os.Getenv(BackgroundEnvKey) != "" {
		args = append([]string{"-g"}, args...)
	}
	switch {
	case os.Getenv(DisabledEnvKey) != "":
		fmt.Printf("open %s\n", strings.Join(args, " "))
		return nil
	default:
		return exec.Command("open", args...)
	}
}

func openWith(uri string, appName string) *exec.Cmd {
	args := []string{"-a", appName, uri}
	if os.Getenv(BackgroundEnvKey) != "" {
		args = append([]string{"-g"}, args...)
	}

	switch {
	case os.Getenv(DisabledEnvKey) != "":
		fmt.Printf("open %s\n", strings.Join(args, " "))
		return nil
	default:
		return exec.Command("open", args...)
	}
}
