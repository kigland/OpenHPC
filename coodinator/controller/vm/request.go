package vm

import (
	"log"
	"math/rand"

	"github.com/gin-gonic/gin"
	"github.com/kigland/HPC-Scheduler/coodinator/controller/mid"
	"github.com/kigland/HPC-Scheduler/coodinator/models/dboper"
	"github.com/kigland/HPC-Scheduler/coodinator/models/openapi"
	"github.com/kigland/HPC-Scheduler/coodinator/shared"
	"github.com/kigland/HPC-Scheduler/coodinator/utils"
	"github.com/kigland/HPC-Scheduler/lib/consts"
	"github.com/kigland/HPC-Scheduler/lib/image"
	"github.com/kigland/HPC-Scheduler/lib/svcTag"
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
		BindPort: consts.LOW_PORT + rand.Intn(consts.HIGH_PORT-consts.LOW_PORT),
	}.Image(image.ImageTorchBook).WithGPU(1)
	img.AutoRemove = true

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
