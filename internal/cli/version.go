package cli

import (
	"fmt"

	"github.com/JahidNishat/docktab/internal/version"
	"github.com/spf13/cobra"
)

func NewVersionCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Display docktab version information",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("docktab %s\n", version.Version)
			fmt.Printf("commit: %s\n", version.Commit)
			fmt.Printf("built: %s\n", version.Date)
		},
	}
}
