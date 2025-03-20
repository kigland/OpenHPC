package common

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/KevinZonda/GoX/pkg/panicx"
)

func InputPort(left int, right int) int {
	if left > right {
		left, right = right, left
	}
	fmt.Printf("Port of the container (%d-%d):\n", left, right)
	portStr, err := Rl.Readline()
	panicx.NotNilErr(err)

	port, err := strconv.Atoi(portStr)
	panicx.NotNilErr(err)
	if port < left || port > right {
		log.Fatalf("Invalid port: %d", port)
		return 0
	}
	return port
}

func InputUsername() string {
	fmt.Println("Username:")
	username, err := Rl.Readline()
	panicx.NotNilErr(err)
	username = strings.TrimSpace(username)
	if username == "" {
		log.Fatalf("Username cannot be empty")
		return ""
	}
	return username
}
