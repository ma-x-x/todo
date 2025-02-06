#!/bin/bash

# 设置备份目录
BACKUP_DIR="/data/backups"
DATE=$(date +%Y%m%d_%H%M%S)

# 创建备份目录
mkdir -p "$BACKUP_DIR"

# MySQL 备份
echo "Backing up MySQL database..."
docker exec todo-api_mysql_1 mysqldump -u root -proot todo_db > "$BACKUP_DIR/mysql_$DATE.sql"

# Redis 备份
echo "Backing up Redis data..."
docker exec todo-api_redis_1 redis-cli SAVE
cp /data/redis/dump.rdb "$BACKUP_DIR/redis_$DATE.rdb"

# 保留最近7天的备份
find "$BACKUP_DIR" -type f -mtime +7 -delete

echo "Backup completed successfully!" 