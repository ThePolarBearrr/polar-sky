package task

import (
	"context"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
	"io"
	"os"
	"polar-sky/log"
)

var Client client.Client

type Docker struct {
	Client *client.Client
	Config *Config
}

type DockerResult struct {
	Action      string
	Error       error
	ContainerID string
	Result      string
}

func NewDocker(cfg *Config) Docker {
	return Docker{
		Client: &Client,
		Config: cfg,
	}
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
		RestartPolicy:   rp,
		Resources:       r,
		PublishAllPorts: true,
	}
	resp, err := d.Client.ContainerCreate(ctx, &cc, &hc, nil, nil, d.Config.Name)
	if err != nil {
		log.Logger.Errorf("Failed create container %s with image %s, error: %v\n", d.Config.Name, d.Config.Image, err)
		return DockerResult{
			Error: err,
		}
	}
	log.Logger.Infof("Container %s has been created\n", resp.ID)

	// start container
	err = d.Client.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{})
	if err != nil {
		log.Logger.Errorf("Failed start container %s, error: %v\n", resp.ID, err)
		return DockerResult{
			Error:       err,
			ContainerID: resp.ID,
		}
	}
	log.Logger.Infof("Container %s has been started\n", resp.ID)

	// stdout container log
	out, err := d.Client.ContainerLogs(ctx, resp.ID, types.ContainerLogsOptions{ShowStdout: true, ShowStderr: true})
	if err != nil {
		log.Logger.Errorf("Failed get container %s log, error: %v\n", resp.ID, err)
		return DockerResult{
			Error:       err,
			ContainerID: resp.ID,
		}
	}
	stdcopy.StdCopy(os.Stdout, os.Stderr, out)

	return DockerResult{
		ContainerID: resp.ID,
		Action:      "Run",
		Result:      "success",
		Error:       nil,
	}
}

func (d *Docker) Stop(ContainerID string) DockerResult {
	log.Logger.Infof("Attempting to stop container %s", ContainerID)

	// stop container
	ctx := context.Background()
	err := d.Client.ContainerStop(ctx, ContainerID, container.StopOptions{})
	if err != nil {
		panic(err)
	}

	// remove container
	removeOptions := types.ContainerRemoveOptions{
		RemoveVolumes: true,
		RemoveLinks:   false,
		Force:         false,
	}
	err = d.Client.ContainerRemove(ctx, ContainerID, removeOptions)
	if err != nil {
		panic(err)
	}

	return DockerResult{
		Action:      "Stop",
		Result:      "success",
		ContainerID: ContainerID,
		Error:       nil,
	}
}
