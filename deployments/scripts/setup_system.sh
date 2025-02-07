#!/bin/bash

# 检查是否具有 root 权限
if [ "$EUID" -ne 0 ]; then
    echo "请使用 root 权限运行此脚本"
    exit 1
fi

# 设置系统参数
echo "正在设置系统参数..."

# 创建系统参数配置文件
cat > /etc/sysctl.d/99-redis.conf << EOF
# Redis 系统优化配置
vm.overcommit_memory = 1
net.core.somaxconn = 512
EOF

# 应用系统参数
sysctl -p /etc/sysctl.d/99-redis.conf

# 禁用透明大页面(THP)
echo never > /sys/kernel/mm/transparent_hugepage/enabled
echo never > /sys/kernel/mm/transparent_hugepage/defrag

# 创建 systemd 服务以持久化 THP 设置
cat > /etc/systemd/system/disable-thp.service << EOF
[Unit]
Description=禁用透明大页面(THP)

[Service]
Type=oneshot
ExecStart=/bin/sh -c 'echo never > /sys/kernel/mm/transparent_hugepage/enabled && echo never > /sys/kernel/mm/transparent_hugepage/defrag'
RemainAfterExit=yes

[Install]
WantedBy=multi-user.target
EOF

# 启用并启动服务
systemctl daemon-reload
systemctl enable disable-thp
systemctl start disable-thp

echo "系统参数设置完成" 