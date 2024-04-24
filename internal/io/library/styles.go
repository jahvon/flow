package library

import (
	"fmt"
	"math"
	"strings"

	"github.com/charmbracelet/lipgloss"
	tuikitStyles "github.com/jahvon/tuikit/styles"
)

var styles = tuikitStyles.BaseTheme()

func renderSelection(s string) string {
	color := styles.PrimaryColor
	style := lipgloss.NewStyle().Foreground(color)
	return style.Render(s)
}

func renderSecondarySelection(s string) string {
	color := styles.SecondaryColor
	style := lipgloss.NewStyle().Foreground(color)
	return style.Render(s)
}

func renderInactive(s string) string {
	color := styles.White
	style := lipgloss.NewStyle().Foreground(color)
	return style.Render(s)
}

func renderDescription(s string) string {
	color := styles.Gray
	style := lipgloss.NewStyle().Foreground(color)
	return style.Render(s)
}

func renderPaneTitle(s string, count int, active bool) string {
	var bg, fg lipgloss.AdaptiveColor
	var title string
	if active {
		bg = styles.PrimaryColor
	} else {
		bg = styles.TertiaryColor
	}
	fg = styles.Black

	if count == 0 {
		title = s
	} else {
		title = fmt.Sprintf("%s (%d)", s, count)
	}
	style := lipgloss.NewStyle().Background(bg).Foreground(fg).Padding(0, 1)
	return style.Render(title) + "\n\n"

}

func renderHeader(ctx string) string {
	prefixBlock := lipgloss.NewStyle().
		Width(1).Height(1).
		Background(styles.Gray).
		Foreground(styles.Gray).
		BorderStyle(lipgloss.RoundedBorder()).Inline(true).Faint(true)
	prefix := strings.Repeat(prefixBlock.Render("+"), 20)
	style := lipgloss.NewStyle().Padding(0, 2).
		Background(styles.PrimaryColor).
		Foreground(styles.TertiaryColor).Inline(true)
	style = style.Align(lipgloss.Center)
	return prefix + style.Bold(true).Render("flow") + style.Render(ctx)
}

func renderFooter(s string) string {
	style := lipgloss.NewStyle().Padding(1, 1).Foreground(styles.Gray)
	return style.Render(s)
}

func paneStyle(active bool, pos int) lipgloss.Style {
	style := lipgloss.NewStyle().Padding(0, 1)
	switch pos {
	case 0:
		style = style.BorderStyle(lipgloss.NormalBorder()).
			BorderTop(true).BorderBottom(true).
			BorderRight(false).BorderLeft(true)
	case 1:
		style = style.BorderStyle(lipgloss.NormalBorder()).
			BorderTop(true).BorderBottom(true).
			BorderRight(false).BorderLeft(false)
	case 2:
		style = style.BorderStyle(lipgloss.NormalBorder()).
			BorderTop(true).BorderBottom(true).
			BorderRight(true).BorderLeft(true)
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
	if len(s) < 10 {
		// Don't truncate very short strings
		return s
	}
	if len(s) > w {
		return s[:w-3] + "..."
	}
	return s
}
