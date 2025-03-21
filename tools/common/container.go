package common

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/KevinZonda/GoX/pkg/panicx"
	kon "github.com/kigland/HPC-Scheduler/coodinator/container"
	"github.com/kigland/HPC-Scheduler/lib/dockerHelper"
	"github.com/kigland/HPC-Scheduler/lib/dockerHelper/image"
	"github.com/kigland/HPC-Scheduler/lib/dockerHelper/rds"
)

var r = rds.RDS{
	BasePath: "/data/rds",
}

func getRDS(username string, imageName image.AllowedImages) (rdsDir string, rdsMountAt string) {
	fmt.Println("RDS Submodule?")
	subfolder, err := Rl.Readline()
	panicx.NotNilErr(err)

	rdsDir, err = r.GetRDSPath(username, subfolder)
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
	sb.WriteString(fmt.Sprintf("URL    : http://127.0.0.2:%d\n", c.Port))
	sb.WriteString(fmt.Sprintf("Token  : %s\n", c.Token))
	sb.WriteString(fmt.Sprintf("CID    : %s\n", c.CID))
	sb.WriteString(fmt.Sprintf("RDS    : %s\n", c.RDSAt))
	sb.WriteString(fmt.Sprintf("SvcTag : %s", c.CName))
	return sb.String()
}

func CreateContainer(dk *dockerHelper.DockerHelper, imageName image.AllowedImages, username, passwd string, port int) (ContainerInfo, error) {
	img := image.Factory{
		Password: passwd,
		BindHost: "127.0.0.2",
		BindPort: strconv.Itoa(port),
	}.Image(imageName).WithGPU(1)
	img.AutoRemove = true

	img.ContainerName = kon.NewContainerName(username)

	rdsDir, rdsMountAt := getRDS(username, imageName)
	img = img.WithRDS(rdsDir, rdsMountAt)

	id, err := dk.StartContainer(img)
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
