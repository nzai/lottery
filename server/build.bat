@echo off
rem 本地构建脚本：交叉编译 linux/amd64 二进制（前端产物已嵌入，单文件部署）

go mod tidy

rem 编译时间作为版本号（前端页面底部显示，用于核对更新是否生效）
for /f "delims=" %%i in ('git log -1 --format^=%%cI') do set BUILD_TIME=%%i
echo 编译时间: %BUILD_TIME%

SET CGO_ENABLED=0
SET GOOS=linux
SET GOARCH=amd64

go build -ldflags "-X github.com/nzai/lottery/server.buildTime=%BUILD_TIME%" -o lottery

SET GOOS=windows
SET GOARCH=amd64
