package executables

import (
	"fmt"

	"github.com/jahvon/flow/config"
	"github.com/jahvon/flow/internal/context"
)

func ExecutableForRef(ctx *context.Context, ref config.Ref) (*config.Executable, error) {
	executableRef := context.ExpandRef(ctx, ref)
	exec, err := ctx.ExecutableCache.GetExecutableByRef(ctx.Logger, executableRef)
	if err != nil {
		return nil, err
	} else if exec == nil {
		return nil, fmt.Errorf("executable missing ref='%s'", ref)
	}

	if exec.Type.Exec != nil {
		fields := map[string]interface{}{
			"executable": exec.ID(),
		}
		exec.Type.Exec.SetLogFields(fields)
	}

	return exec, nil
}

func ExecutableForCmd(parent *config.Executable, cmd string, id int) *config.Executable {
	vis := config.VisibilityInternal
	exec := &config.Executable{
		Verb:       "exec",
		Name:       fmt.Sprintf("%s-cmd-%d", parent.Name, id),
		Visibility: &vis,
		Type: &config.ExecutableTypeSpec{
			Exec: &config.ExecExecutableType{
				Command: cmd,
			},
		},
	}
	fields := map[string]interface{}{"executable": exec.ID()}
	exec.Type.Exec.SetLogFields(fields)
	exec.SetContext(parent.Workspace(), parent.WorkspacePath(), parent.Namespace(), parent.DefinitionPath())
	return exec
}
