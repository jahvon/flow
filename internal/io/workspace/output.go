package workspace

import (
	"fmt"
	"strings"

	tuikitIO "github.com/jahvon/tuikit/io"

	"github.com/jahvon/flow/types/workspace"
)

func PrintWorkspaceList(logger tuikitIO.Logger, format string, workspaces workspace.WorkspaceList) {
	logger.Infof("listing %d workspaces", len(workspaces))
	switch strings.ToLower(format) {
	case "", "yaml", "yml":
		str, err := workspaces.YAML()
		if err != nil {
			logger.Fatalf("Failed to marshal workspace list - %v", err)
		}
		logger.Println(str)
	case "json":
		str, err := workspaces.JSON()
		if err != nil {
			logger.Fatalf("Failed to marshal workspace list - %v", err)
		}
		logger.Println(str)
	default:
		logger.Fatalf("Unsupported output format %s", format)
	}
}

func PrintWorkspaceConfig(logger tuikitIO.Logger, format string, ws *workspace.Workspace) {
	if ws == nil {
		logger.Fatalf("Workspace config is nil")
	}
	logger.Infox(fmt.Sprintf("Workspace %s", ws.AssignedName()), "Location", ws.Location())
	switch strings.ToLower(format) {
	case "", "yaml", "yml":
		str, err := ws.YAML()
		if err != nil {
			logger.Fatalf("Failed to marshal workspace config - %v", err)
		}
		logger.Println(str)
	case "json":
		str, err := ws.JSON()
		if err != nil {
			logger.Fatalf("Failed to marshal workspace config - %v", err)
		}
		logger.Println(str)
	default:
		logger.Fatalf("Unsupported output format %s", format)
	}
}
