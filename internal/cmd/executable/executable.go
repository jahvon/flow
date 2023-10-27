package executable

import (
	"errors"
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/jahvon/flow/internal/cmd/flags"
	"github.com/jahvon/flow/internal/cmd/utils"
	"github.com/jahvon/flow/internal/config"
	flowErrs "github.com/jahvon/flow/internal/errors"
	"github.com/jahvon/flow/internal/executable"
	"github.com/jahvon/flow/internal/executable/consts"
	"github.com/jahvon/flow/internal/io"
	"github.com/jahvon/flow/internal/workspace"
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
		return "", nil, flowErrs.WorkspaceNotFoundError{Workspace: ws}
	}

	definitions, err := workspace.LoadDefinitions(ws, wsPath)
	if err != nil {
		return "", nil, err
	}

	if len(definitions) == 0 {
		log.Debug().Msg("no definitions found in workspace")
		return "", nil, flowErrs.ExecutableNotFoundError{Agent: agentType, Name: name}
	}
	definitions = definitions.FilterByNamespace(ns)

	ns, exec, err := definitions.LookupExecutableByTypeAndName(agentType, name)
	if err != nil {
		return "", nil, err
	} else if exec == nil {
		return "", nil, flowErrs.ExecutableNotFoundError{Agent: agentType, Name: name}
	}

	return ns, exec, nil
}

func FlagsToExecutableList(
	cmd *cobra.Command,
	set flags.FlagSet,
	currentConfig *config.RootConfig,
) (executable.List, error) {
	context, err := utils.ValidateAndGetContext(cmd, set, currentConfig)
	if err != nil {
		return nil, err
	}
	var agentFilter, namespaceFilter string
	var tagsFilter []string

	agentRaw, err := set.ValueFor(cmd, flags.FilterAgentTypeFlag.Name)
	if err != nil {
		return nil, fmt.Errorf("invalid agent flag - %w", err)
	}
	if agentRaw.(string) != "" {
		agentFilter, _ = agentRaw.(string)
		log.Debug().Msgf("only including executables of type %s", agentFilter)
	}

	tagsRaw, err := set.ValueFor(cmd, flags.FilterTagFlag.Name)
	if err != nil {
		return nil, fmt.Errorf("invalid tag flag - %w", err)
	}
	if tagsRaw.([]string) != nil {
		tagsFilter, _ = tagsRaw.([]string)
		log.Debug().Msgf("only including executables with tags %v", tagsFilter)
	}

	namespaceRaw, err := set.ValueFor(cmd, flags.FilterNamespaceFlag.Name)
	if err != nil {
		return nil, fmt.Errorf("invalid namespace flag - %w", err)
	}
	if namespaceRaw.(string) != "" {
		namespaceFilter, _ = namespaceRaw.(string)
		log.Debug().Msgf("only including executables within namespace %s", namespaceFilter)
	}

	var executables executable.List
	if context == "global" {
		for ws, wsPath := range currentConfig.Workspaces {
			collected, err := collectExecutablesInWorkspace(
				ws,
				wsPath,
				agentFilter,
				namespaceFilter,
				tagsFilter,
			)
			if err != nil {
				return nil, err
			}
			executables = append(executables, collected...)
		}
	} else {
		executables, err = collectExecutablesInWorkspace(
			context,
			currentConfig.Workspaces[context],
			agentFilter,
			namespaceFilter,
			tagsFilter,
		)
		if err != nil {
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
	switch len(split) {
	case 1:
		name = split[0]
	case 2:
		ns = split[0]
		name = split[1]
	default:
		return "", "", "", errors.New("invalid executable identifier")
	}

	return ws, ns, name, nil
}

func collectExecutablesInWorkspace(
	ws, wsPath string,
	agentFilter, namespaceFilter string,
	tagsFilter []string,
) (executable.List, error) {
	log.Trace().Msgf("searching for executables in workspace %s", ws)

	var executables executable.List
	definitions, err := workspace.LoadDefinitions(ws, wsPath)
	if err != nil {
		return nil, err
	} else if len(definitions) == 0 {
		log.Debug().Msgf("no definitions found in workspace %s", ws)
		return nil, nil
	}

	if namespaceFilter != "" {
		definitions = definitions.FilterByNamespace(namespaceFilter)
	}
	for _, definition := range definitions {
		defExecutables := definition.Executables
		if agentFilter != "" {
			defExecutables = defExecutables.FilterByType(consts.AgentType(agentFilter))
		}
		if tagsFilter == nil || definition.HasAnyTag(tagsFilter) {
			executables = append(executables, defExecutables...)
		} else {
			defExecutables = defExecutables.FilterByTags(tagsFilter)
			executables = append(executables, defExecutables...)
		}
	}
	return executables, nil
}
