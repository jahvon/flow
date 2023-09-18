package consts

type AgentType string

const (
	AgentTypeOpen AgentType = "open"
	AgentTypeRun  AgentType = "run"
)

var (
	ValidAgentTypes = []AgentType{
		AgentTypeOpen,
		AgentTypeRun,
	}
)
