package ui

import (
	"fmt"
	"sort"
	"strings"
)

type FrameOptions struct {
	CurrentView      View
	CurrentState     State
	CurrentWorkspace string
	CurrentNamespace string
	CurrentFilter    string
	Notice           string

	ObjectContent *TableData
	ObjectList    *TableData
}

type FrameOption func(*FrameOptions)

func MergeFrameOptions(opts ...FrameOption) *FrameOptions {
	o := &FrameOptions{
		CurrentView:      DefaultView,
		CurrentState:     DefaultState,
		CurrentWorkspace: "unk",
		CurrentNamespace: "unk",
		CurrentFilter:    "",
		Notice:           "---",
		ObjectContent:    nil,
		ObjectList:       nil,
	}
	for _, opt := range opts {
		opt(o)
	}
	return o
}

func WithCurrentView(view View) FrameOption {
	return func(o *FrameOptions) {
		o.CurrentView = view
	}
}

func WithCurrentState(state State) FrameOption {
	return func(o *FrameOptions) {
		o.CurrentState = state
	}
}

func WithCurrentWorkspace(workspace string) FrameOption {
	return func(o *FrameOptions) {
		o.CurrentWorkspace = workspace
	}
}

func WithCurrentNamespace(namespace string) FrameOption {
	return func(o *FrameOptions) {
		o.CurrentNamespace = namespace
	}
}

func WithCurrentFilter(tags []string) FrameOption {
	var tagStr string
	sort.Strings(tags)
	for _, tag := range tags {
		tagStr += fmt.Sprintf("<%s> ", tag)
	}
	tagStr = strings.TrimSpace(tagStr)
	return func(o *FrameOptions) {
		o.CurrentFilter = tagStr
	}
}

func WithNotice(notice string) FrameOption {
	return func(o *FrameOptions) {
		o.Notice = notice
	}
}

func WithObjectContent(content *TableData) FrameOption {
	if content == nil {
		return func(o *FrameOptions) {}
	}

	return func(o *FrameOptions) {
		o.ObjectContent = content
		o.ObjectList = nil
	}
}

func WithObjectList(list *TableData) FrameOption {
	if list == nil {
		return func(o *FrameOptions) {}
	}

	return func(o *FrameOptions) {
		o.ObjectList = list
		o.ObjectContent = nil
	}
}
