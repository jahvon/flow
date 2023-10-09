//nolint:cyclop
package executable

import (
	"errors"
	"strings"

	"github.com/jahvon/flow/internal/cmd/flags"
	"github.com/jahvon/flow/internal/cmd/utils"
	"github.com/jahvon/flow/internal/config"
	flowErrs "github.com/jahvon/flow/internal/errors"
	"github.com/jahvon/flow/internal/executable"
	"github.com/jahvon/flow/internal/executable/consts"
	"github.com/jahvon/flow/internal/io"
	"github.com/jahvon/flow/internal/workspace"
	"github.com/spf13/cobra"
)

var log = io.Log()

// ArgsToExecutable parses the given arguments and returns the corresponding executable.
// The first argument is the executable identifier, which can be in the form of: <workspace>/<namespace>:<name>
// If found, the executable is returned along with the namespace it was found in.
func ArgsToExecutable(
	args []string, agentType consts.AgentType, currentConfig *config.RootConfig,
) (string, *executable.Executable, error) {
	if len(args) != 1 {
		return "", nil, errors.New("unexpected number of arguments")
	}
	executableIdentifier := args[0]
	ws, ns, name, err := SplitExecutableIdentifier(executableIdentifier)
	if err != nil {
		return "", nil, err
	}

	if ws == "" {
		log.Debug().Msg("defaulting to the current workspace")
		ws = currentConfig.CurrentWorkspace
	}
	wsPath, found := currentConfig.Workspaces[ws]
	if !found {
		return "", nil, flowErrs.WorkspaceNotFound(ws)
	}

	definitions, err := workspace.LoadDefinitions(ws, wsPath)
	if err != nil {
		return "", nil, err
	}

	if len(definitions) == 0 {
		log.Debug().Msg("no definitions found in workspace")
		return "", nil, flowErrs.ExecutableNotFound(agentType, name)
	}
	definitions = definitions.FilterByNamespace(ns)

	ns, exec, err := definitions.LookupExecutableByTypeAndName(agentType, name)
	if err != nil {
		return "", nil, err
	} else if exec == nil {
		return "", nil, flowErrs.ExecutableNotFound(agentType, name)
	}

	return ns, exec, nil
}

func FlagsToExecutableList(cmd *cobra.Command, currentConfig *config.RootConfig) (executable.List, error) {
	context, err := utils.ValidateAndGetContext(cmd, currentConfig)
	if err != nil {
		return nil, err
	}
	var agentFilter, tagFilter, namespaceFilter *string
	if flag := cmd.Flag(flags.AgentTypeFlagName); flag != nil && flag.Changed {
		*agentFilter = flag.Value.String()
		log.Debug().Msgf("only including executables of type %s", *agentFilter)
	}
	if flag := cmd.Flag(flags.TagFlagName); flag != nil && flag.Changed {
		*tagFilter = flag.Value.String()
		log.Debug().Msgf("only including executables with tag %s", *tagFilter)
	}
	if flag := cmd.Flag(flags.NamespaceFlagName); flag != nil && flag.Changed {
		*namespaceFilter = flag.Value.String()
		log.Debug().Msgf("only including executable within namespace %s", *namespaceFilter)
	}

	var executables executable.List
	collectExecutablesInWorkspace := func(ws, wsPath string) error {
		log.Trace().Msgf("searching for executables in workspace %s", ws)
		definitions, err := workspace.LoadDefinitions(ws, wsPath)
		if err != nil {
			return err
		} else if len(definitions) == 0 {
			log.Debug().Msgf("no definitions found in workspace %s", ws)
			return nil
		}

		if namespaceFilter != nil {
			definitions = definitions.FilterByNamespace(*namespaceFilter)
		}
		for _, definition := range definitions {
			defExecutables := definition.Executables
			if agentFilter != nil {
				defExecutables = defExecutables.FilterByType(consts.AgentType(*agentFilter))
			}
			if tagFilter == nil || definition.HasTag(*tagFilter) {
				executables = append(executables, defExecutables...)
			} else {
				defExecutables = defExecutables.FilterByTag(*tagFilter)
				executables = append(executables, defExecutables...)
			}
		}
		return nil
	}

	if context == "global" {
		for ws, wsPath := range currentConfig.Workspaces {
			if err := collectExecutablesInWorkspace(ws, wsPath); err != nil {
				return nil, err
			}
		}
	} else {
		if err := collectExecutablesInWorkspace(context, currentConfig.Workspaces[context]); err != nil {
			return nil, err
		}
	}

	return executables, nil
}

func SplitExecutableIdentifier(identifier string) (ws, ns, name string, _ error) {
	var split []string
	var remaining string

	if identifier == "" {
		return "", "", "", errors.New("invalid executable identifier")
	}

	split = strings.Split(identifier, "/")
	switch len(split) {
	case 1:
		remaining = split[0]
	case 2:
		ws = split[0]
		remaining = split[1]
	default:
		return "", "", "", errors.New("invalid executable identifier")
	}

	split = strings.Split(remaining, ":")
	if len(split) > 2 {
		return "", "", "", errors.New("invalid executable identifier")
	} else if len(split) == 2 {
		ns = split[0]
		name = split[1]
	} else {
		name = split[0]
	}

	return ws, ns, name, nil
}
