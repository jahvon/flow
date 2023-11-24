package config

import (
	"encoding/json"
	"fmt"

	"gopkg.in/yaml.v3"

	"github.com/jahvon/flow/config"
	"github.com/jahvon/flow/internal/io"
)

var log = io.Log()

func PrintUserConfig(format io.OutputFormat, userConfig *config.UserConfig) {
	if userConfig == nil {
		log.Panic().Msg("Config is nil")
	}

	switch format {
	case io.OutputFormatYAML:
		printUserConfigYAML(userConfig)
	case io.OutputFormatJSON:
		printUserConfigJSON(userConfig, false)
	case io.OutputFormatPrettyJSON:
		printUserConfigJSON(userConfig, true)
	case io.OutputFormatDefault:
		printUserConfigTable(userConfig)
	}
}

func printUserConfigYAML(rootCfg *config.UserConfig) {
	yamlBytes, err := yaml.Marshal(rootCfg)
	if err != nil {
		log.Panic().Msgf("Failed to marshal config - %v", err)
	}
	fmt.Println(string(yamlBytes))
}

func printUserConfigJSON(rootCfg *config.UserConfig, pretty bool) {
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

func printUserConfigTable(rootCfg *config.UserConfig) {
	workspacesYaml, err := yaml.Marshal(rootCfg.Workspaces)
	if err != nil {
		log.Panic().Msgf("Failed to marshal workspaces - %v", err)
	}

	tableRows := [][]string{
		{"Key", "Value"},
		{"Current Workspace", rootCfg.CurrentWorkspace},
		{"Current Namespace", rootCfg.CurrentNamespace},
		{"Workspaces", string(workspacesYaml)},
	}
	io.PrintTableWithHeader(tableRows)
}
