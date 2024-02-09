package executable

import (
	"fmt"

	tuikitIO "github.com/jahvon/tuikit/io"

	"github.com/jahvon/flow/config"
	"github.com/jahvon/flow/internal/io"
)

func PrintExecutableList(logger *tuikitIO.Logger, format io.OutputFormat, executables config.ExecutableList) {
	logger.Infof("listing %d executables", len(executables))
	switch format {
	case io.OutputFormatDocument, io.OutputFormatYAML:
		str, err := executables.YAML()
		if err != nil {
			logger.Fatalf("Failed to marshal executable list - %v", err)
		}
		logger.Println(str)
	case io.OutputFormatJSON:
		str, err := executables.JSON()
		if err != nil {
			logger.Fatalf("Failed to marshal executable list - %v", err)
		}
		logger.Println(str)
	default:
		logger.Fatalf("Unsupported output format %s", format)
	}
}

func PrintExecutable(logger *tuikitIO.Logger, format io.OutputFormat, exec *config.Executable) {
	if exec == nil {
		logger.Fatalf("Executable is nil")
	}
	logger.Infox(fmt.Sprintf("Executable %s", exec.ID()), "Location", exec.DefinitionPath())
	switch format {
	case io.OutputFormatDocument, io.OutputFormatYAML:
		str, err := exec.YAML()
		if err != nil {
			logger.Fatalf("Failed to marshal executable - %v", err)
		}
		logger.Println(str)
	case io.OutputFormatJSON:
		str, err := exec.JSON()
		if err != nil {
			logger.Fatalf("Failed to marshal executable - %v", err)
		}
		logger.Println(str)
	default:
		logger.Fatalf("Unsupported output format %s", format)
	}
}
