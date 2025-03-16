#!/bin/bash

# 解压 tg-card-autosed.zip
unzip tg-card-autosed.zip

# 查找 tg-card-autosed 进程并终止
pkill -f tg-card-autosed

# 等待进程完全退出
sleep 1

# 启动 tg-card-autosed
nohup ./tg-card-autosed > /dev/null 2>&1 &

echo "tg-card-autosed 已重启"