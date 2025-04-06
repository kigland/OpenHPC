package common

import (
	"fmt"
	"strings"

	"github.com/kigland/OpenHPC/lib/svcTag"
	"github.com/kigland/OpenHPC/lib/utils"
)

type VNodeId struct {
	ID     string
	SvcTag svcTag.SvcTag
}

func IDs(cid string) VNodeId {
	summary, ok := DockerHelper.TryGetContainer(cid)
	if ok {
		cid = summary.ID
		svcTag, err := svcTag.Parse(summary.Names[0])
		if err != nil {
			fmt.Println("Failed to parse service tag: ", err)
		}
		return VNodeId{
			ID:     cid,
			SvcTag: svcTag,
		}
	}
	fmt.Println("Container not found or not managed by KHS. Only limited information will be available!")
	svcTag, err := svcTag.Parse(cid)
	if err != nil {
		fmt.Println("Failed to parse service tag: ", err)
		return VNodeId{
			ID:     cid,
			SvcTag: svcTag,
		}
	}
	return VNodeId{
		ID:     cid,
		SvcTag: svcTag,
	}
}

func (v VNodeId) SCID() string {
	return utils.ShortId(v.ID)
}

func (v VNodeId) String() string {
	sb := strings.Builder{}
	sb.WriteString("CID        : ")
	sb.WriteString(v.ID)
	sb.WriteString("\n")
	sb.WriteString("SCID       : ")
	sb.WriteString(utils.ShortId(v.ID))
	sb.WriteString("\n")
	sb.WriteString("SvcTag     : ")
	sb.WriteString(v.SvcTag.String())
	sb.WriteString("\n")
	sb.WriteString("Short Code : ")
	sb.WriteString(v.SvcTag.ShortCode())
	return sb.String()
}
