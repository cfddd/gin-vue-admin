## docker 部署
### 1.从github上拉取源码
```
git clone https://github.com/cfddd/gin-vue-admin.git
```
### 2.拉取镜像
```
docker pull cfddfc/whpu:server
docker pull cfddfc/whpu:web
```

### 3.生成镜像
需要修改名字为Dockerfile里面使用的名字
```
docker tag cfddfc/whpu:server docker-compose-server:latest
docker tag cfddfc/whpu:web docker-compose-web:latest
```

### 4.运行代码
**首先需要进入项目根目录**

启动容器

```
# 使用docker-compose 后台启动
docker-compose -f deploy/docker-compose/docker-compose.yaml up -d
```
### 5.配置数据库信息
就是把数据库文件信息以sql文件导出，然后进入mysql容器里面再导入

数据卷怎么上传？不会(看了官网的文档，使用ORS上传到dockerHub上，然后再拉下来，然后替换就可以了，感觉也没有很方便……)！


下面是把sql文件导入数据库的命令

```
docker stop gva-server 
# 先关闭server容器（反正在数据库迁移后需要重启server容器）

docker cp dump.sql gva-mysql:/
# 复制文件dump.sql到gva-mysql容器里面

docker exec -it gva-mysql /bin/bash
# 进入gva-mysql容器

mysql -u root -p --binary-mode --force gva < ./dump.sql
# 导入sql文件
```
### 6.完成
退出MySQL容器，然后重新启动server容器
```
exit
# 退出MySQL容器
docker-compose -f deploy/docker-compose/docker-compose.yaml up
# 使用docker-compose启动四个容器
docker-compose -f deploy/docker-compose/docker-compose.yaml up -d
# 后台启动
```

接下来就可以访问了
地址为服务器IP:端口
```
http://xxxx:8080
```