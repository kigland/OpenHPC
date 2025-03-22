package handler

import (
	"fmt"
	"log"
	"time"

	"github.com/kigland/HPC-Scheduler/lib/consts"
	"github.com/kigland/HPC-Scheduler/lib/image"
	"github.com/kigland/HPC-Scheduler/lib/utils"
	"github.com/kigland/HPC-Scheduler/tools/common"
)

func Request() {
	port := common.InputPort(consts.LOW_PORT, consts.HIGH_PORT)
	username := common.InputUsername()
	project := common.InputProject()
	dk := common.DockerHelper

	passwd := utils.RndId(32) // 256bit = 32bytes

	imageName := image.ImageJupyterHub

	cinfo, err := common.CreateContainer(dk, imageName, username, passwd, port, project)
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
