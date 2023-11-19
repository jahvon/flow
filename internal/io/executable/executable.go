package executable

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/pterm/pterm"
	"gopkg.in/yaml.v3"

	"github.com/jahvon/flow/config"
	"github.com/jahvon/flow/internal/io"
	"github.com/jahvon/flow/internal/io/utils"
)

var log = io.Log()

type executableListOutput struct {
	Executables []*executableOutput `json:"executables"`
}
type executableOutput struct {
	ID   string             `json:"id"`
	Data *config.Executable `json:"data"`
}

func PrintExecutableList(format io.OutputFormat, executables config.ExecutableList) {
	switch format {
	case io.OutputFormatYAML:
		printExecutableListYAML(executables)
	case io.OutputFormatJSON:
		printExecutableListJSON(executables, false)
	case io.OutputFormatPrettyJSON:
		printExecutableListJSON(executables, true)
	case io.OutputFormatDefault:
		printExecutableListTable(executables)
	default:
		log.Panic().Msgf("Unsupported output format %s", format)
	}
}

func printExecutableListYAML(executables config.ExecutableList) {
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
		log.Panic().Msgf("Failed to marshal executable list - %v", err)
	}
	fmt.Println(string(yamlBytes))
}

func printExecutableListJSON(executables config.ExecutableList, pretty bool) {
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
		log.Panic().Msgf("Failed to marshal executable - %v", err)
	}
	fmt.Println(string(jsonBytes))
}

func printExecutableListTable(executables config.ExecutableList) {
	log.Info().Msgf("Printing %d executables", len(executables))
	tableRows := pterm.TableData{{"ID", "Name", "Verb", "Description", "Tags"}}
	for _, exec := range executables {
		tableRows = append(
			tableRows,
			[]string{
				exec.ID(),
				exec.Name,
				string(exec.Verb),
				utils.WrapLines(exec.Description, 5),
				strings.Join(exec.Tags, ", "),
			},
		)
	}
	io.PrintTableWithHeader(tableRows)
}

func PrintExecutable(format io.OutputFormat, exec *config.Executable) {
	if exec == nil {
		log.Panic().Msg("Executable is nil")
	}

	switch format {
	case io.OutputFormatYAML:
		printExecutableYAML(exec)
	case io.OutputFormatJSON:
		printExecutableJSON(exec, false)
	case io.OutputFormatPrettyJSON:
		printExecutableJSON(exec, true)
	case io.OutputFormatDefault:
		printExecutableTable(exec)
	default:
		log.Panic().Msgf("Unsupported output format %s", format)
	}
}

func printExecutableJSON(exec *config.Executable, pretty bool) {
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
		log.Panic().Msgf("Failed to marshal executable - %v", err)
	}
	fmt.Println(string(jsonBytes))
}

func printExecutableYAML(exec *config.Executable) {
	enriched := &executableOutput{
		ID:   exec.ID(),
		Data: exec,
	}
	yamlBytes, err := yaml.Marshal(enriched)
	if err != nil {
		log.Panic().Msgf("Failed to marshal executable - %v", err)
	}
	fmt.Println(string(yamlBytes))
}

func printExecutableTable(exec *config.Executable) {
	yamlSpec, err := yaml.Marshal(exec.Type)
	if err != nil {
		log.Panic().Msgf("Failed to marshal spec - %v", err)
	}
	tableData := [][]string{
		{"Key", "Value"},
		{"ID", exec.ID()},
		{"Name", exec.Name},
		{"Verb", string(exec.Verb)},
		{"Description", utils.WrapLines(exec.Description, 10)},
		{"Aliases", strings.Join(exec.Aliases, ", ")},
		{"Tags", strings.Join(exec.Tags, ", ")},
		{"Type Spec", string(yamlSpec)},
	}
	io.PrintTableWithHeader(tableData)
}
