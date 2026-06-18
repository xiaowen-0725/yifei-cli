package cmd

import (
	"github.com/spf13/cobra"
)

// Deps carries shared services injected from main. Grows as tasks add packages.
type Deps struct {
	// Schema, Dict, Config-loader, DB-opener added in later tasks.
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
