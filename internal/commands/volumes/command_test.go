package volumes

import (
	"testing"
	"time"

	"github.com/JahidNishat/docktab/internal/docker"
)

func TestSortVolumes(t *testing.T) {
	baseTime := time.Date(2026, 5, 15, 12, 0, 0, 0, time.UTC)

	tests := []struct {
		name      string
		sortBy    string
		input     []docker.Volume
		wantNames []string
	}{
		{
			name:   "sorts by name ascending",
			sortBy: "name",
			input: []docker.Volume{
				{Name: "redis-data"},
				{Name: "api-data"},
				{Name: "postgres-data"},
			},
			wantNames: []string{"api-data", "postgres-data", "redis-data"},
		},
		{
			name:   "sorts by driver ascending",
			sortBy: "driver",
			input: []docker.Volume{
				{Name: "vol-z", Driver: "zfs"},
				{Name: "vol-l", Driver: "local"},
				{Name: "vol-b", Driver: "btrfs"},
			},
			wantNames: []string{"vol-b", "vol-l", "vol-z"},
		},
		{
			name:   "sorts by created newest first",
			sortBy: "created",
			input: []docker.Volume{
				{Name: "old", Created: baseTime.Add(-48 * time.Hour)},
				{Name: "new", Created: baseTime},
				{Name: "middle", Created: baseTime.Add(-12 * time.Hour)},
			},
			wantNames: []string{"new", "middle", "old"},
		},
		{
			name:   "unknown sort value keeps original order",
			sortBy: "unsupported",
			input: []docker.Volume{
				{Name: "first"},
				{Name: "second"},
				{Name: "third"},
			},
			wantNames: []string{"first", "second", "third"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input := cloneVolumes(tt.input)

			got := sortVolumes(input, tt.sortBy)

			if len(got) != len(tt.wantNames) {
				t.Fatalf("expected %d volumes, got %d", len(tt.wantNames), len(got))
			}

			for i, wantName := range tt.wantNames {
				if got[i].Name != wantName {
					t.Fatalf("expected volume %d to be %q, got %q", i, wantName, got[i].Name)
				}
			}
		})
	}
}

func TestGetVolumeColumns(t *testing.T) {
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
			want:    []string{"Name", "Driver", "Created"},
		},
		{
			name:    "compact columns",
			compact: true,
			full:    false,
			want:    []string{"Name", "Driver"},
		},
		{
			name:    "full columns",
			compact: false,
			full:    true,
			want:    []string{"Name", "Driver", "Mountpoint", "Created"},
		},
		{
			name:    "compact takes precedence if both flags are true",
			compact: true,
			full:    true,
			want:    []string{"Name", "Driver"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := getVolumeColumns(tt.compact, tt.full)

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

func cloneVolumes(input []docker.Volume) []docker.Volume {
	cloned := make([]docker.Volume, len(input))
	copy(cloned, input)
	return cloned
}
