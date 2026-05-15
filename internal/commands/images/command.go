package images

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
	return "images"
}

func (c Command) Build() *cobra.Command {
	var (
		all     bool
		compact bool
		full    bool
		sortBy  string
	)

	cmd := &cobra.Command{
		Use:   "images",
		Short: "Display Docker images in a clean table",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			c.log.Debug(
				"listing images",
				"all", all,
				"sort", sortBy,
				"compact", compact,
				"full", full,
			)

			images, err := c.client.ListImages(ctx, all)
			if err != nil {
				c.log.Error("failed to list images", "error", err)
				return err
			}

			sorted := sortImages(images, sortBy)
			columns := getImageColumns(compact, full)
			c.renderer.RenderImages(sorted, columns, c.log)
			return nil
		},
	}

	cmd.Flags().BoolVarP(&all, "all", "a", false, "Show all images (including intermediate)")
	cmd.Flags().BoolVar(&compact, "compact", false, "Compact view")
	cmd.Flags().BoolVar(&full, "full", false, "Full view")
	cmd.Flags().StringVar(&sortBy, "sort", "size", "Sort by: repository, tag, size, created")

	return cmd
}

func sortImages(images []docker.Image, sortBy string) []docker.Image {
	switch sortBy {
	case "repository":
		sort.Slice(images, func(i, j int) bool {
			return images[i].Repository < images[j].Repository
		})
	case "tag":
		sort.Slice(images, func(i, j int) bool {
			return images[i].Tag < images[j].Tag
		})
	case "size":
		sort.Slice(images, func(i, j int) bool {
			return images[i].Size > images[j].Size
		})
	case "created":
		sort.Slice(images, func(i, j int) bool {
			return images[i].Created.After(images[j].Created)
		})
	}
	return images
}

func getImageColumns(compact, full bool) []string {
	if compact {
		return []string{"Repository", "Tag", "Size"}
	}
	if full {
		return []string{"Repository", "Tag", "Image ID", "Created", "Size"}
	}
	return []string{"Repository", "Tag", "Image ID", "Size"}
}
