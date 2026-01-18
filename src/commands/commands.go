package commands

import (
	"CLI_App/src/config"
	internal "CLI_App/src/internals/analysis"
	"time"

	"github.com/spf13/cobra"
)

func FetchCommits() *cobra.Command {
	// Define the cobra command
	command := &cobra.Command{
		Use:   "authors",
		Short: "Scan the whole repository and returns the stats of the authors",
		Run: func(cmd *cobra.Command, args []string) {
			// Get the 'who' flag value
			w, _ := cmd.Flags().GetString("who")
			v, _ := cmd.Flags().GetBool("verbose")
			stats, _ := cmd.Flags().GetBool("stats")
			since, _ := cmd.Flags().GetString("since")
			until, _ := cmd.Flags().GetString("until")
			commitSize, _ := cmd.Flags().GetBool("commit-size")
			internal.FetchCommits(w, v, stats, since, until, commitSize)
		},
	}
	// Non-persistent flags for the command
	command.Flags().StringP("who", "w", "", "Check a person stats on the repo.")
	command.Flags().Bool("verbose", false, "Print the whole information.")
	command.Flags().Bool("stats", false, "Print out the commit stats.")
	command.Flags().Bool("commit-size", false, "Print out the mean of lines changed per commit.")
	command.Flags().String("since", "", "Get the commits from a date on (DD/MM/YYYY).")
	command.Flags().String("until", time.Now().Format("02/01/2006"),
		"Get the commits until a estimated date (DD/MM/YYYY).")

	return command
}

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
		Use:   "scan-config",
		Short: "Configurate the scan variables.",
		Run: func(cmd *cobra.Command, args []string) {
			config.Init()
		},
	}
}
