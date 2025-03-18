package common

import (
	"errors"
	"os"
	"path/filepath"
	"strings"

	"github.com/kigland/HPC-Scheduler/lib/dockerHelper/image"
)

func MountRDSInfo(imageName image.AllowedImages, username string, subfolder string) (rdsDir string, rdsMountAt string, err error) {
	rdsDir = filepath.Join("/data/rds", strings.ToLower(username))
	rdsMountAt = imageName.HomeDir()
	if rdsMountAt != "" {
		rdsMountAt = filepath.Join(rdsMountAt, "rds")
	} else {
		rdsMountAt = "/rds"
	}
	if _, err := os.Stat(rdsDir); err != nil {
		rdsDir = ""
		return "", "", err
	}

	subfolder = strings.TrimSpace(subfolder)
	if strings.Contains(subfolder, "..") {
		return "", "", errors.New("subfolder cannot contain '..'")
	}
	if subfolder != "" {
		rdsDir = filepath.Join(rdsDir, subfolder)
	}
	if _, err := os.Stat(rdsDir); err != nil {
		rdsDir = ""
		return "", "", err
	}
	return rdsDir, rdsMountAt, nil
}
