package dockerHelper

import (
	"context"
	"io"
	"os"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/go-connections/nat"
)

type StartContainerOptions struct {
	ImageName       string
	Env             []string
	User            string
	Volumes         map[string]struct{}
	NetworkDisabled bool `json:",omitempty"`
	ExposedPorts    nat.PortSet
	ContainerName   string

	Binds []string

	Resources container.Resources
}

func (sco StartContainerOptions) toContainerConfig() *container.Config {
	return &container.Config{
		Image:           sco.ImageName,
		Env:             sco.Env,
		User:            sco.User,
		Volumes:         sco.Volumes,
		ExposedPorts:    sco.ExposedPorts,
		NetworkDisabled: sco.NetworkDisabled,
	}
}

func (sco StartContainerOptions) toHostConfig() *container.HostConfig {
	return &container.HostConfig{
		Resources: sco.Resources,
		Binds:     sco.Binds,
	}
}

func (d *DockerHelper) StartContainer(opts StartContainerOptions) (containerID string, err error) {
	cli := d.cli
	out, err := cli.ImagePull(context.Background(), opts.ImageName, image.PullOptions{})
	if err != nil {
		return "", err
	}
	defer out.Close()
	io.Copy(os.Stdout, out)

	resp, err := cli.ContainerCreate(context.Background(), opts.toContainerConfig(), opts.toHostConfig(), nil, nil, opts.ContainerName)
	if err != nil {
		return "", err
	}

	if err := cli.ContainerStart(context.Background(), resp.ID, container.StartOptions{}); err != nil {
		return "", err
	}

	return resp.ID, nil
}
