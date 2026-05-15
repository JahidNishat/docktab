package ps

import (
	"log/slog"
	"sort"
	"strings"

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
	return "ps"
}

func (c Command) Build() *cobra.Command {
	var (
		all        bool
		compact    bool
		full       bool
		nameFilter string
		sortBy     string
	)

	cmd := &cobra.Command{
		Use:   "ps",
		Short: "Display Docker containers in a clean, beautiful table",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			c.log.Debug(
				"listing containers",
				"all", all,
				"name_filter", nameFilter,
				"sort", sortBy,
				"compact", compact,
				"full", full,
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
			sorted := sortContainers(filtered, sortBy)
			columns := getColumns(compact, full)

			c.log.Debug(
				"rendering containers table",
				"count", len(sorted),
				"columns", columns,
			)

			c.renderer.RenderContainers(sorted, columns, c.log)
			return nil
		},
	}

	cmd.Flags().BoolVarP(&all, "all", "a", false, "Show all containers (default shows only running)")
	cmd.Flags().BoolVar(&compact, "compact", false, "Compact view")
	cmd.Flags().BoolVar(&full, "full", false, "Full view with more columns")
	cmd.Flags().StringVar(&nameFilter, "name", "", "Filter containers by name")
	cmd.Flags().StringVar(&sortBy, "sort", "name", "Sort by: name, image, status, created")

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
