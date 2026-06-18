# 只读安全约定

- 工具仅允许单条 `SELECT` / `WITH ... SELECT`。
- 任何 INSERT/UPDATE/DELETE/DROP/ALTER/TRUNCATE/CREATE/EXEC/MERGE/INTO 会被拦截并报错。
- 禁止分号拼接多语句。
- `query` 默认 `--limit 1000` 防止误拉大表;`--limit 0` 解除(谨慎,QMSTB 130 万、INVLA 99 万行)。
- `--translate` 把列头编号翻译成中文(best-effort,跨表同名编号取首个)。
