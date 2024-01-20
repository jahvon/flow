package styles

import (
	_ "embed"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
)

//go:embed markdown.json
var MarkdownStyleJSON string

// Colors

var (
	// Inspirited by Wryan and Jellybeans https://gogh-co.github.io/Gogh/
	PrimaryColor   = lipgloss.AdaptiveColor{Dark: "#477AB3", Light: "#31658C"}
	SecondaryColor = lipgloss.AdaptiveColor{Dark: "#7E62B3", Light: "#5E468C"}
	TertiaryColor  = lipgloss.AdaptiveColor{Dark: "#9E9ECB", Light: "#7C7C99"}

	WarningColor = lipgloss.AdaptiveColor{Dark: "#FFDCA0", Light: "#FFBA7B"}
	InfoColor    = lipgloss.AdaptiveColor{Dark: "#53A6A6", Light: "#287373"}
	ErrorColor   = lipgloss.AdaptiveColor{Dark: "#BF4D80", Light: "#8C4665"}

	White = lipgloss.AdaptiveColor{Dark: "#C0C0C0", Light: "#899CA1"}
	Gray  = lipgloss.AdaptiveColor{Dark: "#3D3D3D", Light: "#333333"}
	Black = lipgloss.AdaptiveColor{Dark: "#2b2a2a", Light: "#000000"}
)

// Styles

func SpinnerStyle() lipgloss.Style {
	return lipgloss.NewStyle().Foreground(SecondaryColor)
}

func EntityViewStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(TertiaryColor).
		MarginLeft(2)
}

func CollectionStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(TertiaryColor).
		MarginLeft(2).
		Padding(0, 1)
}

func TermViewStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		BorderStyle(lipgloss.HiddenBorder()).
		BorderForeground(TertiaryColor).
		Padding(0, 1).
		MarginLeft(2)
}

// Render Functions

func RenderBold(text string) string {
	return lipgloss.NewStyle().Bold(true).Render(text)
}

func RenderInfo(text string) string {
	return lipgloss.NewStyle().Foreground(InfoColor).Render(text)
}

func RenderSuccess(text string) string {
	return lipgloss.NewStyle().Foreground(SecondaryColor).Render(text)
}

func RenderWarning(text string) string {
	return lipgloss.NewStyle().Foreground(WarningColor).Render(text)
}

func RenderError(text string) string {
	return lipgloss.NewStyle().Foreground(ErrorColor).Render(text)
}

func RenderUnknown(text string) string {
	return lipgloss.NewStyle().Foreground(Gray).Render(text)
}

func RenderBrand(text string) string {
	return lipgloss.NewStyle().
		PaddingRight(2).
		PaddingLeft(4).
		Italic(true).
		Bold(true).
		Foreground(Black).
		Background(SecondaryColor).Render(text)
}

func RenderContext(text string, padded bool) string {
	var rightPadding, leftPadding = 0, 2
	if padded {
		rightPadding = 10
		leftPadding = 1
	}
	return lipgloss.NewStyle().
		PaddingRight(rightPadding).
		PaddingLeft(leftPadding).
		Foreground(Black).
		Background(TertiaryColor).Render(text)
}

func RenderInputForm(text string) string {
	return lipgloss.NewStyle().
		PaddingLeft(2).
		Render(text)
}

func RenderHelp(text string) string {
	return lipgloss.NewStyle().
		MarginLeft(2).
		Foreground(Gray).Render(text)
}

// Misc

func ListStyles() list.Styles {
	styles := list.DefaultStyles()
	styles.StatusBar = styles.StatusBar.
		Padding(0, 0, 1, 0).
		Italic(true).
		Foreground(TertiaryColor)
	return styles
}

func ListItemStyles() list.DefaultItemStyles {
	styles := list.NewDefaultItemStyles()
	styles.NormalTitle = styles.NormalTitle.
		Foreground(TertiaryColor).
		Bold(true)
	styles.NormalDesc = styles.NormalDesc.Foreground(White)

	styles.SelectedTitle = styles.SelectedTitle.
		Border(lipgloss.DoubleBorder(), false, false, false, true).
		Foreground(SecondaryColor).
		BorderForeground(SecondaryColor).
		Bold(true)
	styles.SelectedDesc = styles.SelectedDesc.
		Border(lipgloss.HiddenBorder(), false, false, false, true).
		Foreground(White)
	return styles
}

func LoggerStyles() *log.Styles {
	baseStyles := log.DefaultStyles()
	baseStyles.Timestamp = baseStyles.Timestamp.Foreground(lipgloss.AdaptiveColor{Dark: "#505050", Light: "#505050"})
	baseStyles.Key = baseStyles.Key.Foreground(SecondaryColor)
	baseStyles.Value = baseStyles.Value.Foreground(Gray)
	baseStyles.Levels = map[log.Level]lipgloss.Style{
		log.InfoLevel:  baseStyles.Levels[log.InfoLevel].Foreground(InfoColor),
		log.WarnLevel:  baseStyles.Levels[log.WarnLevel].Foreground(WarningColor),
		log.ErrorLevel: baseStyles.Levels[log.ErrorLevel].Foreground(ErrorColor),
		log.DebugLevel: baseStyles.Levels[log.DebugLevel].Foreground(TertiaryColor),
		log.FatalLevel: baseStyles.Levels[log.FatalLevel].Foreground(SecondaryColor).SetString("ERR"),
	}
	return baseStyles
}
