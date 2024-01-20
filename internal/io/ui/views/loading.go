package views

import (
	"fmt"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/jahvon/flow/internal/io/styles"
	"github.com/jahvon/flow/internal/io/ui/types"
)

const LoadingViewType = "loading"

type LoadingView struct {
	msg     string
	spinner spinner.Model
}

func (v *LoadingView) Init() tea.Cmd {
	return v.spinner.Tick
}

func (v *LoadingView) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case error:
		v.msg = msg.Error()
	case string:
		v.msg = msg
	}
	v.spinner, cmd = v.spinner.Update(msg)
	return v, cmd
}

func (v *LoadingView) View() string {
	var txt string
	if v.msg == "" {
		txt = fmt.Sprintf("\n\n %s %s\n\n", v.spinner.View(), styles.RenderInfo("thinking..."))
	} else {
		txt = fmt.Sprintf("\n\n %s %s\n\n", v.spinner.View(), styles.RenderInfo(v.msg))
	}
	return txt
}

func (v *LoadingView) HelpMsg() string {
	return ""
}

func (v *LoadingView) FooterEnabled() bool {
	return false
}

func (v *LoadingView) Type() string {
	return LoadingViewType
}

func NewLoadingView(msg string) types.ViewBuilder {
	spin := spinner.New()
	spin.Style = styles.SpinnerStyle()
	spin.Spinner = spinner.Points
	return &LoadingView{
		msg:     msg,
		spinner: spin,
	}
}
