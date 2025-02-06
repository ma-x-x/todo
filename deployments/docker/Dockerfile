# 第一阶段：构建阶段
# 使用golang:1.21-alpine作为基础镜像
# 选择alpine版本是因为它体积小（约5MB），而完整的golang镜像超过1GB
FROM golang:1.22-alpine AS builder

# 设置容器内的工作目录为/app
# 后续的COPY和RUN命令都将基于此目录执行
WORKDIR /app

# 首先仅复制依赖配置文件
# 这样做是为了利用Docker的缓存机制：
# - 如果go.mod和go.sum没有变化，将使用缓存的依赖
# - 避免源代码改动导致依赖层重新构建
COPY go.mod go.sum ./
RUN go mod download

# 复制所有源代码到容器中
# 这一步放在依赖安装后面，这样源码改动不会影响依赖缓存
COPY . .

# 构建Go应用
# CGO_ENABLED=0: 禁用CGO，生成静态链接的二进制文件，不依赖外部动态库
# GOOS=linux: 明确指定目标系统为Linux，确保二进制文件在Linux环境运行
# -o 参数指定输出文件的路径和名称
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/build/todo-api cmd/server/main.go

# 第二阶段：运行阶段
# 使用alpine:3.19作为基础镜像
# 这是一个极小的Linux发行版，仅包含必要的运行时环境
FROM alpine:3.19

# 设置工作目录
WORKDIR /app

# 从构建阶段复制文件到运行阶段
# 只复制必要的文件：配置文件和编译好的二进制文件
# 这样可以大大减少最终镜像的体积
COPY --from=builder /app/configs/config.yaml /app/configs/
COPY --from=builder /app/build/todo-api /app/

# 声明容器将使用的端口
# 这只是一个文档说明，实际上不会自动打开端口
EXPOSE 8080

# 容器启动时执行的命令
# 使用数组格式可以避免shell解释，直接执行二进制文件
CMD ["/app/todo-api"] 