package docker

import (
	"context"
	"fmt"
	"strings"

	dockertypes "github.com/docker/docker/api/types"
	dockercontainer "github.com/docker/docker/api/types/container"
	dockerclient "github.com/docker/docker/client"
	"github.com/paknahad/docker-log/internal/domain"
)

type containerAPI interface {
	ContainerList(context.Context, dockercontainer.ListOptions) ([]dockertypes.Container, error)
}

type Client struct {
	api containerAPI
}

func NewClient() (*Client, error) {
	api, err := dockerclient.NewClientWithOpts(
		dockerclient.FromEnv,
		dockerclient.WithAPIVersionNegotiation(),
	)
	if err != nil {
		return nil, fmt.Errorf("create Docker client: %w", err)
	}
	return NewClientWithAPI(api), nil
}

func NewClientWithAPI(api containerAPI) *Client {
	return &Client{api: api}
}

func (c *Client) ListRunningContainers(ctx context.Context) ([]domain.Container, error) {
	containers, err := c.api.ContainerList(ctx, dockercontainer.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("list running Docker containers: %w", err)
	}

	discovered := make([]domain.Container, 0, len(containers))
	for _, container := range containers {
		discovered = append(discovered, domain.Container{
			ID:     container.ID,
			Name:   primaryName(container.Names),
			Image:  container.Image,
			Status: container.Status,
		})
	}
	return discovered, nil
}

func primaryName(names []string) string {
	if len(names) == 0 {
		return ""
	}
	return strings.TrimPrefix(names[0], "/")
}
