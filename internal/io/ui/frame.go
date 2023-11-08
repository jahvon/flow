package ui

import (
	"fmt"
	"strings"

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

	brandItem := brandTxt(opts.CurrentState)
	contextItem := contextTxt(opts.CurrentWorkspace, opts.CurrentNamespace, opts.CurrentState)
	filterItem := filterTxt(opts.CurrentFilter, opts.CurrentState)
	noticeItem := noticeTxt(opts.Notice, opts.CurrentState)

	headerRow := tview.NewFlex().SetDirection(tview.FlexColumn).
		AddItem(brandItem, textViewWidth(brandItem)+2, 1, false).
		AddItem(contextItem, textViewWidth(contextItem)+4, 1, false).
		AddItem(filterItem, textViewWidth(filterItem)+4, 1, false).
		AddItem(noticeItem, 0, 2, false)

	contentTitle := fmt.Sprintf(" - %s - ", strings.ToLower(string(opts.CurrentView)))
	flex := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(headerRow, 1, 33, false).
		AddItem(contentBox(contentTitle, *opts), 0, 67, false)

	if err := app.SetRoot(flex, true).EnableMouse(false).Run(); err != nil {
		log.Error().Err(err).Msg("encountered error rendering ui")
	}
}
