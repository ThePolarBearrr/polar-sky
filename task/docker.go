package task

import (
	"context"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"io"
	"os"
	"polar-sky/log"
)

type Docker struct {
	Client      *client.Client
	Config      Config
	ContainerID string
}

type DockerResult struct {
	Action      string
	Error       error
	ContainerID string
	Result      string
}

func (d *Docker) Run() DockerResult {
	// pull image
	ctx := context.Background()
	reader, err := d.Client.ImagePull(ctx, d.Config.Image, types.ImagePullOptions{})
	if err != nil {
		log.Logger.Errorf("Failed pull image %s, error: %v\n", d.Config.Image, err)
		return DockerResult{
			Error: err,
		}
	}
	io.Copy(os.Stdout, reader)

	// create container
	rp := container.RestartPolicy{
		Name: d.Config.RestartPolicy,
	}
	r := container.Resources{
		Memory: d.Config.Memory,
	}
	cc := container.Config{
		Image: d.Config.Image,
		Env:   d.Config.Env,
	}
	hc := container.HostConfig{
		RestartPolicy: rp,
		Resources:     r,
	}
	resp, err := d.Client.ContainerCreate(ctx, &cc, &hc, nil, nil, d.Config.Name)
	if err != nil {
		log.Logger.Errorf("Failed create container %s with image %s, error: %v\n", d.Config.Name, d.Config.Image, err)
		return DockerResult{
			Error: err,
		}
	}

	return DockerResult{}
}
