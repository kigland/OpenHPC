package image

import (
	"path/filepath"
	"slices"
)

type AllowedImages string

const (
	ImageJupyterHub    AllowedImages = "kevinzonda/notebook"
	ImageBase          AllowedImages = "jupyterhub/singleuser"
	ImageJupyterHubIso AllowedImages = "kevinzonda/notebook-iso"
)

var ALLOWED_IMAGES = []AllowedImages{
	ImageJupyterHub,
	ImageJupyterHubIso,
	ImageBase,
}

func (a AllowedImages) IsAllowed() bool {
	return slices.Contains(ALLOWED_IMAGES, a)
}

func (a AllowedImages) HomeDir() string {
	switch a {
	case ImageJupyterHub, ImageBase, ImageJupyterHubIso:
		return "/home/jovyan"
	}
	return ""
}

func (a AllowedImages) BaseURLEnvVar() string {
	switch a {
	case ImageJupyterHub, ImageJupyterHubIso:
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

func (a AllowedImages) SupportSSH() bool {
	switch a {
	case ImageJupyterHub:
		return true
	case ImageJupyterHubIso:
		return false
	}
	return false
}
