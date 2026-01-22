package commands

import (
	"CLI_App/internal/adapter/config"
	"CLI_App/internal/domain"
	"CLI_App/internal/service/scanner"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func ScanCmd() *cobra.Command {
	command := &cobra.Command{
		Use:   "scan",
		Short: "Scan the repository files (must include a .gitignore) and retrieves data from them",
		Run: func(cmd *cobra.Command, args []string) {
			loc, _ := cmd.Flags().GetBool("loc")
			fix, _ := cmd.Flags().GetBool("fix")

			// Load Config
			cfgAdapter := config.NewJSONConfigAdapter()
			cfg, err := cfgAdapter.GetConfig()
			if err != nil {
				// Fallback to default if config not found
				// fmt.Println("Config not found, using default (camelCase). Run 'lit config' to configure.")
				cfg = &domain.Config{NamingConventionIndex: 1}
			}

			pattern := domain.NamingConventions[cfg.NamingConventionIndex]

			service := scanner.NewScanService(cfg, pattern, fix)

			cwd, _ := os.Getwd()

			if loc {
				results, _ := service.CalculateLOC(cwd)
				printLOC(results)
				return
			}

			results, _ := service.ScanProject(cwd)
			printDangerousFunctions(results)
		},
	}
	command.Flags().Bool("loc", false, "Retrieves the languages used with statistics")
	command.Flags().Bool("fix", false, "Fixes up the variables with an invalid naming conventions."+
		"It only one convention to another\nExample: if you have variables snake_case and the active convention is camelCase, it's converted.")

	return command
}

func printDangerousFunctions(results map[string][]*domain.FunctionData) {
	count := 0
	for _, v := range results {
		count += len(v)
	}
	fmt.Printf("\nDangerous functions found in the project: %d\n", count)
	for key, value := range results {
		fmt.Printf("- %s:\n", key)
		for _, item := range value {
			fmt.Printf(" * Function %s (at %d:%d)\n", item.Name, item.StartPosition.Row, item.StartPosition.Column)
			fmt.Printf("   Parameters: %d\n   Total lines of code: %d\n", item.TotalParams, item.Size)
			fmt.Println(item.Feedback)
		}
		fmt.Println()
	}
}

func printLOC(results map[string]int) {
	fmt.Println()
	fmt.Println("Results (language -> total lines of code):")
	total := 0.0
	// Get the total lines of code
	for _, v := range results {
		total += float64(v)
	}

	// Get the results
	for k, v := range results {
		fmt.Printf("%s %d (%.1f%%)\n", k, v, (float64(v)*100)/total)
	}
	fmt.Println("Total lines of code:", total)
}
