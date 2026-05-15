package docker

import "testing"

func TestShortID(t *testing.T) {
	got := shortID("1234567890abcdef", 12)
	want := "1234567890ab"

	if got != want {
		t.Fatalf("expected %q, got %q", want, got)
	}
}

func TestShortIDShortInput(t *testing.T) {
	got := shortID("abc", 12)
	want := "abc"

	if got != want {
		t.Fatalf("expected %q, got %q", want, got)
	}
}

func TestFirstContainerName(t *testing.T) {
	got := firstContainerName([]string{"/redis"})
	want := "redis"

	if got != want {
		t.Fatalf("expected %q, got %q", want, got)
	}
}

func TestSplitRepoTag(t *testing.T) {
	repo, tag := splitRepoTag("localhost:5000/app:latest")

	if repo != "localhost:5000/app" {
		t.Fatalf("unexpected repo: %q", repo)
	}
	if tag != "latest" {
		t.Fatalf("unexpected tag: %q", tag)
	}
}
