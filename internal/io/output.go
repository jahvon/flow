package io

import (
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/rs/zerolog"
)

type OutputFormat string

const (
	OutputFormatUnset      OutputFormat = ""
	OutputFormatJSON       OutputFormat = "json"
	OutputFormatPrettyJSON OutputFormat = "jsonp"
	OutputFormatYAML       OutputFormat = "yaml"
)

func Log() zerolog.Logger {
	writer := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: zerolog.TimeFormatUnix}
	context := zerolog.New(writer).With().Timestamp()
	return context.Logger()
}

func PrintQuestion(question string) {
	color.HiCyan(question)
}

type StdOutWriter struct {
	LogFields   []any
	Logger      *Logger
	AsPlainText bool
}

func (w StdOutWriter) Write(p []byte) (n int, err error) {
	trimmedP := strings.TrimSpace(string(p))
	if trimmedP == "" {
		return len(p), nil
	}
	splitP := strings.Split(trimmedP, "\n")
	for _, line := range splitP {
		if line == "---break" {
			w.AsPlainText = true
			continue
		} else if w.AsPlainText && line == "---endbreak" {
			w.AsPlainText = false
			continue
		}

		switch {
		case w.AsPlainText:
			w.Logger.AsPlainText(func() {
				w.Logger.Infof(line)
			})
		case len(w.LogFields) > 0:
			w.Logger.Infox(line, w.LogFields...)
		default:
			w.Logger.Infof(line)
		}
	}

	return len(p), nil
}

type StdErrWriter struct {
	LogFields   []any
	Logger      *Logger
	AsPlainText bool
}

func (w StdErrWriter) Write(p []byte) (n int, err error) {
	trimmedP := strings.TrimSpace(string(p))
	if trimmedP == "" {
		return len(p), nil
	}

	switch {
	case w.AsPlainText:
		w.Logger.AsPlainText(func() {
			w.Logger.Errorf(trimmedP)
		})
	case len(w.LogFields) > 0:
		w.Logger.Errorx(trimmedP, w.LogFields...)
	default:
		w.Logger.Errorf(trimmedP)
	}

	return len(p), nil
}

func DocsURL(docID string) string {
	return fmt.Sprintf("https://github.com/jahvon/flow/blob/main/docs/%s.md", docID)
}
