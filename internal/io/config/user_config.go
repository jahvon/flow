package config

import (
	"fmt"

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
		fmt.Println(userConfig.YAML())
	case io.OutputFormatJSON:
		fmt.Println(userConfig.JSON(false))
	case io.OutputFormatPrettyJSON:
		fmt.Println(userConfig.JSON(true))
	case io.OutputFormatDefault:
		io.PrintMap(userConfig.Map())
	}
}
