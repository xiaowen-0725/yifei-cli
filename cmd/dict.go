package cmd

import (
	"fmt"
	"sort"
	"strings"

	"github.com/spf13/cobra"
)

func newDictCmd(deps Deps) *cobra.Command {
	dc := &cobra.Command{Use: "dict", Short: "字段含义翻译 (离线,读内嵌字典)"}
	dc.AddCommand(
		newDictTableCmd(deps),
		newDictFieldCmd(deps),
		newDictCodeCmd(deps),
		newDictRelationsCmd(deps),
	)
	return dc
}

func newDictTableCmd(deps Deps) *cobra.Command {
	return &cobra.Command{
		Use:   "table <表名>",
		Short: "某表字段编号→中文对照",
		Args:  cobra.ExactArgs(1),
		RunE: func(c *cobra.Command, args []string) error {
			td, ok := deps.Dict.Table(args[0])
			if !ok {
				return fmt.Errorf("字典中无此表: %s (仅核心表已标注)", args[0])
			}
			out := c.OutOrStdout()
			fmt.Fprintf(out, "%s — %s\n", strings.ToUpper(args[0]), td.Name)
			codes := make([]string, 0, len(td.Fields))
			for code := range td.Fields {
				codes = append(codes, code)
			}
			sort.Strings(codes)
			for _, code := range codes {
				f := td.Fields[code]
				fmt.Fprintf(out, "  %-8s %s  %s\n", code, f.Name, f.Note)
			}
			return nil
		},
	}
}

func newDictFieldCmd(deps Deps) *cobra.Command {
	return &cobra.Command{
		Use:   "field <TABLE.CODE 或 CODE>",
		Short: "单字段含义",
		Args:  cobra.ExactArgs(1),
		RunE: func(c *cobra.Command, args []string) error {
			arg := args[0]
			out := c.OutOrStdout()
			if strings.Contains(arg, ".") {
				parts := strings.SplitN(arg, ".", 2)
				td, ok := deps.Dict.Table(parts[0])
				if ok {
					if f, ok := td.Fields[strings.ToUpper(parts[1])]; ok {
						fmt.Fprintf(out, "%s = %s  %s\n", arg, f.Name, f.Note)
						return nil
					}
				}
				return fmt.Errorf("未找到字段: %s", arg)
			}
			if n, ok := deps.Dict.FieldName(arg); ok {
				fmt.Fprintf(out, "%s = %s\n", strings.ToUpper(arg), n)
				return nil
			}
			return fmt.Errorf("未找到字段: %s", arg)
		},
	}
}

func newDictCodeCmd(deps Deps) *cobra.Command {
	return &cobra.Command{
		Use:   "code <类别> [表]",
		Short: "查参数码表 (如 单别/仓库/部门)",
		Args:  cobra.RangeArgs(1, 2),
		RunE: func(c *cobra.Command, args []string) error {
			key := args[0]
			if len(args) == 2 {
				key = strings.ToUpper(args[1]) + "." + args[0]
			}
			codes, ok := deps.Dict.LookupCodes(key)
			if !ok {
				return fmt.Errorf("未找到码表: %s", key)
			}
			keys := make([]string, 0, len(codes))
			for k := range codes {
				keys = append(keys, k)
			}
			sort.Strings(keys)
			for _, k := range keys {
				fmt.Fprintf(c.OutOrStdout(), "  %-6s %s\n", k, codes[k])
			}
			return nil
		},
	}
}

func newDictRelationsCmd(deps Deps) *cobra.Command {
	return &cobra.Command{
		Use:   "relations <表名>",
		Short: "某表的头身/外联关系",
		Args:  cobra.ExactArgs(1),
		RunE: func(c *cobra.Command, args []string) error {
			rels, ok := deps.Dict.LookupRelations(args[0])
			if !ok {
				return fmt.Errorf("字典中无此表的关系: %s", args[0])
			}
			for _, r := range rels {
				fmt.Fprintln(c.OutOrStdout(), "  "+r)
			}
			return nil
		},
	}
}
