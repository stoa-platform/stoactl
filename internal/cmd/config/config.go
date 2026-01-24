package config

import (
	"github.com/spf13/cobra"
)

// NewConfigCmd creates the config command
func NewConfigCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config",
		Short: "Manage stoactl configuration",
		Long:  `Manage stoactl contexts and configuration settings.`,
	}

	cmd.AddCommand(newSetContextCmd())
	cmd.AddCommand(newUseContextCmd())
	cmd.AddCommand(newGetContextsCmd())
	cmd.AddCommand(newCurrentContextCmd())

	return cmd
}
