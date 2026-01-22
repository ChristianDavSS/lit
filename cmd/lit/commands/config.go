package commands

import (
	"CLI_App/internal/adapter/config"
	"CLI_App/internal/domain"
	"fmt"

	"github.com/spf13/cobra"
)

func ConfigCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "config",
		Short: "Configurate the scan variables.",
		Run: func(cmd *cobra.Command, args []string) {
			index := GetNamingConvention()
			cfg := &domain.Config{NamingConventionIndex: index}

			adapter := config.NewJSONConfigAdapter()
			err := adapter.SaveConfig(cfg)
			if err != nil {
				fmt.Println("Error saving config:", err)
				return
			}
			fmt.Println("Configuration updated.")
		},
	}
}
