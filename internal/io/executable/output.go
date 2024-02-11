package executable

import (
	"fmt"
	"strings"

	tuikitIO "github.com/jahvon/tuikit/io"

	"github.com/jahvon/flow/config"
)

func PrintExecutableList(logger *tuikitIO.Logger, format string, executables config.ExecutableList) {
	logger.Infof("listing %d executables", len(executables))
	switch strings.ToLower(format) {
	case "", "yaml", "yml":
		str, err := executables.YAML()
		if err != nil {
			logger.Fatalf("Failed to marshal executable list - %v", err)
		}
		logger.Println(str)
	case "json":
		str, err := executables.JSON()
		if err != nil {
			logger.Fatalf("Failed to marshal executable list - %v", err)
		}
		logger.Println(str)
	default:
		logger.Fatalf("Unsupported output format %s", format)
	}
}

func PrintExecutable(logger *tuikitIO.Logger, format string, exec *config.Executable) {
	if exec == nil {
		logger.Fatalf("Executable is nil")
	}
	logger.Infox(fmt.Sprintf("Executable %s", exec.ID()), "Location", exec.DefinitionPath())
	switch strings.ToLower(format) {
	case "", "yaml", "yml":
		str, err := exec.YAML()
		if err != nil {
			logger.Fatalf("Failed to marshal executable - %v", err)
		}
		logger.Println(str)
	case "json":
		str, err := exec.JSON()
		if err != nil {
			logger.Fatalf("Failed to marshal executable - %v", err)
		}
		logger.Println(str)
	default:
		logger.Fatalf("Unsupported output format %s", format)
	}
}
