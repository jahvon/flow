package ui

import "github.com/gdamore/tcell/v2"

const DefaultState = IdleState

type State string

const (
	IdleState        State = "idle"
	ProgressingState State = "progressing"
	SuccessState     State = "success"
	ErrorState       State = "error"
)

func (s State) PrimaryFGColor() tcell.Color {
	return tcell.ColorWhite
}

func (s State) PrimaryBGColor() tcell.Color {
	switch s {
	case IdleState:
		return tcell.ColorDarkBlue
	case ProgressingState:
		return tcell.ColorDarkGoldenrod
	case SuccessState:
		return tcell.ColorDarkGreen
	case ErrorState:
		return tcell.ColorDarkRed
	default:
		return tcell.ColorDarkGray
	}
}

func (s State) SecondaryFGColor() tcell.Color {
	return s.PrimaryBGColor()
}

func (s State) SecondaryBGColor() tcell.Color {
	switch s {
	case IdleState:
		return tcell.ColorLightBlue
	case ProgressingState:
		return tcell.ColorLightGoldenrodYellow
	case SuccessState:
		return tcell.ColorLightBlue
	case ErrorState:
		return tcell.ColorLightPink
	default:
		return tcell.ColorLightGray
	}
}
