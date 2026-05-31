package vm

import (
	"log"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/kigland/OpenHPC/coordinator/controller/mid"
	"github.com/kigland/OpenHPC/lib/hypervisor/dockerProv"
	"github.com/kigland/OpenHPC/lib/svcTag"
)

func _verifyACL(aclAllowedList []string, svcTag svcTag.SvcTag) bool {
	if len(aclAllowedList) == 0 {
		return false
	}
	tgtAcl := svcTag.ACLCode()
	for _, s := range aclAllowedList {
		if strings.HasPrefix(tgtAcl, s) {
			return true
		}
	}
	return false
}

func verifyACL(c *gin.Context, svcTag svcTag.SvcTag) bool {
	aclAllowedList := c.GetStringSlice(mid.MID_ACL_ALLOW_LIST)
	return _verifyACL(aclAllowedList, svcTag)
}

func verifyACLFromSvcTagStr(c *gin.Context, svcTagStr string) bool {
	t, err := svcTag.Parse(svcTagStr)
	if err != nil {
		return false
	}
	return verifyACL(c, t)
}

func verifyACLFromContainerSummary(c *gin.Context, summary dockerProv.ContainerSummaryWithSvcTag) bool {
	aclAllowedList := c.GetStringSlice(mid.MID_ACL_ALLOW_LIST)
	if len(aclAllowedList) == 0 {
		return false
	}

	if summary.SvcTag != nil {
		log.Println("Container has svcTag, verifying with svcTag:", summary.SvcTag.String())
		return _verifyACL(aclAllowedList, *summary.SvcTag)
	}
	for _, name := range summary.Names {
		log.Println("Verifying container name with ACL, name:", name)
		if strings.HasPrefix(name, "/") {
			name = name[1:]
		}
		svcTagParse, err := svcTag.Parse(name)
		if err != nil {
			continue
		}
		if _verifyACL(aclAllowedList, svcTagParse) {
			return true
		}
	}
	return false
}
