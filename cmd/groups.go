package cmd

import "github.com/spf13/cobra"

var (
	CrudGroup       = &cobra.Group{ID: "crud", Title: "Configuration CRUD Commands"}
	ExecutableGroup = &cobra.Group{ID: "flow", Title: "Flow Executable Commands"}
)
