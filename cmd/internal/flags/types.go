//nolint:cyclop,lll
package flags

import (
	"fmt"

	"github.com/jahvon/flow/types/executable"
)

type Metadata struct {
	Name      string
	Shorthand string
	Usage     string
	Default   interface{}
	Required  bool
}

var LogLevel = &Metadata{
	Name:      "log-level",
	Usage:     "Log verbosity level (debug, info, fatal)",
	Shorthand: "L",
	Default:   "info",
	Required:  false,
}

var SyncCacheFlag = &Metadata{
	Name:     "sync",
	Usage:    "Sync flow cache and workspaces",
	Default:  false,
	Required: false,
}

var FilterExecSubstringFlag = &Metadata{
	Name:      "filter",
	Shorthand: "f",
	Usage:     "Filter executable by reference substring.",
	Default:   "",
	Required:  false,
}

var FilterWorkspaceFlag = &Metadata{
	Name:      "workspace",
	Shorthand: "w",
	Usage:     "Filter executables by workspace.",
	Default:   "",
	Required:  false,
}

var AllNamespacesFlag = &Metadata{
	Name:      "all",
	Shorthand: "a",
	Usage:     "List from all namespaces.",
	Default:   false,
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
	Usage:     fmt.Sprintf("Filter executables by verb. One of: %s", executable.SortedValidVerbs()),
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
	Usage:     "Output format. One of: yaml, json, or tui.",
	Default:   "tui",
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

var FixedWsModeFlag = &Metadata{
	Name:      "fixed",
	Shorthand: "f",
	Usage:     "Set the workspace mode to fixed",
	Default:   false,
	Required:  false,
}

var ListFlag = &Metadata{
	Name:      "list",
	Shorthand: "l",
	Usage:     "Show a simple list view of executables instead of interactive discovery.",
	Default:   false,
	Required:  false,
}

var CopyFlag = &Metadata{
	Name:     "copy",
	Usage:    "Copy the secret value to the clipboard",
	Default:  false,
	Required: false,
}

var LastLogEntryFlag = &Metadata{
	Name:     "last",
	Usage:    "Print the last execution's logs",
	Default:  false,
	Required: false,
}

var TemplateWorkspaceFlag = &Metadata{
	Name:      "workspace",
	Shorthand: "w",
	Usage:     "Workspace to create the flow file and its artifacts. Defaults to the current workspace.",
	Default:   "",
	Required:  false,
}

var TemplateOutputPathFlag = &Metadata{
	Name:      "output",
	Shorthand: "o",
	Usage:     "Output directory (within the workspace) to create the flow file and its artifacts. If the directory does not exist, it will be created.",
	Default:   "",
	Required:  false,
}

var TemplateFlag = &Metadata{
	Name:      "template",
	Shorthand: "t",
	Usage:     "Registered template name. Templates can be registered in the flow configuration file or with `flow set template`.",
	Default:   "",
	Required:  false,
}

var TemplateFilePathFlag = &Metadata{
	Name:      "file",
	Shorthand: "f",
	Usage:     "Path to the template file. It must be a valid flow file template.",
	Default:   "",
	Required:  false,
}

var SetSoundNotificationFlag = &Metadata{
	Name:    "sound",
	Usage:   "Update completion sound notification setting",
	Default: false,
}

var StoreFullFlag = &Metadata{
	Name:    "full",
	Usage:   "Force clear all stored data",
	Default: false,
}
