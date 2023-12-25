package executable

import (
	"fmt"

	"github.com/jahvon/flow/config"
	"github.com/jahvon/flow/internal/io"
)

var log = io.Log()

func PrintExecutableList(format io.OutputFormat, executables config.ExecutableList) {
	switch format {
	case io.OutputFormatYAML:
		str, err := executables.YAML()
		if err != nil {
			log.Panic().Msgf("Failed to marshal executable list - %v", err)
		}
		fmt.Println(str)
	case io.OutputFormatJSON:
		str, err := executables.JSON(false)
		if err != nil {
			log.Panic().Msgf("Failed to marshal executable list - %v", err)
		}
		fmt.Println(str)
	case io.OutputFormatPrettyJSON:
		str, err := executables.JSON(true)
		if err != nil {
			log.Panic().Msgf("Failed to marshal executable list - %v", err)
		}
		fmt.Println(str)
	case io.OutputFormatDefault:
		io.PrintTableData(executables.TableData())
	default:
		log.Panic().Msgf("Unsupported output format %s", format)
	}
}

func PrintExecutable(format io.OutputFormat, exec *config.Executable) {
	if exec == nil {
		log.Panic().Msg("Executable is nil")
	}

	switch format {
	case io.OutputFormatYAML:
		str, err := exec.YAML()
		if err != nil {
			log.Panic().Msgf("Failed to marshal executable - %v", err)
		}
		fmt.Println(str)
	case io.OutputFormatJSON:
		str, err := exec.JSON(false)
		if err != nil {
			log.Panic().Msgf("Failed to marshal executable - %v", err)
		}
		fmt.Println(str)
	case io.OutputFormatPrettyJSON:
		str, err := exec.JSON(true)
		if err != nil {
			log.Panic().Msgf("Failed to marshal executable - %v", err)
		}
		fmt.Println(str)
	case io.OutputFormatDefault:
		io.PrintMap(exec.Map())
	default:
		log.Panic().Msgf("Unsupported output format %s", format)
	}
}
