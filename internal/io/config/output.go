package config

import (
	"fmt"
	"strings"

	tuikitIO "github.com/jahvon/tuikit/io"

	"github.com/jahvon/flow/types/config"
)

func PrintUserConfig(logger tuikitIO.Logger, format string, userConfig *config.Config) {
	if userConfig == nil {
		logger.Fatalf("FlowFile is nil")
	}

	switch strings.ToLower(format) {
	case "", "yaml", "yml":
		fmt.Println(userConfig.YAML())
	case "json":
		fmt.Println(userConfig.JSON())
	default:
		logger.Fatalf("Unsupported output format %s", format)
	}
}
