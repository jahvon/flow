package executable

import (
	"fmt"
	"strings"

	tuikitIO "github.com/jahvon/tuikit/io"

	"github.com/jahvon/flow/types/executable"
)

const (
	yamlFormat = "yaml"
	ymlFormat  = "yml"
	jsonFormat = "json"
)

func PrintExecutableList(logger tuikitIO.Logger, format string, executables executable.ExecutableList) {
	logger.Infof("listing %d executables", len(executables))
	switch strings.ToLower(format) {
	case "", yamlFormat, ymlFormat:
		str, err := executables.YAML()
		if err != nil {
			logger.Fatalf("Failed to marshal executable list - %v", err)
		}
		logger.Println(str)
	case jsonFormat:
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
		logger.Fatalf("Executable is nil")
	}
	logger.Infox(fmt.Sprintf("Executable %s", exec.ID()), "Location", exec.FlowFilePath())
	switch strings.ToLower(format) {
	case "", yamlFormat, ymlFormat:
		str, err := exec.YAML()
		if err != nil {
			logger.Fatalf("Failed to marshal executable - %v", err)
		}
		logger.Println(str)
	case jsonFormat:
		str, err := exec.JSON()
		if err != nil {
			logger.Fatalf("Failed to marshal executable - %v", err)
		}
		logger.Println(str)
	default:
		logger.Fatalf("Unsupported output format %s", format)
	}
}

func PrintTemplate(logger tuikitIO.Logger, format string, template *executable.Template) {
	if template == nil {
		logger.Fatalf("Template is nil")
	}
	logger.Infof("Template %s", template.Name())
	switch strings.ToLower(format) {
	case "", yamlFormat, ymlFormat:
		str, err := template.YAML()
		if err != nil {
			logger.Fatalf("Failed to marshal template - %v", err)
		}
		logger.Println(str)
	case jsonFormat:
		str, err := template.JSON()
		if err != nil {
			logger.Fatalf("Failed to marshal template - %v", err)
		}
		logger.Println(str)
	default:
		logger.Fatalf("Unsupported output format %s", format)
	}
}

func PrintTemplateList(logger tuikitIO.Logger, format string, templates executable.TemplateList) {
	logger.Infof("listing %d templates", len(templates))
	switch strings.ToLower(format) {
	case "", yamlFormat, ymlFormat:
		str, err := templates.YAML()
		if err != nil {
			logger.Fatalf("Failed to marshal template list - %v", err)
		}
		logger.Println(str)
	case jsonFormat:
		str, err := templates.JSON()
		if err != nil {
			logger.Fatalf("Failed to marshal template list - %v", err)
		}
		logger.Println(str)
	default:
		logger.Fatalf("Unsupported output format %s", format)
	}
}
