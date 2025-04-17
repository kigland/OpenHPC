package image

import (
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

func (f Factory) Image(img AllowedImages) dockerProv.StartContainerOptions {
	switch img {
	case ImageJupyterHub, ImageBase:
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
	if f.BindSSHPort > 0 && id.SupportSSH() {
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
