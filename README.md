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
├── api                 # API 接口定义和处理
│   └── v1
│       ├── dto        # 数据传输对象(Data Transfer Objects)
│       │   ├── auth   # 认证相关DTO
│       │   ├── category   # 分类相关DTO
│       │   ├── reminder   # 提醒相关DTO
│       │   └── todo       # 待办事项相关DTO
│       └── handlers       # HTTP请求处理器
├── cmd                 # 主要应用程序入口
│   └── server         # 服务器启动入口
├── configs            # 配置文件目录
├── deployments        # 部署相关配置
│   ├── docker        # Docker部署配置
│   └── kubernetes    # Kubernetes部署配置
├── docs              # 文档和API文档(Swagger)
├── internal          # 私有应用程序代码
│   ├── middleware    # HTTP中间件
│   ├── models        # 数据模型定义
│   ├── repository    # 数据访问层
│   │   ├── db       # 具体数据库操作实现
│   │   └── ...      # 仓储接口定义
│   ├── router       # 路由配置
│   └── service      # 业务逻辑层
│       ├── impl     # 接口实现
│       └── ...      # 服务接口定义
├── pkg              # 可重用的库代码
│   ├── cache        # 缓存组件
│   ├── config       # 配置管理
│   ├── database     # 数据库连接管理
│   ├── db           # 数据库工具
│   ├── errors       # 错误处理
│   ├── lock         # 分布式锁
│   ├── logger       # 日志组件
│   ├── middleware   # 通用中间件
│   ├── monitor      # 监控组件
│   ├── queue        # 队列组件
│   └── utils        # 通用工具函数
├── logs             # 日志文件目录
├── tmp              # 临时文件
├── Makefile         # 项目管理命令
├── go.mod           # Go模块定义
└── README.md        # 项目说明文档
```

目录结构说明：

1. **api/** - API层
   - 处理HTTP请求响应
   - 数据验证和转换
   - API文档定义

2. **internal/** - 内部应用代码
   - models: 核心数据模型
   - repository: 数据访问层，处理数据持久化
   - service: 业务逻辑层，实现核心功能
   - middleware: 请求处理中间件
   - router: 路由配置和管理

3. **pkg/** - 公共代码包
   - 可被外部项目引用的通用组件
   - 基础设施代码
   - 工具函数和助手方法

4. **configs/** - 配置文件
   - 应用配置
   - 环境变量
   - 部署配置

5. **deployments/** - 部署配置
   - Docker容器化配置
   - Kubernetes编排配置
   - 部署脚本和说明

6. **docs/** - 文档
   - API文档(Swagger)
   - 设计文档
   - 开发指南

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


