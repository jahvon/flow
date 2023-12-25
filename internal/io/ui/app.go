package ui

import (
	"context"
	"fmt"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"

	"github.com/jahvon/flow/internal/io"
)

var log = io.Log().With().Str("scope", "io/ui/frame").Logger()

type PageKeyHandler func(event *tcell.EventKey) *tcell.EventKey

type Application struct {
	app   *tview.Application
	pages *tview.Pages

	activePage      string
	activeView      tview.Primitive
	lastPage        string
	lastView        tview.Primitive
	pageKeyHandlers map[string]PageKeyHandler

	curState                           State
	curWs, curNs, curNotice, curFilter string
}

func StartApplication(cancel context.CancelFunc) *Application {
	tviewApp := tview.NewApplication()
	tviewApp.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyCtrlC:
			tviewApp.Stop()
		}
		return event
	})
	pages := tview.NewPages()
	pages.AddPage("loading", LoadingView(), true, true)
	tviewApp.SetRoot(pages, true).EnableMouse(true).SetFocus(pages)

	go func() {
		if err := tviewApp.Run(); err != nil {
			log.Panic().Err(err).Msg("encountered error rendering ui")
		}
		// Cancel the context when the application is stopped.
		log.Trace().Msg("stopping application")
		cancel()
	}()

	return &Application{
		app:       tviewApp,
		pages:     pages,
		curState:  IdleState,
		curWs:     "unk",
		curNs:     "unk",
		curFilter: "",
		curNotice: "---",
	}
}

func (a *Application) Suspend(f func()) {
	a.app.Suspend(f)
}

func (a *Application) RegisterKeyHandler(pageTitle string, handler PageKeyHandler) {
	if a.pageKeyHandlers == nil {
		a.pageKeyHandlers = make(map[string]PageKeyHandler)
	}
	a.pageKeyHandlers[pageTitle] = handler
}

func (a *Application) SetHeader(opts ...HeaderOption) {
	for _, opt := range opts {
		opt(a)
	}
	a.SetPage(a.activePage, a.activeView)
}

func (a *Application) ShowModal(text string) {
	modal := tview.NewTextView().
		SetText(text)
	// 	AddButtons([]string{"Close"}).
	// 	SetDoneFunc(func(_ int, _ string) {
	// 		a.pages.SwitchToPage(a.activePage)
	// 	})
	// modal.SetBackgroundColor(IdleState.PrimaryBGColor())
	// modal.SetTextColor(IdleState.PrimaryFGColor())

	container := tview.NewFlex().
		AddItem(nil, 0, 1, false).
		AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
			AddItem(nil, 0, 1, false).
			AddItem(modal, 40, 1, true).
			AddItem(nil, 0, 1, false), 10, 1, true).
		AddItem(nil, 0, 1, false)

	a.pages.AddPage("modal", container, true, true)
	a.app.SetFocus(modal)
}

func (a *Application) SetPage(title string, view tview.Primitive) {
	brandItem := brandTxt(a.curState)
	contextItem := contextTxt(a.curWs, a.curNs, a.curState)
	filterItem := filterTxt(a.curFilter, a.curState)
	noticeItem := noticeTxt(a.curNotice, a.curState)

	headerRow := tview.NewFlex().SetDirection(tview.FlexColumn).
		AddItem(brandItem, textViewWidth(brandItem)+2, 1, false).
		AddItem(contextItem, textViewWidth(contextItem)+4, 1, false).
		AddItem(filterItem, textViewWidth(filterItem)+4, 1, false).
		AddItem(noticeItem, 0, 2, false)

	container := tview.NewFlex()
	container.SetDirection(tview.FlexRow).
		SetBorder(true).
		SetBorderColor(tcell.ColorWhite).
		SetTitleAlign(tview.AlignCenter).
		SetTitleColor(a.curState.PrimaryFGColor())
	if title != "" {
		container.SetTitle(fmt.Sprintf(" - %s - ", strings.ToLower(title)))
	}
	container.AddItem(view, 0, 1, false)

	page := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(headerRow, 1, 33, false).
		AddItem(container, 0, 67, false)

	if a.pageKeyHandlers != nil {
		if handler, found := a.pageKeyHandlers[title]; found {
			page.SetInputCapture(handler)
		}
	}

	a.pages.AddAndSwitchToPage(title, page, true)
	a.app.SetFocus(view)
	a.activePage = title
	a.activeView = view
}

func (a *Application) SetErrorPage(err error) {
	a.curState = ErrorState
	a.SetPage("error", ErrorView(err))
}

func (a *Application) Idle() {
	a.SetHeader(WithCurrentState(IdleState))
}

func (a *Application) Progressing() {
	a.SetHeader(WithCurrentState(ProgressingState))
}

func (a *Application) Succeeded() {
	a.SetHeader(WithCurrentState(SuccessState))
}

func (a *Application) Errored() {
	a.SetHeader(WithCurrentState(ErrorState))
}
