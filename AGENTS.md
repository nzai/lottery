# AGENTS.md — 项目上下文（给 AI 助手）

双色球走势图网站：移动端优先（给老人用），展示双色球开奖走势图，支持复选框滑动刷选若干期后实时统计（出现次数/冷热/遗漏/比例）。

## 架构

- **服务端**：Go + Gin + SQLite（`server/`），单二进制部署。前端构建产物（`web/` 构建输出到 `server/static/`）通过 `go:embed` 嵌入二进制（`server/main.go`），所以**部署只有一个可执行文件**。
- **前端**：React 19 + Vite + TS（`web/`），无重型 UI 库，走势图与刷选交互自研。
- **域名**：https://lottery.nzai.me → Nginx → 127.0.0.1:23817 → Go 服务。
- **端口**：23817（用户指定的非默认端口）；SOCKS 隧道 11000（用户指定）。
- **服务器**：境外服务器 `la`（peter 用户，免密 sudo）；境内服务器 `aryana.4apogee.com`（ubuntu 用户）仅作流量转发。

## 关键决策与"为什么"（改代码前必读）

### 数据源：官网对境外 IP 封禁，必须走境内隧道
- 福彩官网 `cwl.gov.cn` 对境外 IP 直接 **403**（IP 级 WAF 封禁，改 UA/Referer 无效）。
- **GitHub Actions 出口也被 WAF 拦截**（返回 200 + HTML 拦截页）——workflow 方案已废弃删除，不要重新尝试。
- 中彩网 `zhcw.com` 页面可达但数据接口在 `jc.zhcw.com`（UCloud WAF 对境外 IP 发 RST）；`chartstatic` 静态数据停更于 2020；新浪/澳客/500 接口废弃或反爬；天行/聚合无彩票接口。
- **现行方案**：`deploy/ssq-tunnel.service` 在 la 上维护一条 SSH SOCKS 隧道（`ssh -D 127.0.0.1:11000 -N ubuntu@aryana.4apogee.com`），lottery 服务配置 `LOTTERY_FETCH_PROXY=socks5://127.0.0.1:11000` 经境内出口抓取官网。隧道监听在 **la 本机**（境内服务器不监听任何端口）。
- 本地开发（国内网络）直连官网即可（代理配置留空）。
- 每日 21:30（Asia/Shanghai）定时增量同步（`LOTTERY_SYNC_CRON`），幂等 upsert（期号唯一键）。

### 静态资源缓存策略
- `/assets/*`：`Cache-Control: public, max-age=2592000, immutable`（30 天）——文件名带内容 hash，修复 bug 重新构建后 hash 变化，浏览器自动拿新文件。
- `/`（index.html）：`no-cache` 每次校验——保证入口永远最新，修复立即生效。
- 页面底部显示编译时间版本号（`/api/version`），用于核对更新是否生效。

### 版本号注入的坑
- `-ldflags "-X main.buildTime=..."` 的符号路径**必须是 `main.buildTime`**（Go 链接器对 main 包符号固定用 `main.` 前缀）；写完整包路径（`github.com/nzai/lottery/server.buildTime`）会**静默失败**（不报错，值保持 "dev"）。
- 验证注入是否生效必须实际运行二进制（`/api/version`），grep 二进制里的字符串不可靠（`-ldflags` 参数文本本身也会出现在构建信息里）。

### 静态文件服务的坑
- 不能用 gin 的 `StaticFS("/")`：它的 `/*filepath` 通配路由与 `/api/*` 冲突。
- 不能用 `http.FileServer` 处理 index.html：标准库会把**任何以 `/index.html` 结尾的路径 301 重定向到 `./`**（目录规范化逻辑）——SPA 回退必须用 `fs.ReadFile` 直接返回（`serveEmbedded`）。
- 空库时 `/api/draws` 必须返回 `[]` 而非 null（store.List 初始化为空切片 + 前端 `?? []` 兜底），否则前端 `draws.map` 白屏。

### 走势图交互设计（移动端优先，老人用）
- **时间正序**：最新一期在底部（用户习惯），初始定位到底部，**向上滚动到顶部加载更早数据**（`before` 分页）。数据层保持 DESC 不变（分页/统计不依赖显示顺序），TrendChart 渲染时 reverse；顶部追加时补偿 scrollTop 防止视口跳动。
- **复选框滑动刷选**：左侧固定列（复选框+期号+日期，sticky left）内手指滑动 = 刷选（`touch-action: none` + pointer capture + `elementFromPoint` 命中行）；数字区域滑动 = 正常滚动。空间分流互不干扰。
- **统计实时更新**：刷选过程中只更新行高亮，300ms debounce 后计算统计面板（计算 <1ms，性能无瓶颈，节奏控制防闪烁）。
- **固定列背景必须不透明**：sticky 列用 `background: inherit`，行背景透明时横向滚动会让数字内容从固定列下方透出（视觉重叠）——所有行显式 `background: #fff`。
- **曲面屏安全边距**：`.chart-scroll` 加 `border-left: 16px`（页面同色）让滚动内容从 16px 处开始，sticky 列自动停在安全区内；工具栏/统计面板 padding 16-18px。不要用 `maximum-scale`/`user-scalable=no`（老人需要双指缩放）。
- 字体基准 16px（老人视力），A−/A+ 调 0.8×~1.5× 记忆在 localStorage。
- 统计面板：次数显示在进度条**左侧**（手机上优先可见），清除按钮紧跟"已选 N 期"（不两端对齐）。

### 部署与发布流程
- **不自动发布**：代码改动只提交（用户明确说"发布"才执行 `./deploy/deploy.sh la`）。
- `deploy/deploy.sh la`：构建前端 → 交叉编译 linux 二进制（注入编译时间）→ scp 到 /tmp → sudo mv 到 `/opt/lottery`（目录属 lottery 用户）→ systemd 重启。
- **独立服务用户**：`lottery`（`useradd -r -s /usr/sbin/nologin`），systemd unit 用 `User=lottery`，与服务器上其他服务隔离。
- 部署脚本用 scp 而非 rsync（Windows 本地无 rsync）；二进制上传后必须 `chmod +x`（scp 保留 644 会导致 203/EXEC）。
- 服务器防火墙只开 80/443（11000/23817 仅回环监听，无需开放）。

### 开发环境
- Windows：`.\dev.ps1`（主推；PowerShell 5.1 对无 BOM UTF-8 会乱码，脚本用纯 ASCII）；Git Bash：`./dev.sh`。
- dev 脚本直接用 `node_modules/.bin/vite`（不用 npm run dev，exec 后 PID 不变便于清理）+ 子进程后台 + `wait`（bash 等待前台命令时信号不触发 trap）。
- npm 12 默认 `allow-remote=none`，`web/.npmrc` 里有 `allow-remote=all`（不要删）。
- `server/static/.gitkeep` 必须保留在 git 里（embed 目录占位；vite 构建会删它，构建后 `touch` 恢复）。

## 常用命令

```bash
# 测试
cd server && go test ./...        # 或仓库根: go test ./server/...
cd web && npx vitest run

# 本地开发
.\dev.ps1                          # Windows；访问 http://localhost:5173

# 部署（用户确认后）
./deploy/deploy.sh la

# 服务器运维
ssh la "sudo journalctl -u lottery --no-pager"   # 服务日志
ssh la "sudo systemctl status ssq-tunnel lottery"
ssh la "sudo -u lottery env LOTTERY_FETCH_PROXY=socks5://127.0.0.1:11000 LOTTERY_DB=/opt/lottery/lottery.db /opt/lottery/lottery-server -sync"  # 手动同步
```

## 服务器拓扑

```
浏览器(手机) → https://lottery.nzai.me → Nginx(:80/443, certbot 证书) → 127.0.0.1:23817
lottery-server(:23817, lottery 用户) ──定时 21:30 抓取──→ socks5://127.0.0.1:11000
ssq-tunnel.service (SSH -D 11000) ──隧道──→ 境内 aryana.4apogee.com (ubuntu, 仅转发)
                                                          ↓
                                                     cwl.gov.cn (官网, 仅境内可达)
```

## 数据流

`/api/draws?limit=N&before=期号`（DESC 分页）→ 前端渲染时 reverse 成正序（最新在底）→ 刷选 → 300ms debounce → 纯函数统计（`web/src/lib/stats.ts`，全部前端计算）。
