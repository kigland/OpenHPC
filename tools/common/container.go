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
	CID     string
	RDSAt   string
	Token   string
	Port    int
	SSHPort int
	SvcTag  svcTag.SvcTag
}

func (c ContainerInfo) String() string {
	sb := strings.Builder{}
	sb.WriteString(fmt.Sprintf("URL    : http://%s:%d\n", consts.CONTAINER_HOST, c.Port))
	if c.SSHPort != 0 {
		sb.WriteString(fmt.Sprintf("SSH    : ssh -p %d jovyan@%s\n", c.SSHPort, consts.CONTAINER_HOST))
	}
	sb.WriteString(fmt.Sprintf("Token  : %s\n", c.Token))
	sb.WriteString(fmt.Sprintf("RDS    : %s\n", c.RDSAt))
	sb.WriteString(fmt.Sprintf("CID    : %s\n", c.CID))
	sb.WriteString(fmt.Sprintf("SCID   : %s\n", dockerHelper.ShortId(c.CID)))
	sb.WriteString(fmt.Sprintf("SvcTag : %s\n", c.SvcTag.String()))
	sb.WriteString(fmt.Sprintf("SC     : %s", c.SvcTag.ShortCode()))
	return sb.String()
}

func CreateContainer(dk *dockerHelper.DockerHelper, imageName image.AllowedImages, username, passwd string, port int, project string, needSSH bool) (ContainerInfo, error) {
	rdsDir, rdsMountAt := getRDS(username, imageName)
	tag := svcTag.New(username).WithProject(project)
	return CreateContainerCustomRDS(dk, imageName, tag, passwd, port, rdsDir, rdsMountAt, needSSH)
}

func CreateContainerCustomRDS(dk *dockerHelper.DockerHelper, imageName image.AllowedImages, tag svcTag.SvcTag, passwd string, port int, rdsDir string, rdsMountAt string, needSSH bool) (ContainerInfo, error) {
	sshPort := 0
	if needSSH {
		sshPort = port + consts.SSH_SHIFT
	}

	img := image.Factory{
		Password:    passwd,
		BindHost:    consts.CONTAINER_HOST,
		BindPort:    port,
		BindSSHHost: consts.SSH_BIND_HOST,
		BindSSHPort: sshPort,
	}.Image(imageName).WithGPU(1).WithAutoRestart()

	img.ContainerName = tag.String()

	img = img.WithMountRW(rdsDir, rdsMountAt)

	shmSize := InputShmSize()
	img = img.WithShmSize(shmSize)

	id, err := dk.StartContainer(img, true)
	if err != nil {
		return ContainerInfo{}, err
	}
	if tag.String() != img.ContainerName {
		fmt.Printf("container name mismatch: %s (SvcTag) != %s (CID)\n", tag.String(), img.ContainerName)
	}
	return ContainerInfo{
		CID:     id,
		RDSAt:   rdsMountAt,
		Token:   passwd,
		Port:    port,
		SSHPort: sshPort,
		SvcTag:  tag,
	}, nil
}
