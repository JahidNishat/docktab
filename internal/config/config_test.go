package config

import "testing"

func TestNew_DefaultValues(t *testing.T) {
	cfg := New()

	if cfg.DockerHost != "" {
		t.Fatalf("expected DockerHost to be empty, got %q", cfg.DockerHost)
	}

	if cfg.Debug {
		t.Fatalf("expected Debug to be false, got %t", cfg.Debug)
	}
}

func TestNew_ReadsDockerHostFromEnv(t *testing.T) {
	t.Setenv("DOCKTAB_DOCKER_HOST", "unix:///var/run/docker.sock")

	cfg := New()
	want := "unix:///var/run/docker.sock"

	if cfg.DockerHost != want {
		t.Fatalf("expected DockerHost to be %q, got %q", want, cfg.DockerHost)
	}
}

func TestNew_ReadsDebugFromEnv(t *testing.T) {
	t.Setenv("DOCKTAB_DEBUG", "true")

	cfg := New()

	if !cfg.Debug {
		t.Fatalf("expected Debug to be true, got %t", cfg.Debug)
	}
}

func TestNew_ReadsMultipleEnvValues(t *testing.T) {
	t.Setenv("DOCKTAB_DOCKER_HOST", "tcp://localhost:2375")
	t.Setenv("DOCKTAB_DEBUG", "true")

	cfg := New()
	wantHost := "tcp://localhost:2375"

	if cfg.DockerHost != wantHost {
		t.Fatalf("expected DockerHost to be %q, got %q", wantHost, cfg.DockerHost)
	}

	if !cfg.Debug {
		t.Fatalf("expected Debug to be true, got %t", cfg.Debug)
	}
}

func TestNew_DebugFalseFromEnv(t *testing.T) {
	t.Setenv("DOCKTAB_DEBUG", "false")

	cfg := New()

	if cfg.Debug {
		t.Fatalf("expected Debug to be false, got %t", cfg.Debug)
	}
}
