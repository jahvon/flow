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
	mkdwn += workspaceDescription(w, true)
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

func workspaceDescription(w *Workspace, withPrefix bool) string {
	if w.Description == "" && w.DescriptionFile == "" {
		return ""
	}
	var mkdwn string

	prefix := ""
	if withPrefix {
		prefix = "> "
	}
	if d := strings.TrimSpace(w.Description); d != "" {
		mkdwn += prefix + "\n"
		mkdwn += addPrefix(d, prefix)
	}
	if w.DescriptionFile != "" {
		wsFile := filepath.Join(w.Location(), w.DescriptionFile)
		mdBytes, err := os.ReadFile(filepath.Clean(wsFile))
		if err != nil {
			mkdwn += addPrefix(fmt.Sprintf("**error rendering description file**: %s", err), prefix)
		} else if d := strings.TrimSpace(string(mdBytes)); d != "" {
			mkdwn += prefix + "\n"
			mkdwn += addPrefix(d, prefix)
		}
	}
	if mkdwn != "" {
		mkdwn += prefix + "\n"
	}
	mkdwn += "\n"
	return mkdwn
}

func addPrefix(s, prefix string) string {
	lines := strings.Split(s, "\n")
	var final string
	for _, line := range lines {
		final += prefix + line + "\n"
	}
	return final
}
