package config

import (
	"fmt"

	tuikitIO "github.com/jahvon/tuikit/io"

	"github.com/jahvon/flow/config"
	"github.com/jahvon/flow/internal/io"
)

func PrintUserConfig(logger *tuikitIO.Logger, format io.OutputFormat, userConfig *config.UserConfig) {
	if userConfig == nil {
		logger.Fatalf("Config is nil")
	}

	switch format {
	case io.OutputFormatDocument, io.OutputFormatYAML:
		fmt.Println(userConfig.YAML())
	case io.OutputFormatJSON:
		fmt.Println(userConfig.JSON())
	default:
		logger.Fatalf("Unsupported output format %s", format)
	}
}
