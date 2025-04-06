package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/kigland/OpenHPC/tools/common"
	"github.com/kigland/OpenHPC/tools/handler"
)

func main() {
	if len(os.Args) < 2 {
		help()
		return
	}
	common.InitRL()
	defer common.Rl.Close()

	var f func()

	switch strings.ToLower(os.Args[1]) {
	case "req", "request", "create", "c", "start":
		f = handler.Request
	case "list", "ls", "ps", "ll", "l", "status":
		f = handler.List
	case "env", "e":
		f = handler.Env
	case "token", "t", "tk":
		f = handler.Token
	case "stop", "s", "rm", "remove":
		f = handler.Stop
	case "ids", "id":
		f = handler.IDs
	case "rds", "r":
		handler.RDS()
		os.Exit(0)
	case "upd", "upgrade", "u":
		f = handler.Upd
	default:
		help()
		os.Exit(1)
	}

	common.InitConfig()
	common.InitDocker()
	f()
}

func help() {
	h := `
Usage hpc [command]:
  - req|request|create|c|start  : create a new VNode
  - stop|s|rm|remove [node_id?] : stop the VNode
  - list|ls|ps|ll|l|status      : list all VNodes
  - env|e         [node_id?]    : show environment variables of the VNode
  - token|t|tk    [node_id?]    : show tokens of the VNode
  - ids|id        [node_id?]    : show CID/SCID/SvcTag/ShortCode of the VNode
  - rds|r         [action]      : manage RDS
  - upd|upgrade|u [node_id?]    : upgrade the VNode
`
	fmt.Println(strings.TrimSpace(h))
}
