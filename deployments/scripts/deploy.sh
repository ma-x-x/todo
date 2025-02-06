#!/bin/bash

# 确保脚本在出错时退出
set -e

# 设置环境变量
export MYSQL_ROOT_PASSWORD=${MYSQL_ROOT_PASSWORD:-root}
export DB_PASSWORD=${DB_PASSWORD:-root}

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

# 停止并删除旧容器
docker-compose down || true

# 等待 MySQL 就绪
wait_for_mysql() {
    echo "Waiting for MySQL to be ready..."
    for i in {1..30}; do
        # 获取 MySQL 容器 ID
        MYSQL_CONTAINER=$(docker-compose ps -q mysql)
        if [ -z "$MYSQL_CONTAINER" ]; then
            echo "MySQL container not found"
            sleep 2
            continue
        fi
        
        if docker exec $MYSQL_CONTAINER mysqladmin ping -h localhost -u"root" -p"${MYSQL_ROOT_PASSWORD}" --silent; then
            echo "MySQL is ready!"
            return 0
        fi
        echo "Waiting for MySQL to be ready... ($i/30)"
        sleep 2
    done
    echo "MySQL did not become ready in time"
    return 1
}

# 构建并启动新容器
docker-compose up -d mysql redis
wait_for_mysql

# 启动应用
docker-compose up -d app

# 检查服务状态
echo "Checking service status..."
docker-compose ps

# 检查应用日志
echo "Checking application logs..."
docker-compose logs --tail=50 app

# 创建应用数据库用户
MYSQL_CONTAINER=$(docker-compose ps -q mysql)
docker exec $MYSQL_CONTAINER mysql -uroot -p"${MYSQL_ROOT_PASSWORD}" -e "
CREATE DATABASE IF NOT EXISTS todo_db;
CREATE USER IF NOT EXISTS 'todo_user'@'%' IDENTIFIED BY '${DB_PASSWORD}';
GRANT ALL PRIVILEGES ON todo_db.* TO 'todo_user'@'%';
FLUSH PRIVILEGES;
"

echo "Deployment completed successfully!" 