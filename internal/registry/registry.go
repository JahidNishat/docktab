package registry

import (
	"log/slog"

	"github.com/JahidNishat/docktab/internal/command"
	"github.com/JahidNishat/docktab/internal/commands/images"
	"github.com/JahidNishat/docktab/internal/commands/networks"
	"github.com/JahidNishat/docktab/internal/commands/ps"
	"github.com/JahidNishat/docktab/internal/commands/volumes"
	"github.com/JahidNishat/docktab/internal/docker"
	"github.com/JahidNishat/docktab/internal/table"
)

func All(client docker.Client, renderer table.Renderer, log *slog.Logger) []command.Command {
	return []command.Command{
		ps.New(client, renderer, log),
		images.New(client, renderer, log),
		volumes.New(client, renderer, log),
		networks.New(client, renderer, log),
	}
}
