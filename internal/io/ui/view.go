package ui

const DefaultView = WorkspaceView

type View string

const (
	WorkspaceView  View = "workspace"
	NamespaceView  View = "namespace"
	ExecutableView View = "executable"
)
