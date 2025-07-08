package executable

import (
	tuikitIO "github.com/flowexec/tuikit/io"

	"github.com/jahvon/flow/internal/io/common"
	"github.com/jahvon/flow/types/executable"
)

func PrintExecutableList(logger tuikitIO.Logger, format string, executables executable.ExecutableList) {
	logger.Debugf("listing %d executables", len(executables))
	switch common.NormalizeFormat(logger, format) {
	case common.YAMLFormat:
		str, err := executables.YAML()
		if err != nil {
			logger.Fatalf("Failed to marshal executable list - %v", err)
		}
		logger.Println(str)
	case common.JSONFormat:
		str, err := executables.JSON()
		if err != nil {
			logger.Fatalf("Failed to marshal executable list - %v", err)
		}
		logger.Println(str)
	default:
		logger.Fatalf("Unsupported output format %s", format)
	}
}

func PrintExecutable(logger tuikitIO.Logger, format string, exec *executable.Executable) {
	if exec == nil {
		logger.Fatalf("Executable type is nil")
	}
	switch common.NormalizeFormat(logger, format) {
	case common.YAMLFormat:
		str, err := exec.YAML()
		if err != nil {
			logger.Fatalf("Failed to marshal executable - %v", err)
		}
		logger.Println(str)
	case common.JSONFormat:
		str, err := exec.JSON()
		if err != nil {
			logger.Fatalf("Failed to marshal executable - %v", err)
		}
		logger.Println(str)
	}
}

func PrintTemplate(logger tuikitIO.Logger, format string, template *executable.Template) {
	if template == nil {
		logger.Fatalf("Template type is nil")
	}
	logger.Debugf("Template %s", template.Name())
	switch common.NormalizeFormat(logger, format) {
	case common.YAMLFormat:
		str, err := template.YAML()
		if err != nil {
			logger.Fatalf("Failed to marshal template - %v", err)
		}
		logger.Println(str)
	case common.JSONFormat:
		str, err := template.JSON()
		if err != nil {
			logger.Fatalf("Failed to marshal template - %v", err)
		}
		logger.Println(str)
	}
}

func PrintTemplateList(logger tuikitIO.Logger, format string, templates executable.TemplateList) {
	logger.Debugf("listing %d templates", len(templates))
	switch common.NormalizeFormat(logger, format) {
	case common.YAMLFormat:
		str, err := templates.YAML()
		if err != nil {
			logger.Fatalf("Failed to marshal template list - %v", err)
		}
		logger.Println(str)
	case common.JSONFormat:
		str, err := templates.JSON()
		if err != nil {
			logger.Fatalf("Failed to marshal template list - %v", err)
		}
		logger.Println(str)
	}
}
