package image

import (
	"path/filepath"
	"strconv"

	"github.com/docker/go-connections/nat"
	"github.com/kigland/OpenHPC/lib/hypervisor/dockerProv"
)

type Factory struct {
	Username    string
	Password    string
	BindHost    string
	BindPort    int
	BindSSHHost string
	BindSSHPort int
	Provider    dockerProv.Provider
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

func (a AllowedImages) BaseURLEnvVar() string {
	switch a {
	case ImageMLBook, ImageJupyterHub, ImageTorchBook:
		return "NB_VAR_BASE_URL"
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

func (f Factory) Image(img AllowedImages) dockerProv.StartContainerOptions {
	switch img {
	case ImageJupyterHub, ImageTorchBook, ImageMLBook, ImageBase:
		return f.jupyterbook(img)
	default:
		return dockerProv.StartContainerOptions{}
	}
}

const (
	JUPYTER_TOKEN = "JUPYTER_TOKEN"
)

func (f Factory) jupyterbook(id AllowedImages) dockerProv.StartContainerOptions {
	port := nat.PortMap{}
	port["8888/tcp"] = []nat.PortBinding{{
		HostIP:   f.BindHost,
		HostPort: strconv.Itoa(f.BindPort),
	}}
	if f.BindSSHPort > 0 {
		port["22/tcp"] = []nat.PortBinding{{
			HostIP:   f.BindSSHHost,
			HostPort: strconv.Itoa(f.BindSSHPort),
		}}
	}
	return dockerProv.StartContainerOptions{
		ImageName: string(id),
		Env: []string{
			JUPYTER_TOKEN + "=" + f.Password,
		},
		PortBindings: port,
		Provider:     f.Provider,
	}
}

func (f Factory) JupyterHub() dockerProv.StartContainerOptions {
	return f.jupyterbook(ImageJupyterHub)
}

func (f Factory) TorchBook() dockerProv.StartContainerOptions {
	return f.jupyterbook(ImageTorchBook)
}

func (f Factory) MLBook() dockerProv.StartContainerOptions {
	return f.jupyterbook(ImageMLBook)
}
