---
name: vimox-hit
description: Use the Vimox Hit CLI to inspect video projects, diagnose AI requests, and create local project folders. Trigger when the user asks to work with Vimox Hit or troubleshoot its model requests.
---

# Vimox Hit

Use `vimox-hit` to read Hit project context and AI request logs. Before the first CLI command, check whether `vimox-hit` is installed. If it is missing, install the latest official release with the platform installer below; the installer verifies the release checksum.

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

Use `vimox-hit project init --name NAME` to create an isolated local project folder. See the [command reference](references/commands.md), [project folder format](references/project-format.md), and [troubleshooting guide](references/troubleshooting.md).

For Simplified Chinese, see [简体中文使用说明](references/zh-CN/SKILL.md).
