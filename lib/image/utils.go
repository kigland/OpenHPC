package image

import "strings"

func PruneImageStr(imgStr string) string {
	imgStr = strings.TrimPrefix(imgStr, "docker.io/")
	imgParts := strings.Split(imgStr, ":")
	return imgParts[0]
}
