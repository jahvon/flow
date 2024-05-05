package common

import (
	"fmt"
	"os"
	"os/exec"
	"slices"

	"github.com/jahvon/tuikit/components"

	"github.com/jahvon/flow/internal/services/open"
)

var termEditors = []string{"vim", "nvim", "emacs", "nano"}

func OpenInEditor(path string) error {
	preferred := os.Getenv("EDITOR")
	if preferred != "" && !slices.Contains(termEditors, preferred) {
		return open.OpenWith(preferred, path, false)
	}
	if preferred == "" {
		preferred = "vim"
	}
	cmd := exec.Command(preferred, path) // #nosec G204
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	return cmd.Run()
}

func DeprecatedOpenInEditor(container *components.ContainerView, path string) {
	preferred := os.Getenv("EDITOR")
	if preferred != "" && !slices.Contains(termEditors, preferred) {
		if err := open.OpenWith(preferred, path, false); err != nil {
			container.HandleError(fmt.Errorf("unable to open editor: %w", err))
		}
	} else {
		if preferred == "" {
			preferred = "vim"
		}
		cmd := exec.Command(preferred, path) // #nosec G204
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		err := cmd.Run()
		if err != nil {
			container.HandleError(fmt.Errorf("unable to open %s: %w", preferred, err))
		}
	}
}
