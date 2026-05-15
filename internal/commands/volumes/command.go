package volumes

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
	return "volumes"
}

func (c Command) Build() *cobra.Command {
	var (
		compact bool
		full    bool
		sortBy  string
	)

	cmd := &cobra.Command{
		Use:   "volumes",
		Short: "Display Docker volumes in a clean table",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			c.log.Debug(
				"listing volumes",
				"sort", sortBy,
				"compact", compact,
				"full", full,
			)

			volumes, err := c.client.ListVolumes(ctx)
			if err != nil {
				return err
			}

			sorted := sortVolumes(volumes, sortBy)
			columns := getVolumeColumns(compact, full)
			c.renderer.RenderVolumes(sorted, columns, c.log)
			return nil
		},
	}

	cmd.Flags().BoolVar(&compact, "compact", false, "Compact view")
	cmd.Flags().BoolVar(&full, "full", false, "Full view")
	cmd.Flags().StringVar(&sortBy, "sort", "name", "Sort by: name, driver, created")

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
