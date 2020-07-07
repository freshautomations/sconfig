package main

import (
	"github.com/freshautomations/sconfig/cmd"
	"github.com/freshautomations/sconfig/exit"
)

func main() {
	if err := cmd.Execute(); err != nil {
		exit.Fail(err)
	}
}
