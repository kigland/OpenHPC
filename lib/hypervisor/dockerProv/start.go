package dockerProv

import (
	"context"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/strslice"
	"github.com/docker/go-connections/nat"
)

type Provider string

const (
	ProviderDocker Provider = ""
	ProviderPodman Provider = "podman"
)

type StartContainerOptions struct {
	Provider Provider

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

	AlwaysRestart bool

	ShmSize int64

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
	restartPolicy := container.RestartPolicy{}
	if sco.AlwaysRestart {
		restartPolicy = container.RestartPolicy{
			Name: "always",
		}
	}
	return &container.HostConfig{
		Resources:     sco.Resources,
		PortBindings:  sco.PortBindings,
		Binds:         sco.Binds,
		AutoRemove:    sco.AutoRemove,
		ShmSize:       sco.ShmSize,
		RestartPolicy: restartPolicy,
	}
}

func (d *DockerHelper) StartContainer(opts StartContainerOptions, pull bool) (containerID string, err error) {
	cli := d.cli
	if pull {
		err := d.Pull(opts.ImageName)
		if err != nil {
			return "", err
		}
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

func (ops StartContainerOptions) WithAutoRestart() StartContainerOptions {
	ops.AlwaysRestart = true
	ops.AutoRemove = false
	return ops
}

func (ops StartContainerOptions) WithAutoRemove() StartContainerOptions {
	ops.AlwaysRestart = false
	ops.AutoRemove = true
	return ops
}
