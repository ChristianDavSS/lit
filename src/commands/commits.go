package commands

import (
	internal "CLI_App/src/internals/analysis"

	"github.com/spf13/cobra"
)

func FetchCommits() *cobra.Command {

	return &cobra.Command{
		Use:   "scan",
		Short: "Scan the whole repository in one single command",
		Run: func(cmd *cobra.Command, args []string) {
			internal.FetchCommits()
		},
	}
}
