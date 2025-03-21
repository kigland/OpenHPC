package handler

import (
	"fmt"
	"log"
	"time"

	kon "github.com/kigland/HPC-Scheduler/coodinator/container"

	"github.com/kigland/HPC-Scheduler/lib/image"
	"github.com/kigland/HPC-Scheduler/tools/common"
)

func Request() {
	port := common.InputPort(40000, 41000)
	username := common.InputUsername()
	dk := common.DockerHelper

	passwd := kon.RndId(32) // 256bit = 32bytes

	imageName := image.ImageJupyterHub

	cinfo, err := common.CreateContainer(dk, imageName, username, passwd, port)
	if err != nil {
		log.Fatalf("Failed to start container: %v", err)
		return
	}
	log.Printf("Container started: %s", cinfo.CID)

	time.Sleep(4 * time.Second)

	logs, err := dk.GetLogs(cinfo.CID, true)
	if err != nil {
		log.Fatalf("Failed to get logs: %v", err)
		return
	}
	fmt.Println(logs)
	fmt.Println("--------------------------------")
	fmt.Println(cinfo.String())
	fmt.Println("--------------------------------")
}
