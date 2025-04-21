package image

import (
	"path/filepath"
	"slices"
)

type HPCImage struct {
	Img  string `json:"img"`
	Home string `json:"home"`
	Rds  string `json:"rds"`

	SSH  int         `json:"ssh"`
	HTTP int         `json:"http"`
	Env  HPCImageEnv `json:"env"`
}

func (i *HPCImage) SupportSSH() bool {
	return i.SSH > 0
}

func (i *HPCImage) SupportHTTP() bool {
	return i.HTTP > 0
}

type HPCImageEnv struct {
	BaseURL string `json:"base_url"`
	Token   string `json:"token"`
}

type AllowedImages string

var allowedImages map[AllowedImages]HPCImage

func InitAllowedImages(images []HPCImage) {
	allowedImages = make(map[AllowedImages]HPCImage)
	allowedImageList := make([]AllowedImages, len(images))
	for i, image := range images {
		allowedImages[AllowedImages(image.Img)] = image
		allowedImageList[i] = AllowedImages(image.Img)
	}
	ALLOWED_IMAGES = allowedImageList
}

func InitDefaultAllowedImages() {
	InitAllowedImages([]HPCImage{
		{
			Img:  "kevinzonda/notebook",
			Home: "/home/jovyan",
			Rds:  "/home/jovyan/rds",
			SSH:  22,
			HTTP: 8888,
			Env: HPCImageEnv{
				BaseURL: "NB_VAR_BASE_URL",
				Token:   "JUPYTER_TOKEN",
			},
		},
		{
			Img:  "kevinzonda/notebook-iso",
			Home: "/home/jovyan",
			Rds:  "/home/jovyan/rds",
			HTTP: 8888,
			Env: HPCImageEnv{
				BaseURL: "NB_VAR_BASE_URL",
				Token:   "JUPYTER_TOKEN",
			},
		},
	})
}

var ALLOWED_IMAGES = []AllowedImages{}

func (a AllowedImages) IsAllowed() bool {
	return slices.Contains(ALLOWED_IMAGES, a)
}

func (a AllowedImages) Cfg() HPCImage {
	return allowedImages[a]
}

func (a AllowedImages) HomeDir() string {
	img, ok := allowedImages[a]
	if ok {
		return img.Home
	}
	return ""
}

func (a AllowedImages) BaseURLEnvVar() string {
	img, ok := allowedImages[a]
	if ok {
		return img.Env.BaseURL
	}
	return ""
}

func (a AllowedImages) RdsDir() string {
	img, ok := allowedImages[a]
	if ok {
		if img.Rds != "" {
			return img.Rds
		}
		if img.Home != "" {
			return filepath.Join(img.Home, "rds")
		}
	}
	return "/rds"
}

func (a AllowedImages) SupportSSH() bool {
	img, ok := allowedImages[a]
	if ok {
		return img.SupportSSH()
	}
	return false
}
