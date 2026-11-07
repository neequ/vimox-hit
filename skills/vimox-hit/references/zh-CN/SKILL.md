# Vimox Hit

[English](../../SKILL.md) | 简体中文

使用 `vimox-hit` 读取 Hit 项目上下文和 AI 请求日志。首次使用前请先阅读[安装说明](install.md)，只安装官方且通过校验的发行版。

先运行 `vimox-hit doctor`，再使用用户提供的项目 ID：

```text
vimox-hit context show --project PROJECT_ID
vimox-hit ai-log list
vimox-hit ai-log show LOG_ID
vimox-hit ai-log diagnose LOG_ID
```

诊断内容会清除凭据和签名 URL 查询参数。使用 `vimox-hit ai-log copy LOG_ID` 可将相同的脱敏诊断包复制到系统剪贴板。分享前请确认提示词和项目数据适合公开。

每次模型调用和最终视频生成前，都必须向用户展示模型、输入素材、关键参数和费用估算，并取得用户明确批准。执行时必须使用服务端针对当前输入指纹签发的短时单次 approval ID。重试、恢复、模型变化或费用变化后不得复用旧批准。当前 CLI 不提供模型执行命令；不得声称未开放的操作已提交。

使用 `vimox-hit project init --name NAME` 创建独立的本地项目目录。更多信息见[命令参考](commands.md)、[项目目录格式](project-format.md)和[故障排查](troubleshooting.md)。
