package io

import (
	"fmt"
	"os"

	"github.com/jahvon/tuikit/components"
	"github.com/rs/zerolog"
)

type OutputFormat components.Format

const (
	OutputFormatDocument = OutputFormat(components.FormatDocument)
	OutputFormatJSON     = OutputFormat(components.FormatJSON)
	OutputFormatYAML     = OutputFormat(components.FormatYAML)
)

func Log() zerolog.Logger {
	writer := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: zerolog.TimeFormatUnix}
	context := zerolog.New(writer).With().Timestamp()
	return context.Logger()
}

func ConfigDocsURL(docID, anchor string) string {
	if anchor != "" {
		anchor = "#" + anchor
	}
	return fmt.Sprintf("https://github.com/jahvon/flow/blob/main/docs/config/%s.md%s", docID, anchor)
}
