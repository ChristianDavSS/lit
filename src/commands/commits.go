package commands

import (
	internal "CLI_App/src/internals/analysis"
	"time"

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
			since, _ := cmd.Flags().GetString("since")
			until, _ := cmd.Flags().GetString("until")
			internal.FetchCommits(w, v, since, until)
		},
	}
	// Non-persistent flags for the command
	command.Flags().StringP("who", "w", "", "Check a person stats on the repo.")
	command.Flags().Bool("verbose", false, "Print the whole information.")
	command.Flags().String("since", "", "Get the commits from a date on (DD/MM/YYYY).")
	command.Flags().String("until", time.Now().Format("02/01/2006"),
		"Get the commits until a estimated date (DD/MM/YYYY).")

	return command
}
