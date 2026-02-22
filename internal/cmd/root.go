// SPDX-License-Identifier: Apache-2.0
// Copyright 2024-2026 CAB Ingénierie / Christophe ABOULICAM
package cmd

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/stoa-platform/stoactl/internal/cmd/apply"
	"github.com/stoa-platform/stoactl/internal/cmd/auth"
	bridgecmd "github.com/stoa-platform/stoactl/internal/cmd/bridge"
	"github.com/stoa-platform/stoactl/internal/cmd/config"
	"github.com/stoa-platform/stoactl/internal/cmd/delete"
	"github.com/stoa-platform/stoactl/internal/cmd/deploy"
	"github.com/stoa-platform/stoactl/internal/cmd/doctor"
	"github.com/stoa-platform/stoactl/internal/cmd/get"
	initcmd "github.com/stoa-platform/stoactl/internal/cmd/init"
	"github.com/stoa-platform/stoactl/internal/cmd/logs"
	"github.com/stoa-platform/stoactl/internal/cmd/token_usage"
)

var (
	// Version is set at build time
	Version = "dev"
	// Commit is set at build time
	Commit = "none"
)

var rootCmd = &cobra.Command{
	Use:   "stoactl",
	Short: "STOA Platform CLI",
	Long: `stoactl is a GitOps-native CLI for managing STOA Platform resources.

It provides a declarative, kubectl-like experience for managing APIs,
subscriptions, and other STOA resources through infrastructure-as-code patterns.

Example usage:
  stoactl config set-context prod --server=https://api.gostoa.dev --tenant=acme
  stoactl config use-context prod
  stoactl auth login
  stoactl get apis
  stoactl apply -f api.yaml`,
	SilenceUsage:  true,
	SilenceErrors: true,
}

// Execute runs the root command
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.AddCommand(config.NewConfigCmd())
	rootCmd.AddCommand(auth.NewAuthCmd())
	rootCmd.AddCommand(get.NewGetCmd())
	rootCmd.AddCommand(apply.NewApplyCmd())
	rootCmd.AddCommand(delete.NewDeleteCmd())
	rootCmd.AddCommand(deploy.NewDeployCmd())
	rootCmd.AddCommand(logs.NewLogsCmd())
	rootCmd.AddCommand(token_usage.NewTokenUsageCmd())
	rootCmd.AddCommand(bridgecmd.NewBridgeCmd())
	rootCmd.AddCommand(initcmd.NewInitCmd())
	rootCmd.AddCommand(doctor.NewDoctorCmd())
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version information",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Printf("stoactl version %s (%s)\n", Version, Commit)
	},
}

// ExitCode constants following ADR-001
const (
	ExitSuccess           = 0
	ExitGeneralError      = 1
	ExitMisuse            = 2
	ExitAuthFailed        = 3
	ExitResourceNotFound  = 4
	ExitConflict          = 5
	ExitValidationError   = 6
)

// Exit exits with the specified code
func Exit(code int) {
	os.Exit(code)
}
