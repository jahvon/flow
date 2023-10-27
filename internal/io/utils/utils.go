package utils

import (
	"strings"

	"github.com/pterm/pterm"
)

var width = pterm.GetTerminalWidth()

// WrapLines Replace every n space with a newline character, leaving at most maxWords words per line
func WrapLines(text string, maxWords int) string {
	trimmed := strings.TrimSpace(text)
	words := strings.Split(trimmed, " ")
	var lines []string
	var line string
	for i, word := range words {
		if i%maxWords == 0 && i != 0 {
			lines = append(lines, line)
			line = ""
		}
		line += word + " "
	}
	lines = append(lines, line)
	return strings.Join(lines, "\n")
}
