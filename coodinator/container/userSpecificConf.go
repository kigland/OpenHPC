package container

import (
	"fmt"
	"path/filepath"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/kigland/OpenHPC/coodinator/models/dbmod"
	"github.com/kigland/OpenHPC/coodinator/shared"
	"github.com/kigland/OpenHPC/lib/hypervisor/dockerHelper"
)

type UserSpecificConf struct {
	dbmod.User
	RequestedMem  int // in MB
	RequestedCPU  int // in Core
	RequestedGPU  int // in Cards
	RequestedTime time.Duration
	ImageName     string
}

func UserToSpecificConf(user dbmod.User) UserSpecificConf {
	return UserSpecificConf{
		User: user,
	}
}

func (u *UserSpecificConf) Normalize() {
	if u.RequestedGPU < 0 {
		if u.MaxVGPU == -1 {
			u.RequestedGPU = 1
		} else {
			u.RequestedGPU = u.MaxVGPU
		}
	}
	if u.RequestedCPU <= 0 {
		if u.MaxVCPU == -1 {
			u.RequestedCPU = 0
		} else {
			u.RequestedCPU = u.MaxVCPU
		}
	}
	if u.RequestedMem <= 0 {
		if u.MaxMemory == -1 {
			u.RequestedMem = 0
		} else {
			u.RequestedMem = u.MaxMemory
		}
	}
}

func (u UserSpecificConf) Validate() error {
	if u.MaxVGPU != -1 && u.RequestedGPU > u.MaxVGPU {
		return fmt.Errorf("requested GPU count is greater than max GPU count")
	}
	if u.MaxVCPU != -1 && u.RequestedCPU > u.MaxVCPU {
		return fmt.Errorf("requested CPU count is greater than max CPU count")
	}
	if u.MaxMemory != -1 && u.RequestedMem > u.MaxMemory {
		return fmt.Errorf("requested memory is greater than max memory")
	}
	return nil
}

func (u UserSpecificConf) GetStoragePath() string {
	return filepath.Join(shared.GetConfig().Storage, "user", u.User.ID)
}

func (u UserSpecificConf) GetDockerOpts() dockerHelper.StartContainerOptions {
	return dockerHelper.StartContainerOptions{
		ImageName: u.ImageName,
		Binds: []string{
			u.GetStoragePath() + ":/rds",
		},
		Resources: container.Resources{
			Memory:         int64(u.RequestedMem * 1024 * 1024), // MB -> B
			NanoCPUs:       int64(u.RequestedCPU * 1000000000),  // Number of CPUs in nanoseconds precision
			DeviceRequests: dockerHelper.GetGPUDeviceRequests(u.RequestedGPU),
		},
	}
}
