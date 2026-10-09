<#
.SYNOPSIS
  vBlog D1 - Build and Deploy to Cloudflare
#>
param(
  [switch]$SkipBuild
)

$ErrorActionPreference = 'Stop'
$Root   = $PSScriptRoot
$Web    = Join-Path $Root 'web'
$Worker = Join-Path $Root 'worker'
$env:WRANGLER_HOME = Join-Path $Root '.wrangler'

if (-not $env:HTTPS_PROXY) {
  try {
    $reg = Get-ItemProperty 'HKCU:\Software\Microsoft\Windows\CurrentVersion\Internet Settings' -ErrorAction SilentlyContinue
    if ($reg.ProxyEnable -eq 1 -and $reg.ProxyServer) {
      $proxyServer = $reg.ProxyServer
      if ($proxyServer -notmatch '^https?://') { $proxyServer = "http://$proxyServer" }
      $env:HTTPS_PROXY = $proxyServer
      $env:HTTP_PROXY  = $proxyServer
      Write-Host "[*] Proxy detected: $proxyServer" -ForegroundColor DarkGray
    }
  } catch {}
}

if (-not (Test-Path (Join-Path $Worker 'node_modules\wrangler'))) {
  Write-Host '[0/3] Installing wrangler locally...' -ForegroundColor Cyan
  Push-Location $Worker
  npm install --no-audit --no-fund
  Pop-Location
}

if (-not $SkipBuild) {
  Write-Host '[1/3] Building frontend...' -ForegroundColor Cyan
  Push-Location $Web
  npm ci --no-audit --no-fund
  if ($LASTEXITCODE -ne 0) { npm install --no-audit --no-fund }
  npm run build
  if ($LASTEXITCODE -ne 0) { Write-Host '[x] Build failed' -ForegroundColor Red; exit 1 }
  Pop-Location
} else {
  Write-Host '[1/3] Skipping frontend build' -ForegroundColor DarkGray
}

Write-Host '[2/3] Deploying Worker + static assets...' -ForegroundColor Cyan
Push-Location $Worker
node node_modules\wrangler\bin\wrangler.js deploy --config "$Worker\wrangler.toml"
if ($LASTEXITCODE -ne 0) { Write-Host '[x] Deploy failed' -ForegroundColor Red; exit 1 }

$secretsJson = node node_modules\wrangler\bin\wrangler.js secret list --config "$Worker\wrangler.toml"
try {
  $secrets = $secretsJson | ConvertFrom-Json
  $names = @($secrets | ForEach-Object { $_.name })
  foreach ($required in @('JWT_SECRET', 'TURNSTILE_SECRET')) {
    if ($names -notcontains $required) {
      Write-Host "[!] Missing secret: $required. Restore with: node node_modules\wrangler\bin\wrangler.js secret put $required" -ForegroundColor Yellow
    }
  }
} catch {
  Write-Host "[!] Note: Could not parse secrets list." -ForegroundColor DarkGray
}
Pop-Location

Write-Host '[3/3] Deployment completed successfully! [OK]' -ForegroundColor Green