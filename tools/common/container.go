package common

import (
	"fmt"
	"strings"

	"github.com/KevinZonda/GoX/pkg/panicx"
	"github.com/kigland/OpenHPC/lib/consts"
	"github.com/kigland/OpenHPC/lib/hypervisor/dockerProv"
	"github.com/kigland/OpenHPC/lib/image"
	"github.com/kigland/OpenHPC/lib/rds"
	"github.com/kigland/OpenHPC/lib/svcTag"
	"github.com/kigland/OpenHPC/lib/utils"
)

var Rds = rds.RDS{
	BasePath: consts.RDS_PATH,
}

func GetRDSWithSubfolder(owner string, subfolder string, imageName image.AllowedImages) (rdsDir string, rdsMountAt string) {
	rdsDir, err := Rds.GetRDSPath(owner, subfolder)
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
		sb.WriteString(fmt.Sprintf("SSH    : ssh -p %d jovyan@%s\n", c.SSHPort, consts.SSH_BIND_HOST))
	}
	sb.WriteString(fmt.Sprintf("Token  : %s\n", c.Token))
	sb.WriteString(fmt.Sprintf("RDS    : %s\n", c.RDSAt))
	sb.WriteString(fmt.Sprintf("CID    : %s\n", c.CID))
	sb.WriteString(fmt.Sprintf("SCID   : %s\n", utils.ShortId(c.CID)))
	sb.WriteString(fmt.Sprintf("SvcTag : %s\n", c.SvcTag.String()))
	sb.WriteString(fmt.Sprintf("SC     : %s", c.SvcTag.ShortCode()))
	return sb.String()
}

func CreateContainerCustomRDS(dk *dockerProv.DockerHelper, imageName image.AllowedImages, tag svcTag.SvcTag, passwd string, port int, rdsDir string, rdsMountAt string, needSSH bool) (ContainerInfo, error) {
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
		Provider:    LoadProvider(),
	}.Image(imageName).WithGPU(1).
		WithAutoRestart().
		WithBaseURL(imageName.BaseURLEnvVar(), consts.BASE_URL(port))

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
