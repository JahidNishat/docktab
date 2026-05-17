package ps

import (
	"log/slog"
	"sort"
	"strings"

	"github.com/JahidNishat/docktab/internal/commands"
	"github.com/spf13/cobra"

	"github.com/JahidNishat/docktab/internal/docker"
	"github.com/JahidNishat/docktab/internal/table"
)

var allowedSortFields = []string{
	"name",
	"image",
	"status",
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
	return "ps"
}

func (c Command) Build() *cobra.Command {
	var (
		all        bool
		nameFilter string
	)
	view := commands.NewViewOptions("name", allowedSortFields)

	cmd := &cobra.Command{
		Use:   "ps",
		Short: "Display Docker containers in a clean, beautiful table",
		Args:  cobra.NoArgs,
		PreRunE: func(cmd *cobra.Command, args []string) error { // NEW
			return view.Validate()
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			c.log.Debug(
				"listing containers",
				"all", all,
				"name_filter", nameFilter,
				"sort", view.Sort,
				"compact", view.Compact,
				"full", view.Full,
			)

			containers, err := c.client.ListContainers(ctx, all)
			if err != nil {
				c.log.Error("failed to list containers", "error", err)
				return err
			}
			c.log.Debug("containers fetched", "count", len(containers))

			// Apply name filter
			filtered := filterByName(containers, nameFilter)
			c.log.Debug("containers filtered", "count", len(filtered))

			// Apply sorting
			sorted := sortContainers(filtered, view.Sort)
			columns := getColumns(view.Compact, view.Full)

			if view.IsJSON() {
				return commands.RenderJSON(cmd.OutOrStdout(), sorted)
			}

			c.renderer.RenderContainers(sorted, columns, c.log)
			return nil
		},
	}

	cmd.Flags().BoolVarP(&all, "all", "a", false, "Show all containers (default shows only running)")
	cmd.Flags().StringVar(&nameFilter, "name", "", "Filter containers by name")
	view.AddFlags(cmd)

	return cmd
}

func filterByName(containers []docker.Container, name string) []docker.Container {
	if name == "" {
		return containers
	}

	var result []docker.Container
	for _, c := range containers {
		if strings.Contains(strings.ToLower(c.Name), strings.ToLower(name)) {
			result = append(result, c)
		}
	}
	return result
}

func sortContainers(containers []docker.Container, sortBy string) []docker.Container {
	switch sortBy {
	case "name":
		sort.Slice(containers, func(i, j int) bool {
			return containers[i].Name < containers[j].Name
		})
	case "image":
		sort.Slice(containers, func(i, j int) bool {
			return containers[i].Image < containers[j].Image
		})
	case "status":
		sort.Slice(containers, func(i, j int) bool {
			return containers[i].Status < containers[j].Status
		})
	case "created":
		sort.Slice(containers, func(i, j int) bool {
			return containers[i].Created.After(containers[j].Created)
		})
	}
	return containers
}

func getColumns(compact, full bool) []string {
	if compact {
		return []string{"ID", "Name", "Image", "Status"}
	}
	if full {
		return []string{"ID", "Name", "Image", "Command", "Created", "Status", "Ports"}
	}
	return []string{"ID", "Name", "Image", "Status", "Ports", "Created"}
}
