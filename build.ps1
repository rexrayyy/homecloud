$ErrorActionPreference = 'Stop'
Set-Location $PSScriptRoot
$oldOS = $env:GOOS; $oldArch = $env:GOARCH; $oldCGO = $env:CGO_ENABLED
try {
    New-Item -ItemType Directory -Force "$PSScriptRoot\internal\client\payload" | Out-Null
    $env:GOOS = 'linux'; $env:GOARCH = 'amd64'; $env:CGO_ENABLED = '0'
    go build -trimpath -o "$PSScriptRoot\internal\client\payload\homecloud-linux-amd64" ./cmd
    if ($LASTEXITCODE -ne 0) { throw 'Linux build failed' }
    $env:GOOS = 'windows'; $env:GOARCH = 'amd64'
    go build -trimpath -tags production -ldflags '-H=windowsgui' -o "$PSScriptRoot\homecloud.exe" ./cmd
    if ($LASTEXITCODE -ne 0) { throw 'Windows build failed. Close HomeCloud before rebuilding.' }
} finally {
    $env:GOOS = $oldOS; $env:GOARCH = $oldArch; $env:CGO_ENABLED = $oldCGO
}