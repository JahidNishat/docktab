package images

import (
	"testing"
	"time"

	"github.com/JahidNishat/docktab/internal/docker"
)

func TestSortImages(t *testing.T) {
	baseTime := time.Date(2026, 5, 15, 12, 0, 0, 0, time.UTC)

	tests := []struct {
		name             string
		sortBy           string
		input            []docker.Image
		wantRepositories []string
	}{
		{
			name:   "sorts by repository ascending",
			sortBy: "repository",
			input: []docker.Image{
				{Repository: "redis"},
				{Repository: "alpine"},
				{Repository: "postgres"},
			},
			wantRepositories: []string{"alpine", "postgres", "redis"},
		},
		{
			name:   "sorts by tag ascending",
			sortBy: "tag",
			input: []docker.Image{
				{Repository: "api", Tag: "v3"},
				{Repository: "api", Tag: "v1"},
				{Repository: "api", Tag: "v2"},
			},
			wantRepositories: []string{"api", "api", "api"},
		},
		{
			name:   "sorts by size descending",
			sortBy: "size",
			input: []docker.Image{
				{Repository: "small", Size: 100},
				{Repository: "large", Size: 900},
				{Repository: "medium", Size: 500},
			},
			wantRepositories: []string{"large", "medium", "small"},
		},
		{
			name:   "sorts by created newest first",
			sortBy: "created",
			input: []docker.Image{
				{Repository: "old", Created: baseTime.Add(-48 * time.Hour)},
				{Repository: "new", Created: baseTime},
				{Repository: "middle", Created: baseTime.Add(-12 * time.Hour)},
			},
			wantRepositories: []string{"new", "middle", "old"},
		},
		{
			name:   "unknown sort keeps original order",
			sortBy: "unsupported",
			input: []docker.Image{
				{Repository: "first"},
				{Repository: "second"},
				{Repository: "third"},
			},
			wantRepositories: []string{"first", "second", "third"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input := cloneImages(tt.input)

			got := sortImages(input, tt.sortBy)

			if len(got) != len(tt.wantRepositories) {
				t.Fatalf("expected %d images, got %d", len(tt.wantRepositories), len(got))
			}

			for i, wantRepository := range tt.wantRepositories {
				if got[i].Repository != wantRepository {
					t.Fatalf(
						"expected image %d repository to be %q, got %q",
						i,
						wantRepository,
						got[i].Repository,
					)
				}
			}
		})
	}
}

func TestSortImagesByTag(t *testing.T) {
	images := []docker.Image{
		{Repository: "api", Tag: "v3"},
		{Repository: "api", Tag: "v1"},
		{Repository: "api", Tag: "v2"},
	}

	got := sortImages(cloneImages(images), "tag")

	wantTags := []string{"v1", "v2", "v3"}

	for i, wantTag := range wantTags {
		if got[i].Tag != wantTag {
			t.Fatalf("expected image %d tag to be %q, got %q", i, wantTag, got[i].Tag)
		}
	}
}

func TestGetImageColumns(t *testing.T) {
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
			want:    []string{"Repository", "Tag", "Image ID", "Size"},
		},
		{
			name:    "compact columns",
			compact: true,
			full:    false,
			want:    []string{"Repository", "Tag", "Size"},
		},
		{
			name:    "full columns",
			compact: false,
			full:    true,
			want:    []string{"Repository", "Tag", "Image ID", "Created", "Size"},
		},
		{
			name:    "compact takes precedence if both flags are true",
			compact: true,
			full:    true,
			want:    []string{"Repository", "Tag", "Size"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := getImageColumns(tt.compact, tt.full)

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

func cloneImages(input []docker.Image) []docker.Image {
	cloned := make([]docker.Image, len(input))
	copy(cloned, input)
	return cloned
}
