package config

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/jahvon/flow/internal/utils"
)

// +docsgen:param
// A parameter is a value that can be passed to an executable and all of its sub-executables.
// Only one of `text`, `secretRef`, or `prompt` must be set. Specifying more than one will result in an error.
type Parameter struct {
	// +docsgen:text
	// A static value to be passed to the executable.
	Text string `yaml:"text"`
	// +docsgen:secretRef
	// A reference to a secret to be passed to the executable.
	Prompt string `yaml:"prompt"`
	// +docsgen:prompt
	// A prompt to be displayed to the user to collect a value to be passed to the executable.
	SecretRef string `yaml:"secretRef"`

	// +docsgen:envKey
	// The name of the environment variable that will be set with the value of the parameter.
	EnvKey string `yaml:"envKey"`
}

type ParameterList []Parameter

const (
	ReservedEnvVarPrefix = "FLOW_"
)

func (p *Parameter) Validate() error {
	if err := utils.ValidateOneOf("parameter type", p.Text, p.SecretRef, p.Prompt); err != nil {
		return err
	}

	if p.EnvKey == "" {
		return errors.New("must specify envKey for parameter")
	} else {
		re := regexp.MustCompile("^[a-zA-Z0-9_]*$")
		if strings.HasPrefix(p.EnvKey, ReservedEnvVarPrefix) {
			return fmt.Errorf("env destination cannot start with reserved prefix '%s'", ReservedEnvVarPrefix)
		} else if !re.MatchString(p.EnvKey) {
			return fmt.Errorf("env destination must be alphanumeric and can only contain underscores characters")
		}
	}

	return nil
}
