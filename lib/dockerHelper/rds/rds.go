package rds

import (
	"os"
	"path/filepath"
	"strings"
)

type RDS struct {
	BasePath string
}

func (r *RDS) GetRDSPath(username string, subfolder string) (path string, err error) {
	path = filepath.Join(r.BasePath, strings.ToLower(username))
	if _, err := os.Stat(path); err != nil {
		return "", err
	}

	subfolder = strings.TrimSpace(subfolder)
	if strings.Contains(subfolder, "..") {
		return "", os.ErrInvalid
	}
	if subfolder == "" {
		return path, nil
	}

	path = filepath.Join(path, subfolder)
	if _, err := os.Stat(path); err != nil {
		return "", err
	}
	return path, nil
}
