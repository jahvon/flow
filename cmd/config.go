package cmd

import "github.com/jahvon/tbox/internal/config"

var currentConfig *config.RootConfig

func init() {
	currentConfig = config.LoadConfig()
}
