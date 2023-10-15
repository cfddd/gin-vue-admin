#!/bin/bash

# 启动容器
docker-compose -f deploy/docker-compose/docker-compose.yaml up --build -d

# 执行 SQL 脚本，导入数据库
docker exec -it gva-mysql bash -c 'exec mysql -uroot -p"$MYSQL_ROOT_PASSWORD"' < deploy/docker-compose/initSQL/dumpAll.sql

# 清理无用的 Docker 资源
docker system prune