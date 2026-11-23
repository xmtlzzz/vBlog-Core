<#
.SYNOPSIS
  备份远程 D1 数据库、历史 KV 图片和 R2 图片；不修改远端数据。
.DESCRIPTION
  pwsh ./backup.ps1
  pwsh ./backup.ps1 -SkipDatabase  # 仅备份图片
  图片目录内的 manifest.json 记录存储来源、原始 key、文件名与元数据。
  使用当前 Wrangler 登录认证；凭据仅保存在进程内存中。
#>
param([switch]$SkipDatabase)
$ErrorActionPreference = 'Stop'
$Root = $PSScriptRoot
$Worker = Join-Path $Root 'worker'
$cli = Join-Path $Worker 'node_modules/wrangler/bin/wrangler.js'
$configPath = Join-Path $Worker 'wrangler.toml'
$config = Get-Content -LiteralPath $configPath -Raw
$accountId = [regex]::Match($config, '(?m)^account_id\s*=\s*"([^"]+)"').Groups[1].Value
if (-not $accountId) { throw '请先在 worker/wrangler.toml 中配置 account_id' }
$env:WRANGLER_HOME = Join-Path $Root '.wrangler'

$requestOptions = @{ TimeoutSec = 120 }
$taskProxy = Get-ItemProperty 'HKCU:\Software\Microsoft\Windows\CurrentVersion\Internet Settings' -ErrorAction SilentlyContinue
$proxyUrl = $env:HTTPS_PROXY
if (-not $proxyUrl -and $taskProxy.ProxyEnable -eq 1) { $proxyUrl = $taskProxy.ProxyServer }
if ($proxyUrl) {
  if ($proxyUrl -notmatch '^https?://') { $proxyUrl = 'http://' + $proxyUrl }
  $env:HTTPS_PROXY = $proxyUrl
  $env:HTTP_PROXY = $proxyUrl
  $requestOptions.Proxy = $proxyUrl
}
$stamp = Get-Date -Format 'yyyyMMdd-HHmmss'
$outDir = Join-Path $Root 'backup'
New-Item -ItemType Directory -Force -Path $outDir | Out-Null

if (-not $SkipDatabase) {
  Write-Host '导出远程 D1 数据库…' -ForegroundColor Cyan
  & node $cli d1 export vblog-d1-db --remote --output (Join-Path $outDir "vblog-$stamp.sql") --config $configPath
  if ($LASTEXITCODE -ne 0) { throw 'D1 备份失败' }
}

function Get-BindingBlock([string]$kind, [string]$binding) {
  foreach ($match in [regex]::Matches($config, "(?ms)^\[\[$kind\]\]\s*(.*?)(?=^\[|\z)")) {
    if ($match.Groups[1].Value -match "(?m)^binding\s*=\s*`"$binding`"") { return $match.Groups[1].Value }
  }
  return ''
}
function Get-BackupPage([string]$uri) {
  try { $page = Invoke-RestMethod -Uri $uri -Headers $authHeaders @requestOptions }
  catch { throw 'Cloudflare 图片列表请求失败，备份不完整；请检查网络和存储读取权限' }
  if (-not $page.success) { throw 'Cloudflare 返回失败，备份不完整' }
  return $page
}
function Save-BackupObject([string]$uri, [string]$key, [string]$directory) {
  # Base64url avoids Windows-invalid filenames, path traversal and key collisions.
  $filename = [Convert]::ToBase64String([Text.Encoding]::UTF8.GetBytes($key)).TrimEnd('=').Replace('+', '-').Replace('/', '_')
  $destination = Join-Path $directory $filename
  $partial = $destination + '.partial'
  try { Invoke-WebRequest -Uri $uri -Headers $authHeaders -OutFile $partial @requestOptions | Out-Null }
  catch { throw "图片下载失败：$key；备份不完整" }
  Move-Item -LiteralPath $partial -Destination $destination -Force
  return $filename
}

$auth = $null
try {
  # Capture stdout; never print or persist the authentication token.
  $authJson = & node $cli auth token --json --config $configPath 2>$null
  if ($LASTEXITCODE -ne 0) { throw 'Wrangler 认证失败，请先运行 wrangler login' }
  $auth = ($authJson -join "`n") | ConvertFrom-Json
  if ($auth.type -notin @('oauth', 'api_token') -or -not $auth.token) { throw '当前认证类型无法备份图片，请使用 Wrangler OAuth 或 API Token' }
  $authHeaders = @{ Authorization = 'Bearer ' + $auth.token }
  $apiRoot = "https://api.cloudflare.com/client/v4/accounts/$accountId"

  foreach ($storage in @('kv', 'r2')) {
    $block = if ($storage -eq 'kv') { Get-BindingBlock 'kv_namespaces' 'IMG' } else { Get-BindingBlock 'r2_buckets' 'R2' }
    if (-not $block) { continue }
    $idKey = if ($storage -eq 'kv') { 'id' } else { 'bucket_name' }
    $resource = [regex]::Match($block, ('(?m)^' + $idKey + '\s*=\s*"([^"]+)"')).Groups[1].Value
    if (-not $resource) { throw "缺少 $storage 存储资源名称" }
    $imageDirectory = Join-Path $outDir "images-$storage-$stamp"
    New-Item -ItemType Directory -Force -Path $imageDirectory | Out-Null
    $entries = [Collections.Generic.List[object]]::new()
    $seen = [Collections.Generic.HashSet[string]]::new()
    $cursor = ''
    Write-Host "导出 $storage 图片…" -ForegroundColor Cyan
    do {
      $encodedResource = [Uri]::EscapeDataString($resource)
      $base = if ($storage -eq 'kv') { "$apiRoot/storage/kv/namespaces/$encodedResource" } else { "$apiRoot/r2/buckets/$encodedResource" }
      $listPath = if ($storage -eq 'kv') { 'keys?limit=1000&prefix=uploads%2F' } else { 'objects?per_page=1000&prefix=uploads%2F' }
      $uri = "$base/$listPath"
      if ($cursor) { $uri += '&cursor=' + [Uri]::EscapeDataString($cursor) }
      $page = Get-BackupPage $uri
      foreach ($item in $page.result) {
        $key = if ($storage -eq 'kv') { $item.name } else { $item.key }
        if (-not $key.StartsWith('uploads/')) { throw '存储列表返回了非图片前缀，备份已停止' }
        $encodedKey = [Uri]::EscapeDataString($key)
        $valuePath = if ($storage -eq 'kv') { "values/$encodedKey" } else { "objects/$encodedKey" }
        $filename = Save-BackupObject "$base/$valuePath" $key $imageDirectory
        $entries.Add(@{ key = $key; file = $filename; metadata = $item })
      }
      $cursor = $page.result_info.cursor
      $more = if ($storage -eq 'kv') { [bool]$cursor } else { [bool]$page.result_info.is_truncated }
      if ($more -and (-not $cursor -or -not $seen.Add($cursor))) { throw '图片列表分页异常，备份不完整' }
    } while ($more)
    @{ storage = $storage; resource = $resource; account_id = $accountId; objects = @($entries.ToArray()) } |
      ConvertTo-Json -Depth 12 | Set-Content -LiteralPath (Join-Path $imageDirectory 'manifest.json') -Encoding utf8
    Write-Host ("已备份 $storage 图片：{0}（{1} 个文件）" -f $imageDirectory, $entries.Count) -ForegroundColor Green
  }
} finally {
  $auth = $null
  $authJson = $null
  $authHeaders = $null
}
