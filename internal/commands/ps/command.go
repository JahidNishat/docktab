package ps

import (
	"github.com/spf13/cobra"
)

type Command struct{}

func New() Command {
	return Command{}
}

func (c Command) Name() string {
	return "ps"
}

func (c Command) Build() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ps",
		Short: "Display Docker containers in a clean, beautiful table",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runPs(cmd)
		},
	}

	cmd.Flags().BoolP("all", "a", false, "Show all containers")
	cmd.Flags().Bool("compact", false, "Compact view")
	cmd.Flags().Bool("full", false, "Full view with all columns")

	return cmd
}

func runPs(cmd *cobra.Command) error {
	cmd.Println("docktab ps - Beautiful table coming in next step!")
	return nil
}
