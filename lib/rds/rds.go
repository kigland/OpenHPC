package rds

import (
	"errors"
	"os"
	"path/filepath"
	"strings"

	"github.com/KevinZonda/GoX/pkg/stringx"
)

type RDS struct {
	BasePath string
}

func (r *RDS) GetRDSPath(username string, subfolder string, createIfNotExists bool) (path string, err error) {
	path, err = r.rdsPath(username, subfolder)
	if err != nil {
		return "", err
	}
	_, err = os.Stat(path)
	if err == nil {
		return path, nil
	}
	if errors.Is(err, os.ErrNotExist) && createIfNotExists {
		err = r.Create(username, subfolder)
		if err != nil {
			return "", err
		}
		return path, nil
	}
	return "", err
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
