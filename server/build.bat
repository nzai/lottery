@echo off

go mod tidy

SET CGO_ENABLED=0
SET GOOS=linux
SET GOARCH=amd64

go build -o lottery

SET GOOS=windows
SET GOARCH=amd64