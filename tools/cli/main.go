package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/kigland/HPC-Scheduler/tools/common"
	"github.com/kigland/HPC-Scheduler/tools/handler"
)

func main() {
	if len(os.Args) < 2 {
		help()
		return
	}
	common.InitRL()
	defer common.Rl.Close()
	common.InitDocker()

	switch strings.ToLower(os.Args[1]) {
	case "req", "request", "create", "c", "start":
		handler.Request()
	case "list", "ls", "ps", "ll", "l":
		handler.List()
	case "env", "e":
		handler.Env()
	case "token", "t":
		handler.Token()
	case "stop", "s", "rm", "remove":
		handler.Stop()
	case "ids", "id":
		handler.IDs()
	case "rds", "r":
		handler.RDS()
	case "upd", "upgrade":
		handler.Upd()
	default:
		help()
	}
}

func help() {
	h := `
Usage hpc [command]:
  - req|request|create|c|start  : create a new VNode
  - stop|s|rm|remove [node_id?] : stop the VNode
  - list|ls|ps|ll|l             : list all VNodes
  - env|e         [node_id?]    : show environment variables of the VNode
  - token|t       [node_id?]    : show tokens of the VNode
  - ids|id        [node_id?]    : show CID/SCID/SvcTag/ShortCode of the VNode
  - rds|r         [action]      : manage RDS
  - upd|upgrade|u [node_id?]    : upgrade the VNode
`
	fmt.Println(strings.TrimSpace(h))
}
