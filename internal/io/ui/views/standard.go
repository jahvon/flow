package views

import (
	"fmt"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/jahvon/flow/internal/io"
	"github.com/jahvon/flow/internal/io/styles"
	"github.com/jahvon/flow/internal/io/ui/types"
)

const TermViewType = "term"

type StandardTermView struct {
	parent     types.ParentView
	viewport   viewport.Model
	logger     *io.Logger
	textInputs []*TextInput
	curInput   int
	navigated  bool
}

type TextInput struct {
	input     *textinput.Model
	submitted chan bool

	Key         string
	Prompt      string
	Placeholder string
	Hidden      bool
}

func (t *TextInput) Value() string {
	return t.input.Value()
}

func (t *TextInput) Render() string {
	t.input.Focus()
	return styles.RenderInputForm(fmt.Sprintf("\n%s\n%s\n", styles.RenderInfo(t.Prompt), t.input.View()))
}

func NewTermView(parent types.ParentView, logger *io.Logger) types.ViewBuilder {
	vp := viewport.New(parent.Width(), parent.Height())
	vp.Style = styles.TermViewStyle()
	return &StandardTermView{
		parent:     parent,
		viewport:   vp,
		textInputs: make([]*TextInput, 0),
		logger:     logger,
	}
}

func (v *StandardTermView) Init() tea.Cmd {
	return v.viewport.Init()
}

func (v *StandardTermView) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if !v.parent.Ready() {
			return v, nil
		}

		switch msg.String() {
		case "up":
			if v.viewport.AtTop() {
				break
			}
			v.viewport.LineUp(1)
			v.navigated = true
		case "down":
			if v.viewport.AtBottom() {
				v.navigated = false
			} else {
				v.viewport.LineDown(1)
				v.navigated = true
			}
		case "enter":
			if v.processingInput() {
				v.textInputs[v.curInput].submitted <- true
				v.curInput++
			}
		}
	case types.TickMsg:
		if v.processingInput() {
			break
		}
		v.viewport.SetContent(v.logger.ReadAllData())
		if !v.navigated {
			v.viewport.GotoBottom()
		}
		_, cmd = v.viewport.Update(msg)
	}
	if v.processingInput() {
		curTextInput := v.textInputs[v.curInput]
		updatedInput, inputCmd := curTextInput.input.Update(msg)
		curTextInput.input = &updatedInput
		cmd = inputCmd
	}
	return v, cmd
}

func (v *StandardTermView) View() string {
	if v.processingInput() && v.textInputs[v.curInput].input != nil {
		return v.textInputs[v.curInput].Render()
	}
	return v.viewport.View()
}

func (v *StandardTermView) StartProcessingUserInputs(inputs ...*TextInput) {
	if v.textInputs == nil {
		v.textInputs = make([]*TextInput, len(inputs))
	}
	for _, input := range inputs {
		if input == nil {
			continue
		}
		input.input = newTextInputView(input)
		input.submitted = make(chan bool)
		v.textInputs = append(v.textInputs, input)
	}
}

func (v *StandardTermView) GetTextInputs() []*TextInput {
	return v.textInputs
}

func (v *StandardTermView) WaitForTextInputs() error {
	if !v.processingInput() {
		return nil
	}
	timeout := time.After(10 * time.Minute)
	select {
	case <-timeout:
		return fmt.Errorf("timed out waiting for user input")
	case <-v.textInputs[len(v.textInputs)-1].submitted:
		return nil
	}
}

func (v *StandardTermView) FooterEnabled() bool {
	return true
}

func (v *StandardTermView) HelpMsg() string {
	return ""
}

func (v *StandardTermView) Type() string {
	return TermViewType
}

func (v *StandardTermView) processingInput() bool {
	return len(v.textInputs) > 0 && v.curInput < len(v.textInputs)
}

func newTextInputView(t *TextInput) *textinput.Model {
	in := textinput.New()
	echoMode := textinput.EchoNormal
	if t.Hidden {
		echoMode = textinput.EchoPassword
	}
	in.EchoMode = echoMode
	if t.Placeholder != "" {
		t.input.Placeholder = t.Placeholder
	}
	return &in
}
