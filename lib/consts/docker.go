package consts

import (
	"github.com/docker/go-connections/nat"
	"github.com/kigland/HPC-Scheduler/lib/dockerHelper"
)

const (
	DOCKER_UNIX_SOCKET = "unix:///var/run/docker.sock"
)

func JupyterHubStartOps(password string, bindHost, bindPort string) dockerHelper.StartContainerOptions {
	return dockerHelper.StartContainerOptions{
		ImageName: "jupyterhub/singleuser",
		Env: []string{
			"JUPYTER_TOKEN=" + password,
		},
		PortBindings: nat.PortMap{
			"8000/tcp": []nat.PortBinding{
				{
					HostIP:   bindHost,
					HostPort: bindPort,
				},
			},
		},
	}
}
