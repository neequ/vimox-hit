$ErrorActionPreference = 'Stop'
$repo = 'neequ/vimox-hit'
$arch = if ([System.Runtime.InteropServices.RuntimeInformation]::OSArchitecture -eq 'Arm64') { 'arm64' } else { 'amd64' }
$release = Invoke-RestMethod "https://api.github.com/repos/$repo/releases/latest" -Headers @{ 'User-Agent' = 'vimox-hit-installer' }
$assetName = "vimox-hit_$($release.tag_name)_windows_$arch.zip"
$asset = $release.assets | Where-Object name -eq $assetName | Select-Object -First 1
$checksums = $release.assets | Where-Object name -eq 'SHA256SUMS' | Select-Object -First 1
if (-not $asset -or -not $checksums) { throw "Release $($release.tag_name) lacks $assetName or SHA256SUMS" }
$temp = Join-Path ([IO.Path]::GetTempPath()) ([guid]::NewGuid().ToString('N'))
$install = Join-Path $env:LOCALAPPDATA 'Vimox\bin'
New-Item -ItemType Directory -Force -Path $temp, $install | Out-Null
try {
  Invoke-WebRequest $asset.browser_download_url -OutFile (Join-Path $temp $assetName)
  Invoke-WebRequest $checksums.browser_download_url -OutFile (Join-Path $temp 'SHA256SUMS')
  $expected = (Select-String -Path (Join-Path $temp 'SHA256SUMS') -Pattern "\s+$([regex]::Escape($assetName))$").Line.Split(' ')[0]
  $actual = (Get-FileHash (Join-Path $temp $assetName) -Algorithm SHA256).Hash.ToLowerInvariant()
  if (-not $expected -or $actual -ne $expected) { throw 'CLI checksum verification failed' }
  Expand-Archive (Join-Path $temp $assetName) -DestinationPath $temp -Force
  Move-Item -Force (Join-Path $temp 'vimox-hit.exe') (Join-Path $install 'vimox-hit.exe')
  $userPath = [Environment]::GetEnvironmentVariable('Path', 'User')
  if (($userPath -split ';') -notcontains $install) { [Environment]::SetEnvironmentVariable('Path', (($userPath.TrimEnd(';') + ';' + $install).Trim(';')), 'User') }
  & (Join-Path $install 'vimox-hit.exe') version
  Write-Host "Installed to $install. Open a new terminal to use vimox-hit."
} finally {
  Remove-Item -LiteralPath $temp -Recurse -Force -ErrorAction SilentlyContinue
}
