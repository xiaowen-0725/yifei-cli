# yifei-cli

> 易飞 ERP 只读数据分析 CLI（鼎捷 Digiwin）— 让人类与 AI Agent 在终端对易飞 ERP 数据库做**只读**查询、结构探索与字段翻译。

[![npm version](https://img.shields.io/npm/v/yifei-cli.svg)](https://www.npmjs.com/package/yifei-cli)

## 安装

```bash
npm install -g yifei-cli
```

安装时会自动从 GitHub Release 下载对应平台的预编译二进制（macOS / Windows / Linux × amd64 / arm64，纯 Go、零运行时依赖）。

也可直接到 [Releases](https://github.com/xiaowen-0725/yifei-cli/releases) 下载单文件二进制。

## 快速开始

```bash
# 1. 配置数据库连接（写入 config.yaml，密码本地存储）
yifei config init --host localhost --port 1433 --user sa --password '***' --database YDSTEST

# 2. 探索表结构（离线，读内嵌 schema）
yifei schema search 客户
yifei schema table COPTC

# 3. 查字段含义（离线，读内嵌字典）
yifei dict table INVMB
yifei dict code 单别 COPTC

# 4. 执行只读查询（Agent 建议加 --format json）
yifei query "SELECT TOP 10 MB001, MB002 FROM INVMB" --format json
yifei query "SELECT * FROM COPTC WHERE TC003 BETWEEN '20220101' AND '20221231'" --translate

# 5. 预置业务分析
yifei analyze list
yifei analyze order-count-by-customer --year 2022
```

## 命令

| 命令 | 说明 |
|------|------|
| `config init/show` | 管理数据库连接配置 |
| `schema tables/search/table/modules` | 离线探索表结构（内嵌 schema） |
| `dict table/field/code/relations` | 离线字段含义翻译（内嵌字典） |
| `query` | 执行**只读** SQL（`--format table\|json\|csv`、`--translate`、`--limit`） |
| `analyze` | 预置业务分析模板 |
| `version` | 显示版本 |

## 只读安全

本工具仅允许单条 `SELECT` / `WITH` 查询，拦截一切 `INSERT/UPDATE/DELETE/DROP/...` 写操作与分号多语句；`query` 默认 `--limit 1000` 防止误拉大表。数据始终只读。

## AI Agent

仓库内附带 [`skills/yifei`](https://github.com/xiaowen-0725/yifei-cli/tree/main/skills/yifei) 一套 Agent Skill（含易飞 + SQL Server 查询避坑指南），Claude / Codex 等 Agent 可据此直接驱动本 CLI 做数据分析。

## License

MIT
