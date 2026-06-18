package main

import (
	"fmt"
	"os"

	"github.com/xiaowen-0725/yifei-cli/cmd"
)

func main() {
	root := cmd.NewRootCmd(cmd.Deps{})
	if err := root.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, "错误:", err)
		os.Exit(1)
	}
}
