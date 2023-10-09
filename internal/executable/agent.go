package executable

import (
	"github.com/jahvon/flow/internal/executable/consts"
)

type Agent interface {
	Name() consts.AgentType
	Exec(executable Executable) error
}
