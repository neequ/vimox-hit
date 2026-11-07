# Vimox Hit CLI

English | [简体中文](README.zh-CN.md)

`vimox-hit` is the local, agent-friendly command line client for [Vimox Hit](https://hit.vimox.cn). It is distributed as a single Go binary and has no runtime dependencies.

## Install the CLI

Windows PowerShell:

```powershell
irm https://raw.githubusercontent.com/neequ/vimox-hit/master/install.ps1 | iex
```

macOS / Linux:

```sh
curl -fsSL https://raw.githubusercontent.com/neequ/vimox-hit/master/install.sh | sh
```

The installers select the current OS and architecture, download the latest GitHub Release, verify its SHA-256 checksum, and install the binary. See [all releases](https://github.com/neequ/vimox-hit/releases).

For users in mainland China, check the public IP country code at `https://api.country.is/`. When it returns `CN`, prefer the `gh-proxy.com` GitHub proxy below; elsewhere, use the official commands above. The installer routes GitHub API and release downloads through the proxy and still verifies `SHA256SUMS`. If detection or the mirror fails, use the official source. The lookup provider receives the public IP; the installer does not store it.

```powershell
$env:VIMOX_HIT_GITHUB_MIRROR='gh-proxy.com'; try { irm 'https://hit.vimox.cn/downloads/cli/install.ps1' | iex } finally { Remove-Item Env:VIMOX_HIT_GITHUB_MIRROR -ErrorAction SilentlyContinue }
```

```sh
curl -fsSL 'https://hit.vimox.cn/downloads/cli/install.sh' | VIMOX_HIT_GITHUB_MIRROR=gh-proxy.com sh
```

## Install the Agent Skill

```sh
npx skills add neequ/vimox-hit --skill vimox-hit -g
```

The Skill is available at [`skills/vimox-hit`](skills/vimox-hit). Its default instructions are in English. Simplified Chinese translations are included under [`references/zh-CN`](skills/vimox-hit/references/zh-CN).

Check whether the CLI is available before use. Remote MCP and the web app are available without installing the CLI. Account commands require an authorized `VIMOX_HIT_TOKEN`; interactive CLI OAuth login is not available yet. Without authorization, use the web app or an already connected MCP client. Do not copy credentials from browser storage.

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
vimox-hit template list --kind video --origin system [--cursor CURSOR]
vimox-hit template upload --kind video --file video-template.json [--media source.mp4] [--private]
vimox-hit template upload --kind suite --file suite-template.json [--private]
vimox-hit template publish --kind video --id TEMPLATE_ID
vimox-hit template stats --kind video --id TEMPLATE_ID
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

Template uploads validate and publish automatically by default; `--private` saves only a personal template. These commands do not invoke models. The template subcommands are available in source, but the current GitHub binary release has not been updated yet; use remote MCP until a binary with template support is published. See [example JSON files](examples/) and [template automation](https://github.com/neequ/vimox-hit/blob/master/docs/template-automation.md) for media indexes, statistics and MCP/REST endpoints.
