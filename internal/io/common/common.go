package common

import (
	"os"
	"os/exec"
	"slices"

	"github.com/jahvon/flow/internal/services/open"
)

const HeaderContextKey = "ctx"

var termEditors = []string{"vim", "nvim", "emacs", "nano"}

func OpenInEditor(path string, stdIn, stdOut *os.File) error {
	preferred := os.Getenv("EDITOR")
	if preferred != "" && !slices.Contains(termEditors, preferred) {
		return open.OpenWith(preferred, path, false)
	}
	if preferred == "" {
		preferred = "vim"
	}
	cmd := exec.Command(preferred, path) // #nosec G204
	cmd.Stdin = stdIn
	cmd.Stdout = stdOut
	return cmd.Run()
}
