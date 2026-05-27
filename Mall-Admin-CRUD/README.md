# Mall Admin CRUD | 商品管理后台系统

## 项目定位

这是一个从零搭建的前后端分离后台管理系统案例，核心业务是“商品管理”的标准 CRUD。

它不是简单的表格 Demo，而是完整跑通了：

- 前端管理页面
- 后端 REST API
- 数据库持久化
- Docker 本地部署
- TiDB Cloud 云数据库
- Hugging Face Spaces 公网部署

这个项目适合作为作品集中的“全栈 CRUD 标准样板”：展示我如何把一个常见后台业务，从页面组件、接口规范、数据库建模、容器化部署，一直推进到公网可访问。

## 访问入口

### 在线演示

- Hugging Face Space: <https://huggingface.co/spaces/zhu592837154/mall-admin>
- Web App: <https://zhu592837154-mall-admin.hf.space>
- 商品管理页面: <https://zhu592837154-mall-admin.hf.space/product-manage>

### 代码仓库

- GitHub: <https://github.com/592837154/scaffold-mall-admin>
- Hugging Face Space Repository: <https://huggingface.co/spaces/zhu592837154/mall-admin/tree/main>

### 本地访问

本地 Docker 启动后访问：

```text
http://localhost:8000/product-manage
```

后端接口测试地址：

```text
http://localhost:8080/api/goods/list?current=1&pageSize=10
```

## 技术栈

### 前端

- React
- TypeScript
- Umi 4
- Ant Design Pro 5
- `@ant-design/pro-components`
- `ProTable`
- `ModalForm`
- Ant Design `message` / `Popconfirm` / `Badge`
- Nginx 静态资源服务与 `/api` 反向代理

### 后端

- Go
- Gin
- GORM
- MySQL 协议数据库
- TiDB Cloud
- RESTful API
- CORS 跨域处理

### 部署与工程化

- Docker
- Docker Compose
- Hugging Face Spaces Docker Runtime
- TiDB Cloud Serverless / Starter
- Git remote 多仓库推送
- Hugging Face Secrets / Variables
- Nginx 单容器前后端整合

## 核心功能

### 商品列表查询

- 使用 `ProTable` 的 `request` 属性请求后端。
- 支持分页参数：
  - `current`
  - `pageSize`
- 支持搜索：
  - 商品名称模糊搜索
  - 商品状态筛选
  - 创建时间范围筛选

后端将 Ant Design Pro 的分页参数转换为 GORM 查询：

```go
offset := (current - 1) * pageSize
db.Offset(offset).Limit(pageSize)
```

### 新建商品

- 前端通过 `ModalForm` 弹出表单。
- 后端 `POST /api/goods/create` 接收 JSON 并写入数据库。
- 成功后刷新 `ProTable`。

### 编辑商品

- 操作列点击“编辑”。
- 表单回显当前行数据。
- 后端 `PUT /api/goods/update` 根据 ID 更新。

### 删除商品

- 操作列使用 `Popconfirm` 二次确认。
- 后端 `DELETE /api/goods/delete?id=xxx` 删除数据。
- 成功后刷新列表。

## 接口规范

为了适配 Ant Design Pro 的数据消费方式，后端统一响应结构：

```json
{
  "success": true,
  "data": {}
}
```

列表接口返回：

```json
{
  "success": true,
  "data": {
    "list": [],
    "total": 0
  }
}
```

核心接口：

```text
GET    /api/goods/list
POST   /api/goods/create
PUT    /api/goods/update
DELETE /api/goods/delete
```

## 主要难点

### 1. ProTable 与后端分页协议对齐

`ProTable` 默认传递的是：

```text
current
pageSize
```

而数据库分页需要的是：

```text
offset
limit
```

这里需要在后端显式转换，否则前端分页看起来正常，实际查询会错位。

### 2. 统一响应结构

Ant Design Pro 的表格组件通常依赖：

```text
success
data
total
```

如果后端直接返回数组，前端表格不会自然识别分页总数。因此后端需要专门为前端组件设计响应格式。

### 3. 跨域与反向代理

本地开发时前端和后端端口不同，需要 CORS。

Docker 部署时，为了让前端调用更稳定，Nginx 将：

```text
/api/*
```

代理到后端服务。

### 4. Docker Compose 到 Hugging Face 的部署模型差异

本地可以使用多容器：

```text
frontend
goods-api
mysql
```

但 Hugging Face Spaces 的 Docker Runtime 更适合单容器 HTTP 服务，而且默认只暴露一个端口。

最终方案是：

```text
同一个容器内：
Nginx 监听 7860
Go 后端监听 8080
Nginx 代理 /api 到 127.0.0.1:8080
```

### 5. 本地与线上数据库一致

为了避免“本地能跑、线上数据不一致”的问题，最终统一使用 TiDB Cloud。

本地 Docker 和 Hugging Face 都连接同一个 TiDB 实例：

```text
MYSQL_HOST
MYSQL_PORT
MYSQL_USER
MYSQL_PASSWORD
MYSQL_DATABASE
MYSQL_TLS
```

密码放在本地 `.env` 或 Hugging Face Secrets 中，不写入仓库。

### 6. Hugging Face 构建限制

部署过程中遇到几个真实问题：

- Hugging Face Git 普通文件不能超过 10 MiB。
- Go 编译产物 `backend/server` 超过限制，不能直接推送。
- 解决方案：不要提交二进制文件，在 Docker 构建阶段编译 Go 后端。

### 7. Node / pnpm 版本问题

构建时 Corepack 默认拉取了较新的 pnpm，引发版本和供应链策略检查问题。

解决方案：

```dockerfile
RUN corepack enable && corepack prepare pnpm@9.15.9 --activate
```

并使用 Node 22 作为前端构建环境。

### 8. Hugging Face Space 元信息

Docker Space 需要在 `README.md` 顶部声明配置：

```yaml
---
title: Mall Admin
emoji: 🛒
colorFrom: blue
colorTo: green
sdk: docker
app_port: 7860
pinned: false
---
```

否则会出现：

```text
Configuration error
Missing configuration in README
```

## 调研与推进路线

### 第一阶段：确定后台页面标准形态

调研 Ant Design Pro 的典型 CRUD 写法，明确：

- 页面使用 `ProTable`。
- 新建/编辑使用 `ModalForm`。
- 删除使用 `Popconfirm`。
- 表格请求统一走 `request`。

目标是让页面符合真实后台系统常见交互，而不是只写一个静态表格。

### 第二阶段：先做前端模拟数据

先用 `setTimeout` 模拟接口延迟，验证：

- 搜索项是否正常
- 分页是否正常
- 新建/编辑弹窗是否可用
- 删除刷新是否符合预期

这一阶段重点是把交互闭环跑通。

### 第三阶段：补后端 API

使用 Gin + GORM 实现真实接口：

- 商品模型定义
- 数据库迁移
- 初始化测试数据
- 分页查询
- 条件过滤
- 增删改接口
- CORS
- 统一响应结构

### 第四阶段：本地 Docker 化

将前端、后端、数据库放进 Docker Compose，解决：

- 服务启动顺序
- Nginx 代理
- 后端环境变量
- 数据库连接
- 本地端口映射

### 第五阶段：迁移到 TiDB Cloud

为了后续公网部署，将本地 MySQL 替换为 TiDB Cloud。

调研重点：

- TiDB Cloud 是否兼容 MySQL 协议
- GORM MySQL Driver 是否能直连
- TLS 参数如何配置
- Hugging Face 如何注入数据库密码

### 第六阶段：部署到 Hugging Face Spaces

Hugging Face 不直接运行 Docker Compose，因此改为单容器部署：

- Dockerfile 多阶段构建前端
- Dockerfile 多阶段构建 Go 后端
- Nginx 统一监听 7860
- `start-hf.sh` 同时启动 Go 后端和 Nginx

### 第七阶段：处理真实构建问题

逐步解决：

- Git 推送认证问题
- Git remote 切换问题
- 10 MiB 文件限制
- README Space 配置缺失
- Node / pnpm 版本冲突
- pnpm lockfile 策略检查

这一阶段是项目最有价值的工程化部分，因为它暴露了从“本地能跑”到“公网能访问”的真实差距。

## 项目收获

- 理解了 Ant Design Pro CRUD 页面与后端接口之间的契约。
- 掌握了 `ProTable` 的分页、搜索、刷新和操作列联动方式。
- 熟悉了 Gin + GORM 的标准 CRUD 写法。
- 理解了 Docker Compose 多容器部署与 Hugging Face 单容器部署的差异。
- 实践了 TiDB Cloud 作为 MySQL 兼容云数据库的接入方式。
- 处理了真实部署中的 Git、构建、依赖版本、平台限制问题。

## 后续优化方向

- 增加登录鉴权与用户权限。
- 增加商品图片上传。
- 增加分类管理页面。
- 增加库存、销量、上下架批量操作。
- 后端增加参数校验与错误码体系。
- 前端抽象统一请求层。
- 使用 CI 自动构建并部署到 Hugging Face。
- 为 TiDB TLS 配置正式 CA 校验，替代 `skip-verify`。

## 展示价值

这个项目可以在作品集中表达三件事：

1. 我能实现标准后台 CRUD 业务。
2. 我能把前端组件规范和后端接口规范对齐。
3. 我能处理从本地开发到公网部署之间的工程问题。

如果只看代码，它是一个商品管理系统；如果看完整过程，它是一个小型全栈项目从 0 到上线的实战记录。
