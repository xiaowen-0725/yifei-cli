package cmd

import (
	"github.com/spf13/cobra"
	"github.com/xiaowen-0725/yifei-cli/internal/dict"
	"github.com/xiaowen-0725/yifei-cli/internal/schema"
)

// Querier is the interface satisfied by *db.DB.
type Querier interface {
	Query(sql string, limit int) ([]string, [][]any, error)
	Close() error
}

// Deps carries shared services injected from main.
type Deps struct {
	Dict   *dict.Dict
	Schema *schema.Schema
	OpenDB func(dsn string) (Querier, error)
}

func NewRootCmd(deps Deps) *cobra.Command {
	root := &cobra.Command{
		Use:          "yifei",
		Short:        "易飞 ERP 只读数据分析 CLI",
		Long:         "yifei-cli — 供人类与 AI Agent 对易飞 ERP 数据库进行只读查询、结构探索与字段翻译。",
		SilenceUsage: true,
	}
	root.PersistentFlags().String("format", "table", "输出格式: table|json|csv")
	root.PersistentFlags().String("config", "", "配置文件路径 (默认: OS 配置目录/yifei-cli/config.yaml)")
	root.PersistentFlags().Bool("quiet", false, "精简输出")
	root.AddCommand(newConfigCmd())
	root.AddCommand(newQueryCmd(deps))
	root.AddCommand(newSchemaCmd(deps))
	root.AddCommand(newDictCmd(deps))
	return root
}

func GlobalFormat(c *cobra.Command) string {
	v, _ := c.Flags().GetString("format")
	return v
}

func GlobalConfigPath(c *cobra.Command) string {
	v, _ := c.Flags().GetString("config")
	return v
}
