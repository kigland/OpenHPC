package vm

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/kigland/OpenHPC/coordinator/controller/mid"
	"github.com/kigland/OpenHPC/lib/svcTag"
)

func verifyACL(c *gin.Context, svcTag svcTag.SvcTag) bool {
	aclAllowedList := c.GetStringSlice(mid.MID_ACL_ALLOW_LIST)
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

func verifyACLFromSvcTagStr(c *gin.Context, svcTagStr string) bool {
	t, err := svcTag.Parse(svcTagStr)
	if err != nil {
		return false
	}
	return verifyACL(c, t)
}
