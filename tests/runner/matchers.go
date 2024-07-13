package runner

import (
	"fmt"

	"go.uber.org/mock/gomock"

	"github.com/jahvon/flow/types/executable"
)

type refMatcher struct {
	ref string
}

func (m *refMatcher) Matches(x any) bool {
	e, ok := x.(*executable.Executable)
	if !ok {
		return false
	}
	return e.Ref().String() == m.ref
}

func (m *refMatcher) String() string {
	return fmt.Sprintf("has ref %s", m.ref)
}

func ExecWithRef(ref executable.Ref) gomock.Matcher {
	return &refMatcher{ref.String()}
}

type cmdMatcher struct {
	cmd string
}

func (m *cmdMatcher) Matches(x any) bool {
	e, ok := x.(*executable.Executable)
	if !ok {
		return false
	}
	if e.Exec == nil {
		return false
	}
	return e.Exec.Cmd == m.cmd
}

func (m *cmdMatcher) String() string {
	return fmt.Sprintf("has command %s", m.cmd)
}

func ExecWithCmd(cmd string) gomock.Matcher {
	return &cmdMatcher{cmd}
}
