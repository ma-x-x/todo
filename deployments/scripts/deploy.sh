#!/bin/bash

# 确保脚本在出错时退出
set -e

# 设置环境变量
export MYSQL_ROOT_PASSWORD=${MYSQL_ROOT_PASSWORD:-root}
export DB_PASSWORD=${DB_PASSWORD:-root}
export DB_HOST=mysql
export REDIS_HOST=redis

# 检查并清理端口占用
check_and_free_port() {
    local port=$1
    if lsof -i :$port > /dev/null; then
        echo "Port $port is in use. Attempting to free it..."
        # 尝试停止使用该端口的进程
        fuser -k $port/tcp || true
        sleep 2
    fi
}

# 检查关键端口
check_and_free_port 8081  # 应用端口
check_and_free_port 3306  # MySQL 端口
check_and_free_port 6379  # Redis 端口

# 创建必要的目录
mkdir -p /data/mysql/conf.d
mkdir -p /data/mysql/logs
mkdir -p /data/redis
mkdir -p logs

# 如果 redis.conf 不存在，创建一个基本配置
if [ ! -f "/data/redis/redis.conf" ]; then
    cat > /data/redis/redis.conf << EOF
bind 0.0.0.0
protected-mode yes
port 6379
tcp-backlog 511
timeout 0
tcp-keepalive 300
daemonize no
supervised no
pidfile /var/run/redis_6379.pid
loglevel notice
logfile ""
databases 16
always-show-logo yes
save 900 1
save 300 10
save 60 10000
stop-writes-on-bgsave-error yes
rdbcompression yes
rdbchecksum yes
dbfilename dump.rdb
dir /data
replica-serve-stale-data yes
replica-read-only yes
repl-diskless-sync no
repl-diskless-sync-delay 5
repl-disable-tcp-nodelay no
replica-priority 100
maxmemory 2gb
maxmemory-policy allkeys-lru
lazyfree-lazy-eviction no
lazyfree-lazy-expire no
lazyfree-lazy-server-del no
replica-lazy-flush no
appendonly yes
appendfilename "appendonly.aof"
appendfsync everysec
no-appendfsync-on-rewrite no
auto-aof-rewrite-percentage 100
auto-aof-rewrite-min-size 64mb
aof-load-truncated yes
aof-use-rdb-preamble yes
EOF
fi

# 等待 MySQL 就绪的函数改进
wait_for_mysql() {
    echo "Waiting for MySQL to be ready..."
    for i in {1..60}; do
        if docker-compose exec mysql mysqladmin ping -h mysql -u"root" -p"${MYSQL_ROOT_PASSWORD}" --silent > /dev/null 2>&1; then
            echo "MySQL is ready!"
            return 0
        fi
        echo "Waiting for MySQL to be ready... ($i/60)"
        sleep 2
    done
    echo "MySQL did not become ready in time"
    return 1
}

# 等待 Redis 就绪的函数
wait_for_redis() {
    echo "Waiting for Redis to be ready..."
    for i in {1..30}; do
        if docker-compose exec redis redis-cli -h redis ping > /dev/null 2>&1; then
            echo "Redis is ready!"
            return 0
        fi
        echo "Waiting for Redis to be ready... ($i/30)"
        sleep 2
    done
    echo "Redis did not become ready in time"
    return 1
}

# 停止并删除旧容器
docker-compose down --volumes --remove-orphans || true

# 清理旧的数据卷（可选，谨慎使用）
# docker volume prune -f

# 创建 docker 网络（如果不存在）
docker network create todo-network || true

# 启动数据库和 Redis
echo "Starting MySQL and Redis..."
docker-compose up -d mysql redis

# 等待服务就绪
wait_for_mysql || exit 1
wait_for_redis || exit 1

# 初始化数据库
echo "Initializing database..."
docker-compose exec -T mysql mysql -uroot -p"${MYSQL_ROOT_PASSWORD}" -e "
CREATE DATABASE IF NOT EXISTS todo_db;
CREATE USER IF NOT EXISTS 'todo_user'@'%' IDENTIFIED BY '${DB_PASSWORD}';
GRANT ALL PRIVILEGES ON todo_db.* TO 'todo_user'@'%';
FLUSH PRIVILEGES;
"

# 导入初始化 SQL
echo "Importing database schema..."
docker-compose exec -T mysql mysql -uroot -p"${MYSQL_ROOT_PASSWORD}" todo_db < scripts/init.sql

# 启动应用
echo "Starting application..."
docker-compose up -d app

# 等待应用就绪
echo "Waiting for application to be ready..."
for i in {1..30}; do
    if curl -s http://localhost:8081/health > /dev/null; then
        echo "Application is ready!"
        break
    fi
    echo "Waiting for application to be ready... ($i/30)"
    sleep 2
done

# 检查服务状态
echo "Checking service status..."
docker-compose ps

# 检查应用日志
echo "Checking application logs..."
docker-compose logs --tail=50 app

echo "Deployment completed successfully!" 