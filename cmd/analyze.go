package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/xiaowen-0725/yifei-cli/internal/analyze"
	"github.com/xiaowen-0725/yifei-cli/internal/config"
	"github.com/xiaowen-0725/yifei-cli/internal/output"
)

func newAnalyzeCmd(deps Deps) *cobra.Command {
	ac := &cobra.Command{Use: "analyze", Short: "预置业务分析"}
	ac.AddCommand(newAnalyzeListCmd())
	for _, tpl := range analyze.All() {
		ac.AddCommand(newAnalyzeRunCmd(deps, tpl))
	}
	return ac
}

func newAnalyzeListCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "列出所有分析模板",
		RunE: func(c *cobra.Command, _ []string) error {
			for _, t := range analyze.All() {
				fmt.Fprintf(c.OutOrStdout(), "%-28s %s\n", t.Name, t.Desc)
			}
			return nil
		},
	}
}

func newAnalyzeRunCmd(deps Deps, tpl analyze.Template) *cobra.Command {
	var year, top string
	cmd := &cobra.Command{
		Use:   tpl.Name,
		Short: tpl.Desc,
		RunE: func(c *cobra.Command, _ []string) error {
			sql, err := tpl.Build(map[string]string{"year": year, "top": top})
			if err != nil {
				return err
			}
			path, err := resolveConfigPath(c)
			if err != nil {
				return err
			}
			cfg, err := config.Load(path)
			if err != nil {
				return err
			}
			conn, err := deps.OpenDB(cfg.DSN())
			if err != nil {
				return err
			}
			defer conn.Close()
			cols, rows, err := conn.Query(sql, 0)
			if err != nil {
				return err
			}
			return output.Render(c.OutOrStdout(), GlobalFormat(c), cols, rows, nil)
		},
	}
	cmd.Flags().StringVar(&year, "year", "", "年份,如 2022")
	cmd.Flags().StringVar(&top, "top", "", "TOP N")
	return cmd
}
