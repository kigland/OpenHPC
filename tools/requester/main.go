package main

import (
	"github.com/kigland/HPC-Scheduler/tools/common"
	"github.com/kigland/HPC-Scheduler/tools/handler"
)

func main() {
	common.InitRL()
	defer common.Rl.Close()
	common.InitDocker()

	handler.Request()
}
