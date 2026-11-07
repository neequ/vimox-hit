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
vimox-hit template list --kind video --origin system [--cursor CURSOR]
vimox-hit template upload --kind video --file video-template.json [--media source.mp4] [--private]
vimox-hit template upload --kind suite --file suite-template.json [--private]
vimox-hit template publish --kind video --id TEMPLATE_ID
vimox-hit template stats --kind video --id TEMPLATE_ID
```

Commands emit a versioned JSON envelope. Provide authentication through the `VIMOX_HIT_TOKEN` environment variable. Never put tokens in command arguments, Skill files, project manifests, or diagnostic bundles. Interactive account authorization will be added after the service authorization flow is available.

See the [Simplified Chinese command reference](zh-CN/commands.md).

Uploads validate and publish automatically unless `--private` is supplied. Repeat `--media` in metadata index order for video/image files. Statistics include views, actual plays, successful clones and shares. Connected MCP clients can use `template_upload`, `template_list`, `template_publish` and `template_stats` without a local CLI.

The template subcommands are available in source, but the current prebuilt GitHub release has not been updated to include them. Use remote MCP until a CLI binary with template support is published.
