package workspace

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func workspaceMarkdown(w *Workspace) string {
	var mkdwn string
	if w.DisplayName != "" {
		mkdwn = fmt.Sprintf("# [Workspace] %s\n", w.DisplayName)
	} else {
		mkdwn = fmt.Sprintf("# [Workspace] %s\n", w.AssignedName())
	}
	mkdwn += workspaceDescription(w)
	if len(w.Tags) > 0 {
		mkdwn += "**Tags**\n"
		for _, tag := range w.Tags {
			mkdwn += fmt.Sprintf("- %s\n", tag)
		}
	}
	if w.Executables != nil {
		mkdwn += "**Executable Filter**\n"
		if len(w.Executables.Included) > 0 {
			mkdwn += "Included\n"
			for _, line := range w.Executables.Included {
				mkdwn += fmt.Sprintf("  %s\n", line)
			}
		}
		if len(w.Executables.Excluded) > 0 {
			mkdwn += "Excluded\n"
			for _, line := range w.Executables.Excluded {
				mkdwn += fmt.Sprintf("  %s\n", line)
			}
		}
	}
	mkdwn += fmt.Sprintf("\n\n_Workspace can be found in_ [%s](%s)\n", w.Location(), w.Location())
	return mkdwn
}

func workspaceDescription(w *Workspace) string {
	var mkdwn string
	const descSpacer = "> \n"
	if w.Description != "" {
		mkdwn += descSpacer
		lines := strings.Split(w.Description, "\n")
		for _, line := range lines {
			mkdwn += fmt.Sprintf("> %s\n", line)
		}
		mkdwn += descSpacer
	}
	if w.DescriptionFile != "" {
		mdBytes, err := os.ReadFile(filepath.Clean(w.DescriptionFile))
		if err != nil {
			mkdwn += fmt.Sprintf("> **error rendering description file**: %s\n", err)
		} else {
			lines := strings.Split(string(mdBytes), "\n")
			for _, line := range lines {
				mkdwn += fmt.Sprintf("> %s\n", line)
			}
		}
		mkdwn += descSpacer
	}
	mkdwn += "\n"
	return mkdwn
}
