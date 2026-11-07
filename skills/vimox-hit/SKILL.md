---
name: vimox-hit
description: 使用 hit.vimox.cn 读取视频项目、诊断 AI 请求并管理独立产物目录；当用户要操作 Vimox Hit 或排查其模型请求时触发。
---

# Vimox Hit

使用 `vimox-hit` CLI 读取 Hit 项目与 AI 请求日志。首次使用时按 [安装说明](references/install.md) 安装官方二进制；不得下载未校验的构建产物。

先检查 `vimox-hit doctor`，再按当前项目 ID 调用：

```text
vimox-hit context show --project PROJECT_ID
vimox-hit ai-log list
vimox-hit ai-log show LOG_ID
vimox-hit ai-log diagnose LOG_ID
```

诊断内容经过凭据与签名 URL 清理，可通过 `vimox-hit ai-log copy LOG_ID` 复制到系统剪贴板。分享或提交诊断前确认内容不含用户不希望暴露的提示词或项目数据。

任何模型调用和最终视频生成前，必须展示模型、输入素材、关键参数和费用估算，并取得用户明确批准。每次执行须消费服务端为当前输入签发的短时单次 approval ID；重试、恢复、模型变化或费用变化都不能沿用旧批准。当前 CLI 尚未提供模型执行命令，不得声称已提交生成。

每项工作可用 `vimox-hit project init --name NAME` 创建独立本地项目目录。参考 [CLI 命令](references/commands.md)、[项目格式](references/project-format.md) 和 [故障恢复](references/troubleshooting.md)。
