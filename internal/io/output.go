package io

import (
	"fmt"
)

func TypesDocsURL(docID, anchor string) string {
	if anchor != "" {
		anchor = "#" + anchor
	}
	return fmt.Sprintf("https://github.com/jahvon/flow/blob/main/docs/types/%s.md%s", docID, anchor)
}
