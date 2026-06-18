package main

import (
	"fmt"
	"os"

	"github.com/xiaowen-0725/yifei-cli/cmd"
	"github.com/xiaowen-0725/yifei-cli/internal/assets"
	"github.com/xiaowen-0725/yifei-cli/internal/db"
	"github.com/xiaowen-0725/yifei-cli/internal/dict"
	"github.com/xiaowen-0725/yifei-cli/internal/schema"
)

func main() {
	sc, err := schema.New(assets.SchemaJSON)
	if err != nil {
		fmt.Fprintln(os.Stderr, "错误:", err)
		os.Exit(1)
	}
	dc, err := dict.New(assets.DictYAML)
	if err != nil {
		fmt.Fprintln(os.Stderr, "错误:", err)
		os.Exit(1)
	}
	deps := cmd.Deps{
		Dict:   dc,
		Schema: sc,
		OpenDB: func(dsn string) (cmd.Querier, error) { return db.Open(dsn) },
	}
	root := cmd.NewRootCmd(deps)
	if err := root.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, "错误:", err)
		os.Exit(1)
	}
}
