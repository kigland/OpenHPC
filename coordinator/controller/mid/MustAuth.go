package mid

import (
	"github.com/gin-gonic/gin"
	"github.com/kigland/OpenHPC/coordinator/models/dboper"
	"github.com/kigland/OpenHPC/coordinator/shared"
	"github.com/kigland/OpenHPC/coordinator/utils"
)

const MID_USER_ID = "mid_user_id"
const MID_ACL_ALLOW_LIST = "mid_acl_allowlist"

func FakeAuth(c *gin.Context) {
	c.Set(MID_USER_ID, "1")
	c.Set(MID_ACL_ALLOW_LIST, []string{""})
	c.Next()
}

func ACLAuth(c *gin.Context) {
	token := c.GetHeader("Authorization")
	var allowList []string
	if token != "" {
		allowList = shared.GetConfig().ACLParsed.GetACLAllowList(token)
	}

	if len(allowList) == 0 {
		utils.Unauthorised(c)
		return
	}
	c.Set(MID_ACL_ALLOW_LIST, allowList)
}

func MustAuth(c *gin.Context) {
	token := c.GetHeader("Authorization")
	if token == "" {
		utils.Unauthorised(c)
		return
	}

	tk, err := dboper.GetTokenByToken(token)
	if err != nil {
		utils.Unauthorised(c)
		return
	}

	c.Set(MID_USER_ID, tk.UserId)
	c.Set(MID_ACL_ALLOW_LIST, []string{""})
	c.Next()
}
