# One API LLM 网关 | 企业级大模型 API 分发系统

## 项目定位

这是一个基于开源项目 One API 搭建的大模型（LLM）接口管理与分发系统。核心业务是将后端的多个大模型（如阿里云百炼的 DeepSeek）统一聚合并转换为标准 OpenAI API 格式，再分发给前端业务、开发工具（如 Cursor）或团队成员使用。

它不仅是一个简单的接口代理，而是完整跑通了：

- 跨云厂商的大模型聚合
- 团队 API Key 的生成与鉴权
- 额度分配与高精度使用量监控
- Hugging Face Spaces 免费云端部署
- UptimeRobot 无人值守自动保活

这个配置记录适合作为团队内部分发 API 或个人工具链基础设施的“标准操作手册（SOP）”：展示如何从云平台申请模型、搭建代理网关，一直推进到客户端无缝接入。

## 访问入口

### 管理后台

- 平台面板: <https://zhu592837154-my-api.hf.space>
- Hugging Face 账号: `zhu592837154`
- UptimeRobot 面板: <https://uptimerobot.com/dashboard>

### 接口调用地址

外部应用（如 Cursor、业务代码）统一接入的基础地址（Base URL）：

```text
https://zhu592837154-my-api.hf.space/v1
```

## 关键账号与凭证记录

> 注意：生产环境中务必妥善保管以下凭证，切勿直接提交到代码仓库。本文只记录凭证类型和配置位置，真实值建议保存到密码管理器或团队密钥库。

### 1. One API 系统账号

- 管理员账号: `root`
- 管理员密码: `<请在私密位置记录，并尽快在后台“用户”面板修改默认密码>`

### 2. 阿里云百炼（上游 Provider）

- 平台地址: <https://bailian.console.aliyun.com/>
- API Key: `<阿里云百炼 API Key，用于 One API 的“渠道”配置>`
- 核心模型: `deepseek-v4-flash`

### 3. 分发给同事/客户端的令牌（下游 Consumer）

- Token 格式: `sk-...`
- 生成位置: One API “令牌”页面
- 使用建议: 一人一 Key，并设置独立额度与名称，便于追踪和停用。

## 技术栈

### 核心架构

- API 网关: One API（统一转化并兼容 OpenAI v1 接口协议）
- 云端宿主机: Hugging Face Spaces（Docker Runtime / CPU Basic 免费实例）
- 上游模型服务: 阿里云百炼（DashScope）/ Google Gemini（备选）
- 保活监控: UptimeRobot（HTTP 探针）

### 客户端消费节点

- 开发辅助: Cursor AI 编程助手

## 核心配置流程

### 1. 渠道配置（对接阿里云）

- 在 One API 中添加新渠道。
- 类型选择 `阿里云百炼`。
- 模型填入 `deepseek-v4-flash`。
- 密钥填入阿里云申请的 `sk-...`。
- 提交并点击“测试”，确认响应时间（如 `2.61s`）即为打通。

### 2. 令牌分发（给同事创建 Key）

- 进入“令牌”菜单，点击“添加新的令牌”。
- 名称使用同事姓名或项目名隔离（如 `Token_For_Cursor_XiaoWang`）。
- 额度可设置固定配额限制（如 1000 万 Token）防刷。
- 生成后复制带有 `sk-` 前缀的长字符串交于对应使用者。

### 3. Cursor 客户端接入

- 打开 Cursor `Settings` -> `Models`。
- 关闭或覆盖原有的 `OpenAI Base URL`，填入网关地址：`https://zhu592837154-my-api.hf.space/v1`。
- 在 `OpenAI API Key` 处填入 One API 生成的令牌。
- 在模型选择框中手动键入 `deepseek-v4-flash` 并回车确认。

## 主要难点与踩坑记录

### 1. 接口协议的路径对齐

- 踩坑: 直接将 `https://zhu592837154-my-api.hf.space` 填入 Cursor，导致调用出现 `404` 错误。
- 解法: 外部客户端依赖 OpenAI 标准协议，必须在 Base URL 末尾显式追加 `/v1`，确保请求正确路由到 `/v1/chat/completions`。

### 2. 模型名称与上游强绑定

- 踩坑: 在 Cursor 下拉列表中随便选了一个模型，或者调用的模型名称与 One API 渠道中配置的不一致，导致报错 `Model not found`。
- 解法: 客户端请求的模型名称，必须与 One API 中配置的名称、以及阿里云后台的模型 Code 三者完全一致。Cursor 找不到时需手动输入名称强行匹配。

### 3. Hugging Face 的容器休眠机制

- 踩坑: HF Spaces 免费实例在长达 48 小时无公网 HTTP 请求访问时，会自动停止容器（Paused），导致 API 接口突然全量超时。
- 解法: 引入 UptimeRobot，创建 `HTTP(s)` 类型的 Monitor，每 5 分钟向 Space 的公网 URL 发送一次探测请求。以此维持 HF 的活跃度检测，实现低成本持续在线。

### 4. 账号额度安全隔离

- 踩坑: 团队多人共用同一个 API Key，一旦发生代码死循环疯狂调用，无法查出是谁消耗了额度，也无法单独阻断。
- 解法: 使用 One API 的令牌体系，实现一用户一 Key。在“日志”面板中可根据 Token 名称精确溯源每一次 API 请求的入参、响应耗时和计费情况。

## 项目收获

- 理解了 LLM API 请求的标准路由协议（OpenAI 兼容规范）。
- 掌握了 API 网关的核心价值：权限控制、聚合转发、费用监控。
- 熟悉了 Hugging Face Spaces 作为 Serverless 容器的运行逻辑及限制绕过技巧。
- 实践了 Cursor 等新一代 AI IDE 接入私有化部署模型的全过程。

## 后续优化方向

- 将系统中其余的 HF Spaces（如 `Mall Admin` 等）全部接入 UptimeRobot 统一保活。
- 为不同职能的同事分配不同价格梯度的模型（如日常使用分配 `qwen-plus`，复杂代码生成分配 `deepseek-v4-flash`）。
- 梳理 One API 面板中的日志记录，制定每月的 Token 账单复盘规范。
