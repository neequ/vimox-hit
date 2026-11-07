---
name: vimox-hit
description: Use the Vimox Hit CLI to inspect video projects, diagnose AI requests, and create local project folders. Trigger when the user asks to work with Vimox Hit or troubleshoot its model requests.
---

# Vimox Hit

Use Vimox Hit's remote MCP tools when they are available. This is the zero-install path for ordinary users: query the project with `video_project_state`, open the returned workbench URL, and use `template_upload` for complete template JSON and keep model approvals in the web app. Do not install a local CLI merely because this Skill was invoked.

If the requested workflow needs the CLI and it is missing, offer the official platform installer below. When the user has asked you to set it up and the host grants terminal/network access, run it and verify with `vimox-hit version`; otherwise show the exact command and stop before claiming installation. Never claim that a prompt can register an MCP connector inside Codex; that remains a client-level setting.

## Remote video project workflow

When the remote Vimox Hit MCP is available, use it for cloud projects even if the local CLI is absent or `vimox-hit doctor` reports `authenticated: false`:

1. If the user supplied a project ID, call `video_project_state` for that ID.
2. Otherwise call `video_projects_list` and let the user choose an existing project. If the list is empty, ask for a project name; after the user provides it, call `video_project_create` without model overrides so the account's global reasoning model is used.
3. After creation, immediately call `video_project_state` with the returned ID and `video_project_open`; report the real state and workbench URL.
4. If creation says no reasoning model is configured, direct the user to configure the global reasoning model in Vimox Hit Settings, then retry. Do not substitute local `project init` for cloud project creation.

If MCP is unavailable, explain that client-side connector registration and account authorization are required, link `https://hit.vimox.cn/mcp` and the web workbench, and name the client's exact settings entry point when known. Do not claim that a chat prompt installed or authorized the connector. Source video upload continues in the returned web workbench. Never call DNA analysis or video generation through an undocumented CLI command.

On Windows PowerShell:

```powershell
$installDir = Join-Path $env:LOCALAPPDATA 'Vimox\bin'
if (-not (Get-Command vimox-hit -ErrorAction SilentlyContinue) -and (Test-Path (Join-Path $installDir 'vimox-hit.exe'))) {
  $env:Path = "$installDir;$env:Path"
}
if (-not (Get-Command vimox-hit -ErrorAction SilentlyContinue)) {
  irm https://raw.githubusercontent.com/neequ/vimox-hit/master/install.ps1 | iex
}
$env:Path = "$installDir;$env:Path"
vimox-hit version
```

On macOS or Linux:

```sh
if ! command -v vimox-hit >/dev/null 2>&1; then
  curl -fsSL https://raw.githubusercontent.com/neequ/vimox-hit/master/install.sh | sh
fi
export PATH="$HOME/.local/bin:$PATH"
vimox-hit version
```

Before selecting a source, check `https://api.country.is/` for the public IP country code. For `CN`, prefer the `gh-proxy.com` GitHub proxy; elsewhere use the official commands. If detection or the proxy fails, fall back to GitHub directly. The installer still verifies `SHA256SUMS`; the lookup provider receives the public IP, and Hit does not store it.

```powershell
$env:VIMOX_HIT_GITHUB_MIRROR='gh-proxy.com'; try { irm 'https://hit.vimox.cn/downloads/cli/install.ps1' | iex } finally { Remove-Item Env:VIMOX_HIT_GITHUB_MIRROR -ErrorAction SilentlyContinue }
```

```sh
curl -fsSL 'https://hit.vimox.cn/downloads/cli/install.sh' | VIMOX_HIT_GITHUB_MIRROR=gh-proxy.com sh
```

Then run `vimox-hit doctor` and inspect `data.authenticated`. If it is `false`, do not use CLI project commands; explain that Hit does not yet provide CLI OAuth authorization. Do not read browser storage or request tokens copied from browser storage. This does not block the remote MCP project workflow above when MCP is connected.

If an authenticated request returns HTTP 401 or 403, stop and report the authorization failure. Do not retry with a token taken from the browser.

When authenticated, use the project ID provided by the user:

```text
vimox-hit context show --project PROJECT_ID
vimox-hit ai-log list
vimox-hit ai-log show LOG_ID
vimox-hit ai-log diagnose LOG_ID
```

Diagnostics redact credentials and signed URL query parameters. Use `vimox-hit ai-log copy LOG_ID` to copy the same redacted bundle to the system clipboard. Before sharing it, check that prompts or project data are appropriate to disclose.

Before any model operation or final video generation, show the model, input assets, key parameters, and cost estimate, then obtain the user's explicit approval. Execution must consume a short-lived, single-use approval ID issued by the service for the exact input fingerprint. Never reuse approval after a retry, resume, model change, or cost change. The current CLI does not provide model execution commands; never claim that an unavailable operation was submitted.

Template automation: use MCP `template_list`, `template_upload`, `template_publish`, and `template_stats` when connected. The commands `template list|upload|publish|stats` require a CLI source build that includes template support; the current prebuilt GitHub release has not been updated, so use MCP until then. In a supported build, uploads validate and publish by default, while `--private` keeps them personal. Repeat `--media` in metadata index order for local video/image files. These operations do not invoke models. Do not claim a local path was uploaded through MCP.

Use `vimox-hit project init --name NAME` to create an isolated local project folder. See the [command reference](references/commands.md), [project folder format](references/project-format.md), and [troubleshooting guide](references/troubleshooting.md).

For Simplified Chinese, see [简体中文使用说明](references/zh-CN/SKILL.md).
