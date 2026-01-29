package commands

import (
	"CLI_App/cmd/adapters/config"
	"CLI_App/cmd/domain"
	"CLI_App/cmd/lit/ui"

	"github.com/spf13/cobra"
)

func Configuration() *cobra.Command {
	return &cobra.Command{
		Use:   "config",
		Short: "Configure the scan variables.",
		Run: func(cmd *cobra.Command, args []string) {
			idx := ui.GetNamingConvention()
			alerts := ui.GetAlertsConfig()
			jsonAdapter := config.NewJSONAdapter()
			newConfig := &domain.Config{NamingConventionIndex: idx, Alerts: alerts}
			jsonAdapter.SaveConfig(newConfig)
		},
	}
}
