#!/bin/bash

# 解压 tg-auto-card-num.zip
unzip tg-auto-card-num.zip

# 查找 tg-auto-card-num 进程并终止
pkill -f tg-auto-card-num

# 等待进程完全退出
sleep 1

# 启动 tg-auto-card-num
nohup ./tg-auto-card-num > /dev/null 2>&1 &

echo "tg-auto-card-num 已重启"