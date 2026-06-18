package cmd

import (
	"fmt"
	"regexp"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/xiaowen-0725/yifei-cli/internal/analyze"
	"github.com/xiaowen-0725/yifei-cli/internal/config"
	"github.com/xiaowen-0725/yifei-cli/internal/output"
	"github.com/xiaowen-0725/yifei-cli/internal/safety"
)

var reYear = regexp.MustCompile(`^\d{4}$`)

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
			// Validate flags before building SQL to prevent injection.
			if year != "" && !reYear.MatchString(year) {
				return fmt.Errorf("--year 必须是 4 位年份,如 2022")
			}
			if top != "" {
				n, err := strconv.Atoi(top)
				if err != nil || n <= 0 {
					return fmt.Errorf("--top 必须是正整数")
				}
			}
			sql, err := tpl.Build(map[string]string{"year": year, "top": top})
			if err != nil {
				return err
			}
			// Defense in depth: validate generated SQL is read-only.
			if err := safety.ValidateReadOnly(sql); err != nil {
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
