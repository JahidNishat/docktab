package cli

import (
	"log/slog"

	"github.com/JahidNishat/docktab/internal/docker"
	"github.com/JahidNishat/docktab/internal/registry"
	"github.com/JahidNishat/docktab/internal/table"
	"github.com/spf13/cobra"
)

func NewRootCommand(
	client docker.Client,
	renderer table.Renderer,
	log *slog.Logger,
) *cobra.Command {
	root := &cobra.Command{
		Use:   "docktab",
		Short: "A clean Docker CLI table viewer",
	}

	for _, cmd := range registry.All(client, renderer, log) {
		root.AddCommand(cmd.Build())
	}

	return root
}
