//nolint:cyclop
package flags

import (
	"fmt"

	"github.com/jahvon/flow/internal/executable/consts"
)

type Metadata struct {
	Name      string
	Shorthand string
	Usage     string
	Default   interface{}
	Required  bool
}

var VerbosityFlag = &Metadata{
	Name:      "verbosity",
	Shorthand: "v",
	Usage:     "Log verbosity level (from 0 to 4 where 4 is most verbose)",
	Default:   2,
	Required:  false,
}

var SyncCacheFlag = &Metadata{
	Name:     "sync",
	Usage:    "Sync flow cache and workspaces",
	Default:  false,
	Required: false,
}

var ListGlobalContextFlag = &Metadata{
	Name:      "global",
	Shorthand: "g",
	Usage:     "List from all workspaces",
	Default:   false,
	Required:  false,
}

var ListWorkspaceContextFlag = &Metadata{
	Name:      "workspace",
	Shorthand: "w",
	Usage:     "List from a specific workspace",
	Default:   "",
	Required:  false,
}

var FilterNamespaceFlag = &Metadata{
	Name:      "namespace",
	Shorthand: "n",
	Usage:     "Filter executables by namespace.",
	Default:   "",
	Required:  false,
}

var FilterTagFlag = &Metadata{
	Name:      "tag",
	Shorthand: "t",
	Usage:     "Filter by tags.",
	Default:   []string{},
	Required:  false,
}

var FilterAgentTypeFlag = &Metadata{
	Name:      "agent",
	Shorthand: "a",
	Usage:     fmt.Sprintf("Filter by executable agent type. One of: %s", consts.ValidAgentTypes),
	Default:   "",
	Required:  false,
}

var SpecificWorkspaceFlag = &Metadata{
	Name:      "workspace",
	Shorthand: "w",
	Usage:     "Specify a workspace to get",
	Default:   "",
	Required:  false,
}

var OutputFormatFlag = &Metadata{
	Name:      "output",
	Shorthand: "o",
	Usage:     "Output format. One of: default, yaml, json, jsonp.",
	Default:   "default",
	Required:  false,
}

var OutputSecretAsPlainTextFlag = &Metadata{
	Name:      "plainText",
	Shorthand: "p",
	Usage:     "Output the secret value as plain text instead of an obfuscated string",
	Default:   false,
	Required:  false,
}

var OutputMetadataFlag = &Metadata{
	Name:      "metadata",
	Shorthand: "m",
	Usage:     "Include metadata in output.",
	Default:   false,
	Required:  false,
}

var AgentTypeFlag = &Metadata{
	Name:      "agent",
	Shorthand: "a",
	Usage:     fmt.Sprintf("Executable agent type. One of: %s", consts.ValidAgentTypes),
	Default:   "",
	Required:  true,
}

var WorkspacePathFlag = &Metadata{
	Name:      "path",
	Shorthand: "p",
	Usage:     "Path to the workspace",
	Default:   "",
	Required:  true,
}

var SetAfterCreateFlag = &Metadata{
	Name:      "set",
	Shorthand: "s",
	Usage:     "Set the newly created workspace as the current workspace",
	Default:   false,
	Required:  false,
}
