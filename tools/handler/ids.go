package handler

import (
	"fmt"
	"log"
	"os"

	"github.com/kigland/HPC-Scheduler/lib/dockerHelper"
	"github.com/kigland/HPC-Scheduler/lib/svcTag"
	"github.com/kigland/HPC-Scheduler/tools/common"
)

func IDs() {
	if len(os.Args) == 3 {
		ids(os.Args[2])
		return
	}
	ids(common.InputWithPrompt("Container ID or Service Tag or Short Code:"))
}

func ids(cid string) {
	summary, ok := common.DockerHelper.TryGetContainer(cid)
	if !ok {
		fmt.Println("Container not found or not managed by KHS. Only limited information will be available!")
		svcTag, err := svcTag.Parse(cid)
		if err != nil {
			fmt.Println("Failed to parse service tag: ", err)
			printIDs(cid, svcTag)
			return
		}
		printIDs("", svcTag)
		return
	}
	cid = summary.ID
	svcTag, err := svcTag.Parse(cid)
	if err != nil {
		log.Fatalf("Failed to parse service tag: %v", err)
	}
	printIDs(cid, svcTag)
}

func printIDs(cid string, svcTag svcTag.SvcTag) {
	fmt.Println("CID        : ", cid)
	fmt.Println("SCID       : ", dockerHelper.ShortId(cid))
	fmt.Println("SvcTag     : ", svcTag.String())
	fmt.Println("Short Code : ", svcTag.ShortCode())
}
