package ui

import (
	"github.com/jahvon/flow/config"
)

type HeaderOption func(application *Application)

func WithCurrentWorkspace(workspace string) HeaderOption {
	return func(a *Application) {
		a.curWs = workspace
	}
}

func WithCurrentNamespace(namespace string) HeaderOption {
	return func(a *Application) {
		a.curNs = namespace
	}
}

func WithCurrentFilter(tags config.Tags) HeaderOption {
	return func(a *Application) {
		a.curFilter = tags
	}
}

func WithNotice(notice string, lvl NoticeLevel) HeaderOption {
	return func(a *Application) {
		a.curNotice = notice
		a.curNoticeLvl = lvl
	}
}
