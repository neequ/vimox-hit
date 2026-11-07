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
```

命令输出带版本号的 JSON envelope。认证使用 `VIMOX_HIT_TOKEN` 环境变量。不要把 token 写入命令参数、Skill 文件、项目清单或诊断包。服务端授权流程就绪后再补充交互式账号授权。

[English](../commands.md)
