package main

import (
	"log"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/kigland/OpenHPC/lib/consts"
	"github.com/kigland/OpenHPC/lib/hypervisor/dockerProv"
)

func main() {
	cli, err := client.NewClientWithOpts(client.WithHost(consts.DOCKER_UNIX_SOCKET), client.WithAPIVersionNegotiation())
	if err != nil {
		log.Fatalf("Failed to create docker client: %v", err)
		return
	}
	dk := dockerProv.NewDockerHelper(cli)
	id, err := dk.StartContainer(dockerProv.StartContainerOptions{
		ImageName: "ubuntu",
		Resources: container.Resources{
			DeviceRequests: dockerProv.GetGPUDeviceRequests(1),
		},
		// AttachStdout: true,
		// AttachStderr: true,
		Cmd: []string{"nvidia-smi"},
	}, true)
	if err != nil {
		log.Fatalf("Failed to start container: %v", err)
		return
	}
	log.Printf("Container started: %s", id)
}
