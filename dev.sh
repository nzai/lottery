#!/usr/bin/env bash
# 本地开发模式：同时启动 Go 服务端与 Vite 前端。
#
#   ./dev.sh
#
# - Go 服务端 :23817（首次启动自动全量回填，数据存 server/lottery.db）
# - Vite 前端 :5173（/api 代理到 23817）
# - Ctrl+C 同时退出前后端（终端信号覆盖整个进程组，脚本负责收尾清理）
set -euo pipefail

ROOT="$(cd "$(dirname "$0")" && pwd)"

cleanup() {
  set +e # 清理过程任何命令失败都不能中断（pkill 无匹配等）
  trap - INT TERM EXIT
  kill $SERVER_PID $VITE_PID 2>/dev/null
  rm -f "$ROOT/.dev-server"
  echo ""
  echo "已停止前后端"
  exit 0
}
trap cleanup INT TERM EXIT

echo "==> 构建并启动 Go 服务端（:23817）"
go build -o "$ROOT/.dev-server" ./server
LOTTERY_DB="$ROOT/server/lottery.db" "$ROOT/.dev-server" &
SERVER_PID=$!

echo "==> 启动 Vite 前端（http://localhost:5173，Ctrl+C 退出）"
cd "$ROOT/web"
# 直接跑 .bin/vite（不用 npm run dev）：vite 经 exec 后 PID 不变，cleanup 能精确终止。
# 放后台 + wait：前台命令会阻塞 trap 执行（bash 等待前台命令时信号不触发 trap），
# wait 则会被信号中断并立即执行 cleanup
"$ROOT/web/node_modules/.bin/vite" &
VITE_PID=$!
wait $VITE_PID
