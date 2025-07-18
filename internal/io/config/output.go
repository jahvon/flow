package config

import (
	tuikitIO "github.com/flowexec/tuikit/io"

	"github.com/flowexec/flow/internal/io/common"
	"github.com/flowexec/flow/types/config"
)

func PrintUserConfig(logger tuikitIO.Logger, format string, userConfig *config.Config) {
	if userConfig == nil {
		logger.Fatalf("Config type is nil")
	}

	switch common.NormalizeFormat(logger, format) {
	case common.YAMLFormat:
		str, err := userConfig.YAML()
		if err != nil {
			logger.Fatalf("Failed to marshal user config - %v", err)
		}
		logger.Println(str)
	case common.JSONFormat:
		str, err := userConfig.JSON()
		if err != nil {
			logger.Fatalf("Failed to marshal user config - %v", err)
		}
		logger.Println(str)
	default:
		logger.Fatalf("Unsupported output format %s", format)
	}
}
