package executable

import "github.com/jahvon/flow/internal/executable/consts"

type Agent interface {
	Name() consts.AgentType
	Exec(
		spec map[string]interface{},
		preferences *Preference,
	) error
}
