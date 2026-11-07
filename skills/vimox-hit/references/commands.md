# CLI Commands

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

Commands emit a versioned JSON envelope. Provide authentication through the `VIMOX_HIT_TOKEN` environment variable. Never put tokens in command arguments, Skill files, project manifests, or diagnostic bundles. Interactive account authorization will be added after the service authorization flow is available.

See the [Simplified Chinese command reference](zh-CN/commands.md).
