# vimox-hit

`vimox-hit` is the local, agent-friendly command line client for [hit.vimox.cn](https://hit.vimox.cn). It is distributed as a single Go binary and has no runtime dependencies.

## Install

Windows PowerShell:

```powershell
irm https://raw.githubusercontent.com/neequ/vimox-hit/main/install.ps1 | iex
```

macOS / Linux:

```sh
curl -fsSL https://raw.githubusercontent.com/neequ/vimox-hit/main/install.sh | sh
```

Install the `vimox-hit` Skill from this repository with the Agent's supported Skill installer. Releases include Windows, macOS and Linux binaries; install scripts verify SHA-256 before installing.

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

Commands emit a versioned JSON envelope. Set `VIMOX_HIT_TOKEN` for authenticated API requests and `VIMOX_HIT_API_URL` only for local development. Never put tokens in command arguments, shell history, project manifests, Skill files or diagnostic bundles.

## Development

```sh
go test ./...
go build ./cmd/vimox-hit
```

Tag a release as `vX.Y.Z`; GitHub Actions cross-compiles the supported OS/architecture matrix and publishes checksums.

## Current API coverage

This first public release reads project context and AI request logs, creates local project folders, and produces redacted diagnostic bundles. Upload, analysis and video generation are not exposed until the service supports the required user approval protocol for every model operation and final video generation.
