# Vimox Hit CLI

English | [简体中文](README.zh-CN.md)

`vimox-hit` is the local, agent-friendly command line client for [Vimox Hit](https://hit.vimox.cn). It is distributed as a single Go binary and has no runtime dependencies.

## Install the CLI

Windows PowerShell:

```powershell
irm https://raw.githubusercontent.com/neequ/vimox-hit/main/install.ps1 | iex
```

macOS / Linux:

```sh
curl -fsSL https://raw.githubusercontent.com/neequ/vimox-hit/main/install.sh | sh
```

The installers select the current OS and architecture, download the latest GitHub Release, verify its SHA-256 checksum, and install the binary. See [all releases](https://github.com/neequ/vimox-hit/releases).

## Install the Agent Skill

```sh
npx skills add neequ/vimox-hit --skill vimox-hit -g
```

The Skill is available at [`skills/vimox-hit`](skills/vimox-hit). Its default instructions are in English. Simplified Chinese translations are included under [`references/zh-CN`](skills/vimox-hit/references/zh-CN).

## Commands

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

Commands emit a versioned JSON envelope. Set `VIMOX_HIT_TOKEN` for authenticated API requests and `VIMOX_HIT_API_URL` only for local development. Never put tokens in command arguments, shell history, project manifests, Skill files, or diagnostic bundles.

## Approval requirements

Every model operation and final video generation must be shown to the user for explicit approval before execution. Approval is short-lived, single-use, and bound to the exact inputs, model, parameters, and cost estimate. A retry or any change to those values requires a new approval. The current CLI does not expose model execution or video generation commands because the service-side approval endpoint is not available yet.

## Development

```sh
go test ./...
go build ./cmd/vimox-hit
```

Push a `vX.Y.Z` tag to build and publish Windows, macOS, and Linux releases for amd64 and arm64. GitHub Actions attaches SHA-256 checksums to each release.
