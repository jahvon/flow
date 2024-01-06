package ui

import (
	_ "embed"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/lipgloss"
)

var (
	primaryColor   = lipgloss.AdaptiveColor{Dark: "#81A2BE", Light: "#5F819D"}
	secondaryColor = lipgloss.AdaptiveColor{Dark: "#B294BB", Light: "#85678F"}
	warningColor   = lipgloss.AdaptiveColor{Dark: "#F0C674", Light: "#DE935F"}
	errorColor     = lipgloss.AdaptiveColor{Dark: "#CC6666", Light: "#A54242"}
	black          = lipgloss.AdaptiveColor{Dark: "#1D1F21", Light: "#FFFFFF"}
	white          = lipgloss.AdaptiveColor{Dark: "#FFFFFF", Light: "#C5C8C6"}
	gray           = lipgloss.AdaptiveColor{Dark: "#373B41", Light: "#707880"}
	lightGray      = lipgloss.AdaptiveColor{Dark: "#C5C8C6", Light: "#373B41"}

	infoStyle    = lipgloss.NewStyle().Foreground(secondaryColor).Render
	warningStyle = lipgloss.NewStyle().Foreground(warningColor).Render
	errorStyle   = lipgloss.NewStyle().Foreground(errorColor).Render
	unknownStyle = lipgloss.NewStyle().Foreground(gray).Render

	spinnerStyle = lipgloss.NewStyle().Foreground(secondaryColor)
	brandStyle   = lipgloss.NewStyle().
			PaddingRight(2).
			PaddingLeft(4).
			Italic(true).
			Bold(true).
			Foreground(white).
			Background(primaryColor).Render
	contextStyle = lipgloss.NewStyle().
			PaddingRight(1).
			PaddingLeft(1).
			Foreground(gray).
			Background(secondaryColor).Render
	filterStyle = lipgloss.NewStyle().
			PaddingRight(1).
			PaddingLeft(0).
			Foreground(gray).
			Background(secondaryColor).Render
	helpStyle = lipgloss.NewStyle().
			MarginLeft(2).
			Foreground(gray).Render
	entityViewStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(primaryColor).
			MarginLeft(2)
	collectionStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(primaryColor).
			MarginLeft(2).
			Padding(0, 1)

	//go:embed markdown.json
	markdownStyleJson string
)

func listStyles() list.Styles {
	styles := list.DefaultStyles()
	styles.StatusBar = styles.StatusBar.
		Padding(0, 0, 1, 0).
		Italic(true).
		Foreground(primaryColor)
	return styles
}

func listItemStyles() list.DefaultItemStyles {
	styles := list.NewDefaultItemStyles()
	styles.NormalTitle = styles.NormalTitle.
		Foreground(lightGray).
		Bold(true)
	styles.NormalDesc = styles.NormalDesc.Foreground(lightGray)

	styles.SelectedTitle = styles.SelectedTitle.
		Border(lipgloss.DoubleBorder(), false, false, false, true).
		Foreground(secondaryColor).
		BorderForeground(secondaryColor).
		Bold(true)
	styles.SelectedDesc = styles.SelectedDesc.
		Border(lipgloss.HiddenBorder(), false, false, false, true).
		Foreground(secondaryColor)
	return styles
}
