package mid

import (
	"github.com/gin-gonic/gin"
	"github.com/kigland/HPC-Scheduler/coodinator/models/dboper"
)

const MID_USER_ID = "mid_user_id"

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
