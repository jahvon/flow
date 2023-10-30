package ui

import (
	"fmt"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"

	"github.com/jahvon/flow/internal/io"
)

var log = io.Log()

const FlowBanner = `
    ______  ____ _      __
   / __/ / / __ \ | /| / /
  / _// /_/ /_/ / |/ |/ / 
 /_/ /____|____/|__/|__/
`

func PrintUiFrame(optList ...FrameOption) {
	opts := MergeFrameOptions(optList...)
	app := tview.NewApplication()

	infoGrid := tview.NewGrid().SetSize(5, 75, 1, 0).
		AddItem(emptyBox(), 0, 0, 2, 75, 0, 0, false).
		AddItem(wsTxt(opts.CurrentWorkspace), 2, 0, 1, 75, 0, 0, false).
		AddItem(nsTxt(opts.CurrentNamespace), 3, 0, 1, 75, 0, 0, false).
		AddItem(noticeTxt(opts.Notice), 4, 0, 1, 75, 0, 0, false)

	headerRow := tview.NewFlex().SetDirection(tview.FlexColumn).
		AddItem(bannerTxt(), 0, 20, false).
		AddItem(infoGrid, 0, 80, false)

	contentBox := tview.NewBox().SetBorder(true).SetBorderColor(tcell.ColorGray).
		SetTitle(string(opts.CurrentView))

	flex := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(headerRow, 5, 33, false).
		AddItem(contentBox, 0, 67, false)

	if err := app.SetRoot(flex, true).EnableMouse(false).Run(); err != nil {
		log.Error().Err(err).Msg("encountered error rendering ui")
	}
}

func bannerTxt() *tview.TextView {
	return tview.NewTextView().
		SetTextColor(tcell.ColorLightGreen).
		SetTextAlign(tview.AlignLeft).
		SetScrollable(false).
		SetText(FlowBanner)
}

func wsTxt(workspace string) *tview.TextView {
	return tview.NewTextView().
		SetDynamicColors(true).
		SetScrollable(false).
		SetText(fmt.Sprintf("\t[lightblue](w) workspace: [white]%s", workspace))
}

func nsTxt(namespace string) *tview.TextView {
	return tview.NewTextView().
		SetDynamicColors(true).
		SetScrollable(false).
		SetText(fmt.Sprintf("\t[lightblue](n) namespace: [white]%s", namespace))
}

func noticeTxt(notice string) *tview.TextView {
	notice = strings.TrimSpace(notice)
	return tview.NewTextView().
		SetTextColor(tcell.ColorPaleVioletRed).
		SetText(fmt.Sprintf("\t%s", notice))
}

func emptyBox() *tview.Box {
	return tview.NewBox().SetBorder(false)
}
