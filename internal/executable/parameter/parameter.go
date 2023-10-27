package parameter

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/samber/lo"

	"github.com/jahvon/flow/internal/io"
	"github.com/jahvon/flow/internal/vault"
)

type Parameter struct {
	// Only one of text, secretRef, or prompt should be set.
	Text      string `yaml:"text"`
	Prompt    string `yaml:"prompt"`
	SecretRef string `yaml:"secretRef"`

	EnvKey string `yaml:"envKey"`
}

const (
	ReservedPrefix = "FLOW_"
)

func (p *Parameter) Validate() error {
	switch {
	case p.Text == "" && p.SecretRef == "" && p.Prompt == "":
		return errors.New("must set either text, secretRef, or prompt for parameter")
	case p.Text != "" && p.SecretRef != "":
		return errors.New("cannot set both text and secretRef for parameter")
	case p.Text != "" && p.Prompt != "":
		return errors.New("cannot set both text and prompt for parameter")
	case p.SecretRef != "" && p.Prompt != "":
		return errors.New("cannot set both secretRef and prompt for parameter")
	}

	if p.EnvKey == "" {
		return errors.New("must specify envKey for parameter")
	} else {
		re := regexp.MustCompile("^[a-zA-Z0-9_]*$")
		if strings.HasPrefix(p.EnvKey, ReservedPrefix) {
			return fmt.Errorf("env destination cannot start with reserved prefix '%s'", ReservedPrefix)
		} else if !re.MatchString(p.EnvKey) {
			return fmt.Errorf("env destination must be alphanumeric and can only contain underscores characters")
		}
	}

	return nil
}

func (p *Parameter) Value() (string, error) {
	switch {
	case p.Text == "" && p.SecretRef == "" && p.Prompt == "":
		return "", nil
	case p.Text != "":
		return p.Text, nil
	case p.Prompt != "":
		response := io.Ask(p.Prompt)
		return response, nil
	case p.SecretRef != "":
		if err := vault.ValidateReference(p.SecretRef); err != nil {
			return "", err
		}
		v := vault.NewVault()
		secret, err := v.GetSecret(p.SecretRef)
		if err != nil {
			return "", err
		}
		return secret.PlainTextString(), nil
	default:
		return "", errors.New("failed to get value for parameter")
	}
}

func ParameterListToEnvList(params []Parameter) ([]string, error) {
	var errs []error
	env := lo.Map(params, func(param Parameter, _ int) string {
		key := param.EnvKey
		value, err := param.Value()
		if err != nil {
			errs = append(errs, err)
			return ""
		}
		return fmt.Sprintf("%s=%s", key, value)
	})

	if len(errs) > 0 {
		return nil, fmt.Errorf("failed to get values for parameters: %v", errs)
	}
	return env, nil
}
