package executable

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/pterm/pterm"
	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v3"

	"github.com/jahvon/flow/internal/executable"
	"github.com/jahvon/flow/internal/io"
)

type executableListOutput struct {
	Executables []*executableOutput `json:"executables"`
}
type executableOutput struct {
	ID   string                 `json:"id"`
	Data *executable.Executable `json:"data"`
}

func PrintExecutableList(format io.OutputFormat, executables executable.List) {
	if format == io.OutputFormatYAML {
		printExecutableListYAML(executables)
	} else if format == io.OutputFormatJSON {
		printExecutableListJSON(executables, false)
	} else if format == io.OutputFormatPrettyJson {
		printExecutableListJSON(executables, true)
	} else if format == io.OutputFormatDefault {
		printExecutableListTable(executables)
	} else {
		log.Fatal().Msgf("Unsupported output format %s", format)
	}
}

func printExecutableListYAML(executables executable.List) {
	log.Info().Msgf("Printing %d executables", len(executables))
	enriched := &executableListOutput{Executables: make([]*executableOutput, 0)}
	for _, exec := range executables {
		enriched.Executables = append(enriched.Executables, &executableOutput{
			ID:   exec.ID(),
			Data: exec,
		})
	}
	yamlBytes, err := yaml.Marshal(enriched)
	if err != nil {
		log.Fatal().Msgf("Failed to marshal executable list - %v", err)
	}
	fmt.Println(string(yamlBytes))
}

func printExecutableListJSON(executables executable.List, pretty bool) {
	log.Info().Msgf("Printing %d executables", len(executables))
	enriched := &executableListOutput{Executables: make([]*executableOutput, 0)}
	for _, exec := range executables {
		enriched.Executables = append(enriched.Executables, &executableOutput{
			ID:   exec.ID(),
			Data: exec,
		})
	}

	var jsonBytes []byte
	var err error
	if pretty {
		jsonBytes, err = json.MarshalIndent(enriched, "", strings.Repeat(" ", 2))
	} else {
		jsonBytes, err = json.Marshal(enriched)
	}
	if err != nil {
		log.Fatal().Msgf("Failed to marshal executable - %v", err)
	}
	fmt.Println(string(jsonBytes))
}

func printExecutableListTable(executables executable.List) {
	log.Info().Msgf("Printing %d executables", len(executables))
	tableRows := pterm.TableData{{"ID", "Name", "Type", "Description", "Tags"}}
	for _, exec := range executables {
		tableRows = append(
			tableRows,
			[]string{exec.ID(), exec.Name, string(exec.Type), exec.Description, strings.Join(exec.Tags, ", ")},
		)
	}

	err := pterm.DefaultTable.WithHasHeader().WithBoxed().WithData(tableRows).Render()
	if err != nil {
		log.Fatal().Msgf("Failed to render executable list - %v", err)
	}
}

func PrintExecutable(format io.OutputFormat, exec *executable.Executable) {
	if format == io.OutputFormatYAML {
		printExecutableYAML(exec)
	} else if format == io.OutputFormatJSON {
		printExecutableJSON(exec, false)
	} else if format == io.OutputFormatPrettyJson {
		printExecutableJSON(exec, true)
	} else if format == io.OutputFormatDefault {
		printExecutableTable(exec)
	} else {
		log.Fatal().Msgf("Unsupported output format %s", format)
	}
}

func printExecutableJSON(exec *executable.Executable, pretty bool) {
	var jsonBytes []byte
	var err error
	enriched := &executableOutput{
		ID:   exec.ID(),
		Data: exec,
	}
	if pretty {
		jsonBytes, err = json.MarshalIndent(enriched, "", strings.Repeat(" ", 2))
	} else {
		jsonBytes, err = json.Marshal(enriched)
	}
	if err != nil {
		log.Fatal().Msgf("Failed to marshal executable - %v", err)
	}
	fmt.Println(string(jsonBytes))
}

func printExecutableYAML(exec *executable.Executable) {
	enriched := &executableOutput{
		ID:   exec.ID(),
		Data: exec,
	}
	yamlBytes, err := yaml.Marshal(enriched)
	if err != nil {
		log.Fatal().Msgf("Failed to marshal executable - %v", err)
	}
	fmt.Println(string(yamlBytes))
}

func printExecutableTable(exec *executable.Executable) {
	yamlSpec, err := yaml.Marshal(exec.Spec)
	if err != nil {
		log.Fatal().Msgf("Failed to marshal spec - %v", err)
	}
	err = pterm.DefaultTable.WithHasHeader().WithBoxed().WithData(pterm.TableData{
		{"Key", "Value"},
		{"ID", exec.ID()},
		{"Name", exec.Name},
		{"Type", string(exec.Type)},
		{"Description", exec.Description},
		{"Aliases", strings.Join(exec.Aliases, ", ")},
		{"Tags", strings.Join(exec.Tags, ", ")},
		{"Spec", string(yamlSpec)},
	}).Render()
	if err != nil {
		log.Fatal().Msgf("Failed to render executable - %v", err)
	}
}
