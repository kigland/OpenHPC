package main

import (
	"os"
	"strings"

	"github.com/kigland/HPC-Scheduler/tools/common"
	"github.com/kigland/HPC-Scheduler/tools/handler"
)

func main() {
	common.InitRL()
	defer common.Rl.Close()
	common.InitDocker()

	switch strings.ToLower(os.Args[1]) {
	case "req", "request":
		handler.Request()
	case "list":
		handler.List()
	case "list-user":
		handler.ListUser(os.Args[2])
	}
}
