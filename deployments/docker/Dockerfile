# 使用 Alpine Linux 3.19 作为基础镜像，它体积小、安全性高
FROM alpine:3.19

# 设置容器内的工作目录
WORKDIR /app

# 配置系统环境：
# 1. 将 Alpine 软件源切换为阿里云镜像，加快国内访问速度
# 2. 安装必要的系统工具和运行时依赖
# 3. 配置时区为中国时区
# 4. 清理安装缓存，减小镜像体积
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories && \
    apk update && \
    apk add --no-cache tzdata curl gettext && \
    cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && \
    echo "Asia/Shanghai" > /etc/timezone && \
    rm -rf /var/cache/apk/* && \
    mkdir -p /app/logs

# 设置容器的环境变量，指定时区
ENV TZ=Asia/Shanghai

# 将编译好的应用程序和配置文件复制到容器中
COPY todo-api .
COPY configs/ ./configs/

# 声明容器对外暴露的端口号
EXPOSE 8081

# 配置容器健康检查
# interval: 检查间隔时间
# timeout: 检查超时时间
# start-period: 启动后多久开始检查
# retries: 连续失败多少次则认为容器不健康
HEALTHCHECK --interval=30s --timeout=10s --start-period=30s --retries=3 \
    CMD curl -f http://localhost:8081/health || exit 1

# 容器启动时执行的命令
CMD ["./todo-api"] 