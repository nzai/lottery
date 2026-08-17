# Local dev mode: start Go server + Vite frontend together.
#
#   .\dev.ps1           (run in PowerShell; if execution policy blocks it:
#                         Set-ExecutionPolicy -Scope Process Bypass)
#
# - Go server :23817 (first start auto-backfills full history to server/lottery.db)
# - Vite frontend :5173 (/api proxied to 23817)
# - Ctrl+C stops both (cleanup via finally + PowerShell.Exiting)
$ErrorActionPreference = 'Stop'
$Root = $PSScriptRoot

# ---------- 1. Build and start Go server ----------
Write-Host '==> Building Go server (:23817)'
Push-Location $Root
go build -o "$Root\.dev-server.exe" ./server
if ($LASTEXITCODE -ne 0) { Pop-Location; Write-Host 'Go build failed'; exit 1 }
$env:LOTTERY_DB = "$Root\server\lottery.db"
$server = Start-Process -FilePath "$Root\.dev-server.exe" -PassThru -NoNewWindow
Pop-Location

# ---------- Fallback cleanup: kill Go process when PowerShell exits ----------
$serverId = $server.Id
Register-EngineEvent -SourceIdentifier PowerShell.Exiting -SupportEvent -Action {
    Stop-Process -Id $serverId -Force -ErrorAction SilentlyContinue
} | Out-Null

# ---------- 2. Run Vite in foreground (Ctrl+C interrupts it) ----------
try {
    Write-Host '==> Starting Vite (http://localhost:5173, Ctrl+C to exit)'
    Push-Location "$Root\web"
    & node "$Root\web\node_modules\vite\bin\vite.js"
}
finally {
    Pop-Location -ErrorAction SilentlyContinue
    if ($server -and -not $server.HasExited) {
        Stop-Process -Id $server.Id -Force -ErrorAction SilentlyContinue
        Write-Host 'Go server stopped'
    }
    # 等文件句柄释放后再删（Stop-Process 后可能有短窗口占用）
    Start-Sleep -Milliseconds 500
    Remove-Item "$Root\.dev-server.exe" -Force -ErrorAction SilentlyContinue
    Write-Host 'Dev environment stopped'
}
