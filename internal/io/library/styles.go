package library

import (
	"fmt"
	"math"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/jahvon/tuikit/themes"
	"github.com/mattn/go-runewidth"

	"github.com/jahvon/flow/types/executable"
)

func renderSelection(s string, theme themes.Theme) string {
	style := lipgloss.NewStyle().Foreground(theme.ColorPalette().PrimaryColor())
	return style.Render(s)
}

func renderSecondarySelection(s string, theme themes.Theme) string {
	style := lipgloss.NewStyle().Foreground(theme.ColorPalette().SecondaryColor())
	return style.Render(s)
}

func renderInactive(s string, theme themes.Theme) string {
	style := lipgloss.NewStyle().Foreground(theme.ColorPalette().GrayColor())
	return style.Render(s)
}

func renderDescription(s string, theme themes.Theme) string {
	style := lipgloss.NewStyle().Foreground(theme.ColorPalette().BodyColor())
	return style.Render(s)
}

func renderPaneTitle(s string, count int, active bool, theme themes.Theme) string {
	var title string
	if count == 0 {
		title = s
	} else {
		title = fmt.Sprintf("%s (%d)", s, count)
	}
	style := lipgloss.NewStyle().Foreground(theme.ColorPalette().SecondaryColor()).Padding(0, 1).Bold(true)
	if active {
		style = style.Underline(true)
	}
	return style.Render(title) + "\n\n"
}

func paneStyle(pos int, theme themes.Theme, splitView bool) lipgloss.Style {
	style := lipgloss.NewStyle().Padding(0, 1)
	if pos == 2 && splitView {
		style = style.BorderStyle(lipgloss.OuterHalfBlockBorder()).
			BorderForeground(theme.ColorPalette().BorderColor()).BorderLeft(true)
	}

	return style
}

func calculateViewportWidths(termWidth int, splitView bool) (int, int, int) {
	if splitView {
		paneOne := math.Floor(float64(termWidth) * 0.20)
		paneTwo := math.Floor(float64(termWidth) * 0.30)
		paneThree := termWidth - int(paneOne) - int(paneTwo)
		return int(paneOne), int(paneTwo), paneThree
	} else {
		paneOne := math.Floor(float64(termWidth) * 0.33)
		paneTwo := math.Floor(float64(termWidth) * 0.67)
		paneThree := termWidth
		return int(paneOne), int(paneTwo), paneThree
	}
}

func shortRef(ref executable.Ref, ws, ns string) string {
	shortID := ref.ID()
	if ws != "" && ref.Workspace() == ws {
		shortID = strings.Replace(shortID, ws+"/", "", 1)
	}
	if ns != "" && ref.Namespace() == ns {
		shortID = strings.Replace(shortID, ns+":", "", 1)
	}
	return executable.NewRef(shortID, ref.Verb()).String()
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
