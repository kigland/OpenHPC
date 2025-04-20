package main

import (
	"fmt"

	"github.com/KevinZonda/GoX/pkg/panicx"
	"github.com/kigland/OpenHPC/lib/nv"
)

func main() {
	log, err := nv.GetNvidiaSmiLog()
	panicx.NotNilErr(err)
	info, err := log.Parse()
	panicx.NotNilErr(err)
	fmt.Println(info)
}
