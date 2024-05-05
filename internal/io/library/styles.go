package library

import (
	"fmt"
	"math"

	"github.com/charmbracelet/lipgloss"
	"github.com/jahvon/tuikit/styles"
	"github.com/mattn/go-runewidth"
)

func renderSelection(s string, theme styles.Theme) string {
	style := lipgloss.NewStyle().Foreground(theme.PrimaryColor)
	return style.Render(s)
}

func renderSecondarySelection(s string, theme styles.Theme) string {
	style := lipgloss.NewStyle().Foreground(theme.TertiaryColor)
	return style.Render(s)
}

func renderInactive(s string, theme styles.Theme) string {
	style := lipgloss.NewStyle().Foreground(theme.Gray)
	return style.Render(s)
}

func renderDescription(s string, theme styles.Theme) string {
	style := lipgloss.NewStyle().Foreground(theme.BodyColor)
	return style.Render(s)
}

func renderPaneTitle(s string, count int, active bool, theme styles.Theme) string {
	var title string
	if count == 0 {
		title = s
	} else {
		title = fmt.Sprintf("%s (%d)", s, count)
	}
	style := lipgloss.NewStyle().Foreground(theme.SecondaryColor).Padding(0, 1).Bold(true)
	if active {
		style = style.Underline(true)
	}
	return style.Render(title) + "\n\n"
}

func paneStyle(pos int, theme styles.Theme) lipgloss.Style {
	style := lipgloss.NewStyle().Padding(0, 1)
	if pos == 2 {
		style = style.BorderStyle(lipgloss.OuterHalfBlockBorder()).
			BorderForeground(theme.BorderColor).BorderLeft(true)
	}

	return style
}

func calculateViewportWidths(termWidth int) (int, int, int) {
	paneOne := math.Floor(float64(termWidth) * 0.20)
	paneTwo := math.Floor(float64(termWidth) * 0.20)
	paneThree := termWidth - int(paneOne) - int(paneTwo)
	return int(paneOne), int(paneTwo), paneThree
}

func truncateText(s string, w int) string {
	padding := 10
	if runewidth.StringWidth(s) <= w-padding {
		// Don't truncate strings that fit
		return s
	}

	runes := []rune(s)
	width := 0
	for i := len(runes) - 1; i >= 0; i-- {
		r := runes[i]
		width += runewidth.RuneWidth(r)
		if width >= w-padding {
			return "..." + string(runes[i+1:])
		}
	}
	return string(runes)
}
