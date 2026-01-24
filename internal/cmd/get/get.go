package get

import (
	"github.com/spf13/cobra"

	"github.com/stoa-platform/stoactl/internal/client"
	"github.com/stoa-platform/stoactl/internal/output"
	"github.com/stoa-platform/stoactl/internal/types"
)

var outputFormat string

// NewGetCmd creates the get command
func NewGetCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get",
		Short: "Display one or many resources",
		Long: `Display one or many resources.

Prints a table of the most important information about the specified resources.
You can filter the list using a NAME or use -o for different output formats.

Examples:
  # List all APIs in table format
  stoactl get apis

  # Get a specific API
  stoactl get api billing-api

  # List APIs in wide format (more columns)
  stoactl get apis -o wide

  # Get API as YAML
  stoactl get api billing-api -o yaml

  # Get APIs as JSON
  stoactl get apis -o json`,
	}

	cmd.PersistentFlags().StringVarP(&outputFormat, "output", "o", "table", "Output format: table, wide, yaml, json")

	cmd.AddCommand(newGetAPIsCmd())

	return cmd
}

func newGetAPIsCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "apis [name]",
		Aliases: []string{"api"},
		Short:   "Display APIs",
		Long: `Display one or many APIs.

Examples:
  stoactl get apis
  stoactl get api billing-api
  stoactl get apis -o yaml`,
		Args: cobra.MaximumNArgs(1),
		RunE: runGetAPIs,
	}
}

func runGetAPIs(cmd *cobra.Command, args []string) error {
	c, err := client.New()
	if err != nil {
		return err
	}

	format := output.ParseFormat(outputFormat)
	printer := output.NewPrinter(format)

	// Single API
	if len(args) == 1 {
		return getAPI(c, printer, args[0])
	}

	// List all APIs
	return listAPIs(c, printer)
}

func getAPI(c *client.Client, printer *output.Printer, name string) error {
	api, err := c.GetAPI(name)
	if err != nil {
		return err
	}

	switch printer.Format {
	case output.FormatYAML:
		resource := apiToResource(api)
		return printer.PrintYAML(resource)
	case output.FormatJSON:
		return printer.PrintJSON(api)
	default:
		headers := []string{"NAME", "VERSION", "STATUS", "PATH"}
		rows := [][]string{{api.Name, api.Version, api.Status, api.Path}}
		printer.PrintTable(headers, rows)
	}

	return nil
}

func listAPIs(c *client.Client, printer *output.Printer) error {
	resp, err := c.ListAPIs()
	if err != nil {
		return err
	}

	if len(resp.Items) == 0 {
		output.Info("No APIs found.")
		return nil
	}

	switch printer.Format {
	case output.FormatYAML:
		var resources []types.Resource
		for _, api := range resp.Items {
			resources = append(resources, apiToResource(&api))
		}
		return printer.PrintYAML(resources)
	case output.FormatJSON:
		return printer.PrintJSON(resp.Items)
	case output.FormatWide:
		headers := []string{"NAME", "VERSION", "STATUS", "PATH", "UPSTREAM", "TENANT", "CREATED"}
		var rows [][]string
		for _, api := range resp.Items {
			rows = append(rows, []string{
				api.Name,
				api.Version,
				api.Status,
				api.Path,
				api.Upstream,
				api.Tenant,
				api.CreatedAt,
			})
		}
		printer.PrintTable(headers, rows)
	default:
		headers := []string{"NAME", "VERSION", "STATUS", "PATH"}
		var rows [][]string
		for _, api := range resp.Items {
			rows = append(rows, []string{api.Name, api.Version, api.Status, api.Path})
		}
		printer.PrintTable(headers, rows)
	}

	return nil
}

func apiToResource(api *types.API) types.Resource {
	return types.Resource{
		APIVersion: "stoa.io/v1",
		Kind:       "API",
		Metadata: types.Metadata{
			Name:      api.Name,
			Namespace: api.Tenant,
		},
		Spec: types.APISpec{
			Version:     api.Version,
			Description: api.Description,
			Upstream: types.UpstreamSpec{
				URL: api.Upstream,
			},
			Routing: types.RoutingSpec{
				Path: api.Path,
			},
		},
	}
}

// GetOutputFormat returns the current output format
func GetOutputFormat() string {
	return outputFormat
}
