package registry

import (
	"github.com/JahidNishat/docktab/internal/command"
	"github.com/JahidNishat/docktab/internal/commands/ps"
)

func All() []command.Command {
	return []command.Command{
		ps.New(),
		// Add more commands here:
		// images.New(),
		// volumes.New(),
	}
}
