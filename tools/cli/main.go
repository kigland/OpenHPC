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
	case "req", "request", "create", "c":
		handler.Request()
	case "list", "ls":
		handler.List()
	case "list-user", "lu", "lsu":
		handler.ListUser(os.Args[2])
	}
}
