package ui

import tea "github.com/charmbracelet/bubbletea"

type ViewBuilder interface {
	tea.Model
	
	FooterEnabled() bool
	HelpMsg() string
	Type() string
}
