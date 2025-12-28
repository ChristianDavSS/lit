package src

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var root = &cobra.Command{
	Use:   "lit",
	Short: "Lit CLI tool",
	Long:  "Lit CLI tool",
}

// Execute function to execute some code
func Execute() {
	root.AddCommand(&cobra.Command{
		Use:   "version",
		Short: "Print the version number",
		Long:  "Command that prints the version number",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("Valor actual del servidor: %d\n", 10)
		},
	})
	// Execute the root, registering all the children commands
	if err := root.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	os.Exit(0)
}
