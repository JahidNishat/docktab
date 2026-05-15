package docker

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/api/types/volume"
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

type Image struct {
	ID         string
	Repository string
	Tag        string
	Created    time.Time
	Size       int64
}

type Volume struct {
	Name       string
	Driver     string
	Mountpoint string
	Created    time.Time
}

type Network struct {
	Name    string
	ID      string
	Driver  string
	Scope   string
	Created time.Time
}

type Client interface {
	ListContainers(ctx context.Context, all bool) ([]Container, error)
	ListImages(ctx context.Context, all bool) ([]Image, error)
	ListVolumes(ctx context.Context) ([]Volume, error)
	ListNetworks(ctx context.Context) ([]Network, error)
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
		panic(fmt.Sprintf("failed to create docker client: %v", err))
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

func (d *dockerClient) ListImages(ctx context.Context, all bool) ([]Image, error) {
	images, err := d.cli.ImageList(ctx, image.ListOptions{All: all})
	if err != nil {
		return nil, err
	}

	result := make([]Image, 0, len(images))
	for _, img := range images {
		repoTags := img.RepoTags
		if len(repoTags) == 0 {
			repoTags = []string{"<none>:<none>"}
		}

		for _, repoTag := range repoTags {
			parts := strings.Split(repoTag, ":")
			repository := parts[0]
			tag := "latest"
			if len(parts) > 1 {
				tag = parts[1]
			}

			result = append(result, Image{
				ID:         img.ID[7:19], // short ID
				Repository: repository,
				Tag:        tag,
				Created:    time.Unix(img.Created, 0),
				Size:       img.Size,
			})
		}
	}

	// Sort by size descending by default
	sort.Slice(result, func(i, j int) bool {
		return result[i].Size > result[j].Size
	})

	return result, nil
}

func (d *dockerClient) ListVolumes(ctx context.Context) ([]Volume, error) {
	volumes, err := d.cli.VolumeList(ctx, volume.ListOptions{})
	if err != nil {
		return nil, err
	}

	result := make([]Volume, 0, len(volumes.Volumes))
	for _, v := range volumes.Volumes {
		created, _ := time.Parse(time.RFC3339, v.CreatedAt)

		result = append(result, Volume{
			Name:       v.Name,
			Driver:     v.Driver,
			Mountpoint: v.Mountpoint,
			Created:    created,
		})
	}
	return result, nil
}

func (d *dockerClient) ListNetworks(ctx context.Context) ([]Network, error) {
	networks, err := d.cli.NetworkList(ctx, types.NetworkListOptions{})
	if err != nil {
		return nil, err
	}

	result := make([]Network, 0, len(networks))
	for _, n := range networks {
		result = append(result, Network{
			ID:      n.ID[:12],
			Name:    n.Name,
			Driver:  n.Driver,
			Scope:   n.Scope,
			Created: n.Created, // may need parsing if it's string
		})
	}
	return result, nil
}
