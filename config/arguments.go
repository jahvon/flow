package config

import (
	"errors"
	"fmt"
	"strconv"
)

type Argument struct {
	// +docsgen:envKey
	// The name of the environment variable that will be assigned the value.
	EnvKey string `yaml:"envKey"`
	// +docsgen:typeArg
	// The type of the argument. This is used to determine how to parse the value of the argument.
	// String is the default type.
	Type string `yaml:"type"`
	// +docsgen:default
	// The default value to use if a value is not set.
	Default string `yaml:"default"`
	// +docsgen:required
	// If true, the argument must be set. If false, the default value will be used if the argument is not set.
	Required bool `yaml:"required"`

	// +docsgen:flag
	// The flag to use when setting the argument from the command line.
	// Either `flag` or `pos` must be set, but not both.
	Flag string `yaml:"flag"`
	// +docsgen:pos
	// The position of the argument in the command line arguments. Values start at 1.
	// Either `flag` or `pos` must be set, but not both.
	Pos int `yaml:"pos"`

	value string
}

func (a *Argument) Set(value string) {
	a.value = value
}

func (a *Argument) Value() string {
	if a.value == "" {
		return a.Default
	}
	return a.value
}

func (a *Argument) Validate() error {
	if a.EnvKey == "" {
		return errors.New("must specify envKey for argument")
	}
	if err := validateArgType(a.Type); err != nil {
		return fmt.Errorf("%s - %w", a.EnvKey, err)
	}
	if a.Flag != "" && a.Pos != 0 {
		return errors.New("either flag or pos must be set, but not both")
	} else if a.Flag == "" && a.Pos == 0 {
		return errors.New("either flag or pos must be set")
	}
	return nil
}

func (a *Argument) ValidateValue() error {
	if a.value == "" && a.Required {
		return fmt.Errorf("required argument not set")
	}

	switch a.Type {
	case "int":
		if _, err := strconv.Atoi(a.value); err != nil {
			return fmt.Errorf("value is not an integer")
		}
	case "bool":
		if _, err := strconv.ParseBool(a.value); err != nil {
			return fmt.Errorf("value is not a boolean")
		}
	}
	return nil
}

func validateArgType(t string) error {
	switch t {
	case "string", "int", "bool":
		return nil
	default:
		return fmt.Errorf("unsupported argument type (%s)", t)
	}
}

type ArgumentList []Argument

func (al *ArgumentList) Validate() error {
	var errs []error
	for _, arg := range *al {
		if err := arg.Validate(); err != nil {
			errs = append(errs, fmt.Errorf("argument %s validation failed - %w", arg.EnvKey, err))
		}
	}
	collectedFlags := make(map[string]struct{})
	collectedPos := make(map[int]struct{})
	for _, arg := range *al {
		if arg.Flag != "" {
			if _, ok := collectedFlags[arg.Flag]; ok {
				errs = append(errs, fmt.Errorf("flag %s is assigned to more than one argument", arg.Flag))
			}
			collectedFlags[arg.Flag] = struct{}{}
		} else if arg.Pos != 0 {
			if _, ok := collectedPos[arg.Pos]; ok {
				errs = append(errs, fmt.Errorf("position %d is assigned to more than one argument", arg.Pos))
			}
			collectedPos[arg.Pos] = struct{}{}
		}
	}
	if len(errs) > 0 {
		return fmt.Errorf("%d argument validation errors: %v", len(errs), errs)
	}
	return nil
}

func (al *ArgumentList) ValidateValues() error {
	var errs []error
	for _, arg := range *al {
		if err := arg.ValidateValue(); err != nil {
			errs = append(errs, fmt.Errorf("argument %s validation failed - %w", arg.EnvKey, err))
		}
	}
	if len(errs) > 0 {
		return fmt.Errorf("%d argument validation errors: %v", len(errs), errs)
	}
	return nil
}

func (al *ArgumentList) ToEnvMap() map[string]string {
	envMap := make(map[string]string)
	for _, arg := range *al {
		envMap[arg.EnvKey] = arg.Value()
	}
	return envMap
}

func (al *ArgumentList) SetValues(flagArgs map[string]string, posArgs []string) error {
	for i, arg := range *al {
		if arg.Flag != "" {
			if val, ok := flagArgs[arg.Flag]; ok {
				arg.Set(val)
				(*al)[i] = arg
			}
		} else if arg.Pos != 0 {
			if arg.Pos <= len(posArgs) {
				arg.Set(posArgs[arg.Pos-1])
				(*al)[i] = arg
			}
		}
	}
	return al.ValidateValues()
}
