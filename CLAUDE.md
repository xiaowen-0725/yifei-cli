# CLAUDE.md — yifei-cli

易飞 ERP（鼎捷 Digiwin）**只读**数据分析 CLI，供人类与 AI Agent 在终端查询 SQL Server 数据库。Go + Cobra，跨平台单文件二进制，通过 npm（postinstall 下载）与 GitHub Release 分发。

仓库: https://github.com/xiaowen-0725/yifei-cli · npm: `yifei-cli` · 参考蓝本: 飞书 lark-cli。

## 铁律
- **数据只读。** 任何走实时 SQL 的路径都必须先过 `internal/safety.ValidateReadOnly`（仅允许单条 `SELECT`/`WITH`，拦截写/DDL/EXEC/分号多语句）。`query` 和 `analyze` 都已接入——新增任何执行 SQL 的命令也必须接入，否则就是破坏核心安全保证。
- **用户可见字符串一律中文。**
- 改了会执行 SQL 的命令后，务必确认它经过 `ValidateReadOnly`（analyze 的回归测试见 `cmd/analyze_test.go` 的 `TestAnalyzeMalicious*`）。

## 架构（cmd 薄，internal 单一职责）
```
cmd/            root(Deps+Querier) · config · query · schema · dict · analyze · version
internal/
  config/   config.yaml 读写 + DSN（密码明文，0600；YIFEI_PASSWORD 可覆盖）
  safety/   只读 SQL 校验（安全红线，正则在 init 编译一次）
  schema/   解析内嵌 schema.json：表/列/模块/搜索
  dict/     解析内嵌 dict.yaml：字段编号→中文、码表、关系
            （方法名是 LookupCodes/LookupRelations，不是 Codes/Relations——后者与字段同名会编译失败）
  output/   table/json/csv 渲染 + 列头翻译
  db/       go-mssqldb 连接；Query(sql, limit) 在扫描时限行（不改写 SQL）
  analyze/  预置 SQL 模板（只用字典已确认的列，禁止臆测金额/价格字段）
  build/    Version/Date（ldflags 注入）
  assets/   //go:embed schema.json + dict.yaml（embed 不能引父目录，故数据放这里）
skills/yifei/   Agent Skill（SKILL.md + references，含 sql-guide 避坑指南）
```
依赖注入：`main.go` 用内嵌资产构造 `cmd.Deps{Schema, Dict, OpenDB}` 注入命令；DB 经 `Querier` 接口，测试用 `fakeDB`（定义在 `cmd/query_test.go`，analyze 测试复用，勿重复定义）。

## 易飞数据约定（写 SQL/模板必读，详见 skills/yifei/references/sql-guide.md）
- T-SQL：限量用 `SELECT TOP n`，**不是 LIMIT**。
- **日期是 `char(8)` 字符串**（`'20220301'`），区间过滤直接字符串比较。
- **无外键**，头身/跨表靠 `XX001(单别)+XX002(单号)` 手动 JOIN；字段编号不连续。
- 金额/数量可能是字符串，聚合前 `CAST(... AS decimal)`。

## 命令速查
`config init/show` · `schema tables/search/table/modules` · `dict table/field/code/relations` · `query "<SQL>" [--format json|csv] [--translate] [--limit N]` · `analyze list|<模板> [--year|--top]` · `version`

## 构建 / 测试 / 交叉编译
```bash
go test ./...                                   # 全部测试（37）
go vet ./...
make build                                      # 本地二进制 -> bin/yifei
GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build -o /tmp/y.exe .   # 交叉编译示例
go test -tags integration ./internal/db/        # 需本地 YDSTEST 容器；默认 build 跳过
```
- **Go floor = 1.25**（go.mod 的 `go` 指令）。注意：`go-mssqldb` v1.10.0 及其 `golang.org/x/*` 依赖树强制 go 1.25；别为了“降低门槛”把它改回 1.23——会触发 `updates to go.mod needed` 而无法构建。
- CGO 始终关闭（纯 Go，单文件零依赖）。

## 分发：npm + GitHub Release（二者强耦合）
- npm 包只是个 **3.9kB 的下载器**：`scripts/postinstall.js` 按平台从 GitHub Release `v<package.json.version>` 拉对应**裸二进制**到 `binaries/`，`bin/yifei.js` 再 exec 它。
- **资产命名必须固定**：`yifei-<os>-<arch>[.exe]`（os∈darwin/windows/linux，arch∈amd64/arm64）。改名就会断掉已发布版本的安装链路——`postinstall.js`、`release.yml`、Release 资产三处命名必须一致。
- npm 默认 registry 可能指向 `registry.npmmirror.com`（不收 npmjs token）；发布须显式 `--registry https://registry.npmjs.org/`（CI 用 setup-node 的 registry-url 已处理）。

## 发版（GitHub Actions 自动化）
```bash
npm version patch          # 改 package.json 版本 + 打 tag vX.Y.Z
git push --follow-tags     # 推 tag -> .github/workflows/release.yml 自动：
                           #   交叉编译6平台 -> 建 Release 上传二进制 -> npm publish
```
- 前置：仓库需配 `NPM_TOKEN` secret（`gh secret set NPM_TOKEN`）。
- `ci.yml`：push/PR 到 main 跑 vet + test + 交叉编译检查。
- 手动发布步骤（如需）与 `release.yml` 里的命令一致。

## 扩展指引
- **加 analyze 模板**：在 `internal/analyze/templates.go` 的 `templates` 切片加一项即可——`cmd/analyze.go` 自动按 `All()` 生成子命令。只用字典已确认的列；给纯函数 SQL builder 写测试。
- **补字典**：`internal/assets/dict.yaml`（首批覆盖核心表，其余字段保留原编号）；schema.json 是只读快照，更新需重新 build（`make sync-data` 从 ../yifei-erp-docs 同步）。
- **经验记忆**:Agent 跨 session 复用的已验证 ERP 经验存 `~/.config/yifei-cli/memory/*.md`（私有，不进 git），由 skill 的 `scripts/recall.mjs` 召回、Agent 按 `skills/yifei/references/experience-format.md` 格式读写。详见 SKILL.md「经验记忆」章节。
- 设计/计划文档在 `../yifei-erp-docs/docs/superpowers/`（spec + 14 任务 plan）。
