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
		data, _ := userConfig.YAML()
		fmt.Println(data)
	case "json":
		data, _ := userConfig.JSON()
		fmt.Println(data)
	default:
		logger.Fatalf("Unsupported output format %s", format)
	}
}
