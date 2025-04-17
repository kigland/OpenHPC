package image

import "strings"

func PruneImageStr(imgStr string) string {
	return pruneImageStr(imgStr, false)
}

func PruneImageStrWithShortID(imgStr string) string {
	return pruneImageStr(imgStr, true)
}

func pruneImageStr(imgStr string, shortId bool) string {
	imgStr = strings.ToLower(imgStr)
	if strings.HasPrefix(imgStr, "sha256:") {
		imgStr = strings.TrimPrefix(imgStr, "sha256:")
		if shortId && len(imgStr) > 12 {
			return imgStr[:12]
		}
		return imgStr
	}
	imgStr = strings.TrimPrefix(imgStr, "docker.io/")
	imgParts := strings.Split(imgStr, ":")
	return imgParts[0]
}
