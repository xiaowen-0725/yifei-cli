package cmd

import (
	"fmt"
	"sort"
	"strings"

	"github.com/spf13/cobra"
	"github.com/xiaowen-0725/yifei-cli/internal/dict"
)

func newSchemaCmd(deps Deps) *cobra.Command {
	sc := &cobra.Command{Use: "schema", Short: "探索表结构 (离线,读内嵌 schema)"}
	sc.AddCommand(
		newSchemaTablesCmd(deps),
		newSchemaSearchCmd(deps),
		newSchemaTableCmd(deps),
		newSchemaModulesCmd(deps),
	)
	return sc
}

func newSchemaTablesCmd(deps Deps) *cobra.Command {
	var module string
	cmd := &cobra.Command{
		Use:   "tables",
		Short: "列出表 (可按模块过滤)",
		RunE: func(c *cobra.Command, _ []string) error {
			for _, t := range deps.Schema.TablesByModule(module) {
				fmt.Fprintln(c.OutOrStdout(), t)
			}
			return nil
		},
	}
	cmd.Flags().StringVar(&module, "module", "", "按模块前缀过滤,如 COP")
	return cmd
}

func newSchemaSearchCmd(deps Deps) *cobra.Command {
	return &cobra.Command{
		Use:   "search <关键词>",
		Short: "按表名/类型/中文名搜表",
		Args:  cobra.ExactArgs(1),
		RunE: func(c *cobra.Command, args []string) error {
			seen := map[string]bool{}
			var hits []string
			add := func(codes []string) {
				for _, t := range codes {
					if !seen[t] {
						seen[t] = true
						hits = append(hits, t)
					}
				}
			}
			add(deps.Schema.SearchTables(args[0])) // 表名 / 后缀类型
			if deps.Dict != nil {
				add(deps.Dict.SearchTables(args[0])) // 中文名 / 字段名 / 备注
			}
			if len(hits) == 0 {
				fmt.Fprintln(c.OutOrStdout(), "无匹配")
				return nil
			}
			sort.Strings(hits)
			for _, t := range hits {
				cname := ""
				if deps.Dict != nil {
					if td, ok := deps.Dict.Table(t); ok {
						cname = td.Name
					}
				}
				if cname != "" {
					fmt.Fprintf(c.OutOrStdout(), "%-10s %s\n", t, cname)
				} else {
					fmt.Fprintln(c.OutOrStdout(), t)
				}
			}
			return nil
		},
	}
}

func newSchemaTableCmd(deps Deps) *cobra.Command {
	return &cobra.Command{
		Use:   "table <表名>",
		Short: "查看某表全部字段",
		Args:  cobra.ExactArgs(1),
		RunE: func(c *cobra.Command, args []string) error {
			tbl, ok := deps.Schema.Table(args[0])
			if !ok {
				return fmt.Errorf("表不存在: %s", args[0])
			}
			out := c.OutOrStdout()
			fmt.Fprintf(out, "表: %s  模块: %s  类型: %s  行数: %d\n",
				strings.ToUpper(args[0]), tbl.Module, tbl.SuffixType, tbl.RowCount)
			var hasDict bool
			var td dict.TableDict
			if deps.Dict != nil {
				td, hasDict = deps.Dict.Table(args[0])
			}
			for _, col := range tbl.Columns {
				cn := ""
				if hasDict {
					if f, ok := td.Fields[col.Name]; ok {
						cn = f.Name
					}
				}
				fmt.Fprintf(out, "  %-12s %-14s %s\n", col.Name, col.DataType, cn)
			}
			return nil
		},
	}
}

func newSchemaModulesCmd(deps Deps) *cobra.Command {
	return &cobra.Command{
		Use:   "modules",
		Short: "列出所有模块",
		RunE: func(c *cobra.Command, _ []string) error {
			for _, m := range deps.Schema.ModuleList() {
				fmt.Fprintf(c.OutOrStdout(), "%-6s %-20s %d 表\n", m.Code, m.Name, m.TableCount)
			}
			return nil
		},
	}
}
