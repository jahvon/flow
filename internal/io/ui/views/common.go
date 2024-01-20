package views

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/jahvon/flow/internal/io/ui/types"
	"github.com/jahvon/flow/internal/services/open"
)

func openInEditor(parent types.ParentView, path string) {
	preferred := os.Getenv("EDITOR")
	if preferred != "" {
		if err := open.OpenWith(preferred, path, false); err != nil {
			parent.HandleInternalError(fmt.Errorf("unable to open editor: %w", err))
		}
	} else {
		cmd := exec.Command("vim", path)
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		err := cmd.Run()
		if err != nil {
			parent.HandleInternalError(fmt.Errorf("unable to open vim: %w", err))
		}
	}
}
