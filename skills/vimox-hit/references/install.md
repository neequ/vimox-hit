# Install

Windows PowerShell:

```powershell
irm https://raw.githubusercontent.com/neequ/vimox-hit/master/install.ps1 | iex
```

macOS / Linux:

```sh
curl -fsSL https://raw.githubusercontent.com/neequ/vimox-hit/master/install.sh | sh
```

The installer downloads the matching build from GitHub Releases and verifies it against `SHA256SUMS`. Install this Skill from `skills/vimox-hit` with the Agent's supported GitHub Skill installer. The Simplified Chinese guide is at [zh-CN/install.md](zh-CN/install.md).

Check `https://api.country.is/` before choosing a source. For country code `CN`, prefer the verified `gh-proxy.com` proxy; otherwise use the official commands. If the lookup or proxy fails, fall back to GitHub directly. The installer verifies `SHA256SUMS`; the lookup provider receives the public IP, and the installer does not store it.

```powershell
$env:VIMOX_HIT_GITHUB_MIRROR='gh-proxy.com'; try { irm 'https://hit.vimox.cn/downloads/cli/install.ps1' | iex } finally { Remove-Item Env:VIMOX_HIT_GITHUB_MIRROR -ErrorAction SilentlyContinue }
```

```sh
curl -fsSL 'https://hit.vimox.cn/downloads/cli/install.sh' | VIMOX_HIT_GITHUB_MIRROR=gh-proxy.com sh
```

When the Skill runs, it checks for `vimox-hit` and installs the CLI if needed. It then runs `vimox-hit doctor`. Project and AI log commands require `VIMOX_HIT_TOKEN`; Hit does not yet provide a CLI OAuth authorization flow, so the Skill stops when the doctor reports `authenticated: false`. Do not copy access tokens from browser storage.
