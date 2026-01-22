package commands

import (
	"fmt"
	"os"

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
	root.AddCommand(ScanCmd())
	root.AddCommand(ConfigCmd())

	// Execute the root, registering all the children commands
	if err := root.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
