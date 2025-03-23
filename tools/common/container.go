package common

import (
	"fmt"
	"strings"

	"github.com/KevinZonda/GoX/pkg/panicx"
	"github.com/kigland/HPC-Scheduler/lib/consts"
	"github.com/kigland/HPC-Scheduler/lib/dockerHelper"
	"github.com/kigland/HPC-Scheduler/lib/image"
	"github.com/kigland/HPC-Scheduler/lib/rds"
	"github.com/kigland/HPC-Scheduler/lib/svcTag"
)

var Rds = rds.RDS{
	BasePath: consts.RDS_PATH,
}

func getRDS(username string, imageName image.AllowedImages) (rdsDir string, rdsMountAt string) {
	fmt.Println("RDS Submodule?")
	subfolder, err := Rl.Readline()
	panicx.NotNilErr(err)

	rdsDir, err = Rds.GetRDSPath(username, subfolder)
	if err == nil {
		return rdsDir, imageName.RdsDir()
	}
	panicx.NotNilErr(err)
	return "", ""
}

type ContainerInfo struct {
	CID   string
	CName string
	RDSAt string
	Token string
	Port  int
}

func (c ContainerInfo) String() string {
	sb := strings.Builder{}
	sb.WriteString(fmt.Sprintf("URL    : http://%s:%d\n", consts.CONTAINER_HOST, c.Port))
	sb.WriteString(fmt.Sprintf("Token  : %s\n", c.Token))
	sb.WriteString(fmt.Sprintf("CID    : %s\n", c.CID))
	sb.WriteString(fmt.Sprintf("RDS    : %s\n", c.RDSAt))
	sb.WriteString(fmt.Sprintf("SvcTag : %s", c.CName))
	return sb.String()
}

func CreateContainer(dk *dockerHelper.DockerHelper, imageName image.AllowedImages, username, passwd string, port int, project string) (ContainerInfo, error) {
	img := image.Factory{
		Password: passwd,
		BindHost: consts.CONTAINER_HOST,
		BindPort: port,
	}.Image(imageName).WithGPU(1)
	img.AutoRemove = true

	svgT := svcTag.New(username).WithProject(project)
	img.ContainerName = svgT.String()

	rdsDir, rdsMountAt := getRDS(username, imageName)
	img = img.WithMountRW(rdsDir, rdsMountAt)

	id, err := dk.StartContainer(img, true)
	if err != nil {
		return ContainerInfo{}, err
	}
	return ContainerInfo{
		CID:   id,
		RDSAt: rdsMountAt,
		Token: passwd,
		Port:  port,
		CName: img.ContainerName,
	}, nil
}
