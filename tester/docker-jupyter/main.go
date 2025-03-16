package main

import (
	"fmt"
	"log"
	"time"

	"github.com/docker/docker/client"
	kon "github.com/kigland/HPC-Scheduler/coodinator/container"
	"github.com/kigland/HPC-Scheduler/lib/consts"
	"github.com/kigland/HPC-Scheduler/lib/dockerHelper"
	"github.com/kigland/HPC-Scheduler/lib/dockerHelper/image"
)

func main() {
	cli, err := client.NewClientWithOpts(client.WithHost(consts.DOCKER_UNIX_SOCKET), client.WithAPIVersionNegotiation())
	if err != nil {
		log.Fatalf("Failed to create docker client: %v", err)
		return
	}

	img := image.Factory{
		Password: "6e644622daba40ba9d0a5999d6dd5734",
		BindHost: "127.0.0.2",
		BindPort: "41000",
	}.Image(image.ImageJupyterHub).WithGPU(1)

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
