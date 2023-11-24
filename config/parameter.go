package config

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/jahvon/flow/internal/utils"
)

type Parameter struct {
	// Only one of text, secretRef, or prompt should be set.
	Text      string `yaml:"text"`
	Prompt    string `yaml:"prompt"`
	SecretRef string `yaml:"secretRef"`

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
