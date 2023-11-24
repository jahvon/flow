package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/jahvon/flow/internal/cmd/flags"
	"github.com/jahvon/flow/internal/io"
)

var Flags *flags.FlagSet

func registerFlagOrPanic(cmd *cobra.Command, flag flags.Metadata) {
	if Flags == nil {
		Flags = &flags.FlagSet{}
	}
	if err := Flags.Register(cmd, flag); err != nil {
		log.Panic().Msgf(err.Error())
	}
}

func getFlagValue[T any](cmd *cobra.Command, flag flags.Metadata) T {
	v, err := Flags.ValueFor(cmd, flag.Name)
	if err != nil {
		io.PrintErrorAndExit(err)
	}
	val, ok := v.(T)
	if !ok {
		io.PrintErrorAndExit(fmt.Errorf("unable to cast flag value to type %T", val))
	}
	return val
}
