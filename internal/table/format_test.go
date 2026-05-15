package table

import (
	"testing"
	"time"
)

func TestTruncate(t *testing.T) {
	tests := []struct {
		name  string
		input string
		max   int
		want  string
	}{
		{
			name:  "returns unchanged string when shorter than max",
			input: "redis",
			max:   10,
			want:  "redis",
		},
		{
			name:  "returns unchanged string when equal to max",
			input: "postgres",
			max:   8,
			want:  "postgres",
		},
		{
			name:  "truncates long string with ellipsis",
			input: "abcdefghijklmnopqrstuvwxyz",
			max:   10,
			want:  "abcdefg...",
		},
		{
			name:  "returns empty string when max is zero",
			input: "abcdef",
			max:   0,
			want:  "",
		},
		{
			name:  "returns empty string when max is negative",
			input: "abcdef",
			max:   -1,
			want:  "",
		},
		{
			name:  "returns first characters when max is too small for ellipsis",
			input: "abcdef",
			max:   3,
			want:  "abc",
		},
		{
			name:  "returns one character when max is one",
			input: "abcdef",
			max:   1,
			want:  "a",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := truncate(tt.input, tt.max)
			if got != tt.want {
				t.Fatalf("truncate(%q, %d) = %q, want %q", tt.input, tt.max, got, tt.want)
			}
		})
	}
}

func TestHumanizeTime(t *testing.T) {
	tests := []struct {
		name  string
		input time.Time
		want  string
	}{
		{
			name:  "minutes ago",
			input: time.Now().Add(-30 * time.Minute),
			want:  "30m ago",
		},
		{
			name:  "hours ago",
			input: time.Now().Add(-2 * time.Hour),
			want:  "2h ago",
		},
		{
			name:  "days ago",
			input: time.Now().Add(-72 * time.Hour),
			want:  "3d ago",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := humanizeTime(tt.input)
			if got != tt.want {
				t.Fatalf("humanizeTime(...) = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestHumanizeSize(t *testing.T) {
	tests := []struct {
		name string
		size int64
		want string
	}{
		{
			name: "bytes",
			size: 512,
			want: "512 B",
		},
		{
			name: "one kilobyte",
			size: 1024,
			want: "1.0 KB",
		},
		{
			name: "fractional kilobyte",
			size: 1536,
			want: "1.5 KB",
		},
		{
			name: "one megabyte",
			size: 1024 * 1024,
			want: "1.0 MB",
		},
		{
			name: "zero bytes",
			size: 0,
			want: "0 B",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := humanizeSize(tt.size)
			if got != tt.want {
				t.Fatalf("humanizeSize(%d) = %q, want %q", tt.size, got, tt.want)
			}
		})
	}
}
