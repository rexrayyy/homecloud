# Builds the Windows client and packs it into a distributable zip for a
# GitHub Release. The exe is intentionally NOT committed to git (43 MB binary);
# it ships as a release asset instead. See README.md.
#
# Usage:  .\package-release.ps1 -Version 1.0.0
param(
    [Parameter(Mandatory = $true)][string]$Version
)

$ErrorActionPreference = 'Stop'
Set-Location $PSScriptRoot

# 1) Rebuild the client (also refreshes the embedded Linux server payload).
& "$PSScriptRoot\build.ps1"

# 2) Stage the release contents.
$stage = Join-Path $PSScriptRoot "dist\HomeCloud-$Version"
if (Test-Path $stage) { Remove-Item -Recurse -Force $stage }
New-Item -ItemType Directory -Force $stage | Out-Null

Copy-Item "$PSScriptRoot\homecloud.exe" (Join-Path $stage 'HomeCloud.exe')
Copy-Item "$PSScriptRoot\README.md" (Join-Path $stage 'README.md')

$notes = @"
HomeCloud $Version — Windows client

Run HomeCloud.exe. No installation required.

Data is stored per-user at %LOCALAPPDATA%\HomeCloud\client.db
(devices + DPAPI-encrypted passwords). Nothing is written next to the exe.
Override the location with the HOMECLOUD_CLIENT_DB environment variable.
"@
Set-Content -Path (Join-Path $stage 'RELEASE-NOTES.txt') -Value $notes -Encoding UTF8

# 3) Zip it.
$zip = Join-Path $PSScriptRoot "dist\HomeCloud-$Version-win64.zip"
if (Test-Path $zip) { Remove-Item -Force $zip }
Compress-Archive -Path (Join-Path $stage '*') -DestinationPath $zip

Write-Host "Release archive ready: $zip"
Write-Host "Attach it to a GitHub Release (do NOT commit it to git)."
