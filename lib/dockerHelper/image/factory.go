package image

import (
	"github.com/docker/go-connections/nat"
	"github.com/kigland/HPC-Scheduler/lib/dockerHelper"
)

type Factory struct {
	Username string
	Password string
	BindHost string
	BindPort string
}

type AllowedImages string

const (
	ImageJupyterHub AllowedImages = "kevinzonda/notebook"
)

func (f Factory) Image(img AllowedImages) dockerHelper.StartContainerOptions {
	switch img {
	case ImageJupyterHub:
		return f.JupyterHub()
	default:
		return dockerHelper.StartContainerOptions{}
	}
}

func (f Factory) JupyterHub() dockerHelper.StartContainerOptions {
	return dockerHelper.StartContainerOptions{
		ImageName: string(ImageJupyterHub),
		Env: []string{
			"JUPYTER_TOKEN=" + f.Password,
		},
		PortBindings: nat.PortMap{
			"8888/tcp": []nat.PortBinding{
				{
					HostIP:   f.BindHost,
					HostPort: f.BindPort,
				},
			},
		},
	}
}
