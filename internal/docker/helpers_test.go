package docker

import "testing"

func TestShortID(t *testing.T) {
	tests := []struct {
		name string
		id   string
		n    int
		want string
	}{
		{
			name: "truncates long id",
			id:   "1234567890abcdef",
			n:    12,
			want: "1234567890ab",
		},
		{
			name: "returns unchanged id when shorter than limit",
			id:   "abc",
			n:    12,
			want: "abc",
		},
		{
			name: "returns unchanged id when equal to limit",
			id:   "1234567890ab",
			n:    12,
			want: "1234567890ab",
		},
		{
			name: "returns empty string for empty id",
			id:   "",
			n:    12,
			want: "",
		},
		{
			name: "returns empty string when max length is zero",
			id:   "abcdef",
			n:    0,
			want: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := shortID(tt.id, tt.n)
			if got != tt.want {
				t.Fatalf("shortID(%q, %d) = %q, want %q", tt.id, tt.n, got, tt.want)
			}
		})
	}
}

func TestFirstContainerName(t *testing.T) {
	tests := []struct {
		name  string
		names []string
		want  string
	}{
		{
			name:  "removes docker leading slash",
			names: []string{"/redis"},
			want:  "redis",
		},
		{
			name:  "keeps name without slash",
			names: []string{"postgres"},
			want:  "postgres",
		},
		{
			name:  "uses first name when multiple names exist",
			names: []string{"/api", "/api-alias"},
			want:  "api",
		},
		{
			name:  "returns empty string for empty slice",
			names: []string{},
			want:  "",
		},
		{
			name:  "returns empty string for nil slice",
			names: nil,
			want:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := firstContainerName(tt.names)
			if got != tt.want {
				t.Fatalf("firstContainerName(%v) = %q, want %q", tt.names, got, tt.want)
			}
		})
	}
}

func TestSplitRepoTag(t *testing.T) {
	tests := []struct {
		name     string
		ref      string
		wantRepo string
		wantTag  string
	}{
		{
			name:     "standard repository with tag",
			ref:      "redis:latest",
			wantRepo: "redis",
			wantTag:  "latest",
		},
		{
			name:     "registry port with tag",
			ref:      "localhost:5000/app:latest",
			wantRepo: "localhost:5000/app",
			wantTag:  "latest",
		},
		{
			name:     "nested repository with explicit tag",
			ref:      "ghcr.io/acme/platform/api:v1.2.3",
			wantRepo: "ghcr.io/acme/platform/api",
			wantTag:  "v1.2.3",
		},
		{
			name:     "repository without tag defaults to latest",
			ref:      "nginx",
			wantRepo: "nginx",
			wantTag:  "latest",
		},
		{
			name:     "registry repository without tag defaults to latest",
			ref:      "localhost:5000/app",
			wantRepo: "localhost:5000/app",
			wantTag:  "latest",
		},
		{
			name:     "docker dangling image",
			ref:      "<none>:<none>",
			wantRepo: "<none>",
			wantTag:  "<none>",
		},
		{
			name:     "empty reference",
			ref:      "",
			wantRepo: "<none>",
			wantTag:  "<none>",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotRepo, gotTag := splitRepoTag(tt.ref)

			if gotRepo != tt.wantRepo {
				t.Fatalf("splitRepoTag(%q) repo = %q, want %q", tt.ref, gotRepo, tt.wantRepo)
			}

			if gotTag != tt.wantTag {
				t.Fatalf("splitRepoTag(%q) tag = %q, want %q", tt.ref, gotTag, tt.wantTag)
			}
		})
	}
}
