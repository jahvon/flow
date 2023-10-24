package cmd

import "github.com/spf13/cobra"

var (
	DataGroup       = &cobra.Group{ID: "data", Title: "Data and Metadata"}
	ExecutableGroup = &cobra.Group{ID: "exec", Title: "Executables"}
)
