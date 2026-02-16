// SPDX-License-Identifier: Apache-2.0
// Copyright 2024-2026 CAB Ingénierie / Christophe ABOULICAM
package deploy

import (
	"fmt"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/stoa-platform/stoactl/internal/client"
	"github.com/stoa-platform/stoactl/internal/output"
	"github.com/stoa-platform/stoactl/internal/types"
)

var (
	outputFormat string
	apiID        string
	environment  string
	version      string
	gatewayID    string
	commitSHA    string
	watch        bool
	statusFilter string
	page         int
	pageSize     int
)

// NewDeployCmd creates the deploy command group
func NewDeployCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "deploy",
		Short: "Manage API deployments",
		Long: `Manage API deployments across environments and gateways.

Examples:
  # Deploy an API to dev
  stoactl deploy create --api-id api-1 --env dev --version 1.0.0

  # List deployments
  stoactl deploy list

  # List deployments for a specific environment
  stoactl deploy list --env production

  # Get deployment details
  stoactl deploy get <deployment-id>

  # Rollback a deployment
  stoactl deploy rollback <deployment-id>`,
	}

	cmd.PersistentFlags().StringVarP(&outputFormat, "output", "o", "table", "Output format: table, wide, yaml, json")

	cmd.AddCommand(newDeployCreateCmd())
	cmd.AddCommand(newDeployListCmd())
	cmd.AddCommand(newDeployGetCmd())
	cmd.AddCommand(newDeployRollbackCmd())

	return cmd
}

func newDeployCreateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new deployment",
		Long: `Create a new API deployment to a target environment.

Examples:
  stoactl deploy create --api-id api-1 --env dev --version 1.0.0
  stoactl deploy create --api-id api-1 --env production --version 2.0.0 --gateway gw-1
  stoactl deploy create --api-id api-1 --env dev --version 1.0.0 --watch`,
		RunE: runDeployCreate,
	}

	cmd.Flags().StringVar(&apiID, "api-id", "", "API ID to deploy (required)")
	cmd.Flags().StringVar(&environment, "env", "", "Target environment (required)")
	cmd.Flags().StringVar(&version, "version", "", "Version to deploy (required)")
	cmd.Flags().StringVar(&gatewayID, "gateway", "", "Target gateway ID (optional)")
	cmd.Flags().StringVar(&commitSHA, "commit", "", "Git commit SHA (optional)")
	cmd.Flags().BoolVarP(&watch, "watch", "w", false, "Watch deployment until completion")

	_ = cmd.MarkFlagRequired("api-id")
	_ = cmd.MarkFlagRequired("env")
	_ = cmd.MarkFlagRequired("version")

	return cmd
}

func newDeployListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "List deployments",
		Long: `List deployments with optional filters.

Examples:
  stoactl deploy list
  stoactl deploy list --env production
  stoactl deploy list --api-id api-1 --status success
  stoactl deploy list -o wide`,
		RunE: runDeployList,
	}

	cmd.Flags().StringVar(&apiID, "api-id", "", "Filter by API ID")
	cmd.Flags().StringVar(&environment, "env", "", "Filter by environment")
	cmd.Flags().StringVar(&statusFilter, "status", "", "Filter by status (pending, in_progress, success, failed, rolled_back)")
	cmd.Flags().IntVar(&page, "page", 1, "Page number")
	cmd.Flags().IntVar(&pageSize, "page-size", 20, "Items per page")

	return cmd
}

func newDeployGetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "get <deployment-id>",
		Short: "Get deployment details",
		Long: `Get detailed information about a specific deployment.

Examples:
  stoactl deploy get deploy-abc123`,
		Args: cobra.ExactArgs(1),
		RunE: runDeployGet,
	}
}

func newDeployRollbackCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "rollback <deployment-id>",
		Short: "Rollback a deployment",
		Long: `Rollback a deployment to its previous version.

Examples:
  stoactl deploy rollback deploy-abc123`,
		Args: cobra.ExactArgs(1),
		RunE: runDeployRollback,
	}
}

func runDeployCreate(cmd *cobra.Command, _ []string) error {
	c, err := client.New()
	if err != nil {
		return err
	}

	create := &types.DeploymentCreate{
		APIID:       apiID,
		Environment: environment,
		Version:     version,
		GatewayID:   gatewayID,
		CommitSHA:   commitSHA,
	}

	dep, err := c.CreateDeployment(create)
	if err != nil {
		return err
	}

	output.Success("Deployment %s created (status: %s)", dep.ID, dep.Status)

	if watch {
		return watchDeployment(c, dep.ID)
	}

	return nil
}

func runDeployList(_ *cobra.Command, _ []string) error {
	c, err := client.New()
	if err != nil {
		return err
	}

	resp, err := c.ListDeployments(apiID, environment, statusFilter, page, pageSize)
	if err != nil {
		return err
	}

	if len(resp.Items) == 0 {
		output.Info("No deployments found.")
		return nil
	}

	format := output.ParseFormat(outputFormat)
	printer := output.NewPrinter(format)

	switch format {
	case output.FormatYAML:
		return printer.PrintYAML(resp.Items)
	case output.FormatJSON:
		return printer.PrintJSON(resp.Items)
	case output.FormatWide:
		headers := []string{"ID", "API", "ENV", "VERSION", "STATUS", "DEPLOYED BY", "CREATED", "COMPLETED", "ATTEMPTS"}
		var rows [][]string
		for _, d := range resp.Items {
			completed := "-"
			if d.CompletedAt != "" {
				completed = formatTime(d.CompletedAt)
			}
			rows = append(rows, []string{
				d.ID, d.APIName, strings.ToUpper(d.Environment), d.Version,
				statusIcon(d.Status), d.DeployedBy, formatTime(d.CreatedAt),
				completed, fmt.Sprintf("%d", d.AttemptCount),
			})
		}
		printer.PrintTable(headers, rows)
	default:
		headers := []string{"ID", "API", "ENV", "VERSION", "STATUS", "CREATED"}
		var rows [][]string
		for _, d := range resp.Items {
			rows = append(rows, []string{
				d.ID, d.APIName, strings.ToUpper(d.Environment), d.Version,
				statusIcon(d.Status), formatTime(d.CreatedAt),
			})
		}
		printer.PrintTable(headers, rows)
	}

	if resp.Total > resp.PageSize {
		totalPages := (resp.Total + resp.PageSize - 1) / resp.PageSize
		output.Info("Page %d/%d (total: %d)", resp.Page, totalPages, resp.Total)
	}

	return nil
}

func runDeployGet(_ *cobra.Command, args []string) error {
	c, err := client.New()
	if err != nil {
		return err
	}

	dep, err := c.GetDeployment(args[0])
	if err != nil {
		return err
	}

	format := output.ParseFormat(outputFormat)
	printer := output.NewPrinter(format)

	switch format {
	case output.FormatYAML:
		return printer.PrintYAML(dep)
	case output.FormatJSON:
		return printer.PrintJSON(dep)
	default:
		headers := []string{"FIELD", "VALUE"}
		rows := [][]string{
			{"ID", dep.ID},
			{"API", dep.APIName},
			{"Environment", strings.ToUpper(dep.Environment)},
			{"Version", dep.Version},
			{"Status", statusIcon(dep.Status)},
			{"Deployed By", dep.DeployedBy},
			{"Created", dep.CreatedAt},
			{"Updated", dep.UpdatedAt},
			{"Attempts", fmt.Sprintf("%d", dep.AttemptCount)},
		}
		if dep.CompletedAt != "" {
			rows = append(rows, []string{"Completed", dep.CompletedAt})
		}
		if dep.ErrorMessage != "" {
			rows = append(rows, []string{"Error", dep.ErrorMessage})
		}
		if dep.RollbackOf != "" {
			rows = append(rows, []string{"Rollback Of", dep.RollbackOf})
		}
		if dep.RollbackVersion != "" {
			rows = append(rows, []string{"Rollback Version", dep.RollbackVersion})
		}
		if dep.GatewayID != "" {
			rows = append(rows, []string{"Gateway", dep.GatewayID})
		}
		if dep.CommitSHA != "" {
			rows = append(rows, []string{"Commit", dep.CommitSHA})
		}
		printer.PrintTable(headers, rows)
	}

	return nil
}

func runDeployRollback(_ *cobra.Command, args []string) error {
	c, err := client.New()
	if err != nil {
		return err
	}

	dep, err := c.RollbackDeployment(args[0])
	if err != nil {
		return err
	}

	output.Success("Rollback initiated: %s → v%s (status: %s)", dep.ID, dep.Version, dep.Status)
	return nil
}

func watchDeployment(c *client.Client, deploymentID string) error {
	output.Info("Watching deployment %s...", deploymentID)

	for i := 0; i < 60; i++ {
		time.Sleep(2 * time.Second)

		dep, err := c.GetDeployment(deploymentID)
		if err != nil {
			return err
		}

		fmt.Printf("  %s %s (attempt %d)\n", statusIcon(dep.Status), dep.Status, dep.AttemptCount)

		switch dep.Status {
		case "success":
			output.Success("Deployment completed successfully!")
			return nil
		case "failed":
			if dep.ErrorMessage != "" {
				return fmt.Errorf("deployment failed: %s", dep.ErrorMessage)
			}
			return fmt.Errorf("deployment failed")
		case "rolled_back":
			return fmt.Errorf("deployment was rolled back")
		}
	}

	return fmt.Errorf("deployment watch timed out after 2 minutes")
}

func statusIcon(status string) string {
	switch status {
	case "success":
		return "OK " + status
	case "failed":
		return "!! " + status
	case "in_progress":
		return ".. " + status
	case "pending":
		return "-- " + status
	case "rolled_back":
		return "<< " + status
	default:
		return "?? " + status
	}
}

func formatTime(t string) string {
	parsed, err := time.Parse(time.RFC3339, t)
	if err != nil {
		return t
	}
	return parsed.Format("2006-01-02 15:04")
}
