#!/bin/bash

# 确保脚本在出错时退出
set -e

# 设置环境变量
export MYSQL_ROOT_PASSWORD=${MYSQL_ROOT_PASSWORD:-root}
export DB_PASSWORD=${DB_PASSWORD:-root}

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

# 构建并启动新容器
docker-compose up -d --build

# 等待服务启动
echo "Waiting for services to start..."
sleep 10

# 检查服务状态
echo "Checking service status..."
docker-compose ps

# 检查应用日志
echo "Checking application logs..."
docker-compose logs --tail=50 app

# 等待 MySQL 启动
echo "Waiting for MySQL to start..."
sleep 20

# 创建应用数据库用户
docker exec todo-api_mysql_1 mysql -uroot -p"${MYSQL_ROOT_PASSWORD}" -e "
CREATE DATABASE IF NOT EXISTS todo_db;
CREATE USER IF NOT EXISTS 'todo_user'@'%' IDENTIFIED BY '${DB_PASSWORD}';
GRANT ALL PRIVILEGES ON todo_db.* TO 'todo_user'@'%';
FLUSH PRIVILEGES;
"

echo "Deployment completed successfully!" 