package main

import (
	"fmt"
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
	case "list", "ls", "ps":
		handler.List()
	case "env", "e":
		handler.Env()
	case "token", "t":
		handler.Token()
	default:
		help()
	}
}

func help() {
	h := `
Usage hpc [command]:
  - req|request|create|c: create a new VNode
  - list|ls|ps          : list all VNodes
  - env|e      [cid?]   : show environment variables of the VNode
  - token|t    [cid?]   : show tokens of the VNode
`
	fmt.Println(strings.TrimSpace(h))
}
