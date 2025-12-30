package commands

import (
	internal "CLI_App/src/internals/analysis"

	"github.com/spf13/cobra"
)

func FetchCommits() *cobra.Command {
	// Define the cobra command
	command := &cobra.Command{
		Use:   "scan",
		Short: "Scan the whole repository in one single command",
		Run: func(cmd *cobra.Command, args []string) {
			// Get the 'who' flag value
			w, _ := cmd.Flags().GetString("who")
			v, _ := cmd.Flags().GetBool("verbose")
			internal.FetchCommits(w, v)
		},
	}
	// Non-persistent flags for the command
	command.Flags().StringP("who", "w", "", "Check a person stats on the repo.")
	command.Flags().BoolP("verbose", "", false, "Print the whole information.")

	return command
}
