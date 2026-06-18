package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/xiaowen-0725/yifei-cli/internal/config"
	"github.com/xiaowen-0725/yifei-cli/internal/output"
	"github.com/xiaowen-0725/yifei-cli/internal/safety"
)

func newQueryCmd(deps Deps) *cobra.Command {
	var file string
	var translate bool
	var limit int
	cmd := &cobra.Command{
		Use:   "query [SQL]",
		Short: "执行只读 SQL 查询",
		RunE: func(c *cobra.Command, args []string) error {
			sql, err := readSQL(args, file)
			if err != nil {
				return err
			}
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

			cols, rows, err := conn.Query(sql, limit)
			if err != nil {
				return err
			}
			var tr func(string) (string, bool)
			if translate && deps.Dict != nil {
				tr = deps.Dict.FieldName
			}
			return output.Render(c.OutOrStdout(), GlobalFormat(c), cols, rows, tr)
		},
	}
	cmd.Flags().StringVarP(&file, "file", "f", "", "从文件读取 SQL")
	cmd.Flags().BoolVar(&translate, "translate", false, "列头按字典翻译为中文")
	cmd.Flags().IntVar(&limit, "limit", 1000, "最大返回行数 (0=不限制)")
	return cmd
}

func readSQL(args []string, file string) (string, error) {
	if file != "" {
		b, err := os.ReadFile(file)
		if err != nil {
			return "", err
		}
		return string(b), nil
	}
	if len(args) == 0 {
		return "", fmt.Errorf("请提供 SQL (参数) 或 -f 文件")
	}
	return args[0], nil
}
