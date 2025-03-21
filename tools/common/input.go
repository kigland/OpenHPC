package common

import (
	"fmt"
	"log"
	"strconv"

	"github.com/KevinZonda/GoX/pkg/panicx"
)

func InputPort(left int, right int) int {
	if left > right {
		left, right = right, left
	}
	fmt.Printf("Port of the container (%d-%d):\n", left, right)
	portStr := rlStr()

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
	username := rlStr()
	if username == "" {
		log.Fatalf("Username cannot be empty")
		return ""
	}
	return username
}
