package docker

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

type Container struct {
	ID      string
	Name    string
	Image   string
	Command string
	Created time.Time
	Status  string
	State   string
	Ports   string
}

type Client interface {
	ListContainers(ctx context.Context, all bool) ([]Container, error)
}

type dockerClient struct {
	cli *client.Client
}

func NewClient(host string) Client {
	opts := []client.Opt{client.FromEnv}
	if host != "" {
		opts = append(opts, client.WithHost(host))
	}

	cli, err := client.NewClientWithOpts(opts...)
	if err != nil {
		panic(fmt.Sprintf("Failed to create Docker client: %v", err))
	}
	return &dockerClient{cli: cli}
}

func (d *dockerClient) ListContainers(ctx context.Context, all bool) ([]Container, error) {
	containers, err := d.cli.ContainerList(ctx, container.ListOptions{All: all})
	if err != nil {
		return nil, err
	}

	result := make([]Container, 0, len(containers))
	for _, c := range containers {
		result = append(result, Container{
			ID:      c.ID[:12],                           // Shorten ID for display
			Name:    strings.TrimPrefix(c.Names[0], "/"), // Remove leading slash from name
			Image:   c.Image,
			Command: c.Command,
			Created: time.Unix(c.Created, 0),
			Status:  c.Status,
			State:   c.State,
			Ports:   formatPorts(c.Ports),
		})
	}
	return result, nil
}

func formatPorts(ports []types.Port) string {
	if len(ports) == 0 {
		return ""
	}

	var parts []string
	for _, p := range ports {
		if p.PublicPort != 0 {
			parts = append(parts, fmt.Sprintf("%d->%d", p.PublicPort, p.PrivatePort))
		} else {
			parts = append(parts, fmt.Sprintf("%d", p.PrivatePort))
		}
	}
	return strings.Join(parts, ", ")
}
