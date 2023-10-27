package errors

import (
	"fmt"

	"github.com/jahvon/flow/internal/executable/consts"
)

type ExecutableNotFoundError struct {
	Agent consts.AgentType
	Name  string
}

func (e ExecutableNotFoundError) Error() string {
	return fmt.Sprintf("%s executable %s not found", e.Agent, e.Name)
}

type WorkspaceNotFoundError struct {
	Workspace string
}

func (e WorkspaceNotFoundError) Error() string {
	return fmt.Sprintf("workspace %s not found", e.Workspace)
}
