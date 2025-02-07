# 第一阶段：构建阶段
FROM golang:1.22-alpine AS builder

# 设置工作目录
WORKDIR /app

# 设置 Go 环境变量
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

# 设置 Alpine 镜像源为阿里云镜像（构建阶段也需要）
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories

# 首先仅复制依赖配置文件
COPY go.mod go.sum ./
RUN go mod download

# 复制所有源代码到容器中
COPY . .

# 构建Go应用
RUN go build -ldflags="-s -w" -o /app/build/todo-api cmd/server/main.go

# 第二阶段：运行阶段
FROM alpine:3.19

# 设置工作目录
WORKDIR /app

# 设置 Alpine 镜像源为阿里云镜像
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories && \
    # 安装必要的运行时依赖
    apk update && \
    apk add --no-cache tzdata curl && \
    # 设置时区
    cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && \
    echo "Asia/Shanghai" > /etc/timezone && \
    # 清理缓存
    rm -rf /var/cache/apk/* && \
    # 创建日志目录
    mkdir -p /app/logs

# 设置环境变量
ENV TZ=Asia/Shanghai

# 从构建阶段复制编译好的二进制文件和配置文件
COPY --from=builder /app/build/todo-api .
COPY configs/ ./configs/

# 声明容器将使用的端口
EXPOSE 8081

# 健康检查
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD curl -f http://localhost:8081/health || exit 1

# 容器启动时执行的命令
CMD ["./todo-api"] 