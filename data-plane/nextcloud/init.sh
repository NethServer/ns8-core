#!/bin/bash

# Test mysql access: podman run --pod nextcloud -it --rm mariadb:10 mysql -h127.0.0.1 -unextcloud -pnextcloud

N=nextcloud
HOST=nextcloud.$(hostname -f)
PORT=8181

mkdir -p .config/systemd/user/
pushd .config/systemd/user/

cat <<EOF > nextcloud.service
[Unit]
Description=Podman nextcloud.service
Documentation=man:podman-generate-systemd(1)
Wants=network.target
After=network-online.target
Requires=nextcloud-app.service nextcloud-db.service nextcloud-redis.service
Before=nextcloud-app.service nextcloud-db.service nextcloud-redis.service

[Service]
Environment=PODMAN_SYSTEMD_UNIT=%n
Restart=on-failure
TimeoutStopSec=70
ExecStartPre=/bin/rm -f %t/nextcloud.pid %t/nextcloud.pod-id
ExecStartPre=/usr/bin/podman pod create --infra-conmon-pidfile %t/nextcloud.pid --pod-id-file %t/nextcloud.pod-id --name nextcloud -p 127.0.0.1:$PORT:80 --replace
ExecStart=/usr/bin/podman pod start --pod-id-file %t/nextcloud.pod-id
ExecStop=/usr/bin/podman pod stop --ignore --pod-id-file %t/nextcloud.pod-id -t 10
ExecStopPost=/usr/bin/podman pod rm --ignore -f --pod-id-file %t/nextcloud.pod-id
PIDFile=%t/nextcloud.pid
Type=forking

[Install]
WantedBy=multi-user.target default.target
EOF

cat <<EOF > nextcloud-app.service

[Unit]
Description=Podman nextcloud-app.service
Documentation=man:podman-generate-systemd(1)
Wants=network.target
After=network-online.target
BindsTo=nextcloud.service
After=nextcloud.service

[Service]
Environment=PODMAN_SYSTEMD_UNIT=%n
Restart=on-failure
TimeoutStopSec=70
ExecStartPre=/bin/rm -f %t/nextcloud-app.pid %t/nextcloud-app.ctr-id
ExecStart=/usr/bin/podman run --conmon-pidfile %t/nextcloud-app.pid --cidfile %t/nextcloud-app.ctr-id --cgroups=no-conmon --pod-id-file %t/nextcloud.pod-id --replace -d --name nextcloud-app -e MYSQL_DATABASE=nextcloud -e MYSQL_HOST=127.0.0.1 -e MYSQL_USER=nextcloud -e MYSQL_PASSWORD=nextcloud -e NEXTCLOUD_ADMIN_USER=admin -e NEXTCLOUD_ADMIN_PASSWORD=Nethesis,1234 -e REDIS_HOST=127.0.0.1 -e REDIS_HOST_PORT=6379 -e "NEXTCLOUD_TRUSTED_DOMAINS=nextcloud.debian.neth.loc localhost" -e APACHE_DISABLE_REWRITE_IP=1 -e TRUSTED_PROXIES=10.0.0.0/8 -v nextcloud-app-data:/var/www/html/data -v nextcloud-app-source:/var/www/html docker.io/nextcloud:latest
ExecStop=/usr/bin/podman stop --ignore --cidfile %t/nextcloud-app.ctr-id -t 10
ExecStopPost=/usr/bin/podman rm --ignore -f --cidfile %t/nextcloud-app.ctr-id
PIDFile=%t/nextcloud-app.pid
Type=forking

[Install]
WantedBy=multi-user.target default.target
EOF

cat <<EOF > nextcloud-db.service
[Unit]
Description=Podman nextcloud-db.service
Documentation=man:podman-generate-systemd(1)
Wants=network.target
After=network-online.target
BindsTo=nextcloud.service
After=nextcloud.service

[Service]
Environment=PODMAN_SYSTEMD_UNIT=%n
Restart=on-failure
TimeoutStopSec=70
ExecStartPre=/bin/rm -f %t/nextcloud-db.pid %t/nextcloud-db.ctr-id
ExecStartPre=-/bin/touch %t/dump-sql
ExecStart=/usr/bin/podman run --conmon-pidfile %t/nextcloud-db.pid --cidfile %t/nextcloud-db.ctr-id --cgroups=no-conmon --pod-id-file %t/nextcloud.pod-id --replace -d --name nextcloud-db -e MYSQL_ROOT_PASSWORD=nextcloud -e MYSQL_DATABASE=nextcloud -e MYSQL_USER=nextcloud -e MYSQL_PASSWORD=nextcloud -v nextcloud-db-data:/var/lib/mysql -v %t/dump.sql:/docker-entrypoint-initdb.d/dump.sql docker.io/mariadb:10
ExecStop=/usr/bin/podman stop --ignore --cidfile %t/nextcloud-db.ctr-id -t 10
ExecStopPost=/usr/bin/podman rm --ignore -f --cidfile %t/nextcloud-db.ctr-id
PIDFile=%t/nextcloud-db.pid
Type=forking

[Install]
WantedBy=multi-user.target default.target
EOF

cat <<EOF > nextcloud-redis.service
[Unit]
Description=Podman nextcloud-redis.service
Documentation=man:podman-generate-systemd(1)
Wants=network.target
After=network-online.target
BindsTo=nextcloud.service
After=nextcloud.service

[Service]
Environment=PODMAN_SYSTEMD_UNIT=%n
Restart=on-failure
TimeoutStopSec=70
ExecStartPre=/bin/rm -f %t/nextcloud-redis.pid %t/nextcloud-redis.ctr-id
ExecStart=/usr/bin/podman run --conmon-pidfile %t/nextcloud-redis.pid --cidfile %t/nextcloud-redis.ctr-id --cgroups=no-conmon --pod-id-file %t/nextcloud.pod-id --replace -d --name nextcloud-redis -v nextcloud-redis-data:/data docker.io/redis:6-alpine --appendonly yes
ExecStop=/usr/bin/podman stop --ignore --cidfile %t/nextcloud-redis.ctr-id -t 10
ExecStopPost=/usr/bin/podman rm --ignore -f --cidfile %t/nextcloud-redis.ctr-id
PIDFile=%t/nextcloud-redis.pid
Type=forking

[Install]
WantedBy=multi-user.target default.target
EOF

XDG_RUNTIME_DIR=/run/user/$(id -u) systemctl --user enable --now $N

# The command below should be executed by the API server
podman run -i --network host --rm docker.io/redis redis-cli <<EOF
SET traefik/http/services/$N/loadbalancer/servers/0/url http://127.0.0.1:8181
SET traefik/http/routers/$N-http/service $N
SET traefik/http/routers/$N-http/entrypoints http,https
SET traefik.http/routers/$N-http/rule "Host(\`$HOST\`)"
SET traefik/http/routers/$N-https/entrypoints http,https
SET traefik/http/routers/$N-https/rule "Host(\`$HOST\`)"
SET traefik/http/routers/$N-https/tls true
SET traefik/http/routers/$N-https/service $N
SET traefik/http/routers/$N-https/tls/domains/0/main $HOST  
EOF
# enable let's encrypt
# SET traefik/http/routers/$N-https/tls/certresolver letsencrypt
