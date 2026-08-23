<#
.SYNOPSIS
  vBlog D1 — 一键构建 + 部署到 Cloudflare

.DESCRIPTION
  构建前端（npm ci + vite build），随后 wrangler deploy 发布 Worker + 静态资源。
  前置：已 wrangler login（见 README）、worker/node_modules 里已装 wrangler
  （首次可先跑：cd worker; npm install）

  用法：
    pwsh ./deploy.ps1            # 构建 + 部署
    pwsh ./deploy.ps1 -SkipBuild # 只部署（前端没改动时更快）
#>
param(
  [switch]$SkipBuild
)

$ErrorActionPreference = 'Stop'
$Root   = $PSScriptRoot
$Web    = Join-Path $Root 'web'
$Worker = Join-Path $Root 'worker'
$env:WRANGLER_HOME = Join-Path $Root '.wrangler'

if (-not (Test-Path (Join-Path $Worker 'node_modules\wrangler'))) {
  Write-Host '[0/3] 安装 wrangler（本地依赖）…' -ForegroundColor Cyan
  Push-Location $Worker
  npm install --no-audit --no-fund
  Pop-Location
}

if (-not $SkipBuild) {
  Write-Host '[1/3] 构建前端…' -ForegroundColor Cyan
  Push-Location $Web
  npm ci --no-audit --no-fund
  if ($LASTEXITCODE -ne 0) { npm install --no-audit --no-fund }
  npm run build
  if ($LASTEXITCODE -ne 0) { Write-Host '[x] 构建失败' -ForegroundColor Red; exit 1 }
  Pop-Location
} else {
  Write-Host '[1/3] 跳过前端构建' -ForegroundColor DarkGray
}

Write-Host '[2/3] 部署 Worker + 静态资源…' -ForegroundColor Cyan
Push-Location $Worker
node node_modules\wrangler\bin\wrangler.js deploy
if ($LASTEXITCODE -ne 0) { Write-Host '[x] 部署失败' -ForegroundColor Red; exit 1 }
Pop-Location

Write-Host '[3/3] 部署完成 ✅' -ForegroundColor Green