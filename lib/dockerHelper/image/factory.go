package image

import (
	"path/filepath"

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
	ImageTorchBook  AllowedImages = "kevinzonda/torchbook"
	ImageMLBook     AllowedImages = "kevinzonda/mlbook"
)

func (a AllowedImages) HomeDir() string {
	switch a {
	case ImageMLBook, ImageJupyterHub, ImageTorchBook:
		return "/home/jovyan"
	}
	return ""
}

func (a AllowedImages) RdsDir() string {
	home := a.HomeDir()
	if home != "" {
		return filepath.Join(home, "rds")
	}
	return "/rds"
}

func (f Factory) Image(img AllowedImages) dockerHelper.StartContainerOptions {
	switch img {
	case ImageJupyterHub, ImageTorchBook, ImageMLBook:
		return f.jupyterbook(img)
	default:
		return dockerHelper.StartContainerOptions{}
	}
}

func (f Factory) jupyterbook(id AllowedImages) dockerHelper.StartContainerOptions {
	return dockerHelper.StartContainerOptions{
		ImageName: string(id),
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

func (f Factory) JupyterHub() dockerHelper.StartContainerOptions {
	return f.jupyterbook(ImageJupyterHub)
}

func (f Factory) TorchBook() dockerHelper.StartContainerOptions {
	return f.jupyterbook(ImageTorchBook)
}

func (f Factory) MLBook() dockerHelper.StartContainerOptions {
	return f.jupyterbook(ImageMLBook)
}
