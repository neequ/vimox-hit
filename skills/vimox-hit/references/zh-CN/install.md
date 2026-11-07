# 安装

Windows PowerShell：

```powershell
irm https://raw.githubusercontent.com/neequ/vimox-hit/main/install.ps1 | iex
```

macOS / Linux：

```sh
curl -fsSL https://raw.githubusercontent.com/neequ/vimox-hit/main/install.sh | sh
```

安装程序会从 GitHub Releases 下载当前平台的构建，并使用 `SHA256SUMS` 校验。使用 Agent 支持的 GitHub Skill 安装工具，从 `skills/vimox-hit` 安装本 Skill。

[English](../install.md)
