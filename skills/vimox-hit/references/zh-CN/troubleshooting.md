# 故障排查

运行 `vimox-hit doctor` 检查 API 连通性，再用 `vimox-hit ai-log list` 找到请求。`show` 读取当前账号可见的详情，`diagnose` 生成脱敏诊断包，`copy` 将相同内容复制到剪贴板。凭据和签名 URL 查询参数会被清除。

如果 doctor 返回 `authenticated: false`，当前 CLI 尚不能通过 OAuth 读取项目和 AI 日志。请在 Hit 网页端查看账户数据；不要从浏览器存储复制访问令牌。

只对只读查询进行有限重试。重试任务前先检查服务端状态。每次模型调用和最终视频生成都要求用户明确批准；重试或模型、输入、参数、费用估算发生变化后，必须重新申请批准。

[English](../troubleshooting.md)
