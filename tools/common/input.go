package common

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/KevinZonda/GoX/pkg/intx"
	"github.com/KevinZonda/GoX/pkg/panicx"
	"github.com/KevinZonda/GoX/pkg/stringx"
	"github.com/kigland/OpenHPC/lib/image"
	"github.com/kigland/OpenHPC/lib/utils"
)

func InputWithPrompt(prompt string) string {
	prompt = strings.TrimSpace(prompt)
	if prompt != "" {
		prompt += " "
	}
	return rlStrWithPrompt(prompt)
}

func ValidateInputAsSvgTagPart(input string) {
	invalid := []string{"-", " ", "@", ";", "/"}
	for _, v := range invalid {
		if strings.Contains(input, v) {
			log.Fatalf("Input cannot contain '%s': %s", v, input)
			return
		}
	}
}

func InputNeedSSH() bool {
	ssh := stringx.TrimLower(InputWithPrompt("Need SSH? (y/n):"))
	switch ssh {
	case "y", "yes", "true", "1":
		return true
	case "n", "no", "false", "0":
		return false
	default:
		log.Fatalf("Invalid input: %s", ssh)
		return false
	}
}

func InputPort(left int, right int) int {
	if left > right {
		left, right = right, left
	}
	portStr := InputWithPrompt(fmt.Sprintf("Port of the container (%d-%d):", left, right))

	port, err := strconv.Atoi(portStr)
	panicx.NotNilErr(err)
	if !intx.InRange(port, left, right) {
		log.Fatalf("Invalid port: %d", port)
		return 0
	}
	return port
}

func InputOwner() string {
	owner := InputWithPrompt("Enter VNode's Owner:")
	if owner == "" {
		log.Fatalf("Owner cannot be empty")
		return ""
	}
	owner = stringx.TrimLower(owner)
	ValidateInputAsSvgTagPart(owner)
	return owner
}

func InputProject() string {
	project := InputWithPrompt("Enter VNode's Project (for SvcTag only):")
	project = stringx.TrimLower(project)
	ValidateInputAsSvgTagPart(project)
	return project
}

func InputTokenOrGenerate(minLen int) string {
	token := InputWithPrompt("Enter VNode's Token (Keep empty to generate):")
	if token == "" {
		goto generate
	}
	if len(token) < minLen {
		log.Fatalf("Token must be at least %d characters long", minLen)
		return ""
	}
generate:
	if token == "" {
		if minLen < 8 {
			minLen = 32
		}
		token = utils.RndId(minLen)
	}
	return token
}

func InputImage() image.AllowedImages {
	fmt.Println("Select allowed Images:")
	for i, img := range image.ALLOWED_IMAGES {
		fmt.Println(i, ")", img)
	}
	idx := InputWithPrompt("Index (default [0]):")
	if idx == "" {
		return image.ALLOWED_IMAGES[0]
	}
	idxInt, err := strconv.Atoi(idx)
	panicx.NotNilErr(err)
	if !intx.InRange(idxInt, 0, len(image.ALLOWED_IMAGES)-1) {
		log.Fatalf("Invalid index: %d", idxInt)
		return ""
	}
	return image.ALLOWED_IMAGES[idxInt]
}

func InputShmSize() int64 {
	shmSize := InputWithPrompt("ShmSize (MB) or (_GB):")
	if shmSize == "" {
		return 64
	}
	factor := 1
	if strings.HasSuffix(shmSize, "GB") {
		factor = 1024
		shmSize = strings.TrimSuffix(shmSize, "GB")
	}
	shmSizeInt, err := strconv.Atoi(shmSize)
	panicx.NotNilErr(err)
	return int64(shmSizeInt) * int64(factor)
}
