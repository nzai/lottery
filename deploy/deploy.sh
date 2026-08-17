#!/usr/bin/env bash
# 本地构建并部署到服务器。用法: ./deploy/deploy.sh user@server-ip
set -euo pipefail

SERVER="${1:?用法: deploy.sh user@server-ip}"
ROOT="$(cd "$(dirname "$0")/.." && pwd)"
TARGET=/opt/lottery

echo "==> 构建前端"
(cd "$ROOT/web" && npm ci && npm run build)

echo "==> 交叉编译 Go 服务端（linux/amd64）"
mkdir -p "$ROOT/dist"
(cd "$ROOT/server" && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o "$ROOT/dist/lottery-server" .)

echo "==> 上传到 $SERVER:$TARGET"
rsync -avz --delete "$ROOT/dist/" "$SERVER:$TARGET/"
rsync -avz "$ROOT/deploy/lottery.service" "$SERVER:/tmp/lottery.service"

echo "==> 安装 systemd 单元并重启"
ssh "$SERVER" "
  sudo cp /tmp/lottery.service /etc/systemd/system/lottery.service &&
  sudo systemctl daemon-reload &&
  sudo systemctl enable lottery &&
  sudo systemctl restart lottery &&
  sleep 2 &&
  sudo systemctl status lottery --no-pager | head -10
"
