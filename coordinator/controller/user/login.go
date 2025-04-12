package user

import (
	"github.com/gin-gonic/gin"
	"github.com/kigland/OpenHPC/coordinator/models/apimod"
	"github.com/kigland/OpenHPC/coordinator/models/dboper"
	"github.com/kigland/OpenHPC/coordinator/utils"
)

func login(c *gin.Context) {
	body := utils.BodyAsF[apimod.LoginReq](c)

	user, err := dboper.GetUserByID(body.Username)
	if err != nil {
		utils.ErrorMsg(c, 400, "user not found")
		return
	}

	if user.Password != body.Password {
		utils.ErrorMsg(c, 400, "password is incorrect")
		return
	}

	token := utils.RndId(32)
	err = dboper.CreateToken(user.ID, token)
	if err != nil {
		utils.ErrorMsg(c, 400, "failed to create token")
		return
	}

	c.JSON(200, apimod.Token{
		Token: token,
	})
}
