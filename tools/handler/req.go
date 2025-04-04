package handler

import (
	"fmt"
	"log"
	"time"

	"github.com/kigland/OpenHPC/lib/consts"
	"github.com/kigland/OpenHPC/tools/common"
)

func Request() {
	port := common.InputPort(consts.LOW_PORT, consts.HIGH_PORT)
	username := common.InputUsername()
	project := common.InputProject()
	dk := common.DockerHelper

	token := common.InputTokenOrGenerate(32)

	image := common.InputImage()

	cinfo, err := common.CreateContainer(dk, image, username, token, port, project)

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
