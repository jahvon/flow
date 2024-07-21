package executables

import (
	"fmt"

	"github.com/jahvon/flow/internal/context"
	"github.com/jahvon/flow/types/common"
	"github.com/jahvon/flow/types/executable"
)

func ExecutableForRef(ctx *context.Context, ref executable.Ref) (*executable.Executable, error) {
	executableRef := context.ExpandRef(ctx, ref)
	exec, err := ctx.ExecutableCache.GetExecutableByRef(ctx.Logger, executableRef)
	if err != nil {
		return nil, err
	} else if exec == nil {
		return nil, fmt.Errorf("executable missing ref='%s'", ref)
	}

	if exec.Exec != nil {
		fields := map[string]interface{}{
			"executable": exec.ID(),
		}
		exec.Exec.SetLogFields(fields)
	}

	return exec, nil
}

func ExecutableForCmd(parent *executable.Executable, cmd string, id int) *executable.Executable {
	vis := executable.ExecutableVisibility(common.VisibilityInternal)
	exec := &executable.Executable{
		Verb:       "exec",
		Name:       fmt.Sprintf("%s-cmd-%d", parent.Name, id),
		Visibility: &vis,
		Exec: &executable.ExecExecutableType{
			Cmd: cmd,
		},
	}
	fields := map[string]interface{}{"executable": exec.ID()}
	exec.Exec.SetLogFields(fields)
	exec.SetContext(parent.Workspace(), parent.WorkspacePath(), parent.Namespace(), parent.FlowFilePath())
	return exec
}
