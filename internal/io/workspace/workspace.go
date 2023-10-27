package workspace

import (
	"encoding/json"
	"fmt"
	"strings"

	"gopkg.in/yaml.v3"

	"github.com/jahvon/flow/internal/io"
	"github.com/jahvon/flow/internal/io/utils"
	"github.com/jahvon/flow/internal/workspace"
)

var log = io.Log()

func PrintWorkspaceList(format io.OutputFormat, workspaces []workspace.Config) {
	switch format {
	case io.OutputFormatYAML:
		printWorkspaceListYAML(workspaces)
	case io.OutputFormatJSON:
		printWorkspaceListJSON(workspaces, false)
	case io.OutputFormatPrettyJSON:
		printWorkspaceListJSON(workspaces, true)
	case io.OutputFormatDefault:
		printWorkspaceListTable(workspaces)
	default:
		log.Panic().Msgf("Unsupported output format %s", format)
	}
}

func printWorkspaceListYAML(workspaces []workspace.Config) {
	log.Info().Msgf("Printing %d workspaces", len(workspaces))
	yamlBytes, err := yaml.Marshal(workspaces)
	if err != nil {
		log.Panic().Msgf("Failed to marshal workspace list - %v", err)
	}
	fmt.Println(string(yamlBytes))
}

func printWorkspaceListJSON(workspaces []workspace.Config, pretty bool) {
	log.Info().Msgf("Printing %d workspaces", len(workspaces))
	var jsonBytes []byte
	var err error
	if pretty {
		jsonBytes, err = json.MarshalIndent(workspaces, "", "  ")
	} else {
		jsonBytes, err = json.Marshal(workspaces)
	}
	if err != nil {
		log.Panic().Msgf("Failed to marshal workspace list - %v", err)
	}
	fmt.Println(string(jsonBytes))
}

func printWorkspaceListTable(workspaces []workspace.Config) {
	log.Info().Msgf("Printing %d workspaces", len(workspaces))
	tableRows := [][]string{{"Name", "Location", "Description", "Tags"}}
	for _, ws := range workspaces {
		name := ws.DisplayName
		if name == "" {
			log.Debug().Msg("Workspace config has no display name, using assigned name")
			name = ws.AssignedName()
		}
		tableRows = append(
			tableRows,
			[]string{
				name,
				ws.Location(),
				utils.WrapLines(ws.Description, 5),
				strings.Join(ws.Tags, ", "),
			},
		)
	}
	io.PrintTableWithHeader(tableRows)
}

func PrintWorkspaceConfig(format io.OutputFormat, ws *workspace.Config) {
	if ws == nil {
		log.Panic().Msg("Workspace config is nil")
	}

	switch format {
	case io.OutputFormatYAML:
		printWorkspaceConfigsYAML(ws)
	case io.OutputFormatJSON:
		printWorkspaceConfigJSON(ws, false)
	case io.OutputFormatPrettyJSON:
		printWorkspaceConfigJSON(ws, true)
	case io.OutputFormatDefault:
		printWorkspaceConfigTable(ws)
	default:
		log.Panic().Msgf("Unsupported output format %s", format)
	}
}

func printWorkspaceConfigsYAML(ws *workspace.Config) {
	log.Info().Msgf("Printing workspace config for %s", ws.DisplayName)
	yamlBytes, err := yaml.Marshal(ws)
	if err != nil {
		log.Panic().Msgf("Failed to marshal workspace config - %v", err)
	}
	fmt.Println(string(yamlBytes))
}

func printWorkspaceConfigJSON(ws *workspace.Config, pretty bool) {
	log.Info().Msgf("Printing workspace config for %s", ws.DisplayName)
	var jsonBytes []byte
	var err error
	if pretty {
		jsonBytes, err = json.MarshalIndent(ws, "", "  ")
	} else {
		jsonBytes, err = json.Marshal(ws)
	}
	if err != nil {
		log.Panic().Msgf("Failed to marshal workspace config - %v", err)
	}
	fmt.Println(string(jsonBytes))
}

func printWorkspaceConfigTable(ws *workspace.Config) {
	tableRows := [][]string{
		{"Key", "Value"},
		{"Name", ws.AssignedName()},
		{"Location", ws.Location()},
	}
	if ws.AssignedName() != ws.DisplayName && ws.DisplayName != "" {
		tableRows = append(tableRows, []string{"Display Name", ws.DisplayName})
	}
	if ws.Description != "" {
		tableRows = append(tableRows, []string{"Description", utils.WrapLines(ws.Description, 10)})
	}
	if ws.Tags != nil {
		tableRows = append(tableRows, []string{"Tags", strings.Join(ws.Tags, ", ")})
	}
	if ws.Git != nil {
		gitConfig, err := yaml.Marshal(ws.Git)
		if err != nil {
			log.Error().Err(err).Msg("failed to marshal git config")
		} else {
			tableRows = append(tableRows, []string{"Git Config", string(gitConfig)})
		}
	}
	if ws.Executables != nil {
		execs, err := yaml.Marshal(ws.Executables)
		if err != nil {
			log.Error().Err(err).Msg("failed to marshal executables")
		} else {
			tableRows = append(tableRows, []string{"Executables", string(execs)})
		}
	}
	io.PrintTableWithHeader(tableRows)
}
