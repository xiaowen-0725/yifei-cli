---
name: yifei
description: 易飞 ERP 只读数据分析。当用户要查易飞/YDSTEST 数据库的品号/客户/供应商/销售订单/采购订单/工单/库存/BOM,或要做销售/库存/生产分析时使用。触发词:易飞、ERP 数据、查品号、查订单、查工单、库存异动、yifei、销售分析。
---

# 易飞 ERP 数据分析 (yifei-cli)

只读分析易飞 ERP(YDSTEST,SQL Server)。所有数据只读,禁止写操作。

## 取数前:凭据检查
先确认配置存在: `yifei config show`。若报「配置文件不存在」,引导用户运行 `yifei config init --host <h> --port 1433 --user sa --password <pw> --database YDSTEST`。

## 推荐流程(先懂结构,再查数)
1. 不确定查哪张表 → `yifei schema search <关键词>` 或读 references/table-map.md
2. 不懂字段含义 → `yifei dict table <表名>` / `yifei dict field <TABLE.CODE>`
3. 看头身/外联关系 → `yifei dict relations <表名>`
4. 写并执行查询 → `yifei query "<SQL>" --format json`
   - Agent 一律加 `--format json`
   - 列头要中文加 `--translate`
   - 默认最多返回 1000 行,需要更多用 `--limit`

## 关键约定
- 仅 SELECT/WITH,工具会拦截一切写操作(见 references/safety.md)
- T-SQL 用 `SELECT TOP n`,不是 LIMIT;日期是 char(8) 字符串(见 references/sql-guide.md)
- 成品分析: `yifei analyze list` 查看可用模板

## 命令速查
| 目的 | 命令 |
|------|------|
| 搜表 | `yifei schema search 客户` |
| 看表结构 | `yifei schema table COPTC` |
| 字段翻译 | `yifei dict table COPTC` |
| 码表 | `yifei dict code 单别 COPTC` |
| 查询 | `yifei query "SELECT TOP 10 ..." --format json` |
| 分析 | `yifei analyze order-count-by-customer --year 2022` |
