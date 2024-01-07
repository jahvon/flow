package ui

import (
	"os"
	"os/exec"

	"github.com/jahvon/flow/internal/services/open"
)

func openInEditor(app *Application, path string) {
	preferred := os.Getenv("EDITOR")
	if preferred != "" {
		if err := open.OpenWith(preferred, path, false); err != nil {
			log.Err(err).Msg("unable to open editor")
			app.HandleInternalError(err)
		}
	} else {
		cmd := exec.Command("vim", path)
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		err := cmd.Run()
		if err != nil {
			log.Err(err).Msg("unable to open vim")
			app.HandleInternalError(err)
		}
	}
}
