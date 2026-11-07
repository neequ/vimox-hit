# CLI 命令

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

命令输出带版本号的 JSON envelope。认证使用 `VIMOX_HIT_TOKEN` 环境变量。不要把 token 写入命令参数、Skill 文件、项目清单或诊断包。服务端授权流程就绪后再补充交互式账号授权。

[English](../commands.md)

模板上传默认校验后自动发布，无需人工审核。`--private` 只保存个人模板；`--media` 按索引顺序重复传入视频/图片。统计包含浏览、播放、成功克隆和分享。没有 CLI 授权时使用已连接的远程 MCP：`template_upload`、`template_list`、`template_publish`、`template_stats`。

模板子命令已在源码中实现，当前预编译 GitHub Release 尚未包含这些命令。CLI 二进制更新前，请使用远程 MCP 完成模板操作。
