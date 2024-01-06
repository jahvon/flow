package workspace

import (
	"fmt"

	"github.com/jahvon/flow/config"
	"github.com/jahvon/flow/internal/io"
)

var log = io.Log()

func PrintWorkspaceList(format io.OutputFormat, workspaces config.WorkspaceConfigList) {
	log.Info().Msgf("%d workspaces", len(workspaces))
	switch format {
	case io.OutputFormatUnset, io.OutputFormatYAML:
		str, err := workspaces.YAML()
		if err != nil {
			log.Panic().Msgf("Failed to marshal workspace list - %v", err)
		}
		fmt.Println(str)
	case io.OutputFormatJSON:
		str, err := workspaces.JSON(false)
		if err != nil {
			log.Panic().Msgf("Failed to marshal workspace list - %v", err)
		}
		fmt.Println(str)
	case io.OutputFormatPrettyJSON:
		str, err := workspaces.JSON(true)
		if err != nil {
			log.Panic().Msgf("Failed to marshal workspace list - %v", err)
		}
		fmt.Println(str)
	default:
		log.Panic().Msgf("Unsupported output format %s", format)
	}
}

func PrintWorkspaceConfig(format io.OutputFormat, ws *config.WorkspaceConfig) {
	if ws == nil {
		log.Panic().Msg("Workspace config is nil")
	}
	log.Info().Str("Location", ws.Location()).Msgf("Workspace %s", ws.AssignedName())
	switch format {
	case io.OutputFormatUnset, io.OutputFormatYAML:
		str, err := ws.YAML()
		if err != nil {
			log.Panic().Msgf("Failed to marshal workspace config - %v", err)
		}
		fmt.Println(str)
	case io.OutputFormatJSON:
		str, err := ws.JSON(false)
		if err != nil {
			log.Panic().Msgf("Failed to marshal workspace config - %v", err)
		}
		fmt.Println(str)
	case io.OutputFormatPrettyJSON:
		str, err := ws.JSON(true)
		if err != nil {
			log.Panic().Msgf("Failed to marshal workspace config - %v", err)
		}
		fmt.Println(str)
	default:
		log.Panic().Msgf("Unsupported output format %s", format)
	}
}
