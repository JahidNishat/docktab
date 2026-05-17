package volumes

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
	return "volumes"
}

func (c Command) Build() *cobra.Command {
	view := commands.NewViewOptions("name", allowedSortFields)

	cmd := &cobra.Command{
		Use:   "volumes",
		Short: "Display Docker volumes in a clean table",
		Args:  cobra.NoArgs,
		PreRunE: func(cmd *cobra.Command, args []string) error { // NEW
			return view.Validate()
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			c.log.Debug(
				"listing volumes",
				"sort", view.Sort,
				"compact", view.Compact,
				"full", view.Full,
			)

			volumes, err := c.client.ListVolumes(ctx)
			if err != nil {
				return err
			}

			sorted := sortVolumes(volumes, view.Sort)
			columns := getVolumeColumns(view.Compact, view.Full)

			if view.IsJSON() {
				return commands.RenderJSON(cmd.OutOrStdout(), sorted)
			}

			c.renderer.RenderVolumes(sorted, columns, c.log)
			return nil
		},
	}

	view.AddFlags(cmd)

	return cmd
}

func sortVolumes(volumes []docker.Volume, sortBy string) []docker.Volume {
	switch sortBy {
	case "name":
		sort.Slice(volumes, func(i, j int) bool {
			return volumes[i].Name < volumes[j].Name
		})
	case "driver":
		sort.Slice(volumes, func(i, j int) bool {
			return volumes[i].Driver < volumes[j].Driver
		})
	case "created":
		sort.Slice(volumes, func(i, j int) bool {
			return volumes[i].Created.After(volumes[j].Created)
		})
	}
	return volumes
}

func getVolumeColumns(compact, full bool) []string {
	if compact {
		return []string{"Name", "Driver"}
	}
	if full {
		return []string{"Name", "Driver", "Mountpoint", "Created"}
	}
	return []string{"Name", "Driver", "Created"}
}
