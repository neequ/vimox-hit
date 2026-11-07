# Install

Windows PowerShell:

```powershell
irm https://raw.githubusercontent.com/neequ/vimox-hit/main/install.ps1 | iex
```

macOS / Linux:

```sh
curl -fsSL https://raw.githubusercontent.com/neequ/vimox-hit/main/install.sh | sh
```

The installer downloads the matching build from GitHub Releases and verifies it against `SHA256SUMS`. Install this Skill from `skills/vimox-hit` with the Agent's supported GitHub Skill installer. The Simplified Chinese guide is at [zh-CN/install.md](zh-CN/install.md).

When the Skill runs, it checks for `vimox-hit` and installs the CLI if needed. It then runs `vimox-hit doctor`. Project and AI log commands require `VIMOX_HIT_TOKEN`; Hit does not yet provide a CLI OAuth authorization flow, so the Skill stops when the doctor reports `authenticated: false`. Do not copy access tokens from browser storage.
