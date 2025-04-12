package mid

import (
	"slices"

	"github.com/gin-gonic/gin"
	"github.com/kigland/OpenHPC/coordinator/models/dboper"
	"github.com/kigland/OpenHPC/coordinator/shared"
)

const MID_USER_ID = "mid_user_id"

func FakeAuth(c *gin.Context) {
	c.Set(MID_USER_ID, "1")
	c.Next()
}

func ACLAuth(c *gin.Context) {
	token := c.GetHeader("Authorization")
	if token == "" || !slices.Contains(shared.GetConfig().ACL.APIKeys, token) {
		c.JSON(401, gin.H{
			"message": "Unauthorised",
		})
		c.Abort()
		return
	}
}

func MustAuth(c *gin.Context) {
	token := c.GetHeader("Authorization")
	if token == "" {
		c.JSON(401, gin.H{
			"message": "Unauthorized",
		})
		c.Abort()
		return
	}

	tk, err := dboper.GetTokenByToken(token)
	if err != nil {
		c.JSON(401, gin.H{
			"message": "Unauthorized",
		})
		c.Abort()
		return
	}

	c.Set(MID_USER_ID, tk.UserId)
	c.Next()
}
