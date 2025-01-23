# Docker 部署指南

本文档详细介绍如何使用 Docker 部署待办事项管理系统。

## 目录结构

```bash
deployments/docker/
├── Dockerfile # Docker镜像构建文件
├── docker-compose.yml # Docker编排配置文件
└── README.md # 部署说明文档
```

## 前置要求

在开始部署之前，请确保您的系统已经安装以下软件：

- Docker (版本 20.10.0 或更高)
- Docker Compose (版本 2.0.0 或更高)

可以使用以下命令检查版本：

```bash
docker --version
docker-compose --version
```

## 快速开始

### 1. 环境配置

1. 复制项目根目录下的环境变量示例文件：
```bash
cp .env.example .env
```

2. 根据实际情况修改 .env 文件中的配置：
   - 修改数据库连接信息
   - 修改 JWT 密钥（生产环境必须修改）
   - 根据需要调整日志配置

### 2. 构建和启动

1. 进入 Docker 部署目录：
```bash
cd deployments/docker
```

2. 构建 Docker 镜像：
```bash
docker-compose build
```

3. 启动服务：
```bash
docker-compose up -d
```

现在可以通过 `http://localhost:8080` 访问服务了。

### 3. 查看运行状态

- 查看容器状态：
```bash
docker-compose ps
```

- 查看应用日志：
```bash
docker-compose logs -f app
```

- 查看数据库日志：
```bash
docker-compose logs -f db
```

## 常用操作

### 停止服务
```bash
docker-compose down
```

### 重启服务
```bash
docker-compose restart
```

### 清理数据（慎用）
```bash
docker-compose down -v
```

## 生产环境部署注意事项

1. 安全配置
   - 修改默认的数据库密码
   - 设置复杂的 JWT 密钥
   - 关闭调试模式（SERVER_MODE=release）
   - 配置 HTTPS

2. 性能优化
   - 根据服务器配置调整数据库连接池
   - 配置合适的日志级别
   - 考虑使用 Redis 缓存

3. 数据备份
   - 定期备份数据库
   - 配置日志轮转
   - 建议使用数据卷持久化存储

## 常见问题

1. 端口冲突
   - 检查 8080 端口是否被占用
   - 可以在 .env 文件中修改 SERVER_PORT

2. 数据库连接失败
   - 确认数据库容器是否正常运行
   - 检查数据库连接配置是否正确
   - 等待数据库完全启动

3. 权限问题
   - 确保数据卷目录具有正确的权限
   - 检查日志目录的写入权限

## 故障排查

1. 查看容器日志
```bash
docker-compose logs -f [服务名]
```

2. 进入容器内部
```bash
docker-compose exec [服务名] sh
```

3. 检查网络连接
```bash
docker network ls
docker network inspect [网络名]
```

## 技术支持

如果遇到问题，可以：
1. 查看项目 [Issues](https://github.com/your-repo/issues)
2. 提交新的 Issue
3. 通过 Pull Request 贡献代码

## 参考文档

- [Docker 官方文档](https://docs.docker.com/)
- [Docker Compose 文档](https://docs.docker.com/compose/)
- [项目技术文档](../docs)
