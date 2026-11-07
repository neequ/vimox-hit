# 安装

Windows PowerShell:

```powershell
irm https://raw.githubusercontent.com/neequ/vimox-hit/main/install.ps1 | iex
```

macOS / Linux:

```sh
curl -fsSL https://raw.githubusercontent.com/neequ/vimox-hit/main/install.sh | sh
```

安装脚本从 GitHub Releases 下载当前平台构建并校验 `SHA256SUMS`。安装 Skill 时使用 Agent 支持的 GitHub Skill 安装器，源目录是仓库中的 `skills/vimox-hit`。
