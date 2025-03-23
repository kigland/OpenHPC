package rds

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/KevinZonda/GoX/pkg/stringx"
)

type RDS struct {
	BasePath string
}

func (r *RDS) GetRDSPath(username string, subfolder string) (path string, err error) {
	path, err = r.rdsPath(username, subfolder)
	if err != nil {
		return "", err
	}
	if _, err := os.Stat(path); err != nil {
		return "", err
	}
	return path, nil
}

func (r *RDS) rdsPath(username string, subfolder string) (string, error) {
	username = stringx.TrimLower(username)
	if username == "" {
		return "", os.ErrInvalid
	}
	subfolder = strings.TrimSpace(subfolder)
	if strings.Contains(subfolder, ".") {
		return "", os.ErrInvalid
	}
	if subfolder == "" {
		return filepath.Join(r.BasePath, username), nil
	}
	return filepath.Join(r.BasePath, username, subfolder), nil
}
