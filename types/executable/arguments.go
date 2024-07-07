package executable

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/jahvon/flow/internal/utils"
)

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
	if err := utils.ValidateOneOf("argument type", a.Flag, a.Pos); err != nil {
		return err
	}

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
	case "float":
		if _, err := strconv.ParseFloat(a.value, 64); err != nil {
			return fmt.Errorf("value is not a float")
		}
	case "bool":
		if _, err := strconv.ParseBool(a.value); err != nil {
			return fmt.Errorf("value is not a boolean")
		}
	case "string":
		// no-op
	default:
		return fmt.Errorf("unsupported argument type (%s)", a.Type)
	}
	return nil
}

func validateArgType(t ArgumentType) error {
	switch t {
	case "string", "int", "bool", "float":
		return nil
	default:
		return fmt.Errorf("unsupported argument type (%s)", t)
	}
}

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
