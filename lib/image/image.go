package image

import (
	"path/filepath"
	"slices"
)

type AllowedImages string

const (
	ImageJupyterHub AllowedImages = "kevinzonda/notebook"
	ImageBase       AllowedImages = "jupyterhub/singleuser"
)

var ALLOWED_IMAGES = []AllowedImages{
	ImageJupyterHub,
	ImageBase,
}

func (a AllowedImages) IsAllowed() bool {
	return slices.Contains(ALLOWED_IMAGES, a)
}

func (a AllowedImages) HomeDir() string {
	switch a {
	case ImageJupyterHub, ImageBase:
		return "/home/jovyan"
	}
	return ""
}

func (a AllowedImages) BaseURLEnvVar() string {
	switch a {
	case ImageJupyterHub:
		return "NB_VAR_BASE_URL"
	}
	return ""
}

func (a AllowedImages) RdsDir() string {
	home := a.HomeDir()
	if home != "" {
		return filepath.Join(home, "rds")
	}
	return "/rds"
}
