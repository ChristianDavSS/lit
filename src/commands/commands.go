package commands

import (
	"CLI_App/src/config"
	internal "CLI_App/src/internals/analysis"

	"github.com/spf13/cobra"
)

func Files() *cobra.Command {
	command := &cobra.Command{
		Use:   "files",
		Short: "Scan the repository files (must include a .gitignore) and retrieves data from them",
		Run: func(cmd *cobra.Command, args []string) {
			loc, _ := cmd.Flags().GetBool("loc")
			internal.Files(loc)
		},
	}
	command.Flags().Bool("loc", false, "Retrieves the languages used with statistics")

	return command
}

func Configuration() *cobra.Command {
	return &cobra.Command{
		Use:   "config",
		Short: "Configurate the scan variables.",
		Run: func(cmd *cobra.Command, args []string) {
			config.Init()
		},
	}
}
