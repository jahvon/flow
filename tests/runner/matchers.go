package runner

import (
	"fmt"

	"go.uber.org/mock/gomock"

	"github.com/jahvon/flow/config"
)

type refMatcher struct {
	ref string
}

func (m *refMatcher) Matches(x any) bool {
	e, ok := x.(*config.Executable)
	if !ok {
		return false
	}
	return e.Ref().String() == m.ref
}

func (m *refMatcher) String() string {
	return fmt.Sprintf("has ref %s", m.ref)
}

func ExecWithRef(ref config.Ref) gomock.Matcher {
	return &refMatcher{ref.String()}
}

type cmdMatcher struct {
	cmd string
}

func (m *cmdMatcher) Matches(x any) bool {
	e, ok := x.(*config.Executable)
	if !ok {
		return false
	}
	if e.Type == nil || e.Type.Exec == nil {
		return false
	}
	return e.Type.Exec.Command == m.cmd
}

func (m *cmdMatcher) String() string {
	return fmt.Sprintf("has command %s", m.cmd)
}

func ExecWithCmd(cmd string) gomock.Matcher {
	return &cmdMatcher{cmd}
}
