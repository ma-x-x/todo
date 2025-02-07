#!/bin/bash

# 设置备份目录
BACKUP_DIR="/data/backups"
DATE=$(date +%Y%m%d_%H%M%S)

# 创建备份目录
mkdir -p "$BACKUP_DIR"

# MySQL 数据库备份
echo "正在备份 MySQL 数据库..."
docker exec todo-api_mysql_1 mysqldump -u root -proot todo_db > "$BACKUP_DIR/mysql_$DATE.sql"

# Redis 数据备份
echo "正在备份 Redis 数据..."
docker exec todo-api_redis_1 redis-cli SAVE
cp /data/redis/dump.rdb "$BACKUP_DIR/redis_$DATE.rdb"

# 删除7天前的备份文件
find "$BACKUP_DIR" -type f -mtime +7 -delete

echo "备份完成！" 