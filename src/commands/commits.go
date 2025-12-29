package commands

import (
	internal "CLI_App/src/internals/commits"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func FetchCommits() *cobra.Command {
	path, err := os.Getwd()
	if err != nil {
		fmt.Println("Error: ", err)
		os.Exit(1)
	}

	return &cobra.Command{
		Use:   "commits",
		Short: "Get the commits",
		Run: func(cmd *cobra.Command, args []string) {
			internal.FetchCommits(path)
		},
	}
}
