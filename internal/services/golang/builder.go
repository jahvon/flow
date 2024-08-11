package golang

type Builder struct {
	Binary    string
	Command   string
	Directory string
	Main      string
	Ldflags   []string
}

func NewBuilder() *Builder {
	b := Builder{}
	return b
}
