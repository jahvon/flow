package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/jahvon/flow/internal/cmd/flags"
)

var Flags *flags.FlagSet

func registerFlagOrPanic(cmd *cobra.Command, flag flags.Metadata) {
	if Flags == nil {
		Flags = &flags.FlagSet{}
	}
	if err := Flags.Register(cmd, flag, false); err != nil {
		panic(err)
	}
}

func getFlagValue[T any](cmd *cobra.Command, flag flags.Metadata) T {
	v, err := Flags.ValueFor(cmd, flag.Name, false)
	if err != nil {
		panic(err)
	}
	val, ok := v.(T)
	if !ok {
		panic(fmt.Errorf("unable to cast flag value to type %T", val))
	}
	return val
}

func registerPersistentFlagOrPanic(cmd *cobra.Command, flag flags.Metadata) {
	if Flags == nil {
		Flags = &flags.FlagSet{}
	}
	if err := Flags.Register(cmd, flag, true); err != nil {
		panic(err)
	}
}

func getPersistentFlagValue[T any](cmd *cobra.Command, flag flags.Metadata) T {
	v, err := Flags.ValueFor(cmd, flag.Name, true)
	if err != nil {
		panic(err)
	}
	val, ok := v.(T)
	if !ok {
		panic(fmt.Errorf("unable to cast flag value to type %T", val))
	}
	return val
}
