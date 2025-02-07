# 使用 alpine 作为基础镜像
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

# 复制预编译的二进制文件和配置文件
COPY todo-api .
COPY configs/ ./configs/

# 声明容器将使用的端口
EXPOSE 8081

# 健康检查
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD curl -f http://localhost:8081/health || exit 1

# 容器启动时执行的命令
CMD ["./todo-api"] 