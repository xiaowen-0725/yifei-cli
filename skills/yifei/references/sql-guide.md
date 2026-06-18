# 易飞 + SQL Server 查询指南(避坑)

## T-SQL 方言
- 限量用 `SELECT TOP n ...`,**不是 LIMIT**;分页用 `OFFSET n ROWS FETCH NEXT m ROWS ONLY`(需 ORDER BY)。
- 空值 `ISNULL(col, 0)`;字符串拼接用 `+`;长度 `LEN()`。

## 易飞数据陷阱(最关键)
- **日期是 char(8) 字符串**(如 '20220301')。区间过滤直接字符串比较:
  `WHERE TC003 BETWEEN '20220101' AND '20221231'`。需转日期: `CONVERT(date, TC003, 112)`。
- **金额/数量可能是字符串或带尾零**。聚合前转型:`SUM(CAST(NULLIF(LTRIM(RTRIM(col)),'') AS decimal(18,4)))`。
- **无外键**。头身/跨表全靠 `XX001(单别)+XX002(单号)` 手动 JOIN。
- **字段编号不连续**(如 INVLA 从 LA001 跳到 LA004),不要假设连续。

## 常用统计片段
- 按月: `GROUP BY SUBSTRING(TA003,1,6)`。
- TOP N: `SELECT TOP 20 ... ORDER BY 计数 DESC`。
- 计数前先想清楚:大表(QMSTB/INVLA/MOCTB)务必带过滤或 TOP。

## 性能/安全
- 大表避免无过滤 `COUNT(*)` / 全表扫描。
- 一律先 `schema table` / `dict table` 确认字段再写 SQL。
