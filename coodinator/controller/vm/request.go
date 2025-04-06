package vm

import (
	"log"

	"github.com/KevinZonda/GoX/pkg/randx"
	"github.com/gin-gonic/gin"
	"github.com/kigland/OpenHPC/coodinator/controller/mid"
	"github.com/kigland/OpenHPC/coodinator/models/dboper"
	"github.com/kigland/OpenHPC/coodinator/models/openapi"
	"github.com/kigland/OpenHPC/coodinator/shared"
	"github.com/kigland/OpenHPC/coodinator/utils"
	"github.com/kigland/OpenHPC/lib/consts"
	"github.com/kigland/OpenHPC/lib/hypervisor/dockerProv"
	"github.com/kigland/OpenHPC/lib/image"
	"github.com/kigland/OpenHPC/lib/svcTag"
)

func request(c *gin.Context) {
	_ = utils.BodyAsF[openapi.VmReq](c)
	uid := c.GetString(mid.MID_USER_ID)
	user, err := dboper.GetUserByID(uid)
	if err != nil {
		c.JSON(400, gin.H{
			"message": "user not found",
		})
		return
	}

	passwd := utils.RndId(32)

	img := image.Factory{
		Password: passwd,
		BindHost: consts.CONTAINER_HOST,
		BindPort: randx.RndRange(consts.LOW_PORT, consts.HIGH_PORT),
		Provider: dockerProv.ProviderDocker,
	}.Image(image.ImageTorchBook).WithGPU(1).WithAutoRestart()

	svgT := svcTag.New(user.ID)
	img.ContainerName = svgT.String()

	id, err := shared.DockerHelper.StartContainer(img, true)
	if err != nil {
		log.Fatalf("Failed to start container: %v", err)
		return
	}

	c.JSON(200, gin.H{
		"id": id,
	})
}
