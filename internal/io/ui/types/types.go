package types

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

type TickMsg time.Time
type NoticeLevel string

const (
	NoticeLevelInfo    NoticeLevel = "info"
	NoticeLevelWarning NoticeLevel = "warning"
	NoticeLevelError   NoticeLevel = "error"
)

type ParentView interface {
	Ready() bool
	Height() int
	Width() int
	HandleInternalError(err error)
	HandleUserError(err error)
	SetContext(ws, ns string)
	GetContext() (ws, ns string)
	SetNotice(notice string, lvl NoticeLevel)
	BuildAndSetView(ViewBuilder)
	Finalize()
}

type ViewBuilder interface {
	tea.Model

	FooterEnabled() bool
	HelpMsg() string
	Type() string
}
