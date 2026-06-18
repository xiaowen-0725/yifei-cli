package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/xiaowen-0725/yifei-cli/internal/config"
)

func resolveConfigPath(c *cobra.Command) (string, error) {
	if p := GlobalConfigPath(c); p != "" {
		return p, nil
	}
	return config.DefaultPath()
}

func newConfigCmd() *cobra.Command {
	cc := &cobra.Command{Use: "config", Short: "管理数据库连接配置"}
	cc.AddCommand(newConfigInitCmd(), newConfigShowCmd())
	return cc
}

func newConfigInitCmd() *cobra.Command {
	var host, user, password, database, encrypt string
	var port int
	cmd := &cobra.Command{
		Use:   "init",
		Short: "生成 config.yaml",
		RunE: func(c *cobra.Command, _ []string) error {
			path, err := resolveConfigPath(c)
			if err != nil {
				return err
			}
			cfg := &config.Config{Host: host, Port: port, User: user,
				Password: password, Database: database, Encrypt: encrypt}
			if err := config.Save(path, cfg); err != nil {
				return err
			}
			fmt.Fprintf(c.OutOrStdout(), "已写入配置: %s\n", path)
			return nil
		},
	}
	cmd.Flags().StringVar(&host, "host", "localhost", "数据库主机")
	cmd.Flags().IntVar(&port, "port", 1433, "端口")
	cmd.Flags().StringVar(&user, "user", "sa", "用户名")
	cmd.Flags().StringVar(&password, "password", "", "密码")
	cmd.Flags().StringVar(&database, "database", "YDSTEST", "数据库名")
	cmd.Flags().StringVar(&encrypt, "encrypt", "disable", "TLS: disable|true|false")
	return cmd
}

func newConfigShowCmd() *cobra.Command {
	var reveal bool
	cmd := &cobra.Command{
		Use:   "show",
		Short: "查看当前配置 (密码默认脱敏)",
		RunE: func(c *cobra.Command, _ []string) error {
			path, err := resolveConfigPath(c)
			if err != nil {
				return err
			}
			cfg, err := config.Load(path)
			if err != nil {
				return err
			}
			pw := "****"
			if reveal {
				pw = cfg.Password
			}
			out := c.OutOrStdout()
			fmt.Fprintf(out, "配置文件: %s\n", path)
			fmt.Fprintf(out, "host: %s\nport: %d\nuser: %s\npassword: %s\ndatabase: %s\nencrypt: %s\n",
				cfg.Host, cfg.Port, cfg.User, pw, cfg.Database, cfg.Encrypt)
			return nil
		},
	}
	cmd.Flags().BoolVar(&reveal, "reveal", false, "明文显示密码")
	return cmd
}
