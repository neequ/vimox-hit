# 模板自动化指南

Vimox Hit 的模板操作覆盖视频模板和套图模板：浏览系统模板、上传个人模板、自动发布以及读取浏览/播放/克隆/分享统计。模板上传和发布只做格式校验与存储，不会调用模型，也不会创建生成任务。

## 选用 CLI 或 MCP

- 已连接 Vimox Hit 远程 MCP 的 Agent 可直接使用 `template_list`、`template_upload`、`template_publish` 和 `template_stats`。
- 需要从本机上传视频或图片时使用 CLI 的重复 `--media` 参数。远程 MCP 无法读取你电脑上的文件路径。
- CLI 使用授权渠道提供的 `VIMOX_HIT_TOKEN`。当前 GitHub Releases 预编译版尚未包含模板子命令；先运行 `vimox-hit version` 和 `vimox-hit doctor` 检查环境。若命令提示未知 `template`，请使用远程 MCP 或包含模板功能的 CLI 源码构建。
- 不要从浏览器存储复制 token，也不要把凭据写进命令参数、JSON、Skill 或日志。

## 浏览模板

```sh
# 系统视频模板
vimox-hit template list --kind video --origin system --limit 24

# 个人视频模板和套图模板
vimox-hit template list --kind video --origin mine
vimox-hit template list --kind suite --origin mine

# 搜索并读取下一页
vimox-hit template list --kind video --origin system --query "商品展示" --limit 24
vimox-hit template list --kind video --origin system --cursor NEXT_CURSOR
```

列表输出包含 `items`、`total`、`nextCursor`、分类、标签和统计。视频分类必须使用服务端已有分类，可先从视频模板列表的 `categories` 读取；REST 客户端也可调用 `GET /api/v1/video-template-categories`。

## 上传视频模板

先准备符合 `ImportVideoTemplateInput` 的 JSON。视频模板需要 `sourceKey`、名称、描述、有效分类、标签和至少一个分镜。分镜的图型、中文静帧提示词和中文视频提示词必须有效；画幅、分辨率和时长写在 `payload`。

```json
{
  "sourceKey": "docs:product-demo:v1",
  "name": "商品展示模板",
  "description": "可复用的商品近景展示分镜。",
  "category": "按摩养生",
  "tags": ["商品展示"],
  "payload": {
    "campaignStyleLock": "自然柔光，商品外观保持真实",
    "videoAspectRatio": "9:16",
    "videoResolution": "720p",
    "slots": [],
    "shots": [
      {
        "assetType": "hero-image",
        "displayName": "商品主镜头",
        "promptInstruction": "商品置于干净台面，柔光表现真实材质和轮廓。",
        "videoPrompt": "镜头缓缓推近商品，保持商品外观一致，背景简洁。",
        "durationSec": 6,
        "videoAspectRatio": "9:16",
        "videoResolution": "720p",
        "sourceVideoMediaIndex": 0
      }
    ]
  }
}
```

`sourceKey` 是账号范围内的幂等键：相同账号再次提交相同 key 会返回已有模板，不会覆盖修订；要上传新版本，请换一个 key。

视频和图片按 `--media` 出现顺序编号，从 0 开始。分镜的 `sourceVideoMediaIndex` 指向视频文件，DNA 槽位的 `assetMediaIndex` 指向图片。每个分镜不写 `sourceVideoMediaIndex` 时，服务端会关联本次上传的第一个视频；没有上传媒体时，不能保留非空媒体索引。

```sh
# sourceVideoMediaIndex 为 0，上传一段原片
vimox-hit template upload --kind video --file video-template.json --media ./source.mp4

# 同一命令可上传多个视频/图片，索引按参数顺序计算
vimox-hit template upload --kind video --file video-template.json --media ./source.mp4 --media ./product.png
```

上传默认自动发布，无需人工审核。成功响应中的 `data.template.id` 是模板 ID。仅保存为个人模板时加 `--private`；之后可以用该 ID 单独发布：

```sh
vimox-hit template upload --kind video --file video-template.json --private
vimox-hit template publish --kind video --id TEMPLATE_ID
```

若保存成功但发布失败，CLI 会在错误中返回模板 ID，并提示重试命令。无需重新上传，可直接运行 `template publish`。

## 上传套图模板

套图 JSON 使用 `EcomSuiteFile` 格式。提示词模式至少提供一个 `promptTemplate`；资产模式还必须提供完整、双向匹配的资产身份与逐图绑定矩阵。下面是提示词模式的最小结构，可按需要增加分镜和说明字段：

```json
{
  "name": "商品柔光套图",
  "category": { "l1": "商品展示", "l2": "通用", "leaf": "商品" },
  "styleLock": { "lockText": "白底柔光，保持商品真实外观" },
  "shots": [
    {
      "shotId": "hero",
      "order": 1,
      "shotRole": "HERO",
      "displayName": "商品主图",
      "promptTemplate": "在白底柔光中展示 {product}，突出商品真实材质与轮廓。"
    }
  ]
}
```

```sh
# 默认发布
vimox-hit template upload --kind suite --file suite-template.json

# 只保存为个人模板
vimox-hit template upload --kind suite --file suite-template.json --private
vimox-hit template publish --kind suite --id TEMPLATE_ID
```

CLI 不接受套图 `--media`。套图原图使用 `provenance.sourceImages` 中可访问的 HTTPS URL；CLI 不会读取或上传本机套图媒体。

## 使用远程 MCP

模板工具需要 Vimox Hit MCP 已连接并完成 OAuth 授权。`template_upload` 接收完整模板 JSON；`publish` 省略时默认为 `true`，传 `false` 则只保存个人模板。

```text
template_list({ kind: "VIDEO", origin: "system", limit: 24 })
template_upload({ kind: "VIDEO", template: <完整视频模板 JSON> })
template_upload({ kind: "SUITE", template: <完整套图 JSON>, publish: false })
template_publish({ kind: "VIDEO", templateId: "TEMPLATE_ID" })
template_stats({ kind: "VIDEO", templateId: "TEMPLATE_ID" })
```

MCP 不会从 JSON 中的本机路径读取或上传媒体。视频模板需要本机原片或 DNA 图片时，改用 CLI `--media`；Agent 不应声称本地文件已经通过 MCP 上传。

## 查看统计

```sh
vimox-hit template stats --kind video --id TEMPLATE_ID
vimox-hit template stats --kind suite --id TEMPLATE_ID
```

统计字段为 `viewCount`、`playCount`、`cloneCount` 和 `shareCount`：打开模板详情计一次浏览；视频实际开始播放计一次播放；成功创建模板项目计一次克隆；成功复制模板链接计一次分享。浏览、播放和分享按账号与模板在 30 分钟内去重；克隆仅在项目创建成功时记录，相同项目 ID 不重复计数。读取列表或统计本身不会增加次数。统计从功能上线后累积，不回填历史数据。

## REST 接口

所有业务接口都需要当前 Vimox Hit 账号的授权。浏览器不要直接访问 Redis、SQLite、Provider 或文件存储。

| 用途 | 接口 |
| --- | --- |
| 读取视频分类 | `GET /api/v1/video-template-categories` |
| 浏览视频模板 | `GET /api/v1/video-templates` |
| 导入视频模板 | `POST /api/v1/video-templates/import`，JSON 或 multipart |
| 浏览/导入套图模板 | `GET /api/v1/suites` / `POST /api/v1/suites` |
| 发布个人模板 | `POST /api/v1/template-publications`，请求体 `{ "kind": "VIDEO", "templateId": "..." }` |
| 读取统计 | `GET /api/v1/template-stats?kind=VIDEO&templateId=...` |
| 记录浏览/播放/分享 | `POST /api/v1/template-stats/events` |

视频 multipart 使用 `metadata` 字段传 JSON 字符串，并按索引顺序重复传 `media` 字段。客户端不能自行增加克隆数；克隆统计由成功创建项目的服务端流程记录。

## 上传限制与自动发布

- 视频模板 JSON 不超过 200 KB；名称最多 120 字，描述最多 1000 字，标签最多 32 个。
- 视频媒体最多 20 个文件；单文件最多 128 MiB，总大小最多 192 MiB。接受视频和图片 MIME 类型；上传媒体要求服务端已配置 COS。
- 视频模板的 `category` 必须是当前有效分类；每个分镜必须通过当前模板、提示词及媒体引用校验。
- 套图模板须通过 `EcomSuiteFile` 和资产矩阵校验；格式不通过时不会发布。
- 上传后默认立即创建系统模板快照。重复发布会复用已有快照；修改个人原稿不会静默覆盖已发布快照，修订应重新导入并使用新的 `sourceKey`。

字段和限制会由服务端当前契约校验；格式不匹配时，接口会返回校验错误且不会发布模板。
