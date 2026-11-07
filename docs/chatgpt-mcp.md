# Vimox Hit MCP 接入

Vimox Hit 提供远程 MCP 地址：

```text
https://hit.vimox.cn/mcp
```

连接后可用自然语言列出视频项目、检查项目状态、创建空项目并打开网页工作台。MCP 不会直接读取用户电脑上的视频文件，也不会绕过确认触发 DNA 解析或视频生成。

## 普通用户使用方式

在 ChatGPT、Codex 或其他 Agent 客户端中复制下面的任务文本。它会先检查依赖；缺少时给出官方连接/安装方式，并在有终端权限时执行安装后验证：

```text
请按 Vimox Hit 视频创作工作流协助我完成当前项目。
先检查本轮可用的 Vimox Hit MCP、Skill 和 CLI。缺少 MCP 时展示地址 https://hit.vimox.cn/mcp，并告诉我在当前客户端的 MCP/Connector 设置中添加服务、完成 Vimox 登录和安全配置授权凭据；无法由当前 Agent 完成的步骤要具体指出。缺少 Skill 时展示并执行（若有终端权限）npx skills add neequ/vimox-hit --skill vimox-hit -g；缺少 CLI 时按系统展示官方安装命令：Windows 使用 irm https://raw.githubusercontent.com/neequ/vimox-hit/master/install.ps1 | iex，macOS/Linux 使用 curl -fsSL https://raw.githubusercontent.com/neequ/vimox-hit/master/install.sh | sh。若任务需要 CLI 且我已要求配置，终端/网络权限可用时可执行安装；否则展示命令。安装后验证：vimox-hit version。安装或连接后再继续，不要让我猜下一步。
目标流程：导入原片、解析视频 DNA、检查并补充素材、确认分镜后生成成片。
每次模型调用或视频生成前，展示实际模型、输入素材、关键参数和费用估算，等待我明确批准；不要伪造完成结果。
```

如果当前客户端允许 Agent 使用终端和网络，Agent 可以执行上面的官方安装命令；MCP 注册和 Vimox 登录仍需按客户端授权界面完成。当前 MCP 入口使用 Vimox Bearer 身份认证，客户端若没有 OAuth 登录按钮，必须在客户端的安全凭据设置中配置 Vimox OAuth Bearer token；不要把 token 粘贴到对话或项目文件。远程客户端不能直接读取本机视频，原片继续通过网页工作台上传。

## 工具边界

- `video_projects_list`：列出当前账户的视频项目。
- `video_project_state`：读取原片、图片素材、分镜、任务、视频计划和成片状态，并返回下一步。
- `video_project_create`：创建空项目，不上传素材、不调用模型。
- `video_project_open`：返回网页工作台地址。

模型调用、DNA 解析、素材生成和成片生成必须在展示模型、素材、参数与费用后取得用户明确批准。HTTP 202 只表示任务已受理，不表示任务完成。
