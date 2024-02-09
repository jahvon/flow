package workspace

import (
	"fmt"

	tuikitIO "github.com/jahvon/tuikit/io"

	"github.com/jahvon/flow/config"
	"github.com/jahvon/flow/internal/io"
)

func PrintWorkspaceList(logger *tuikitIO.Logger, format io.OutputFormat, workspaces config.WorkspaceConfigList) {
	logger.Infof("listing %d workspaces", len(workspaces))
	switch format {
	case io.OutputFormatDocument, io.OutputFormatYAML:
		str, err := workspaces.YAML()
		if err != nil {
			logger.Fatalf("Failed to marshal workspace list - %v", err)
		}
		logger.Println(str)
	case io.OutputFormatJSON:
		str, err := workspaces.JSON()
		if err != nil {
			logger.Fatalf("Failed to marshal workspace list - %v", err)
		}
		logger.Println(str)
	default:
		logger.Fatalf("Unsupported output format %s", format)
	}
}

func PrintWorkspaceConfig(logger *tuikitIO.Logger, format io.OutputFormat, ws *config.WorkspaceConfig) {
	if ws == nil {
		logger.Fatalf("Workspace config is nil")
	}
	logger.Infox(fmt.Sprintf("Workspace %s", ws.AssignedName()), "Location", ws.Location())
	switch format {
	case io.OutputFormatDocument, io.OutputFormatYAML:
		str, err := ws.YAML()
		if err != nil {
			logger.Fatalf("Failed to marshal workspace config - %v", err)
		}
		logger.Println(str)
	case io.OutputFormatJSON:
		str, err := ws.JSON()
		if err != nil {
			logger.Fatalf("Failed to marshal workspace config - %v", err)
		}
		logger.Println(str)
	default:
		logger.Fatalf("Unsupported output format %s", format)
	}
}
