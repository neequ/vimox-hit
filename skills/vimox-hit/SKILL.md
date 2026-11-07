---
name: vimox-hit
description: Use the Vimox Hit CLI to inspect video projects, diagnose AI requests, and create local project folders. Trigger when the user asks to work with Vimox Hit or troubleshoot its model requests.
---

# Vimox Hit

Use Vimox Hit's remote MCP tools when they are available. This is the zero-install path for ordinary users: query the project with `video_project_state`, open the returned workbench URL, and use `template_upload` for complete template JSON and keep model approvals in the web app. Do not install a local CLI merely because this Skill was invoked.

If the requested workflow needs the CLI and it is missing, offer the official platform installer below. When the user has asked you to set it up and the host grants terminal/network access, run it and verify with `vimox-hit version`; otherwise show the exact command and stop before claiming installation. Never claim that a prompt can register an MCP connector inside Codex; that remains a client-level setting.

On Windows PowerShell:

```powershell
$installDir = Join-Path $env:LOCALAPPDATA 'Vimox\bin'
if (-not (Get-Command vimox-hit -ErrorAction SilentlyContinue) -and (Test-Path (Join-Path $installDir 'vimox-hit.exe'))) {
  $env:Path = "$installDir;$env:Path"
}
if (-not (Get-Command vimox-hit -ErrorAction SilentlyContinue)) {
  irm https://raw.githubusercontent.com/neequ/vimox-hit/main/install.ps1 | iex
}
$env:Path = "$installDir;$env:Path"
vimox-hit version
```

On macOS or Linux:

```sh
if ! command -v vimox-hit >/dev/null 2>&1; then
  curl -fsSL https://raw.githubusercontent.com/neequ/vimox-hit/main/install.sh | sh
fi
export PATH="$HOME/.local/bin:$PATH"
vimox-hit version
```

Then run `vimox-hit doctor` and inspect `data.authenticated`. If it is `false`, stop and explain that Hit does not yet provide CLI OAuth authorization. Do not read browser storage, request tokens copied from browser storage, or continue with project commands. The user can still inspect the project in the Hit web app.

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

Template automation: use MCP `template_list`, `template_upload`, `template_publish`, and `template_stats` when connected. For an authenticated local CLI, use `template upload --kind video|suite --file FILE`; uploads validate and publish by default, while `--private` keeps them personal. Repeat `--media` in metadata index order for local video/image files. These operations do not invoke models. Do not claim a local path was uploaded through MCP.

Use `vimox-hit project init --name NAME` to create an isolated local project folder. See the [command reference](references/commands.md), [project folder format](references/project-format.md), and [troubleshooting guide](references/troubleshooting.md).

For Simplified Chinese, see [简体中文使用说明](references/zh-CN/SKILL.md).
