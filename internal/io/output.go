package io

import (
	"fmt"
)

func ConfigDocsURL(docID, anchor string) string {
	if anchor != "" {
		anchor = "#" + anchor
	}
	return fmt.Sprintf("https://github.com/jahvon/flow/blob/main/docs/config/%s.md%s", docID, anchor)
}
