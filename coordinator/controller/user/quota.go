package user

import (
	"github.com/gin-gonic/gin"
	"github.com/kigland/OpenHPC/coordinator/models/dboper"
)

func quota(c *gin.Context) {

	uid := c.GetString("uid")
	user, err := dboper.GetUserByID(uid)
	if err != nil {
		c.JSON(400, gin.H{
			"message": "user not found",
		})
		return
	}

	c.JSON(200, gin.H{
		"vcpu":   user.MaxVCPU,
		"vgpu":   user.MaxVGPU,
		"memory": user.MaxMemory,
	})
}
