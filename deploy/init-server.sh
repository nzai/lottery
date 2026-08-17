#!/usr/bin/env bash
# 服务器一次性初始化：nginx + certbot 证书 + 反代配置（80 端口）。
# 用法：拷到服务器上执行  sudo bash init-server.sh
# 前置条件：lottery.nzai.me 的 DNS A 记录已指向本服务器。
set -euo pipefail

DOMAIN=lottery.nzai.me

echo "==> 安装 nginx / certbot / rsync"
sudo apt update
sudo apt install -y nginx certbot python3-certbot-nginx rsync

echo "==> 写入 80 端口反代配置（443 证书申请后再由 certbot 补全）"
cat > /tmp/lottery-http.conf << 'EOF'
server {
    listen 80;
    server_name lottery.nzai.me;
    location / {
        proxy_pass http://127.0.0.1:23817;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
EOF
sudo cp /tmp/lottery-http.conf /etc/nginx/sites-available/lottery.conf
sudo ln -sf /etc/nginx/sites-available/lottery.conf /etc/nginx/sites-enabled/lottery.conf
sudo nginx -t
sudo systemctl reload nginx

echo ""
echo "==> 下一步：申请证书（交互式，会问邮箱并自动补全 443 + 跳转配置）"
echo "    在服务器上执行：  sudo certbot --nginx -d $DOMAIN --redirect"
echo "    证书申请完成后，用 deploy/nginx.conf 中的完整配置核对即可（证书路径相同）。"
