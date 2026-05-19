package networks

import (
	"log/slog"
	"sort"

	"github.com/JahidNishat/docktab/internal/commands"
	"github.com/spf13/cobra"

	"github.com/JahidNishat/docktab/internal/docker"
	"github.com/JahidNishat/docktab/internal/table"
)

var allowedSortFields = []string{
	"name",
	"driver",
	"created",
}

type Command struct {
	client   docker.Client
	renderer table.Renderer
	log      *slog.Logger
}

func New(client docker.Client, renderer table.Renderer, log *slog.Logger) Command {
	return Command{
		client:   client,
		renderer: renderer,
		log:      log,
	}
}

func (c Command) Name() string {
	return "networks"
}

func (c Command) Build() *cobra.Command {
	view := commands.NewViewOptions("name", allowedSortFields)

	cmd := &cobra.Command{
		Use:   "networks",
		Short: "Display Docker networks in a clean table",
		Args:  cobra.NoArgs,
		PreRunE: func(cmd *cobra.Command, args []string) error { // NEW
			return view.Validate()
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			c.log.Debug(
				"listing networks",
				"sort", view.Sort,
				"compact", view.Compact,
				"full", view.Full,
			)

			networks, err := c.client.ListNetworks(ctx)
			if err != nil {
				c.log.Error("failed to list networks", "error", err)
				return err
			}

			sorted := sortNetworks(networks, view.Sort)
			columns := getNetworkColumns(view.Compact, view.Full)

			if view.IsJSON() {
				return commands.RenderJSON(cmd.OutOrStdout(), sorted)
			}

			c.renderer.RenderNetworks(sorted, columns, c.log)
			return nil
		},
	}

	view.AddFlags(cmd)

	return cmd
}

func sortNetworks(networks []docker.Network, sortBy string) []docker.Network {
	switch sortBy {
	case "name":
		sort.Slice(networks, func(i, j int) bool {
			return networks[i].Name < networks[j].Name
		})
	case "driver":
		sort.Slice(networks, func(i, j int) bool {
			return networks[i].Driver < networks[j].Driver
		})
	case "created":
		sort.Slice(networks, func(i, j int) bool {
			return networks[i].Created.After(networks[j].Created)
		})
	}
	return networks
}

func getNetworkColumns(compact, full bool) []string {
	if compact {
		return []string{"Name", "Driver", "Scope"}
	}
	if full {
		return []string{"Name", "ID", "Driver", "Scope", "Created"}
	}
	return []string{"Name", "Driver", "Scope"}
}
