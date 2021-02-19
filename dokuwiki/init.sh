#!/bin/bash

# DOC: https://github.com/bitnami/bitnami-docker-dokuwiki#configuration

N=dokuwiki

useradd dokuwiki
loginctl enable-linger dokuwiki
DOKUWIKI_HOME=$(getent passwd dokuwiki | cut -d':' -f6)
mkdir -p $DOKUWIKI_HOME/.config/systemd/user

# configure persistente dir
#mkdir -p $HOME/dokuwiki/data
#podman unshare chown 1001:1001 /home/dokuwiki/dokuwiki/data/
#vols=" --volume $HOME/dokuwiki/data:/bitnami/dokuwiki:Z"

# set virtual host name
HOST=dokuwiki.$(hostname -f)


cat <<EOF > $DOKUWIKI_HOME/.config/systemd/user/dokuwiki.service
[Unit]
Description=Podman dokuwiki.service
Documentation=man:podman-generate-systemd(1)
Wants=network.target
After=network-online.target

[Service]
Environment=PODMAN_SYSTEMD_UNIT=%n
Restart=on-failure
ExecStartPre=/bin/rm -f %t/dokuwiki.pid %t/dokuwiki.ctr-id
ExecStart=/usr/bin/podman run --conmon-pidfile %t/dokuwiki.pid --cidfile %t/dokuwiki.ctr-id --cgroups=no-conmon --replace --log-driver journald --name dokuwiki -d -p 127.0.0.1:8080:8080 -v dokuwiki-data:/bitnami/dokuwiki docker.io/bitnami/dokuwiki:latest 
ExecStop=/usr/bin/podman stop --ignore --cidfile %t/dokuwiki.ctr-id -t 10
ExecStopPost=/usr/bin/podman rm --ignore -f --cidfile %t/dokuwiki.ctr-id
PIDFile=%t/dokuwiki.pid
KillMode=none
Type=forking

[Install]
WantedBy=multi-user.target default.target
EOF
chown -R dokuwiki:dokuwiki $DOKUWIKI_HOME/.config
restorecon -R $DOKUWIKI_HOME/.config
grep -q 'XDG_RUNTIME_DIR' $DOKUWIKI_HOME/.bashrc || echo "export XDG_RUNTIME_DIR=/run/user/$(id -u dokuwiki)" >>  $DOKUWIKI_HOME/.bashrc
XDG_RUNTIME_DIR=/run/user/$(id -u dokuwiki) runuser -w XDG_RUNTIME_DIR -u dokuwiki -- systemctl --user enable --now dokuwiki.service

exit
podman run -it --network host --rm docker.io/redis redis-cli SET traefik/http/services/$N/loadbalancer/servers/0/url http://127.0.0.1:8080
podman run -it --network host --rm docker.io/redis redis-cli SET traefik/http/routers/$N-http/service $N
podman run -it --network host --rm docker.io/redis redis-cli SET traefik/http/routers/$N-http/entrypoints http,https
podman run -it --network host --rm docker.io/redis redis-cli SET traefik.http/routers/$N-http/rule "Host(\`$HOST\`)"
podman run -it --network host --rm docker.io/redis redis-cli SET traefik/http/routers/$N-https/entrypoints http,https
podman run -it --network host --rm docker.io/redis redis-cli SET traefik/http/routers/$N-https/rule "Host(\`$HOST\`)"
podman run -it --network host --rm docker.io/redis redis-cli SET traefik/http/routers/$N-https/tls true
podman run -it --network host --rm docker.io/redis redis-cli SET traefik/http/routers/$N-https/service $N
podman run -it --network host --rm docker.io/redis redis-cli SET traefik/http/routers/$N-https/tls/certresolver letsencrypt
podman run -it --network host --rm docker.io/redis redis-cli SET traefik/http/routers/$N-https/tls/domains/0/main $HOST

