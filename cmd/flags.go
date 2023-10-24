package cmd

import (
	"github.com/spf13/cobra"

	"github.com/jahvon/flow/internal/cmd/flags"
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
