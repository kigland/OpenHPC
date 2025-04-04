package common

import (
	"fmt"
	"strings"

	"github.com/KevinZonda/GoX/pkg/panicx"
	"github.com/kigland/OpenHPC/lib/consts"
	"github.com/kigland/OpenHPC/lib/dockerHelper"
	"github.com/kigland/OpenHPC/lib/image"
	"github.com/kigland/OpenHPC/lib/rds"
	"github.com/kigland/OpenHPC/lib/svcTag"
)

var Rds = rds.RDS{
	BasePath: consts.RDS_PATH,
}

func getRDS(username string, imageName image.AllowedImages) (rdsDir string, rdsMountAt string) {
	subfolder := InputWithPrompt("RDS Submodule (default \"\")")

	rdsDir, err := Rds.GetRDSPath(username, subfolder)
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
	rdsDir, rdsMountAt := getRDS(username, imageName)
	tag := svcTag.New(username).WithProject(project)
	return CreateContainerCustomRDS(dk, imageName, tag, passwd, port, rdsDir, rdsMountAt)
}

func CreateContainerCustomRDS(dk *dockerHelper.DockerHelper, imageName image.AllowedImages, tag svcTag.SvcTag, passwd string, port int, rdsDir string, rdsMountAt string) (ContainerInfo, error) {
	img := image.Factory{
		Password: passwd,
		BindHost: consts.CONTAINER_HOST,
		BindPort: port,
	}.Image(imageName).WithGPU(false, 1).WithAutoRestart()

	img.ContainerName = tag.String()

	img = img.WithMountRW(rdsDir, rdsMountAt)

	shmSize := InputShmSize()
	img = img.WithShmSize(shmSize)

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
