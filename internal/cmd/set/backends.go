package set

import (
	"errors"
	"fmt"
	"strconv"
	"syscall"
	"time"

	"github.com/spf13/cobra"
	"golang.org/x/term"

	"github.com/jahvon/flow/internal/backend"
	"github.com/jahvon/flow/internal/backend/consts"
	"github.com/jahvon/flow/internal/cmd/flags"
	"github.com/jahvon/flow/internal/io"
	"github.com/jahvon/flow/internal/parameter"
)

var log = io.Log()

func FlagsToAuthConfig(
	cmd *cobra.Command,
	baseAuthConfig *backend.AuthConfig,
) (*backend.AuthConfig, error) {
	if cmd == nil || baseAuthConfig == nil {
		return nil, errors.New("unexpected empty input")
	} else if cmd.Flags().NFlag() == 0 {
		_ = cmd.Usage()
		return nil, errors.New("no flags set")
	}

	var backendName consts.BackendName
	if cmd.Flag(flags.BackendFlagName).Value.String() != "" {
		log.Trace().Msg("backend flag set; using specified backend")
		n := consts.BackendName(cmd.Flag(flags.BackendFlagName).Value.String())
		backendName = n
	} else {
		log.Trace().Msg("no backend flag set; using existing backend")
		backendName = baseAuthConfig.Backend
	}
	var preferredMode consts.AuthMode
	if setVal := cmd.Flag(flags.PreferredModeFlagName).Value.String(); setVal != "" {
		log.Trace().Msg("preferred mode flag set; using specified mode")
		setMode := consts.AuthMode(setVal)
		preferredMode = setMode
	} else {
		log.Trace().Msg("no preferred mode flag set; using existing mode")
		preferredMode = baseAuthConfig.PreferredMode
	}

	var rememberMe bool
	if setVal := cmd.Flag(flags.RememberMeFlagName).Value.String(); setVal != "" {
		log.Trace().Msg("remember me flag set; using new value")
		setBool, err := strconv.ParseBool(setVal)
		if err != nil {
			return nil, fmt.Errorf("invalid remember me flag - %v", err)
		}
		rememberMe = setBool
	} else {
		log.Trace().Msg("no remember me flag set; using existing value")
		rememberMe = baseAuthConfig.RememberMe
	}

	var rememberMeDuration time.Duration
	if setVal := cmd.Flag(flags.RememberMeDurationFlagName).Value.String(); setVal != "" {
		log.Trace().Msg("remember me duration flag set; using specified duration")
		setDuration, err := time.ParseDuration(setVal)
		if err != nil {
			return nil, fmt.Errorf("invalid duration - %v", err)
		}
		rememberMeDuration = setDuration
	} else {
		log.Trace().Msg("no remember me duration flag set; using existing duration")
		rememberMeDuration = baseAuthConfig.RememberMeDuration
	}

	authConfig := &backend.AuthConfig{
		Backend:            backendName,
		PreferredMode:      preferredMode,
		RememberMe:         rememberMe,
		RememberMeDuration: rememberMeDuration,
	}
	return authConfig, nil
}

func FlagsToSecretConfig(
	cmd *cobra.Command,
	baseSecretConfig *backend.SecretConfig,
) (*backend.SecretConfig, error) {
	if cmd == nil || baseSecretConfig == nil {
		return nil, errors.New("unexpected empty input")
	} else if cmd.Flags().NFlag() == 0 {
		_ = cmd.Usage()
		return nil, errors.New("no flags set")
	}

	var backendName *consts.BackendName
	if cmd.Flag(flags.BackendFlagName).Value.String() != "" {
		n := consts.BackendName(cmd.Flag(flags.BackendFlagName).Value.String())
		backendName = &n
	} else {
		backendName = &baseSecretConfig.Backend
	}

	secretConfig := &backend.SecretConfig{
		Backend: *backendName,
	}

	return secretConfig, nil
}

func FlagsToParameter(cmd *cobra.Command, existingParamList []*parameter.Parameter) (*parameter.Parameter, error) {
	key := cmd.Flag(flags.KeyFlagName).Value.String()
	// exposeAs := parameter.Destination(cmd.Flag(flags.ExposeAsFlagName).Value.String())
	txtVal := cmd.Flag(flags.TextValueFlagName).Value.String()
	secretRef := cmd.Flag(flags.SecretRefFlagName).Value.String()

	param := &parameter.Parameter{}
	existingParam, found := parameter.LookupParameter(existingParamList, key)
	if found && existingParam != nil {
		param = existingParam
	}

	if key == "" {
		return nil, errors.New("key cannot be empty")
	}
	param.Key = parameter.NormalizeKey(key)

	if txtVal == "" && secretRef == "" {
		return nil, errors.New("must set either text or secretRef")
	} else if txtVal != "" && secretRef != "" {
		return nil, errors.New("cannot set both text and secretRef")
	}

	if param.Value == nil {
		param.Value = &parameter.Value{}
	}

	if txtVal == "-" {
		param.Value.Text = io.Ask("Value:")
		param.Value.Ref = ""
	} else if txtVal != "" {
		param.Value.Text = txtVal
		param.Value.Ref = ""
	}

	if secretRef != "" {
		param.Value.Ref = secretRef
	}

	// todo
	// if exposeAs == "" && param.ExposeAs == nil {
	// 	exposeAs = parameter.DestinationEnv
	// }

	return param, nil
}

func ArgsToSecretKV(args []string) (string, backend.Secret, error) {
	if len(args) != 2 {
		return "", "", errors.New("unexpected number of arguments")
	}
	key := args[0]
	value := args[1]

	if key == "" {
		return "", "", errors.New("key cannot be empty")
	}

	if value == "-" {
		rawValue, err := term.ReadPassword(syscall.Stdin)
		if err != nil {
			return "", "", fmt.Errorf("failed to read secret value - %v", err)
		}
		value = string(rawValue)
	}
	if value == "" {
		return "", "", errors.New("value cannot be empty")
	}

	return parameter.NormalizeKey(key), backend.Secret(value), nil
}
