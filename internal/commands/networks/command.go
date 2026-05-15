package networks

import (
	"log/slog"
	"sort"

	"github.com/spf13/cobra"

	"github.com/JahidNishat/docktab/internal/docker"
	"github.com/JahidNishat/docktab/internal/table"
)

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
	var (
		compact bool
		full    bool
		sortBy  string
	)

	cmd := &cobra.Command{
		Use:   "networks",
		Short: "Display Docker networks in a clean table",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			networks, err := c.client.ListNetworks(ctx)
			if err != nil {
				return err
			}

			sorted := sortNetworks(networks, sortBy)
			columns := getNetworkColumns(compact, full)
			c.renderer.RenderNetworks(sorted, columns, c.log)
			return nil
		},
	}

	cmd.Flags().BoolVar(&compact, "compact", false, "Compact view")
	cmd.Flags().BoolVar(&full, "full", false, "Full view")
	cmd.Flags().StringVar(&sortBy, "sort", "name", "Sort by: name, driver, created")

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
