#!/bin/bash

# Test mysql access: podman run --pod nextcloud -it --rm mariadb:10 mysql -h127.0.0.1 -unextcloud -pnextcloud

N=nextcloud
HOST=nextcloud.$(hostname -f)

podman pod create --name $N -p 127.0.0.1:8181:80

podman run --pod $N \
  --name $N-db \
  -e MYSQL_ROOT_PASSWORD=nextcloud  -e MYSQL_DATABASE=nextcloud -e MYSQL_USER=nextcloud -e MYSQL_PASSWORD=nextcloud \
  -v $N-db-data:/var/lib/mysql \
  docker.io/mariadb:10

podman run -d --pod $N \
  --name $N-redis \
  -v $N-redis-data:/data \
  docker.io/redis:6-alpine --appendonly yes

podman run -d --pod $N \
  --name $N-app \
  -e MYSQL_DATABASE=nextcloud -e MYSQL_HOST=127.0.0.1 -e MYSQL_USER=nextcloud -e MYSQL_PASSWORD=nextcloud \
  -e NEXTCLOUD_ADMIN_USER=admin -e NEXTCLOUD_ADMIN_PASSWORD=Nethesis,1234 \
  -e REDIS_HOST=127.0.0.1 -e REDIS_HOST_PORT=6379 \
  -e NEXTCLOUD_TRUSTED_DOMAINS="$HOST localhost" \
  -e APACHE_DISABLE_REWRITE_IP=1 -e TRUSTED_PROXIES=10.0.0.0/8 \
  -v $N-app-data:/var/www/html/data -v $N-app-source:/var/www/html \
  docker.io/nextcloud:latest

podman generate systemd --new --name --files $N
