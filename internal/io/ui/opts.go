package ui

import "github.com/rivo/tview"

type FrameOptions struct {
	CurrentView      View
	CurrentWorkspace string
	CurrentNamespace string
	Notice           string

	ObjectContent *tview.Box
	ObjectList    *tview.List
}

type FrameOption func(*FrameOptions)

func MergeFrameOptions(opts ...FrameOption) *FrameOptions {
	o := &FrameOptions{
		CurrentView:      DefaultView,
		CurrentWorkspace: "[red]unk",
		CurrentNamespace: "[red]unk",
		Notice:           "testing. this is a placeholder notice!",
		ObjectContent:    nil,
		ObjectList:       nil,
	}
	for _, opt := range opts {
		opt(o)
	}
	return o
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

func WithObjectContent(content *tview.Box) FrameOption {
	return func(o *FrameOptions) {
		o.ObjectContent = content
	}
}

func WithObjectList(list *tview.List) FrameOption {
	return func(o *FrameOptions) {
		o.ObjectList = list
	}
}
