package main

import (
	"fmt"
	"log"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	kon "github.com/kigland/HPC-Scheduler/coodinator/container"
	"github.com/kigland/HPC-Scheduler/lib/consts"
	"github.com/kigland/HPC-Scheduler/lib/dockerHelper"
)

func main() {
	cli, err := client.NewClientWithOpts(client.WithHost(consts.DOCKER_UNIX_SOCKET), client.WithAPIVersionNegotiation())
	if err != nil {
		log.Fatalf("Failed to create docker client: %v", err)
		return
	}
	img := consts.JupyterHubStartOps("test", "127.0.0.1", "8000")
	img.Resources = container.Resources{
		DeviceRequests: dockerHelper.GetGPUDeviceRequests(1),
	}
	dk := dockerHelper.NewDockerHelper(cli)
	img.ContainerName = kon.NewContainerName("KevinZonda")
	id, err := dk.StartContainer(img)
	if err != nil {
		log.Fatalf("Failed to start container: %v", err)
		return
	}
	log.Printf("Container started: %s", id)

	time.Sleep(4 * time.Second)

	logs, err := dk.GetLogs(id, true)
	if err != nil {
		log.Fatalf("Failed to get logs: %v", err)
		return
	}
	fmt.Println(logs)
}
