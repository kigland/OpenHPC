package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/KevinZonda/GoX/pkg/panicx"
	"github.com/chzyer/readline"
	"github.com/docker/docker/client"
	kon "github.com/kigland/HPC-Scheduler/coodinator/container"
	"github.com/kigland/HPC-Scheduler/lib/consts"
	"github.com/kigland/HPC-Scheduler/lib/dockerHelper"
	"github.com/kigland/HPC-Scheduler/lib/dockerHelper/image"
)

func getRDS(username string, imageName image.AllowedImages) (rdsDir string, rdsMountAt string) {
	rdsDir = filepath.Join("/data/rds", strings.ToLower(username))
	rdsMountAt = imageName.HomeDir()
	if rdsMountAt != "" {
		rdsMountAt = filepath.Join(rdsMountAt, "rds")
	} else {
		rdsMountAt = "/rds"
	}
	if _, err := os.Stat(rdsDir); err != nil {
		fmt.Println("RDS directory not found, skipping...")
		rdsDir = ""
	}
	return rdsDir, rdsMountAt
}

var rl *readline.Instance

func InputPort() int {
	fmt.Println("Port of the container (40000-41000):")
	portStr, err := rl.Readline()
	panicx.NotNilErr(err)

	port, err := strconv.Atoi(portStr)
	panicx.NotNilErr(err)
	if port < 40000 || port > 41000 {
		log.Fatalf("Invalid port: %d", port)
		return 0
	}
	return port
}

func InputUsername() string {
	fmt.Println("Username:")
	username, err := rl.Readline()
	panicx.NotNilErr(err)
	username = strings.TrimSpace(username)
	if username == "" {
		log.Fatalf("Username cannot be empty")
		return ""
	}
	return username
}

func main() {
	var err error
	rl, err = readline.New("> ")
	panicx.NotNilErr(err)
	defer rl.Close()

	port := InputPort()

	username := InputUsername()

	cli, err := client.NewClientWithOpts(client.WithHost(consts.DOCKER_UNIX_SOCKET), client.WithAPIVersionNegotiation())
	panicx.NotNilErr(err)

	passwd := kon.RndId(32) // 256bit = 32bytes

	imageName := image.ImageJupyterHub

	img := image.Factory{
		Password: passwd,
		BindHost: "127.0.0.2",
		BindPort: strconv.Itoa(port),
	}.Image(imageName).WithGPU(1)
	img.AutoRemove = true

	dk := dockerHelper.NewDockerHelper(cli)
	img.ContainerName = kon.NewContainerName(username)

	rdsDir, rdsMountAt := getRDS(username, imageName)
	img = img.WithRDS(rdsDir, rdsMountAt)

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
	fmt.Println("--------------------------------")
	fmt.Println("URL  : http://127.0.0.2:" + strconv.Itoa(port))
	fmt.Println("Token: " + passwd)
	fmt.Println("CID  : " + id)
	fmt.Println("RDS  : " + rdsMountAt)
	fmt.Println("--------------------------------")
}
