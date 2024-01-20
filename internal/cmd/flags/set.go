package flags

import (
	"fmt"
	"reflect"

	"github.com/spf13/cobra"
)

type FlagSet struct {
	registeredFlags map[string]Metadata
}

//nolint:gocognit
func (f *FlagSet) Register(cmd *cobra.Command, flag Metadata, persistent bool) error {
	flagSet := cmd.Flags()
	if persistent {
		flagSet = cmd.PersistentFlags()
	}
	if flag.Default == nil {
		return fmt.Errorf("flag default must be defined using explicit type")
	}
	switch reflect.TypeOf(flag.Default).Kind() { //nolint:exhaustive
	case reflect.String:
		if flag.Shorthand != "" {
			flagSet.StringP(flag.Name, flag.Shorthand, flag.Default.(string), flag.Usage)
		} else {
			flagSet.String(flag.Name, flag.Default.(string), flag.Usage)
		}
	case reflect.Bool:
		if flag.Shorthand != "" {
			flagSet.BoolP(flag.Name, flag.Shorthand, flag.Default.(bool), flag.Usage)
		} else {
			flagSet.Bool(flag.Name, flag.Default.(bool), flag.Usage)
		}
	case reflect.Slice:
		var def []string
		var ok bool
		if def, ok = flag.Default.([]string); !ok {
			return fmt.Errorf("unexpected type received for %s flag", flag.Name)
		}

		if len(def) == 0 {
			def = nil
		}

		if flag.Shorthand != "" {
			flagSet.StringArrayP(flag.Name, flag.Shorthand, def, flag.Usage)
		} else {
			flagSet.StringArray(flag.Name, def, flag.Usage)
		}
	case reflect.Int:
		if flag.Shorthand != "" {
			flagSet.IntP(flag.Name, flag.Shorthand, flag.Default.(int), flag.Usage)
		} else {
			flagSet.Int(flag.Name, flag.Default.(int), flag.Usage)
		}
	default:
		return fmt.Errorf("unexpected flag default type (%v)", reflect.TypeOf(flag.Default).Kind())
	}

	if flag.Required {
		if err := cmd.MarkFlagRequired(flag.Name); err != nil {
			return err
		}
	}

	if f.registeredFlags == nil {
		f.registeredFlags = make(map[string]Metadata)
	}
	f.registeredFlags[flag.Name] = flag
	return nil
}

func (f *FlagSet) ValueFor(cmd *cobra.Command, flagName string, persistent bool) (interface{}, error) {
	metadata, found := f.registeredFlags[flagName]
	if !found {
		return nil, fmt.Errorf("flag %s not registered for command", flagName)
	}

	flagSet := cmd.Flags()
	if persistent {
		flagSet = cmd.PersistentFlags()
	}
	flag := cmd.Flag(flagName)
	if flag == nil {
		return "", nil
	}

	switch reflect.TypeOf(metadata.Default).Kind() { //nolint:exhaustive
	case reflect.String:
		val, err := flagSet.GetString(flagName)
		if err != nil {
			return nil, err
		}
		return val, nil
	case reflect.Bool:
		val, err := flagSet.GetBool(flagName)
		if err != nil {
			return nil, err
		}
		return val, nil
	case reflect.Slice:
		val, err := flagSet.GetStringArray(flagName)
		if err != nil {
			return nil, err
		}
		return val, nil
	case reflect.Int:
		val, err := flagSet.GetInt(flagName)
		if err != nil {
			return nil, err
		}
		return val, nil
	default:
		return nil, fmt.Errorf("unexpected flag default type (%v)", reflect.TypeOf(metadata.Default).Kind())
	}
}
