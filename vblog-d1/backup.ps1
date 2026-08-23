<#
.SYNOPSIS
  vBlog D1 — 备份远程 D1 数据库

.DESCRIPTION
  把线上 vblog-d1-db 完整导出为 SQL 文件到 backup/ 目录。
  建议：每次重大改动/删库前先跑一次；也可加入计划任务定期备份。

  用法：
    pwsh ./backup.ps1
#>

$ErrorActionPreference = 'Stop'
$Root   = $PSScriptRoot
$Worker = Join-Path $Root 'worker'
$env:WRANGLER_HOME = Join-Path $Root '.wrangler'

$stamp = Get-Date -Format 'yyyyMMdd-HHmmss'
$outDir = Join-Path $Root 'backup'
New-Item -ItemType Directory -Force -Path $outDir | Out-Null

Write-Host '导出远程数据库…' -ForegroundColor Cyan
Push-Location $Worker
node node_modules\wrangler\bin\wrangler.js d1 export vblog-d1-db --remote --output (Join-Path $outDir "vblog-$stamp.sql")
if ($LASTEXITCODE -ne 0) { Pop-Location; Write-Host '[x] 备份失败' -ForegroundColor Red; exit 1 }
Pop-Location

Write-Host ("已备份: backup\vblog-{0}.sql" -f $stamp) -ForegroundColor Green