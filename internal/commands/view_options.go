package commands

import (
	"encoding/json"
	"io"
	"strings"

	"github.com/spf13/cobra"
)

const (
	OutputTable = "table"
	OutputJSON  = "json"
)

var allowedOutputFormats = []string{
	OutputTable,
	OutputJSON,
}

type ViewOptions struct {
	Compact bool
	Full    bool
	Output  string
	Sort    string

	allowedSortFields []string
}

func NewViewOptions(defaultSort string, allowedSortFields []string) *ViewOptions {
	return &ViewOptions{
		Output:            OutputTable,
		Sort:              defaultSort,
		allowedSortFields: allowedSortFields,
	}
}

func (o *ViewOptions) AddFlags(cmd *cobra.Command) {
	cmd.Flags().BoolVar(&o.Compact, "compact", false, "Compact view")
	cmd.Flags().BoolVar(&o.Full, "full", false, "Full view")
	cmd.Flags().StringVar(&o.Sort, "sort", o.Sort, "Sort by: "+strings.Join(o.allowedSortFields, ", "))
	cmd.Flags().StringVarP(&o.Output, "output", "o", OutputTable, "Output format: table, json")
}

func (o *ViewOptions) Validate() error {
	if err := ValidateAllowedValue("sort field", o.Sort, o.allowedSortFields); err != nil {
		return err
	}

	if err := ValidateAllowedValue("output format", o.Output, allowedOutputFormats); err != nil {
		return err
	}

	return ValidateMutuallyExclusive(o.Compact, "--compact", o.Full, "--full")
}

func (o *ViewOptions) IsJSON() bool {
	return o.Output == OutputJSON
}

func RenderJSON(w io.Writer, value any) error {
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")
	return encoder.Encode(value)
}
