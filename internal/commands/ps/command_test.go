package ps

import (
	"testing"
	"time"

	"github.com/JahidNishat/docktab/internal/docker"
)

func TestFilterByName(t *testing.T) {
	containers := []docker.Container{
		{Name: "api-service"},
		{Name: "postgres-db"},
		{Name: "redis-cache"},
		{Name: "worker-api"},
	}

	tests := []struct {
		name      string
		filter    string
		wantNames []string
	}{
		{
			name:      "empty filter returns all containers",
			filter:    "",
			wantNames: []string{"api-service", "postgres-db", "redis-cache", "worker-api"},
		},
		{
			name:      "matches single container",
			filter:    "redis",
			wantNames: []string{"redis-cache"},
		},
		{
			name:      "matches multiple containers",
			filter:    "api",
			wantNames: []string{"api-service", "worker-api"},
		},
		{
			name:      "matching is case insensitive",
			filter:    "POSTGRES",
			wantNames: []string{"postgres-db"},
		},
		{
			name:      "returns empty slice when no containers match",
			filter:    "nginx",
			wantNames: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := filterByName(containers, tt.filter)

			if len(got) != len(tt.wantNames) {
				t.Fatalf("expected %d containers, got %d", len(tt.wantNames), len(got))
			}

			for i, wantName := range tt.wantNames {
				if got[i].Name != wantName {
					t.Fatalf("expected container %d to be %q, got %q", i, wantName, got[i].Name)
				}
			}
		})
	}
}

func TestSortContainers(t *testing.T) {
	baseTime := time.Date(2026, 5, 15, 12, 0, 0, 0, time.UTC)

	tests := []struct {
		name      string
		sortBy    string
		input     []docker.Container
		wantNames []string
	}{
		{
			name:   "sorts by name ascending",
			sortBy: "name",
			input: []docker.Container{
				{Name: "redis"},
				{Name: "api"},
				{Name: "postgres"},
			},
			wantNames: []string{"api", "postgres", "redis"},
		},
		{
			name:   "sorts by image ascending",
			sortBy: "image",
			input: []docker.Container{
				{Name: "worker", Image: "node:22"},
				{Name: "db", Image: "postgres:16"},
				{Name: "cache", Image: "redis:7"},
			},
			wantNames: []string{"worker", "db", "cache"},
		},
		{
			name:   "sorts by status ascending",
			sortBy: "status",
			input: []docker.Container{
				{Name: "stopped", Status: "Exited"},
				{Name: "running", Status: "Up"},
				{Name: "created", Status: "Created"},
			},
			wantNames: []string{"created", "stopped", "running"},
		},
		{
			name:   "sorts by created newest first",
			sortBy: "created",
			input: []docker.Container{
				{Name: "old", Created: baseTime.Add(-48 * time.Hour)},
				{Name: "new", Created: baseTime},
				{Name: "middle", Created: baseTime.Add(-12 * time.Hour)},
			},
			wantNames: []string{"new", "middle", "old"},
		},
		{
			name:   "unknown sort value keeps original order",
			sortBy: "unsupported",
			input: []docker.Container{
				{Name: "first"},
				{Name: "second"},
				{Name: "third"},
			},
			wantNames: []string{"first", "second", "third"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input := cloneContainers(tt.input)

			got := sortContainers(input, tt.sortBy)

			if len(got) != len(tt.wantNames) {
				t.Fatalf("expected %d containers, got %d", len(tt.wantNames), len(got))
			}

			for i, wantName := range tt.wantNames {
				if got[i].Name != wantName {
					t.Fatalf("expected container %d to be %q, got %q", i, wantName, got[i].Name)
				}
			}
		})
	}
}

func TestGetColumns(t *testing.T) {
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
			want:    []string{"ID", "Name", "Image", "Status", "Ports", "Created"},
		},
		{
			name:    "compact columns",
			compact: true,
			full:    false,
			want:    []string{"ID", "Name", "Image", "Status"},
		},
		{
			name:    "full columns",
			compact: false,
			full:    true,
			want:    []string{"ID", "Name", "Image", "Command", "Created", "Status", "Ports"},
		},
		{
			name:    "compact takes precedence if both flags are true",
			compact: true,
			full:    true,
			want:    []string{"ID", "Name", "Image", "Status"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := getColumns(tt.compact, tt.full)

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

func cloneContainers(input []docker.Container) []docker.Container {
	cloned := make([]docker.Container, len(input))
	copy(cloned, input)
	return cloned
}
