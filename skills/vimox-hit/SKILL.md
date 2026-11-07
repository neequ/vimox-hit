---
name: vimox-hit
description: Use the Vimox Hit CLI to inspect video projects, diagnose AI requests, and create local project folders. Trigger when the user asks to work with Vimox Hit or troubleshoot its model requests.
---

# Vimox Hit

Use `vimox-hit` to read Hit project context and AI request logs. Follow [the installation guide](references/install.md) before first use. Install only official, checksum-verified releases.

Start with `vimox-hit doctor`, then use the project ID provided by the user:

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
