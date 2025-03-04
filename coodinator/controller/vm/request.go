package vm

import (
	"github.com/gin-gonic/gin"
	"github.com/kigland/HPC-Scheduler/coodinator/controller/mid"
	"github.com/kigland/HPC-Scheduler/coodinator/models/dboper"
)

func request(c *gin.Context) {

	uid := c.GetString(mid.MID_USER_ID)
	user, err := dboper.GetUserByID(uid)
	if err != nil {
		c.JSON(400, gin.H{
			"message": "user not found",
		})
		return
	}

	if user.MaxVCPU == -1 {
		c.JSON(400, gin.H{
			"message": "user not found",
		})
		return
	}

}
