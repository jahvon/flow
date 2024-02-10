//nolint:cyclop
package flags

import (
	"fmt"

	"github.com/jahvon/flow/config"
)

type Metadata struct {
	Name      string
	Shorthand string
	Usage     string
	Default   interface{}
	Required  bool
}

var VerbosityFlag = &Metadata{
	Name:     "verbosity",
	Usage:    "Log verbosity level (-1 to 1)",
	Default:  0,
	Required: false,
}

var SyncCacheFlag = &Metadata{
	Name:     "sync",
	Usage:    "Sync flow cache and workspaces",
	Default:  false,
	Required: false,
}

var FilterWorkspaceFlag = &Metadata{
	Name:      "workspace",
	Shorthand: "w",
	Usage:     "Filter executables by workspace.",
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

var FilterVerbFlag = &Metadata{
	Name:      "verb",
	Shorthand: "v",
	Usage:     fmt.Sprintf("Filter executables by verb. One of: %s", config.SortedValidVerbs()),
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

var OutputFormatFlag = &Metadata{
	Name:      "output",
	Shorthand: "o",
	Usage:     "Output format. One of: yaml, json, jsonp.",
	Default:   "",
	Required:  false,
}

var OutputSecretAsPlainTextFlag = &Metadata{
	Name:      "plainText",
	Shorthand: "p",
	Usage:     "Output the secret value as plain text instead of an obfuscated string",
	Default:   false,
	Required:  false,
}

var SetAfterCreateFlag = &Metadata{
	Name:      "set",
	Shorthand: "s",
	Usage:     "Set the newly created workspace as the current workspace",
	Default:   false,
	Required:  false,
}

var NonInteractiveFlag = &Metadata{
	Name:      "non-interactive",
	Shorthand: "x",
	Usage: "Disable displaying flow output via terminal UI rendering. " +
		"This is only needed if the interactive output is enabled by default in flow's configuration.",
	Default:  false,
	Required: false,
}

var LastLogEntryFlag = &Metadata{
	Name:     "last",
	Usage:    "Print the last execution's logs",
	Default:  false,
	Required: false,
}
