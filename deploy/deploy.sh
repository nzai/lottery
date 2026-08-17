#!/usr/bin/env bash
# 本地构建并部署到服务器。用法: ./deploy/deploy.sh la（ssh 别名或 user@host）
set -euo pipefail

SERVER="${1:?用法: deploy.sh la 或 deploy.sh user@host}"
ROOT="$(cd "$(dirname "$0")/.." && pwd)"
TARGET=/opt/lottery

echo "==> 构建前端（产物嵌入二进制；touch .gitkeep 保持 embed 源目录在 git 中）"
(cd "$ROOT/web" && npm ci && npm run build)
touch "$ROOT/server/static/.gitkeep"

echo "==> 交叉编译 Go 服务端（linux/amd64，前端产物已嵌入，单文件部署）"
mkdir -p "$ROOT/dist"
BUILD_TIME="$(date -u '+%Y-%m-%dT%H:%M:%SZ')"
echo "    编译时间: $BUILD_TIME（前端页面底部显示，用于核对更新是否生效）"
(cd "$ROOT/server" && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
  -ldflags "-X main.buildTime=$BUILD_TIME" \
  -o "$ROOT/dist/lottery-server" .)

echo "==> 初始化服务器：创建独立服务用户 lottery（幂等）"
ssh "$SERVER" "
  id -u lottery >/dev/null 2>&1 || sudo useradd -r -s /usr/sbin/nologin lottery
  sudo mkdir -p $TARGET
  sudo chown lottery:lottery $TARGET
"

echo "==> 上传到 $SERVER（先传 /tmp，再 sudo 移入，避免目录权限冲突）"
scp "$ROOT/dist/lottery-server" "$SERVER:/tmp/lottery-server"
scp "$ROOT/deploy/lottery.service" "$SERVER:/tmp/lottery.service"

echo "==> 安装 systemd 单元（lottery 用户运行）并重启"
ssh "$SERVER" "
  sudo mv /tmp/lottery-server $TARGET/lottery-server &&
  sudo chmod +x $TARGET/lottery-server &&
  sudo chown lottery:lottery $TARGET/lottery-server &&
  sudo cp /tmp/lottery.service /etc/systemd/system/lottery.service &&
  sudo systemctl daemon-reload &&
  sudo systemctl enable lottery &&
  sudo systemctl restart lottery &&
  sleep 2 &&
  sudo systemctl status lottery --no-pager | head -10
"
