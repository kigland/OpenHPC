package common

import (
	"fmt"
	"log"
	"strconv"

	"github.com/KevinZonda/GoX/pkg/intx"
	"github.com/KevinZonda/GoX/pkg/panicx"
	"github.com/KevinZonda/GoX/pkg/stringx"
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
	if !intx.InRange(port, left, right) {
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
	return stringx.TrimLower(project)
}

func InputTokenOrGenerate(minLen int) string {
	token := InputWithPrompt("Token:")
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
