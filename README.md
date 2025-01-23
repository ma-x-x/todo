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

3. 安装依赖
```bash
make deps
```

4. 启动开发服务器
```bash
make dev
```

5. 访问服务
- API 服务: http://localhost:8080
- API 文档: http://localhost:8080/swagger/index.html

## 项目结构

项目采用标准的 Go 项目布局，主要目录说明如下：

```
todo-demo
├── api                 # API 接口定义
│   └── v1             # API 版本 1
│       ├── dto        # 数据传输对象
│       │   ├── auth        # 认证相关 DTO
│       │   │   ├── login.go      # 登录请求/响应
│       │   │   └── register.go   # 注册请求/响应
│       │   ├── category    # 分类相关 DTO
│       │   │   ├── create.go     # 创建分类
│       │   │   ├── list.go       # 分类列表
│       │   │   └── update.go     # 更新分类
│       │   ├── reminder    # 提醒相关 DTO
│       │   │   ├── create.go     # 创建提醒
│       │   │   └── update.go     # 更新提醒
│       │   └── todo        # 待办事项 DTO
│       │       ├── create.go     # 创建待办
│       │       ├── list.go       # 待办列表
│       │       └── update.go     # 更新待办
│       ├── handlers        # HTTP 处理器
│       │   ├── auth.go          # 认证处理
│       │   ├── category.go      # 分类处理
│       │   ├── health.go        # 健康检查
│       │   ├── reminder.go      # 提醒处理
│       │   └── todo.go          # 待办处理
│       └── routes          # 路由定义
│           └── routes.go        # API 路由注册
├── cmd                 # 主程序入口
│   └── server         # 服务器程序
│       └── main.go         # 主程序入口点
├── configs            # 配置文件目录
│   └── config.yaml        # 应用配置文件
├── internal          # 内部代码包
│   ├── middleware    # 中间件
│   │   ├── auth.go        # 认证中间件
│   │   ├── cors.go        # 跨域中间件
│   │   └── logger.go      # 日志中间件
│   ├── models       # 数据模型
│   │   ├── base.go        # 基础模型
│   │   ├── category.go    # 分类模型
│   │   ├── reminder.go    # 提醒模型
│   │   ├── todo.go        # 待办模型
│   │   └── user.go        # 用户模型
│   ├── repository   # 数据访问层
│   │   ├── db            # 数据库实现
│   │   │   ├── category.go   # 分类存储
│   │   │   ├── reminder.go   # 提醒存储
│   │   │   ├── todo.go       # 待办存储
│   │   │   └── user.go       # 用户存储
│   │   ├── category.go   # 分类仓储接口
│   │   ├── reminder.go   # 提醒仓储接口
│   │   ├── todo.go       # 待办仓储接口
│   │   └── user.go       # 用户仓储接口
│   ├── router       # 路由配置
│   │   └── router.go     # 主路由配置
│   └── service      # 业务逻辑层
│       ├── impl          # 接口实现
│       │   ├── auth.go       # 认证服务实现
│       │   ├── auth_test.go  # 认证测试
│       │   ├── category.go   # 分类服务实现
│       │   ├── reminder.go   # 提醒服务实现
│       │   ├── todo.go       # 待办服务实现
│       │   └── todo_test.go  # 待办测试
│       ├── auth.go       # 认证服务接口
│       ├── category.go   # 分类服务接口
│       ├── reminder.go   # 提醒服务接口
│       └── todo.go       # 待办服务接口
├── pkg              # 公共代码包
│   ├── cache        # 缓存组件
│   │   └── redis.go      # Redis 实现
│   ├── config      # 配置管理
│   │   └── config.go     # 配置加载
│   ├── database    # 数据库组件
│   │   └── mysql.go      # MySQL 连接
│   ├── db          # 数据库工具
│   │   ├── db.go         # 数据库操作
│   │   └── pool.go       # 连接池管理
│   ├── errors      # 错误处理
│   │   └── errors.go     # 错误定义
│   ├── lock        # 分布式锁
│   │   └── distributed_lock.go  # 锁实现
│   ├── logger      # 日志组件
│   │   └── logger.go     # 日志实现
│   ├── middleware  # 通用中间件
│   │   ├── cache.go      # 缓存中间件
│   │   ├── db.go         # 数据库中间件
│   │   ├── performance.go # 性能监控
│   │   └── ratelimit.go  # 限流中间件
│   ├── monitor     # 监控组件
│   │   └── prometheus.go # 指标收集
│   ├── queue       # 队列组件
│   │   └── task_queue.go # 任务队列
│   └── utils       # 工具函数
│       └── jwt.go        # JWT 工具
├── deployments     # 部署配置
│   ├── docker          # Docker 部署
│   │   ├── Dockerfile       # 容器构建
│   │   └── docker-compose.yml  # 容器编排
│   └── kubernetes     # K8s 部署
│       ├── configmap.yaml    # 配置映射
│       ├── deploy.sh         # 部署脚本
│       ├── deployment.yaml   # 部署配置
│       ├── grafana.yaml      # 监控面板
│       ├── ingress.yaml      # 入口配置
│       ├── mysql.yaml        # 数据库配置
│       ├── prometheus.yaml   # 监控配置
│       ├── redis-config.yaml # Redis配置
│       ├── redis.yaml        # Redis部署
│       ├── secret.yaml       # 密钥配置
│       └── service.yaml      # 服务配置
├── docs            # 文档目录
│   ├── docs.go          # Swagger 文档
│   ├── swagger.json     # API 文档(JSON)
│   └── swagger.yaml     # API 文档(YAML)
├── logs            # 日志目录
│   └── app.log          # 应用日志
├── DESIGN.md       # 设计文档
├── Makefile        # 构建脚本
├── README.md       # 项目说明
├── START.md        # 开发指南
├── go.mod          # Go 模块定义
├── go.sum          # 依赖版本锁定
└── tmp             # 临时文件目录
```

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


