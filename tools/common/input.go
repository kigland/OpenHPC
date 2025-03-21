package common

import (
	"fmt"
	"log"
	"strconv"

	"github.com/KevinZonda/GoX/pkg/panicx"
	"github.com/kigland/HPC-Scheduler/lib/utils"
)

func InputWithPrompt(prompt string) string {
	fmt.Println(prompt)
	return rlStr()
}

func InputPort(left int, right int) int {
	if left > right {
		left, right = right, left
	}
	portStr := InputWithPrompt(fmt.Sprintf("Port of the container (%d-%d):", left, right))

	port, err := strconv.Atoi(portStr)
	panicx.NotNilErr(err)
	if port < left || port > right {
		log.Fatalf("Invalid port: %d", port)
		return 0
	}
	return port
}

func InputUsername() string {
	username := InputWithPrompt("Username:")
	if username == "" {
		log.Fatalf("Username cannot be empty")
		return ""
	}
	return username
}

func InputProject() string {
	project := InputWithPrompt("Project:")
	return utils.TrimLower(project)
}
