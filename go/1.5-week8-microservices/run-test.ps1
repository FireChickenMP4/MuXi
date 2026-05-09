param(
    [switch]$SkipBuild,
    [switch]$KeepEtcd
)

$ErrorActionPreference = "Stop"
$Root = Split-Path -Parent $MyInvocation.MyCommand.Path

Write-Host "=== 启动 etcd ===" -ForegroundColor Cyan
docker compose -f "$Root\docker-compose.yml" up -d

Write-Host "=== 运行集成测试 ===" -ForegroundColor Cyan
Set-Location $Root
go test -tags=integration -v -timeout 60s .\integration\

if (-not $KeepEtcd) {
    Write-Host "=== 停止 etcd ===" -ForegroundColor Cyan
    docker compose -f "$Root\docker-compose.yml" down
}

Write-Host "=== 完成 ===" -ForegroundColor Green
