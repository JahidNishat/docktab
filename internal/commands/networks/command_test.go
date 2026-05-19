package networks

import (
	"strings"
	"testing"
	"time"

	"github.com/JahidNishat/docktab/internal/commands"
	"github.com/JahidNishat/docktab/internal/docker"
)

func TestSortNetworks(t *testing.T) {
	baseTime := time.Date(2026, 5, 15, 12, 0, 0, 0, time.UTC)

	tests := []struct {
		name      string
		sortBy    string
		input     []docker.Network
		wantNames []string
	}{
		{
			name:   "sorts by name ascending",
			sortBy: "name",
			input: []docker.Network{
				{Name: "zeta-net"},
				{Name: "alpha-net"},
				{Name: "beta-net"},
			},
			wantNames: []string{"alpha-net", "beta-net", "zeta-net"},
		},
		{
			name:   "sorts by driver ascending",
			sortBy: "driver",
			input: []docker.Network{
				{Name: "overlay-net", Driver: "overlay"},
				{Name: "bridge-net", Driver: "bridge"},
				{Name: "macvlan-net", Driver: "macvlan"},
			},
			wantNames: []string{"bridge-net", "macvlan-net", "overlay-net"},
		},
		{
			name:   "sorts by created newest first",
			sortBy: "created",
			input: []docker.Network{
				{Name: "old", Created: baseTime.Add(-48 * time.Hour)},
				{Name: "new", Created: baseTime},
				{Name: "middle", Created: baseTime.Add(-12 * time.Hour)},
			},
			wantNames: []string{"new", "middle", "old"},
		},
		{
			name:   "unknown sort value keeps original order",
			sortBy: "unsupported",
			input: []docker.Network{
				{Name: "first"},
				{Name: "second"},
				{Name: "third"},
			},
			wantNames: []string{"first", "second", "third"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input := cloneNetworks(tt.input)

			got := sortNetworks(input, tt.sortBy)

			if len(got) != len(tt.wantNames) {
				t.Fatalf("expected %d networks, got %d", len(tt.wantNames), len(got))
			}

			for i, wantName := range tt.wantNames {
				if got[i].Name != wantName {
					t.Fatalf("expected network %d to be %q, got %q", i, wantName, got[i].Name)
				}
			}
		})
	}
}

func TestGetNetworkColumns(t *testing.T) {
	tests := []struct {
		name    string
		compact bool
		full    bool
		want    []string
	}{
		{
			name:    "default columns",
			compact: false,
			full:    false,
			want:    []string{"Name", "Driver", "Scope"},
		},
		{
			name:    "compact columns",
			compact: true,
			full:    false,
			want:    []string{"Name", "Driver", "Scope"},
		},
		{
			name:    "full columns",
			compact: false,
			full:    true,
			want:    []string{"Name", "ID", "Driver", "Scope", "Created"},
		},
		{
			name:    "compact takes precedence if both flags are true",
			compact: true,
			full:    true,
			want:    []string{"Name", "Driver", "Scope"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := getNetworkColumns(tt.compact, tt.full)

			if len(got) != len(tt.want) {
				t.Fatalf("expected %d columns, got %d", len(tt.want), len(got))
			}

			for i, wantColumn := range tt.want {
				if got[i] != wantColumn {
					t.Fatalf("expected column %d to be %q, got %q", i, wantColumn, got[i])
				}
			}
		})
	}
}

func cloneNetworks(input []docker.Network) []docker.Network {
	cloned := make([]docker.Network, len(input))
	copy(cloned, input)
	return cloned
}

func TestCommandRejectsUnexpectedArgs(t *testing.T) {
	cmd := New(nil, nil, nil).Build()

	_, _, err := commands.ExecuteCommandForTest(t, cmd, "hello")
	if err == nil {
		t.Fatal("expected error for unexpected argument")
	}

	if !strings.Contains(err.Error(), `unknown command "hello"`) {
		t.Fatalf("expected unexpected argument error, got: %v", err)
	}
}

func TestCommandRejectsInvalidSortField(t *testing.T) {
	cmd := New(nil, nil, nil).Build()

	_, _, err := commands.ExecuteCommandForTest(t, cmd, "--sort", "banana")
	if err == nil {
		t.Fatal("expected error for invalid sort field")
	}

	if !strings.Contains(err.Error(), `invalid sort field "banana"`) {
		t.Fatalf("expected invalid sort error, got: %v", err)
	}
}

func TestCommandRejectsCompactAndFullTogether(t *testing.T) {
	cmd := New(nil, nil, nil).Build()

	_, _, err := commands.ExecuteCommandForTest(t, cmd, "--compact", "--full")
	if err == nil {
		t.Fatal("expected error for conflicting flags")
	}

	if !strings.Contains(err.Error(), "--compact and --full cannot be used together") {
		t.Fatalf("expected conflicting flags error, got: %v", err)
	}
}

func TestCommandRejectsInvalidOutputFormat(t *testing.T) {
	cmd := New(nil, nil, nil).Build()

	_, _, err := commands.ExecuteCommandForTest(t, cmd, "--output", "xml")
	if err == nil {
		t.Fatal("expected error for invalid output format")
	}

	if !strings.Contains(err.Error(), `invalid output format "xml"`) {
		t.Fatalf("expected invalid output format error, got: %v", err)
	}
}

func TestCommandRejectsInvalidOutputFormatShortFlag(t *testing.T) {
	cmd := New(nil, nil, nil).Build()

	_, _, err := commands.ExecuteCommandForTest(t, cmd, "-o", "xml")
	if err == nil {
		t.Fatal("expected error for invalid output format")
	}

	if !strings.Contains(err.Error(), `invalid output format "xml"`) {
		t.Fatalf("expected invalid output format error, got: %v", err)
	}
}
