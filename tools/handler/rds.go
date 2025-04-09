package handler

import (
	"fmt"
	"os"
	"strings"

	"github.com/KevinZonda/GoX/pkg/ruby"
	"github.com/kigland/OpenHPC/tools/common"
)

func RDS() {
	verb := PopFst()
	// [exec] [rds] [action]
	switch verb {
	case "list", "ls", "ps":
		rdsList()
	case "create", "new":
		rdsCreate()
	default:
		rdsHelp()
		os.Exit(1)
	}
}

func rdsHelp() {
	help := `
Usage:
  hpc rds [action]

Actions:
  list|ls|ps : List all RDS
  create|new : Create a new RDS
`
	fmt.Println(strings.TrimSpace(help))
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
	owner := common.InputOwner()
	subfolder := common.InputProject()
	err := common.Rds.Create(owner, subfolder)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("RDS created successfully")
}
