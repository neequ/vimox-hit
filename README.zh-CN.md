# Vimox Hit CLI

[English](README.md) | 简体中文

`vimox-hit` 是面向本地终端和 Agent 的 [Vimox Hit](https://hit.vimox.cn) 命令行工具。它以单个 Go 二进制文件分发，无需额外运行时依赖。

## 安装 CLI

Windows PowerShell：

```powershell
irm https://raw.githubusercontent.com/neequ/vimox-hit/main/install.ps1 | iex
```

macOS / Linux：

```sh
curl -fsSL https://raw.githubusercontent.com/neequ/vimox-hit/main/install.sh | sh
```

安装脚本会识别系统和架构，从 GitHub Release 下载最新版本、校验 SHA-256，并安装二进制文件。版本列表见 [Releases](https://github.com/neequ/vimox-hit/releases)。

## 安装 Agent Skill

```sh
npx skills add neequ/vimox-hit --skill vimox-hit -g
```

Skill 位于 [`skills/vimox-hit`](skills/vimox-hit)。默认说明使用英文；完整简体中文译文位于 [`references/zh-CN`](skills/vimox-hit/references/zh-CN)。

## 命令

```text
vimox-hit version
vimox-hit doctor
vimox-hit context show --project PROJECT_ID
vimox-hit project init --name NAME
vimox-hit ai-log list [--status FAILED]
vimox-hit ai-log show LOG_ID
vimox-hit ai-log diagnose LOG_ID
vimox-hit ai-log copy LOG_ID
```

命令输出带版本号的 JSON envelope。认证使用 `VIMOX_HIT_TOKEN` 环境变量；`VIMOX_HIT_API_URL` 仅供本地开发。不要把 token 写入命令参数、Shell 历史、项目清单、Skill 文件或诊断包。

## 审批要求

每次模型调用和最终视频生成都必须先向用户展示并取得明确批准。批准为短时单次授权，并绑定准确的输入、模型、参数和费用估算；重试或这些信息发生变化时必须重新申请。目前服务端尚未提供审批接口，因此 CLI 不开放模型执行和视频生成命令。

## 开发

```sh
go test ./...
go build ./cmd/vimox-hit
```

推送 `vX.Y.Z` 标签会触发 GitHub Actions，构建并发布 Windows、macOS、Linux 的 amd64 与 arm64 版本，同时附带 SHA-256 校验文件。
