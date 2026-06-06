# Sub2API 接入 Cursor 过程记录

记录时间：2026-06-06 20:18:36 +08:00

本记录整理了本次在 Codex 窗口里完成的 Sub2API、Cloudflare Tunnel、Cursor 自定义 OpenAI Base URL 接入过程，包括中途遇到的问题、原因、修复动作和最终验证证据。

源码 fork：

```text
https://github.com/592837154/sub2api
```

该 fork 用于后续保留本地部署、配置修复、Cursor 接入经验和可能的二次开发改动。当前实操环境仍运行在本机已下载的 Sub2API 项目目录中，后续如需长期维护，建议把本地目录切换到该 fork 的 Git remote，或重新 clone fork 后迁移配置。

## 1. 环境信息

Sub2API 项目目录：

```text
C:\Users\Administrator\Desktop\learn\sub2api-main\sub2api-main
```

GitHub fork：

```text
https://github.com/592837154/sub2api
```

Sub2API 后端配置文件：

```text
C:\Users\Administrator\Desktop\learn\sub2api-main\sub2api-main\backend\config.yaml
```

后端监听配置：

```yaml
server:
  host: 127.0.0.1
  port: 8080
```

数据库配置：

```yaml
database:
  host: 127.0.0.1
  port: 5432
  user: postgres
  password: sub2api_dev_password
  dbname: sub2api
```

Redis：

```yaml
redis:
  host: 127.0.0.1
  port: 6379
```

本机实际监听确认：

```text
127.0.0.1:8080  Sub2API backend
127.0.0.1:5432  PostgreSQL
127.0.0.1:6379  Redis
```

后端进程：

```text
server.exe
PID: 12372
```

PostgreSQL 路径：

```text
C:\Users\Administrator\Desktop\learn\sub2api-tools\pgsql\bin\postgres.exe
```

## 2. Sub2API 初始状态

用户已经在 Sub2API 后台成功添加了一个 OpenAI OAuth 账号。

账号信息：

```text
账号名称：11231
邮箱：zhukaixue2026@126.com
平台：OpenAI
类型：OAuth
状态：正常 / active
调度：开启
```

API Key：

```text
名称：111
Key：sk-47c1fbe...413f69
```

出于安全原因，这里只记录掩码，不记录完整 API Key。

## 3. API Key 创建建议

创建 API Key 时建议：

```text
名称：按用途命名，例如 cursor、test、myapp
分组：后续要和上游账号同平台同分组
自定义密钥：关闭
IP 限制：测试阶段关闭
额度限制：自己测试可设 1、5、10 美元，也可 0 表示无限制
速率限制：测试阶段关闭
密钥有效期：自己长期用可关闭，给别人临时用建议开启
```

注意：生成多个 API Key 并不会生成多份上游额度。多个 Key 共享所绑定账号池的可用能力和 Sub2API 用户余额。

## 4. Cursor 不能直接使用本地地址

最开始尝试在 Cursor Settings 里填：

```text
Override OpenAI Base URL:
http://127.0.0.1:8080/v1
```

Cursor 返回错误：

```text
Provider Error
Access to private networks is forbidden
```

原因：

Cursor 内置 Ask / Agent 对自定义 Provider 禁止访问私有网络地址，包括：

```text
localhost
127.0.0.1
192.168.x.x
10.x.x.x
172.16.x.x ~ 172.31.x.x
```

因此即使本机 Sub2API 正常运行，Cursor 内置聊天也不能直接访问本机或局域网 IP。

## 5. Cloudflare Tunnel 方案

为解决 Cursor 不能访问本地地址的问题，使用 Cloudflare Tunnel 把本机后端临时暴露为公网 HTTPS。

### 5.1 winget 安装失败

尝试：

```powershell
winget install --id Cloudflare.cloudflared --accept-package-agreements --accept-source-agreements
```

失败信息：

```text
No applicable app licenses found
```

于是改为直接下载 Cloudflare 官方 GitHub 发行版。

### 5.2 下载 cloudflared

下载到：

```text
C:\Users\Administrator\Desktop\learn\sub2api-main\sub2api-main\tools\cloudflared.exe
```

下载地址：

```text
https://github.com/cloudflare/cloudflared/releases/latest/download/cloudflared-windows-amd64.exe
```

版本确认：

```text
cloudflared version 2026.5.2
```

### 5.3 启动 Tunnel

启动命令：

```powershell
.\tools\cloudflared.exe tunnel --url http://127.0.0.1:8080 --logfile .\tools\cloudflared-tunnel.log --loglevel info
```

实际是后台启动的，进程信息：

```text
cloudflared PID: 8296
```

日志文件：

```text
C:\Users\Administrator\Desktop\learn\sub2api-main\sub2api-main\tools\cloudflared-tunnel.log
```

生成的临时公网地址：

```text
https://entirely-bible-etc-households.trycloudflare.com
```

Cursor 里最终使用的 Base URL：

```text
https://entirely-bible-etc-households.trycloudflare.com/v1
```

注意：这是 Cloudflare Quick Tunnel 临时地址。重启电脑、停止 cloudflared 进程或隧道断开后，地址可能失效，需要重新启动并获取新地址。

停止 Tunnel：

```powershell
Stop-Process -Id 8296
```

### 5.4 Tunnel 可用性验证

请求：

```powershell
Invoke-WebRequest -Uri "https://entirely-bible-etc-households.trycloudflare.com/health"
```

返回：

```json
{"status":"ok"}
```

说明 Cloudflare Tunnel 已成功转发到本地 Sub2API 后端。

## 6. Cursor 配置

Cursor Settings 中配置：

```text
OpenAI API Key:
sk-47c1fbe...413f69
```

```text
Override OpenAI Base URL:
https://entirely-bible-etc-households.trycloudflare.com/v1
```

最开始使用 `gpt-5`，后续发现该模型在当前 Sub2API/OpenAI OAuth 账号下不可用或不在模型列表内，因此改为自定义模型：

```text
gpt-5.5
```

也测试成功过：

```text
gpt-5.4-mini
```

## 7. 遇到的问题和修复

### 7.1 地址被 Cursor 拦截

错误：

```text
Access to private networks is forbidden
```

原因：

Cursor 不允许自定义 Provider 访问 `127.0.0.1` 或局域网 IP。

解决：

使用 Cloudflare Tunnel，改填公网 HTTPS：

```text
https://entirely-bible-etc-households.trycloudflare.com/v1
```

### 7.2 用户余额不足

错误：

```text
INSUFFICIENT_BALANCE
Insufficient account balance
```

原因：

Sub2API 用户余额为 `$0.00`。即使上游账号可用，Sub2API 自身也会先做余额判断。

解决：

在 Sub2API 后台用户管理中给 `admin` 用户充值，最终余额显示：

```text
$10.00
```

### 7.3 No available accounts

错误：

```text
No available accounts: no available accounts
```

核心原因：

API Key 和 OpenAI 账号没有在同一个 OpenAI 分组里。

最开始系统里有默认分组：

```text
name: default
platform: anthropic
```

这个默认分组是 Anthropic 平台的，不是 OpenAI 平台。API Key 一开始绑到了 `default/anthropic`，但 OpenAI 账号不能通过这个分组调度。

后来创建/确认了 OpenAI 分组：

```text
id: 2
name: openai
platform: openai
status: active
rate_multiplier: 1
```

账号绑定：

```text
account_id: 1
account_name: 11231
platform: openai
group: openai
```

API Key 原状态：

```text
api_key_id: 1
name: 111
group_id: 1
group_name: default
group_platform: anthropic
```

修复动作：

直接在数据库中把 API Key 改到 OpenAI 分组：

```sql
update api_keys
set group_id = (
  select id
  from groups
  where name = 'openai'
    and platform = 'openai'
    and deleted_at is null
  limit 1
),
updated_at = now()
where id = 1;
```

修复后：

```text
api_key_id: 1
name: 111
group_id: 2
group_name: openai
group_platform: openai
```

### 7.4 gpt-5 不可用

使用 `gpt-5` 时出现过：

```text
Service temporarily unavailable
```

查询 `/v1/models` 后发现可用模型列表里有：

```text
gpt-5.2
gpt-5.2-chat-latest
gpt-5.3-codex
gpt-5.4
gpt-5.4-mini
gpt-5.5
```

但没有直接的：

```text
gpt-5
```

测试结果：

```text
gpt-5.2      失败：The 'gpt-5.2' model is not supported when using Codex with a ChatGPT account.
gpt-5.4-mini 成功
gpt-5.5      成功
```

最终 Cursor 使用：

```text
gpt-5.5
```

## 8. 后台验证命令

使用 PostgreSQL 客户端：

```text
C:\Users\Administrator\Desktop\learn\sub2api-tools\pgsql\bin\psql.exe
```

### 8.1 查看分组

```powershell
& 'C:\Users\Administrator\Desktop\learn\sub2api-tools\pgsql\bin\psql.exe' `
  -U postgres -h 127.0.0.1 -d sub2api `
  -c "select id,name,platform,status,rate_multiplier from groups where deleted_at is null order by id;"
```

结果：

```text
 id |  name   | platform  | status | rate_multiplier
----+---------+-----------+--------+-----------------
  1 | default | anthropic | active |          1.0000
  2 | openai  | openai    | active |          1.0000
```

### 8.2 查看账号和分组绑定

```powershell
& 'C:\Users\Administrator\Desktop\learn\sub2api-tools\pgsql\bin\psql.exe' `
  -U postgres -h 127.0.0.1 -d sub2api `
  -c "select a.id,a.name,a.platform,a.type,a.status,a.schedulable,a.priority,a.rate_multiplier,string_agg(g.name || ':' || g.platform, ', ') as groups from accounts a left join account_groups ag on ag.account_id=a.id left join groups g on g.id=ag.group_id where a.deleted_at is null group by a.id order by a.id;"
```

结果：

```text
 id | name  | platform | type  | status | schedulable | priority | rate_multiplier |    groups
----+-------+----------+-------+--------+-------------+----------+-----------------+---------------
  1 | 11231 | openai   | oauth | active | t           |        1 |          1.0000 | openai:openai
```

### 8.3 查看 API Key 和分组绑定

```powershell
& 'C:\Users\Administrator\Desktop\learn\sub2api-tools\pgsql\bin\psql.exe' `
  -U postgres -h 127.0.0.1 -d sub2api `
  -c "select k.id,k.name,left(k.key,10)||'...'||right(k.key,6) as key,k.status,k.group_id,g.name as group_name,g.platform as group_platform from api_keys k left join groups g on g.id=k.group_id where k.deleted_at is null order by k.id;"
```

最终结果：

```text
 id | name |         key         | status | group_id | group_name | group_platform
----+------+---------------------+--------+----------+------------+----------------
  1 | 111  | sk-47c1fbe...413f69 | active |        2 | openai     | openai
```

## 9. 直接接口测试

### 9.1 查看模型

```powershell
$key='sk-47c1fbe...413f69'
$base='https://entirely-bible-etc-households.trycloudflare.com/v1'
Invoke-WebRequest -Uri "$base/models" -Headers @{Authorization="Bearer $key"} -UseBasicParsing
```

返回模型包含：

```text
gpt-5.4-mini
gpt-5.5
```

### 9.2 测试 gpt-5.4-mini

请求：

```powershell
$body = @{
  model = 'gpt-5.4-mini'
  messages = @(@{ role = 'user'; content = 'Return only 156' })
  stream = $false
} | ConvertTo-Json -Depth 8

Invoke-WebRequest `
  -Uri "$base/chat/completions" `
  -Method POST `
  -Headers @{Authorization="Bearer $key"; 'Content-Type'='application/json'} `
  -Body $body `
  -UseBasicParsing
```

结果：

```json
{
  "model": "gpt-5.4-mini",
  "choices": [
    {
      "message": {
        "role": "assistant",
        "content": "156"
      }
    }
  ]
}
```

### 9.3 测试 gpt-5.5

请求模型：

```text
gpt-5.5
```

结果：

```json
{
  "model": "gpt-5.5",
  "choices": [
    {
      "message": {
        "role": "assistant",
        "content": "156"
      }
    }
  ]
}
```

## 10. 如何证明 Cursor 走的是 Sub2API

最终用户问：如何证明 `gpt-5.5` 不是 Cursor Pro 自带的，而是 Sub2API？

最硬证据是 Sub2API 的 `usage_logs` 表里出现了 Cursor 的请求记录。

查询命令：

```powershell
& 'C:\Users\Administrator\Desktop\learn\sub2api-tools\pgsql\bin\psql.exe' `
  -U postgres -h 127.0.0.1 -d sub2api `
  -c "select u.id,u.request_id,u.api_key_id,k.name as key_name,u.account_id,a.name as account_name,u.group_id,g.name as group_name,g.platform as group_platform,u.model,u.requested_model,u.upstream_model,u.input_tokens,u.output_tokens,u.total_cost,u.user_agent,u.created_at from usage_logs u join api_keys k on k.id=u.api_key_id join accounts a on a.id=u.account_id left join groups g on g.id=u.group_id order by u.created_at desc limit 10;"
```

实际记录显示：

```text
user_agent: Cursor/1.0
api_key_id: 1
key_name: 111
account_id: 1
account_name: 11231
group_id: 2
group_name: openai
group_platform: openai
model: gpt-5.5
requested_model: gpt-5.5
```

部分成功调用记录：

```text
id: 12
request_id: client:f69e5b67-8116-4335-a7b5-bf5650091495
api_key_id: 1
key_name: 111
account_id: 1
account_name: 11231
group_id: 2
group_name: openai
group_platform: openai
model: gpt-5.5
requested_model: gpt-5.5
input_tokens: 677
output_tokens: 41
total_cost: 0.0176710000
user_agent: Cursor/1.0
created_at: 2026-06-06 20:11:17.867244+08
```

结论：

如果 Cursor 走的是 Cursor Pro 自带额度，Sub2API 数据库不会出现：

```text
Cursor/1.0
api_key_id = 1
key_name = 111
account_id = 1
account_name = 11231
group_name = openai
model = gpt-5.5
total_cost > 0
```

因此可以确认实际链路是：

```text
Cursor
  -> Cloudflare Tunnel
  -> Sub2API
  -> API Key 111
  -> OpenAI OAuth 账号 11231
  -> gpt-5.5
```

## 11. 当前最终可用配置

Cursor：

```text
OpenAI API Key:
sk-47c1fbe...413f69

Override OpenAI Base URL:
https://entirely-bible-etc-households.trycloudflare.com/v1

Model:
gpt-5.5
```

Sub2API：

```text
用户：admin
余额：$10.00
API Key：111
API Key 分组：openai
OpenAI 账号：11231
账号分组：openai
```

Tunnel：

```text
cloudflared PID: 8296
public URL: https://entirely-bible-etc-households.trycloudflare.com
```

## 12. 下次重启后的注意事项

1. 确认 Sub2API 后端仍在运行：

```powershell
Test-NetConnection 127.0.0.1 -Port 8080
```

2. 如果 Cloudflare Tunnel 失效，重新启动：

```powershell
cd C:\Users\Administrator\Desktop\learn\sub2api-main\sub2api-main
.\tools\cloudflared.exe tunnel --url http://127.0.0.1:8080
```

3. 复制新的 `https://xxxx.trycloudflare.com` 地址。

4. Cursor 里更新：

```text
Override OpenAI Base URL:
https://新的地址.trycloudflare.com/v1
```

5. 模型继续使用：

```text
gpt-5.5
```

或：

```text
gpt-5.4-mini
```

6. 如果再次出现 `No available accounts`，优先检查：

```text
API Key 分组是否为 openai
OpenAI 账号是否绑定 openai 分组
账号是否 active
账号是否 schedulable
用户余额是否大于 0
```

7. 如果 Cursor 又报本地私有网络禁止访问，说明 Base URL 填成了 localhost、127.0.0.1 或 192.168.x.x，需要重新使用 Cloudflare Tunnel 的 HTTPS 公网地址。

## 13. Fork 仓库维护建议

已 fork 源码到：

```text
https://github.com/592837154/sub2api
```

建议后续维护方式：

1. 将当前本地代码目录关联到 fork，方便记录自己的修复和配置说明：

```powershell
cd C:\Users\Administrator\Desktop\learn\sub2api-main\sub2api-main
git remote -v
git remote set-url origin https://github.com/592837154/sub2api.git
```

2. 如果还想保留上游仓库，增加 upstream：

```powershell
git remote add upstream https://github.com/原作者/sub2api.git
```

3. 后续自己的配置记录、Cursor 接入说明、Windows 本地运行脚本，可以提交到 fork：

```powershell
git status
git add .
git commit -m "docs: record Cursor gateway setup"
git push origin main
```

4. 注意不要把真实 API Key、OAuth token、数据库密码、Cloudflare Tunnel 临时地址长期作为敏感配置提交。文档中密钥应保持掩码，例如：

```text
sk-47c1fbe...413f69
```

5. 本次记录里的 Cloudflare Quick Tunnel 地址是临时地址，不适合作为 fork 里的固定生产配置。长期使用建议配置 Cloudflare Named Tunnel 和固定域名。
