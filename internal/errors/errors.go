package errors

import (
	"fmt"

	"github.com/jahvon/flow/internal/executable/consts"
)

func ExecutableNotFound(agent consts.AgentType, name string) error {
	return fmt.Errorf("%s executable %s not found", agent, name)
}

func WorkspaceNotFound(ws string) error {
	return fmt.Errorf("workspace %s not found", ws)
}
