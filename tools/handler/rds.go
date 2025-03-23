package handler

import (
	"fmt"
	"os"
	"strings"

	"github.com/KevinZonda/GoX/pkg/ruby"
	"github.com/kigland/HPC-Scheduler/tools/common"
)

func RDS() {
	// [exec] [rds] [action]
	switch os.Args[2] {
	case "list", "ls", "ps":
		rdsList()
	case "create", "new":
		rdsCreate()
	default:
		help := `
Usage:
  hpc rds [action]

Actions:
  list|ls|ps : List all RDS
  create|new : Create a new RDS
`
		fmt.Println(strings.TrimSpace(help))
	}
}

func rdsList() {
	names := common.Rds.List()
	if len(names) == 0 {
		fmt.Println("No RDS found")
		return
	}
	ruby.Apply(names, ruby.Println)

}
func rdsCreate() {
	username := common.InputUsername()
	subfolder := common.InputProject()
	err := common.Rds.Create(username, subfolder)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("RDS created successfully")
}
