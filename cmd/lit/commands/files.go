package commands

import (
	"CLI_App/cmd/adapters/analysis/languages"
	"CLI_App/cmd/adapters/config"
	"CLI_App/cmd/domain"
	"CLI_App/cmd/service"

	"github.com/spf13/cobra"
)

func Files() *cobra.Command {
	command := &cobra.Command{
		Use:   "scan",
		Short: "Scan the repository files (must include a .gitignore) and retrieves data from them",
		Run: func(cmd *cobra.Command, args []string) {
			loc, _ := cmd.Flags().GetBool("loc")
			fix, _ := cmd.Flags().GetBool("fix")

			jsonAdapter := config.NewJSONAdapter()
			scanner := service.NewScannerService(
				languages.NewFileAnalyzer(domain.Conventions[jsonAdapter.GetConfig().NamingConventionIndex]),
			)

			switch {
			case loc:
				scanner.ExecuteLOC()
				scanner.PrintLOCResults()
			case fix:
				scanner.FixFile()
			default:
				scanner.ScanFiles()
				scanner.PrintScanningResults()
			}
		},
	}
	command.Flags().Bool("loc", false, "Retrieves the languages used with statistics")
	command.Flags().Bool("fix", false, "Fixes up the variables with an invalid naming conventions."+
		"It only one convention to another\nExample: if you have variables snake_case and the active convention is camelCase, it's converted.")

	return command
}
