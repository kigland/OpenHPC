package main

import (
	"log"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/kigland/HPC-Scheduler/lib/dockerHelper"
)

func main() {
	cli, err := client.NewClientWithOpts(client.WithHost("unix:///var/run/docker.sock"), client.WithAPIVersionNegotiation())
	if err != nil {
		log.Fatalf("Failed to create docker client: %v", err)
		return
	}
	dk := dockerHelper.NewDockerHelper(cli)
	id, err := dk.StartContainer(dockerHelper.StartContainerOptions{
		ImageName: "ubuntu",
		Resources: container.Resources{
			DeviceRequests: dockerHelper.GetGPUDeviceRequests(1),
		},
		AttachStdout: true,
		AttachStderr: true,
		Cmd:          []string{"nvidia-smi"},
	})
	if err != nil {
		log.Fatalf("Failed to start container: %v", err)
		return
	}
	log.Printf("Container started: %s", id)
}
