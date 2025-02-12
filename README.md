# Todo Demo

一个基于 Go + Gin + GORM 的待办事项管理系统。

## 功能特性

### 核心功能
- 完整的用户认证系统 (JWT)
- 待办事项的 CRUD 操作
- 分类管理
- 提醒通知
- Swagger API 文档
- 日志管理
- 配置管理

### 技术特性
- RESTful API 设计
- 基于 JWT 的无状态认证
- 完善的数据验证
- 统一的错误处理
- 详细的 API 文档
- 完整的测试覆盖
- 容器化部署支持
- 云原生支持

### 性能特性
- 多级缓存架构
- 连接池管理
- 请求限流
- 异步处理
- 批量操作

### 安全特性
- 密码加密存储
- JWT Token 认证
- 请求频率限制
- SQL 注入防护
- XSS 防护
- CSRF 防护

## 系统要求

- Go 1.21+
- MySQL 8.0+
- Redis 6.0+ (可选,用于缓存)
- Docker 20.10+ (可选)
- Kubernetes 1.20+ (可选)

## 快速开始

1. 克隆项目
```bash
git clone https://github.com/ma-x-x/todo.git
cd todo
```

2. 配置环境
```bash
# 复制环境变量示例文件
cp .env.example .env
# 编辑配置文件
vim .env
```

3. 初始化数据库
```bash
# 确保 MySQL 服务已启动
# 使用 root 用户执行初始化脚本
mysql -u root -p < scripts/init.sql

# 或者使用 make 命令初始化
make init-db
```

4. 安装依赖
```bash
make deps
```

5. 启动开发服务器
```bash
make dev
```

6. 访问服务
# API 服务
http://localhost:8080

# API 文档
http://localhost:8080/swagger/index.html

# 测试 API
# 1. 先注册用户
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"username":"test","password":"test123","email":"test@example.com"}'

# 2. 登录获取 token
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"test","password":"test123"}'

# 3. 使用 token 创建待办事项
curl -X POST http://localhost:8080/api/v1/todos \
  -H "Authorization: Bearer <your-token>" \
  -H "Content-Type: application/json" \
  -d '{"title":"测试待办事项","description":"这是一个测试","priority":2}'
```

## 项目结构

项目采用标准的 Go 项目布局，主要目录说明如下：

```
todo
├── api                     # API 层，处理 HTTP 请求和响应
│   └── v1                 # API 版本控制
│       ├── dto            # 数据传输对象，定义请求和响应结构
│       │   ├── auth       # 认证相关 DTO
│       │   ├── category   # 分类相关 DTO
│       │   ├── reminder   # 提醒相关 DTO
│       │   └── todo       # 待办事项相关 DTO
│       └── handlers       # HTTP 请求处理器
├── cmd                    # 应用程序入口
│   └── server            # API 服务器
│       └── main.go       # 主程序入口点
├── configs               # 配置文件目录
│   ├── config.dev.yaml   # 开发环境配置
│   ├── config.prod.yaml  # 生产环境配置
│   └── config.yaml       # 基础配置
├── deployments           # 部署相关配置和脚本
│   ├── docker           # Docker 容器化配置
│   ├── kubernetes       # Kubernetes 编排配置
│   └── scripts          # 部署自动化脚本
├── docs                 # 项目文档
│   ├── docs.go          # Swagger 自动生成的文档
│   ├── swagger.json     # Swagger API 定义
│   └── swagger.yaml     # Swagger API 配置
├── internal             # 私有应用代码
│   ├── middleware      # HTTP 中间件
│   ├── models         # 数据模型定义
│   ├── repository     # 数据访问层实现
│   ├── router         # 路由注册和管理
│   └── service        # 业务逻辑层
├── pkg                 # 可复用的公共代码包
│   ├── cache          # 缓存实现
│   ├── config         # 配置管理
│   ├── database       # 数据库连接和管理
│   ├── errors         # 错误处理
│   ├── jwt            # JWT 认证
│   ├── logger         # 日志处理
│   ├── monitor        # 监控指标
│   ├── queue          # 消息队列
│   ├── response       # HTTP 响应处理
│   └── utils          # 通用工具函数
├── scripts            # 维护脚本和工具
│   └── init.sql      # 数据库初始化脚本
├── build             # 编译构建产物
├── logs              # 应用日志文件
├── tmp               # 临时文件
├── Makefile          # 项目管理命令
├── README.md         # 项目说明文档
└── go.mod            # Go 模块依赖定义
```

目录结构详细说明：

1. **api/** - API 层
   - 处理 HTTP 请求和响应
   - 实现 RESTful API 接口
   - 请求参数验证和响应格式化
   - 版本化 API 管理

2. **internal/** - 内部应用代码
   - models: 定义核心数据结构和业务实体
   - repository: 实现数据持久化和访问逻辑
   - service: 封装核心业务逻辑
   - middleware: 处理横切关注点（认证、日志等）
   - router: 管理 API 路由和处理器映射

3. **pkg/** - 公共代码包
   - cache: 实现多级缓存策略
   - config: 处理配置加载和管理
   - database: 数据库连接池和操作封装
   - errors: 统一错误处理机制
   - jwt: 用户认证和授权
   - logger: 结构化日志记录
   - monitor: 性能监控和指标收集
   - queue: 异步任务处理
   - response: 统一响应格式
   - utils: 通用辅助函数

4. **configs/** - 配置文件
   - 支持多环境配置
   - 敏感配置分离
   - 配置热重载

5. **deployments/** - 部署配置
   - docker: 容器化部署配置
   - kubernetes: 容器编排配置
   - scripts: 自动化部署脚本

6. **docs/** - 项目文档
   - API 文档 (Swagger)
   - 架构设计文档
   - 开发指南

7. **scripts/** - 工具脚本
   - 数据库初始化
   - 维护工具
   - 自动化任务

主要业务模块：

1. 用户管理模块
   - 用户注册
   - 用户登录
   - 用户认证
   - 用户信息管理

2. 待办事项模块
   - 创建待办事项
   - 查询待办事项
   - 更新待办事项
   - 删除待办事项

3. 分类管理模块
   - 创建分类
   - 修改分类
   - 删除分类
   - 按分类查询

4. 提醒管理模块
   - 设置提醒
   - 提醒通知
   - 重复提醒
   - 提醒方式选择

性能优化设计：

1. 缓存层设计
   - 多级缓存架构
   - 缓存预热和更新
   - 防止缓存穿透

2. 数据库优化
   - 连接池管理
   - 读写分离支持
   - 索引优化

3. 并发处理
   - 异步任务队列
   - 分布式锁
   - 限流降级

4. 监控指标
   - 性能指标采集
   - 资源使用监控
   - 业务指标统计

## 开发命令

- `make deps`: 安装项目依赖
- `make dev`: 启动开发服务器
- `make build`: 构建应用
- `make test`: 运行测试
- `make test-coverage`: 生成测试覆盖率报告
- `make lint`: 运行代码检查
- `make swagger`: 生成 API 文档
- `make docker-build`: 构建 Docker 镜像
- `make docker-run`: 运行 Docker 容器
- `make clean`: 清理构建文件

## 部署方式

### Docker 部署

```bash
# 构建镜像
make docker-build

# 运行容器
make docker-run
```

### Kubernetes 部署

```bash
# 部署应用
cd deployments/kubernetes
./deploy.sh

# 查看状态
kubectl get pods -n todo-app
```

## 监控

- Prometheus 指标采集
- Grafana 监控面板
- ELK 日志聚合
- Jaeger 链路追踪

## 开发指南

详细的开发指南请查看 [START.md](START.md)。
设计文档请查看 [DESIGN.md](DESIGN.md)。

## 贡献指南

1. Fork 项目
2. 创建功能分支 (`git checkout -b feature/xxx`)
3. 提交更改 (`git commit -am 'feat: add xxx'`)
4. 推送分支 (`git push origin feature/xxx`)
5. 创建 Pull Request

## 版本历史

- v1.0.0 - 基础功能实现
- v1.1.0 - 添加缓存支持
- v1.2.0 - 添加监控功能

## 许可证

MIT License


