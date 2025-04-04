package handler

import (
	"fmt"
	"log"
	"time"

	"github.com/kigland/OpenHPC/lib/consts"
	"github.com/kigland/OpenHPC/lib/svcTag"
	"github.com/kigland/OpenHPC/tools/common"
)

func Request() {
	port := common.InputPort(consts.LOW_PORT, consts.HIGH_PORT)
	owner := common.InputOwner()
	project := common.InputProject()

	enableRDS := common.InputEnableRDS()

	subfolder := ""
	if enableRDS {
		subfolder = common.InputWithPrompt("RDS Submodule (default \"\"):")
	}

	dk := common.DockerHelper

	token := common.InputTokenOrGenerate(32)

	image := common.InputImage()

	needSSH := common.InputNeedSSH()

	rdsDir, rdsMountAt := "", ""
	if enableRDS {
		rdsDir, rdsMountAt = common.GetRDSWithSubfolder(owner, subfolder, image)
	}
	tag := svcTag.New(owner).WithProject(project)
	cinfo, err := common.CreateContainerCustomRDS(dk, image, tag, token, port, rdsDir, rdsMountAt, needSSH)

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
