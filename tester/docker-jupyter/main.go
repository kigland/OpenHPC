package main

import (
	"fmt"
	"log"
	"time"

	"github.com/docker/docker/client"
	"github.com/kigland/HPC-Scheduler/coodinator/utils"
	"github.com/kigland/HPC-Scheduler/lib/consts"
	"github.com/kigland/HPC-Scheduler/lib/dockerHelper"
	"github.com/kigland/HPC-Scheduler/lib/image"
	"github.com/kigland/HPC-Scheduler/lib/svcTag"
)

func main() {
	cli, err := client.NewClientWithOpts(client.WithHost(consts.DOCKER_UNIX_SOCKET), client.WithAPIVersionNegotiation())
	if err != nil {
		log.Fatalf("Failed to create docker client: %v", err)
		return
	}

	passwd := utils.RndId(32) // 256bit = 32bytes

	img := image.Factory{
		Password: passwd,
		BindHost: consts.CONTAINER_HOST,
		BindPort: 41000,
	}.Image(image.ImageTorchBook).WithGPU(1)
	img.AutoRemove = true

	dk := dockerHelper.NewDockerHelper(cli)
	svgT := svcTag.New("KevinZonda")
	img.ContainerName = svgT.String()
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
	fmt.Println(passwd)
}
