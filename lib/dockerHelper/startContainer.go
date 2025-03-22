package dockerHelper

import (
	"context"
	"io"
	"os"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/api/types/strslice"
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
	Cmd             []string
	// Tty             bool
	// AttachStdout    bool
	// AttachStderr    bool
	// AttachStdin     bool
	Binds        []string
	PortBindings nat.PortMap
	AutoRemove   bool
	Labels       map[string]string

	Resources container.Resources
}

func (sco StartContainerOptions) ToContainerConfig() *container.Config {
	return &container.Config{
		Image:           sco.ImageName,
		Env:             sco.Env,
		User:            sco.User,
		Volumes:         sco.Volumes,
		ExposedPorts:    sco.ExposedPorts,
		NetworkDisabled: sco.NetworkDisabled,
		Cmd:             strslice.StrSlice(sco.Cmd),

		// Tty:             sco.Tty,
		// AttachStdout:    sco.AttachStdout,
		// AttachStderr:    sco.AttachStderr,
		// AttachStdin:     sco.AttachStdin,
		Labels: sco.Labels,
	}
}

func (sco StartContainerOptions) ToHostConfig() *container.HostConfig {
	return &container.HostConfig{
		Resources:    sco.Resources,
		PortBindings: sco.PortBindings,
		Binds:        sco.Binds,
		AutoRemove:   sco.AutoRemove,
	}
}

func (d *DockerHelper) StartContainer(opts StartContainerOptions, pull bool) (containerID string, err error) {
	cli := d.cli
	if pull {
		out, err := cli.ImagePull(context.Background(), opts.ImageName, image.PullOptions{})
		if err != nil {
			return "", err
		}
		defer out.Close()
		io.Copy(os.Stdout, out)
	}

	resp, err := cli.ContainerCreate(context.Background(), opts.ToContainerConfig(), opts.ToHostConfig(), nil, nil, opts.ContainerName)
	if err != nil {
		return "", err
	}

	if err := cli.ContainerStart(context.Background(), resp.ID, container.StartOptions{}); err != nil {
		return "", err
	}

	return resp.ID, nil
}
