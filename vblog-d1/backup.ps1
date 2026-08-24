<#
.SYNOPSIS
  vBlog D1 — 备份远程 D1 数据库 + KV 图床图片

.DESCRIPTION
  把线上 vblog-d1-db 完整导出为 SQL 文件到 backup/ 目录；
  同时把 KV 图床（IMG 绑定）里的全部图片下载到 backup/images-<stamp>/。
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
if ($LASTEXITCODE -ne 0) { Pop-Location; Write-Host '[x] D1 备份失败' -ForegroundColor Red; exit 1 }
Pop-Location

# ── KV 图床备份（IMG 绑定；无 R2 时图片全在这里，漏了就没法完整恢复）──
Write-Host '导出 KV 图床图片…' -ForegroundColor Cyan
$imgDir = Join-Path $outDir "images-$stamp"
New-Item -ItemType Directory -Force -Path $imgDir | Out-Null
Push-Location $Worker
$keys = node node_modules\wrangler\bin\wrangler.js kv key list --binding IMG --remote | ConvertFrom-Json
foreach ($k in $keys) {
  $safeName = $k.name -replace '/', '__'
  $target = Join-Path $imgDir $safeName
  # 借 cmd 做字节级重定向（PS 的 > 会按文本解码编码，损坏二进制图片且不报错）
  cmd /c "node node_modules\wrangler\bin\wrangler.js kv key get --binding IMG --remote $($k.name) > `"$target`""
  if ($LASTEXITCODE -ne 0) { Write-Host ("[!] 跳过 {0}" -f $k.name) -ForegroundColor Yellow }
}
Pop-Location

Write-Host ("已备份: backup\vblog-{0}.sql" -f $stamp) -ForegroundColor Green
Write-Host ("已备份图床: backup\images-{0}（{1} 个文件）" -f $stamp, @($keys).Count) -ForegroundColor Green
