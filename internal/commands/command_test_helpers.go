package commands

import (
	"bytes"
	"testing"

	"github.com/spf13/cobra"
)

func ExecuteCommandForTest(t *testing.T, cmd *cobra.Command, args ...string) (string, string, error) {
	t.Helper()

	var stdout bytes.Buffer
	var stderr bytes.Buffer

	cmd.SetArgs(args)
	cmd.SetOut(&stdout)
	cmd.SetErr(&stderr)

	err := cmd.Execute()

	return stdout.String(), stderr.String(), err
}
