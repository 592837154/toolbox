# Remote Debug Log Relay

面向前端测试环境和 H5 真机问题排查的远程日志回流工作流。目标是在本机终端收集页面运行日志，再把结构化日志交给 AI 分析调用链、状态分支和异常原因，减少人工复现、截图和反复猜测。

## 适用场景

- 测试环境页面能复现问题，但本地开发环境难以复现。
- H5 真机、微信/QQ/支付宝 WebView 等环境无法方便打开完整 DevTools。
- bug 需要多轮加日志、部署、复现、分析。
- 希望把前端运行时状态、接口结果、分支判断整理成结构化日志后交给 AI 排查。

## 架构

```text
测试环境页面 / H5 页面
  -> debugLog(tag, msg, data)
  -> new Image().src = LOG_URL/pixel.gif?payload=...
  -> Cloudflare Tunnel HTTPS 域名
  -> cloudflared 转发到本机 http://127.0.0.1:3695
  -> server-debug-log.cjs
  -> 解析 payload
  -> console.log 打到本机终端
  -> 复制日志给 AI 分析
```

## 组件职责

### debugLog.ts

前端日志 SDK。负责把 `tag`、`msg`、`data` 序列化成 JSON，并通过图片 GET 请求上报。

使用图片请求的原因：

- 接入轻，不需要业务接口支持。
- 对 H5 和 WebView 兼容性较好。
- 避免部分场景下 `fetch`、`sendBeacon`、CORS、预检请求带来的干扰。

### debug-log-endpoint.json

上报配置文件。

```json
{
  "protocol": "https",
  "host": "example.trycloudflare.com",
  "listenPort": 3695
}
```

- `protocol + host + port`：前端实际请求的公网地址。
- `listenPort`：本机 Node 日志服务监听端口。

注意：如果使用 Cloudflare quick tunnel，`host` 每次重启 tunnel 后都可能变化。

### server-debug-log.cjs

本机日志接收服务。监听 `listenPort`，收到 `/pixel.gif?payload=...` 后解析并打印 JSON，同时返回 1x1 gif，避免浏览器资源加载报错。

### cloudflared

公网入口。把测试环境页面可访问的 HTTPS 地址转发到本机日志服务。

```text
https://xxx.trycloudflare.com -> http://127.0.0.1:3695
```

## 标准工作流

### 1. 启动本机日志服务

```bash
cd D:\work\phoenix-share\debug-log
node server-debug-log.cjs
```

期望输出：

```text
log server listening: http://127.0.0.1:3695
```

### 2. 启动 Cloudflare Tunnel

```bash
cloudflared tunnel --url http://127.0.0.1:3695
```

记录输出里的 quick tunnel 域名：

```text
Your quick Tunnel has been created! Visit it:
https://xxx.trycloudflare.com
```

### 3. 验证本机服务

```bash
curl -s -o /dev/null "http://localhost:3695/pixel.gif?payload=%7B%22tag%22%3A%22local-test%22%7D"
```

`server-debug-log.cjs` 终端应打印：

```json
{
  "tag": "local-test"
}
```

### 4. 验证 tunnel

```bash
curl -s -o /dev/null "https://xxx.trycloudflare.com/pixel.gif?payload=%7B%22tag%22%3A%22tunnel-test%22%7D"
```

`server-debug-log.cjs` 终端应打印：

```json
{
  "tag": "tunnel-test"
}
```

如果返回 `530` 或终端没有输出，说明 Cloudflare Tunnel 没有转发到本机。先修 tunnel，不要急着部署测试环境。

### 5. 更新前端上报地址

把新 tunnel 域名写入 `debug-log-endpoint.json`：

```json
{
  "protocol": "https",
  "host": "xxx.trycloudflare.com",
  "listenPort": 3695
}
```

`host` 不要带 `https://`。

### 6. 部署测试环境

`debug-log-endpoint.json` 会被打进前端包。改完配置后，必须重新构建并部署测试环境，否则页面仍会请求旧 tunnel 域名。

### 7. 复现问题并收集日志

打开测试环境页面，复现 bug。本机 `server-debug-log.cjs` 终端会输出结构化日志。

建议复制包含以下信息的完整日志给 AI：

- 页面初始化日志。
- 关键按钮点击日志。
- 接口入参和响应 code。
- 关键状态值和分支判断。
- 异常信息、浏览器 Network 状态码、截图描述。

## 排查清单

### 本机 curl 有输出，tunnel curl 无输出

说明本机 Node 服务正常，Cloudflare Tunnel 不通。

处理方式：

- 检查 cloudflared 是否仍在运行。
- 使用 `http://127.0.0.1:3695` 重新启动 tunnel。
- 重新复制 quick tunnel 新域名。
- 再次执行 tunnel curl 验证。

### 浏览器 Network 显示 530

说明请求到达 Cloudflare，但没有成功转发到本机服务。

处理方式：

- 先用 curl 验证 tunnel。
- 如果 curl 也 530，重启 cloudflared 获取新域名。
- 新域名验证成功后再更新配置并部署。

### 浏览器 Network 显示 ERR_NAME_NOT_RESOLVED

说明域名解析失败。

处理方式：

- 检查 `debug-log-endpoint.json` 里的 host 是否和 cloudflared 输出完全一致。
- 注意单词拼写，不要把 quick tunnel 旧域名写进配置。
- 确认 dev3 已部署最新包。

### 浏览器控制台有 debugLog，但本机无输出

说明前端调用触发了，但请求没有进入本机 Node 服务。

检查顺序：

```text
Network 请求地址是否是最新 tunnel 域名
-> 请求状态是否 200
-> curl tunnel 是否能打印 tunnel-test
-> server-debug-log.cjs 是否监听正确端口
```

## 日志设计建议

日志要面向 AI 分析，尽量结构化。

推荐格式：

```ts
debugLog('video-flow-limit', 'close traffic limit modal', {
  shareKey,
  fileId,
  showType,
  isMobile,
  pageType,
  trafficSwitch,
  totalFileSize,
  reason: 'video_flow_limit_modal_closed'
});
```

推荐字段：

- `tag`：业务链路，例如 `video-flow-limit`。
- `msg`：当前阶段，例如 `show modal`、`close modal`。
- `shareKey`：分享标识。
- `traceId`：一次页面访问的链路 ID。
- `href`：当前页面 URL。
- `userAgent`：终端环境。
- `isMobile`：是否移动端。
- `pageType`：页面类型。
- `request` / `response`：关键接口信息，注意脱敏。
- `reason`：分支命中的原因。

## 安全与清理

- 不要上报 token、手机号、身份证、完整 cookie、支付二维码等敏感信息。
- 临时业务日志用完后要搜索 `debugLog(` 和业务 tag 清理。
- `debugLog.ts` 工具可以保留，业务调用应按需添加和删除。
- Quick tunnel 是临时公网入口，不适合长期稳定依赖。

## 后续优化

当前 quick tunnel 方案需要把域名写入配置并重新部署。更好的方向是运行时配置：

```js
localStorage.setItem('__debug_log_url__', 'https://xxx.trycloudflare.com');
location.reload();
```

或 URL 参数：

```text
?debugLogUrl=https%3A%2F%2Fxxx.trycloudflare.com
```

这样可以避免每次 tunnel 域名变化后都重新部署测试环境。

