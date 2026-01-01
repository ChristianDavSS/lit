package cmd

import (
	"fmt"
	"os"

	"CLI_App/src/commands"

	"github.com/spf13/cobra"
)

// root: commands entry point. Every command is a subcommand of root
var root = &cobra.Command{
	Use:     "lit",
	Short:   "Lit CLI tool for your git projects",
	Long:    "Lit CLI is a tool made for with love for developers.\nScan your repository and get feedback now.",
	Version: "0.01 beta",
}

// Execute function to execute some code
func Execute() {
	// Add commands to the root
	root.AddCommand(commands.FetchCommits())
	root.AddCommand(commands.Files())

	// Execute the root, registering all the children commands
	if err := root.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	// End the execution
	os.Exit(0)
}
