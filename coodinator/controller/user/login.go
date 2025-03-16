package user

import (
	"github.com/gin-gonic/gin"
	"github.com/kigland/HPC-Scheduler/coodinator/models/dboper"
	"github.com/kigland/HPC-Scheduler/coodinator/models/openapi"
	"github.com/kigland/HPC-Scheduler/coodinator/utils"
)

func login(c *gin.Context) {
	body := utils.BodyAsF[openapi.LoginReq](c)

	user, err := dboper.GetUserByID(body.Username)
	if err != nil {
		c.JSON(400, gin.H{
			"message": "user not found",
		})
		return
	}

	if user.Password != body.Password {
		c.JSON(400, gin.H{
			"message": "password is incorrect",
		})
		return
	}

	token := utils.RndId(32)
	err = dboper.CreateToken(user.ID, token)
	if err != nil {
		c.JSON(400, gin.H{
			"message": "failed to create token",
		})
		return
	}

	c.JSON(200, openapi.Token{
		Token: token,
	})
}
