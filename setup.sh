#!/bin/bash

#
# Install NS8 control plane:
# - redis
# - traefik
#

HOST=$(hostname -f)

cat <<EOF > /etc/sysctl.d/podman.conf
user.max_user_namespaces=28633
net.ipv4.ip_unprivileged_port_start=0
EOF
sysctl -p /etc/sysctl.d/podman.conf

# Install podman
dnf -y install podman

# Install redis
useradd redis
loginctl enable-linger redis
REDIS_HOME=$(getent passwd redis | cut -d':' -f6)
mkdir -p $REDIS_HOME/.config/systemd/user/{default.target.wants,multi-user.target.wants}
cat <<EOF > $REDIS_HOME/.config/systemd/user/redis.service
[Unit]
Description=Podman redis.service
Documentation=man:podman-generate-systemd(1)
Wants=network.target
After=network-online.target

[Service]
Environment=PODMAN_SYSTEMD_UNIT=%n
Restart=on-failure
ExecStartPre=/bin/rm -f %t/redis.pid %t/redis.ctr-id
ExecStart=/usr/bin/podman run --conmon-pidfile %t/redis.pid --cidfile %t/redis.ctr-id --cgroups=no-conmon --replace --name redis -d --log-driver journald -p 127.0.0.1:6379:6379 --volume redis-data:/data:Z docker.io/redis:6-alpine --appendonly yes
ExecStop=/usr/bin/podman stop --ignore --cidfile %t/redis.ctr-id -t 10
ExecStopPost=/usr/bin/podman rm --ignore -f --cidfile %t/redis.ctr-id
PIDFile=%t/redis.pid
KillMode=none
Type=forking

[Install]
WantedBy=multi-user.target default.target

EOF
chown -R redis:redis $REDIS_HOME/.config
restorecon -R $REDIS_HOME/.config
grep -q 'XDG_RUNTIME_DIR' $REDIS_HOME/.bashrc || echo "export XDG_RUNTIME_DIR=/run/user/$(id -u redis)" >>  $REDIS_HOME/.bashrc
runuser -l redis -c "systemctl --user enable --now redis.service"

# Install traefik
podman run -it --network host  --rm redis redis-cli SET traefik '' # prepare traefik root key
useradd traefik
loginctl enable-linger traefik
TRAEFIK_HOME=$(getent passwd traefik | cut -d':' -f6)
cat <<EOF > $TRAEFIK_HOME/traefik.yaml
defaultEntryPoints:
  - http
  - https

file: {}

log:
  level: DEBUG

accessLog: {}

entryPoints:
  http:
   address: ":80"
  https:
   address: ":443"

providers:
  redis:
    endpoints:
      - "127.0.0.1:6379"

tls:
  certResolver: letsencrypt
  options: {}

certificatesResolvers:
  letsencrypt:
    acme:
      email: root@$HOST
      storage: /etc/traefik/acme/acme.json
      httpChallenge:
        entryPoint: http
      tlsChallenge: false
EOF
chown -R traefik:traefik $TRAEFIK_HOME/traefik.yaml

mkdir -p $TRAEFIK_HOME/.config/systemd/user
cat <<EOF > $TRAEFIK_HOME/.config/systemd/user/traefik.service
[Unit]
Description=Podman traefik.service
Documentation=man:podman-generate-systemd(1)
Wants=network.target
After=network-online.target

[Service]
Environment=PODMAN_SYSTEMD_UNIT=%n
Restart=on-failure
ExecStartPre=/bin/rm -f %t/traefik.pid %t/traefik.ctr-id
ExecStart=/usr/bin/podman run --conmon-pidfile %t/traefik.pid --cidfile %t/traefik.ctr-id --cgroups=no-conmon --replace --network=host --log-driver journald --name traefik -d -v traefik-acme:/etc/traefik/acme:Z -v $TRAEFIK_HOME/traefik.yaml:/etc/traefik/traefik.yaml:Z docker.io/traefik:v2.4
ExecStop=/usr/bin/podman stop --ignore --cidfile %t/traefik.ctr-id -t 10
ExecStopPost=/usr/bin/podman rm --ignore -f --cidfile %t/traefik.ctr-id
PIDFile=%t/traefik.pid
KillMode=none
Type=forking

[Install]
WantedBy=multi-user.target default.target
EOF
chown -R traefik:traefik $TRAEFIK_HOME/.config
restorecon -R $TRAEFIK_HOME/.config
grep -q 'XDG_RUNTIME_DIR' $TRAEFIK_HOME/.bashrc || echo "export XDG_RUNTIME_DIR=/run/user/$(id -u traefik)" >>  $TRAEFIK_HOME/.bashrc
runuser -l traefik -c "systemctl --user enable --now traefik.service"
