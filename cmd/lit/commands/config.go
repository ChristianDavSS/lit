package commands

import (
	"CLI_App/cmd/adapters/config"
	"CLI_App/cmd/lit/ui"

	"github.com/spf13/cobra"
)

func Configuration() *cobra.Command {
	return &cobra.Command{
		Use:   "config",
		Short: "Configure the scan variables.",
		Run: func(cmd *cobra.Command, args []string) {
			idx := ui.GetNamingConvention()
			jsonAdapter := config.NewJSONAdapter()
			newConfig := &config.ConfigDto{NamingConventionIndex: idx}
			jsonAdapter.SaveConfig(newConfig)
		},
	}
}
