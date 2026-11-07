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

命令默认输出版本化 JSON。认证通过 `VIMOX_HIT_TOKEN` 环境变量提供；不要把 token 放进命令参数、Skill、项目 manifest 或诊断包。生产登录交互和设备授权会在服务端授权协议落地后补齐。
