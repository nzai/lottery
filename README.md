# lottery — 双色球走势图

移动端优先的双色球走势图工具：定时从福彩官网同步开奖数据，支持复选框滑动刷选若干期，
实时统计号码出现次数、冷热、遗漏与比例。

## 功能

- 双色球走势图（33 红球 + 16 蓝球），竖屏/横屏/桌面自适应
- 复选框滑动刷选：手指在左侧复选框列划过即勾选，数字区域滑动正常滚动
- 实时统计（300ms 防抖）：出现次数（条形图）、冷热号、遗漏值、奇偶/大小比例
- 双指捏合缩放 + A−/A+ 字号调节（记忆在 localStorage）
- 每日 21:30（北京时间）自动同步最新开奖，首次启动全量回填历史

## 技术栈

- 服务端：Go + Gin + SQLite（modernc.org/sqlite，纯 Go 免 CGO）
- 前端：React 19 + Vite + TypeScript + Vitest

## 开发

**一键启动前后端（推荐）**：

```powershell
# Windows（PowerShell）
.\dev.ps1

# 或 Git Bash
./dev.sh
```

启动后访问 http://localhost:5173（Vite，`/api` 代理到 Go），Go 服务端在 `:23817`，首次启动自动全量回填历史数据，Ctrl+C 同时退出。

**分开启动**：

```bash
# 服务端（首次启动会自动全量回填，需要能访问福彩官网）
cd server && go run .

# 前端（dev 代理 /api → localhost:23817）
cd web && npm install && npm run dev

# 测试
cd server && go test ./...
cd web && npx vitest run
```

## 部署

前端构建产物在 `go build` 时**嵌入进二进制**（go:embed），部署只有一个可执行文件：

```bash
./deploy/deploy.sh user@server-ip
```

**服务器首次初始化**（一次性，在服务器上执行）：

```bash
sudo apt install -y nginx certbot python3-certbot-nginx rsync
# DNS 的 A 记录先指向服务器，然后申请证书
sudo certbot --nginx -d lottery.nzai.me
# 配置 Nginx（deploy/nginx.conf 放到 /etc/nginx/sites-available/lottery.conf 并启用）
sudo ln -sf /etc/nginx/sites-available/lottery.conf /etc/nginx/sites-enabled/
sudo nginx -t && sudo systemctl reload nginx
```

部署后验证：

```bash
curl -s http://127.0.0.1:23817/api/health
curl -s 'http://127.0.0.1:23817/api/draws?limit=5'
# 查看回填与每日同步日志
sudo journalctl -u lottery --no-pager
```

## 配置（环境变量）

| 变量 | 默认 | 说明 |
|---|---|---|
| `LOTTERY_ADDR` | `:23817` | HTTP 监听地址 |
| `LOTTERY_DB` | `lottery.db` | SQLite 文件路径 |
| `LOTTERY_SYNC_CRON` | `30 21 * * *` | 每日同步 cron（Asia/Shanghai） |
| `LOTTERY_UA` | Chrome UA | 抓取请求 User-Agent |
| `LOTTERY_FETCH_DELAY_MS` | `1500` | 分页抓取间隔 |
| `LOTTERY_FETCH_ENABLE` | `true` | 是否启用抓取与定时同步 |

> 前端静态资源不再有 `LOTTERY_STATIC` 配置——构建产物在编译期嵌入二进制。
> 缓存策略：`/assets/*`（文件名带内容 hash）一年 immutable；`index.html` no-cache，
> 保证修复 bug 后用户立即拿到新版本。

## 手动同步

```bash
cd server && go run . -sync   # 触发一次同步后退出
```
