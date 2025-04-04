package shared

import (
	"github.com/KevinZonda/GoX/pkg/panicx"
	"github.com/docker/docker/client"
	"github.com/kigland/OpenHPC/lib/consts"
	"github.com/kigland/OpenHPC/lib/dockerHelper"
	"github.com/kigland/OpenHPC/lib/rds"
)

var Docker *client.Client
var DockerHelper *dockerHelper.DockerHelper
var Rds *rds.RDS

func initDocker() {
	var err error
	Docker, err = client.NewClientWithOpts(client.WithHost(GetConfig().DockerHost), client.WithAPIVersionNegotiation())
	panicx.NotNilErr(err)
	DockerHelper = dockerHelper.NewDockerHelper(Docker)
	Rds = &rds.RDS{
		BasePath: consts.RDS_PATH,
	}
}
