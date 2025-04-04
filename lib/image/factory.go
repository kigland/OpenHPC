package image

import (
	"path/filepath"
	"strconv"

	"github.com/docker/go-connections/nat"
	"github.com/kigland/OpenHPC/lib/dockerHelper"
)

type Factory struct {
	Username    string
	Password    string
	BindHost    string
	BindPort    int
	BindSSHPort int
}

type AllowedImages string

const (
	ImageJupyterHub AllowedImages = "kevinzonda/notebook"
	ImageTorchBook  AllowedImages = "kevinzonda/torchbook"
	ImageMLBook     AllowedImages = "kevinzonda/mlbook"
	ImageBase       AllowedImages = "jupyterhub/singleuser"
)

var ALLOWED_IMAGES = []AllowedImages{
	ImageJupyterHub,
	ImageBase,
}

func (a AllowedImages) HomeDir() string {
	switch a {
	case ImageMLBook, ImageJupyterHub, ImageTorchBook, ImageBase:
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
	case ImageJupyterHub, ImageTorchBook, ImageMLBook, ImageBase:
		return f.jupyterbook(img)
	default:
		return dockerHelper.StartContainerOptions{}
	}
}

const (
	JUPYTER_TOKEN = "JUPYTER_TOKEN"
)

func (f Factory) jupyterbook(id AllowedImages) dockerHelper.StartContainerOptions {
	port := nat.PortMap{}
	port["8888/tcp"] = []nat.PortBinding{{
		HostIP:   f.BindHost,
		HostPort: strconv.Itoa(f.BindPort),
	}}
	if f.BindSSHPort > 0 {
		port["22/tcp"] = []nat.PortBinding{{
			HostIP:   f.BindHost,
			HostPort: strconv.Itoa(f.BindSSHPort),
		}}
	}
	return dockerHelper.StartContainerOptions{
		ImageName: string(id),
		Env: []string{
			JUPYTER_TOKEN + "=" + f.Password,
		},
		PortBindings: port,
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
