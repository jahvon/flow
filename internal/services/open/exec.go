//go:build !windows && !darwin
// +build !windows,!darwin

package open

import (
	"fmt"
	"os"
	"os/exec"
)

func open(uri string) *exec.Cmd {
	switch {
	case os.Getenv(DisabledEnvKey) != "":
		fmt.Println(fmt.Sprintf("xdg-open %s", uri))
		return nil
	default:
		return exec.Command("xdg-open", uri)
	}
}

func openWith(input string, appName string) *exec.Cmd {
	switch {
	case os.Getenv(DisabledEnvKey) != "":
		fmt.Println(fmt.Sprintf("%s %s", appName, input))
		return nil
	default:
		return exec.Command(appName, input)
	}
}
