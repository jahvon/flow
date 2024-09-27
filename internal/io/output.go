package io

import (
	"fmt"
	"strings"
)

func TypesDocsURL(docID, anchor string) string {
	if anchor != "" {
		anchor = "#" + anchor
	} else {
		anchor = "?id=" + strings.ToLower(anchor)
	}
	return fmt.Sprintf("https://flowexec.io/#/types/%s%s", docID, anchor)
}
