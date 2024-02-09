package io

import (
	"fmt"

	"github.com/jahvon/tuikit/components"
)

type OutputFormat components.Format

const (
	OutputFormatDocument = OutputFormat(components.FormatDocument)
	OutputFormatJSON     = OutputFormat(components.FormatJSON)
	OutputFormatYAML     = OutputFormat(components.FormatYAML)
)

func ConfigDocsURL(docID, anchor string) string {
	if anchor != "" {
		anchor = "#" + anchor
	}
	return fmt.Sprintf("https://github.com/jahvon/flow/blob/main/docs/config/%s.md%s", docID, anchor)
}
