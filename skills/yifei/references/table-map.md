# 核心表速查 — 查什么业务用哪张表

| 业务 | 主档/单头 | 单身/明细 | 关键字段 |
|------|-----------|-----------|----------|
| 品号 | INVMB | — | MB001 品号, MB002 品名 |
| 客户 | COPMA | — | MA001 代号, MA002 简称 |
| 供应商 | PURMA | — | MA001 代号, MA002 简称 |
| 销售订单 | COPTC | COPTD | TC004 客户, TD004 品号 |
| 采购订单 | PURTA | PURTB | TA004 供应商, TB004 品号 |
| 生产工单 | MOCTA | MOCTB | TA006 制造品号, TB003 物料 |
| BOM | BOMCA | BOMCB | CA003 母件, CB005 子件 |
| 库存异动 | INVLA | — | LA001 品号, LA005 方向(1入/-1出) |

模块前缀: COP 销售 / PUR 采购 / MOC 工单 / INV 库存 / BOM 物料清单 / CMS 基础参数。
完整 58 模块见 yifei-erp-docs/database_overview.md。
