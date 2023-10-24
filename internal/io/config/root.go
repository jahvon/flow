package config

import (
	"encoding/json"
	"fmt"

	"gopkg.in/yaml.v3"

	"github.com/jahvon/flow/internal/config"
	"github.com/jahvon/flow/internal/io"
)

var log = io.Log()

func PrintRootConfig(format io.OutputFormat, rootCfg *config.RootConfig) {
	if rootCfg == nil {
		log.Panic().Msg("Config is nil")
	}

	switch format {
	case io.OutputFormatYAML:
		printRootConfigYAML(rootCfg)
	case io.OutputFormatJSON:
		printRootConfigJSON(rootCfg, false)
	case io.OutputFormatPrettyJSON:
		printRootConfigJSON(rootCfg, true)
	case io.OutputFormatDefault:
		printRootConfigTable(rootCfg)
	}
}

func printRootConfigYAML(rootCfg *config.RootConfig) {
	yamlBytes, err := yaml.Marshal(rootCfg)
	if err != nil {
		log.Panic().Msgf("Failed to marshal config - %v", err)
	}
	fmt.Println(string(yamlBytes))
}

func printRootConfigJSON(rootCfg *config.RootConfig, pretty bool) {
	var jsonBytes []byte
	var err error
	if pretty {
		jsonBytes, err = json.MarshalIndent(rootCfg, "", "  ")
	} else {
		jsonBytes, err = json.Marshal(rootCfg)
	}
	if err != nil {
		log.Panic().Msgf("Failed to marshal config - %v", err)
	}
	fmt.Println(string(jsonBytes))
}

func printRootConfigTable(rootCfg *config.RootConfig) {
	workspacesYaml, err := yaml.Marshal(rootCfg.Workspaces)
	if err != nil {
		log.Panic().Msgf("Failed to marshal workspaces - %v", err)
	}

	tableRows := [][]string{
		{"Key", "Value"},
		{"Current Workspace", rootCfg.CurrentWorkspace},
		{"Workspaces", string(workspacesYaml)},
	}
	io.PrintTableWithHeader(tableRows)
}
