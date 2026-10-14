package server

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

type DockerController struct{ client *client.Client }

func NewDockerController() (*DockerController, error) {
	c, e := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if e != nil {
		return nil, e
	}
	return &DockerController{c}, nil
}
func (d *DockerController) List(ctx context.Context) ([]types.Container, error) {
	return d.client.ContainerList(ctx, types.ContainerListOptions{All: true})
}
func (d *DockerController) Control(ctx context.Context, id, action string) error {
	switch action {
	case "start":
		return d.client.ContainerStart(ctx, id, types.ContainerStartOptions{})
	case "stop":
		return d.client.ContainerStop(ctx, id, container.StopOptions{})
	case "restart":
		return d.client.ContainerRestart(ctx, id, container.StopOptions{})
	default:
		return fmt.Errorf("unsupported action %q", action)
	}
}
