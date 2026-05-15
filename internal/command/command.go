package command

import "github.com/spf13/cobra"

type Command interface {
	Name() string
	Build() *cobra.Command
}
