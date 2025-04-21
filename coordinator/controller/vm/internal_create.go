package vm

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/kigland/OpenHPC/coordinator/models/apimod"
	"github.com/kigland/OpenHPC/coordinator/shared"
	"github.com/kigland/OpenHPC/lib/consts"
	"github.com/kigland/OpenHPC/lib/hypervisor/dockerProv"
	"github.com/kigland/OpenHPC/lib/image"
	"github.com/kigland/OpenHPC/lib/svcTag"
)

type CreateRequest struct {
	Provider dockerProv.Provider
	Dk       *dockerProv.DockerHelper
	Image    image.AllowedImages
	Tag      svcTag.SvcTag
	Passwd   string

	BindHost    string
	BindPort    int
	BindSSHHost string
	BindSSHPort int

	RdsDir     string
	RdsMountAt string
	ShmSize    int

	NeedGPU    bool
	MaxMemByte int64
}

type createdInfo struct {
	CID     string
	Image   image.AllowedImages
	RDSAt   string
	Token   string
	Port    int
	SSHPort int
	SvcTag  svcTag.SvcTag
}

func createdInfoToVMInfo(info createdInfo) apimod.VmCreatedInfo {
	resp := apimod.VmCreatedInfo{
		Cid:   info.CID,
		Image: string(info.Image),
		RdsAt: info.RDSAt,
		Token: info.Token,
		Http:  parseHTTPVisitURL(info.Port),

		SvcTag: info.SvcTag.String(),
		Sc:     info.SvcTag.ShortCode(),
	}
	if info.SSHPort != 0 && info.Image.SupportSSH() {
		resp.Ssh = shared.GetConfig().VisitSSHHost + ":" + strconv.Itoa(info.SSHPort)

	}
	return resp
}

func CreateContainerCustomRDS(req CreateRequest) (createdInfo, error) {
	sshPort := 0
	if req.BindSSHPort > 0 { //&& req.Image.SupportSSH() {
		sshPort = req.BindSSHPort
	}

	img := image.Factory{
		Password:    req.Passwd,
		BindHost:    req.BindHost,
		BindPort:    req.BindPort,
		BindSSHHost: req.BindSSHHost,
		BindSSHPort: sshPort,
		Provider:    req.Provider,
	}.Image(req.Image).
		WithAutoRestart().
		WithBaseURL(req.Image.BaseURLEnvVar(), consts.BASE_URL(req.BindPort)).
		WithShmSize(int64(req.ShmSize)).
		WithMaxMemoryByte(req.MaxMemByte)

	if req.NeedGPU {
		img = img.WithGPU(1)
	}

	img.ContainerName = req.Tag.String()

	img = img.WithMountRW(req.RdsDir, req.RdsMountAt)

	id, err := req.Dk.StartContainer(img, true)
	if err != nil {
		return createdInfo{}, err
	}
	if req.Tag.String() != img.ContainerName {
		log.Printf("container name mismatch: %s (SvcTag) != %s (CID)", req.Tag.String(), img.ContainerName)
	}
	return createdInfo{
		CID:     id,
		Image:   req.Image,
		RDSAt:   req.RdsMountAt,
		Token:   req.Passwd,
		Port:    req.BindPort,
		SSHPort: sshPort,
		SvcTag:  req.Tag,
	}, nil
}

type VNodeId struct {
	ID     string
	SvcTag svcTag.SvcTag
}

func IDs(dk *dockerProv.DockerHelper, cid string) (VNodeId, error) {
	summary, ok := dk.TryGetContainer(cid)
	if !ok {
		return VNodeId{}, fmt.Errorf("container not found or not managed by OHPC")
	}
	cid = summary.ID
	svcTag, err := svcTag.Parse(summary.Names[0])
	if err != nil {
		return VNodeId{}, fmt.Errorf("failed to parse service tag: %w", err)
	}
	return VNodeId{
		ID:     cid,
		SvcTag: svcTag,
	}, nil
}

func parseHTTPVisitURL(newPort int) string {
	return strings.ReplaceAll(shared.GetConfig().VisitHTTPHost, "$PORT", strconv.Itoa(newPort))
}
