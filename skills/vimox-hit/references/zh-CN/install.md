# 安装

Windows PowerShell：

```powershell
irm https://raw.githubusercontent.com/neequ/vimox-hit/master/install.ps1 | iex
```

macOS / Linux：

```sh
curl -fsSL https://raw.githubusercontent.com/neequ/vimox-hit/master/install.sh | sh
```

安装程序会从 GitHub Releases 下载当前平台的构建，并使用 `SHA256SUMS` 校验。使用 Agent 支持的 GitHub Skill 安装工具，从 `skills/vimox-hit` 安装本 Skill。

选择下载源前查询 `https://api.country.is/` 的公网 IP 国家代码。结果为 `CN` 时优先使用已验证的 `gh-proxy.com` GitHub 代理；其他地区使用官方命令。检测或镜像不可用时回退官方源。安装程序仍校验 `SHA256SUMS`；归属查询服务会收到公网 IP，安装程序不会保存它。

```powershell
$env:VIMOX_HIT_GITHUB_MIRROR='gh-proxy.com'; try { irm 'https://hit.vimox.cn/downloads/cli/install.ps1' | iex } finally { Remove-Item Env:VIMOX_HIT_GITHUB_MIRROR -ErrorAction SilentlyContinue }
```

```sh
curl -fsSL 'https://hit.vimox.cn/downloads/cli/install.sh' | VIMOX_HIT_GITHUB_MIRROR=gh-proxy.com sh
```

Skill 首次运行时会检查 `vimox-hit`，缺失时自动安装 CLI，然后运行 `vimox-hit doctor`。项目和 AI 日志命令需要 `VIMOX_HIT_TOKEN`；Hit 目前尚未提供 CLI OAuth 授权流程，因此 doctor 返回 `authenticated: false` 时 Skill 会停止。不要从浏览器存储复制访问令牌。

[English](../install.md)
