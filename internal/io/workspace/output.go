package workspace

import (
	tuikitIO "github.com/flowexec/tuikit/io"

	"github.com/flowexec/flow/internal/io/common"
	"github.com/flowexec/flow/types/workspace"
)

func PrintWorkspaceList(logger tuikitIO.Logger, format string, workspaces workspace.WorkspaceList) {
	logger.Debugf("listing %d workspaces", len(workspaces))
	switch common.NormalizeFormat(logger, format) {
	case common.YAMLFormat:
		str, err := workspaces.YAML()
		if err != nil {
			logger.Fatalf("Failed to marshal workspace list - %v", err)
		}
		logger.Println(str)
	case common.JSONFormat:
		str, err := workspaces.JSON()
		if err != nil {
			logger.Fatalf("Failed to marshal workspace list - %v", err)
		}
		logger.Println(str)
	}
}

func PrintWorkspaceConfig(logger tuikitIO.Logger, format string, ws *workspace.Workspace) {
	if ws == nil {
		logger.Fatalf("Workspace type is nil")
	}
	switch common.NormalizeFormat(logger, format) {
	case common.YAMLFormat:
		str, err := ws.YAML()
		if err != nil {
			logger.Fatalf("Failed to marshal workspace config - %v", err)
		}
		logger.Println(str)
	case common.JSONFormat:
		str, err := ws.JSON()
		if err != nil {
			logger.Fatalf("Failed to marshal workspace config - %v", err)
		}
		logger.Println(str)
	}
}
