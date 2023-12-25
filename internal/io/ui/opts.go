package ui

import (
	"fmt"
	"sort"
	"strings"

	"github.com/jahvon/flow/config"
)

type HeaderOption func(application *Application)

func WithCurrentState(state State) HeaderOption {
	return func(a *Application) {
		a.curState = state
	}
}

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
	var tagStr string
	sort.Strings(tags)
	for _, tag := range tags {
		tagStr += fmt.Sprintf("<%s> ", tag)
	}
	tagStr = strings.TrimSpace(tagStr)
	return func(a *Application) {
		a.curFilter = tagStr
	}
}

func WithNotice(notice string) HeaderOption {
	return func(a *Application) {
		a.curNotice = notice
	}
}
