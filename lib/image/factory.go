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
	return f.image(img)
}

func (f Factory) image(id AllowedImages) dockerProv.StartContainerOptions {
	img, ok := allowedImages[id]
	if !ok {
		return dockerProv.StartContainerOptions{}
	}

	port := nat.PortMap{}
	if img.SupportHTTP() {
		_p := nat.Port(strconv.Itoa(img.HTTP) + "/tcp")
		port[_p] = []nat.PortBinding{{
			HostIP:   f.BindHost,
			HostPort: strconv.Itoa(f.BindPort),
		}}
	}

	if img.SupportSSH() {
		_p := nat.Port(strconv.Itoa(img.SSH) + "/tcp")
		port[_p] = []nat.PortBinding{{
			HostIP:   f.BindSSHHost,
			HostPort: strconv.Itoa(f.BindSSHPort),
		}}
	}

	ops := dockerProv.StartContainerOptions{
		ImageName:    string(id),
		PortBindings: port,
		Provider:     f.Provider,
	}

	if img.Env.Token != "" {
		ops.Env = append(ops.Env, img.Env.Token+"="+f.Password)
	}

	return ops
}
