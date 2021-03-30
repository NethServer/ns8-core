# Install

**NOTE**: this is just a prototype, *do not* use in production.

## Core installation

Execute as root on each node:
```
# curl https://raw.githubusercontent.com/NethServer/ns8-scratchpad/main/core/install.sh | bash
```

When installing on Debian 10 Buster, first make sure to have the latest running kernel:
```
apt-get update
apt-get --with-new-pkgs upgrade -y
apt-get install curl
reboot
```

After installation, to initialize the cluster you need 3 different steps:

- master node initialization
- join a worker to cluster
- grant access to the worker

To initialize the cluster, execute on the master node:
```
nethserver init
```

This command will initialize the master node as WireGuard VPN server with private IP `10.5.4.1`, listening on port `55820`.

Then, move to the worker node and execute:
```
nethserver join <master_pubkey> <master_public_address:vpn_port> <worker_vpn_ip>
```

Where:
- `<master_pubkey>` is the master WireGuard publick key, execute `cat /etc/wireguard/privatekey | wg pubkey` to access it
- `<master_public_address:vpn_port>` is the master public address with WireGuard listening port, something like `1.2.3.4:55820`
- `<worker_vpn_ip>` is the VPN private IP of the worker, something like `10.5.4.2`

Finally, execute on the master node:
```
nethserver grant <worker_hostname> <worker_pubkey> <worker_vpn_ip>
```

Where:
- `<worker_hostname>` is the worker hostname from `hostname -s`
- `<worker_pubkey>` is the worker WireGuard publick key, execute `cat /etc/wireguard/privatekey | wg pubkey` to access it
- `<worker_vpn_ip>` is the VPN worker IP set in the previous command


### Developer configuration

If you're a developer and you need to push images to the registry, you must configure the authentication.
Create a [GitHub PAT](https://docs.github.com/en/github/authenticating-to-github/creating-a-personal-access-token)
for the **ghcr.io** registry (for read-only access `read:packages private` scope should be enough) then run the following command, specifying
your GitHub user name and providing the generated PAT as password:
```
podman login --authfile /usr/local/etc/registry.json ghcr.io
export REGISTRY_AUTH_FILE=/usr/local/etc/registry.json
chmod -c 644 /usr/local/etc/registry.json
```

The core is composed also by the following components:

- traefik, running with `traefik0` user
- restic rest-server, running with user `restic0`


### Redis

Once the core has been initialized, you can access Redis with one of the following command:

    podman run -i --network host --rm docker.io/redis:6-alpine redis-cli <<EOF
    PING
    EOF

As alternative, use `nc` command:

    # nc 127.0.0.1 6379 <<EOF
    ...
    EOF

Or even shorter in Bash:

    # cat >/dev/tcp/127.0.0.1/6379 <<EOF
    ...
    EOF

### Traefik

To inspect and modify the module start a SSH session. SSH is preferred to `su - traefik0` because the latter
does not properly initialize the Systemd session environment. Check the services are running with:

    # ssh traefik0@localhost
    # systemctl --user status

To uninstall the `traefik0` module run

    # loginctl disable-linger traefik0
    # userdel -r traefik0

#### Default Let's Encrypt certificate

To request a Let's Encrypt certificate for the server FQDN, just execute:
```
N=default HOST=$(hostname -f); podman run -i --network host --rm docker.io/redis:6-alpine redis-cli <<EOF
SET traefik0/http/routers/$N-http/service $N
SET traefik0/http/routers/$N-http/entrypoints http,https
SET traefik0/http/routers/$N-http/rule "Host(\`$HOST\`)"
SET traefik0/http/routers/$N-https/entrypoints http,https
SET traefik0/http/routers/$N-https/rule "Host(\`$HOST\`)"
SET traefik0/http/routers/$N-https/tls true
SET traefik0/http/routers/$N-https/service $N
SET traefik0/http/routers/$N-https/tls/certresolver letsencrypt
SET traefik0/http/routers/$N-https/tls/domains/0/main $HOST
EOF
```

Traefik will generate the certificate without exposing any new service.

### Ldapproxy

The LDAP account provider service can be local or remote, and can require TLS or not. To help modules using the LDAP
service, a LDAP proxy is always available at 127.0.0.1:3890 without TLS.

Rootless containers can establish a connection from their private network to the loopback interface with the
following Podman arguments:

    --network=slirp4netns:allow_host_loopback=true --add-host=accountprovider:10.0.2.2

To create a Ldapproxy instance run the following commands, adjusting the LDAPHOST value if needed.

```
cat >/dev/tcp/127.0.0.1/6379 <<EOF
HSET module/ldapproxy0/module.env EVENTS_IMAGE ghcr.io/nethserver/ldapproxy:latest PROXYPORT <proxy_listen_port> LDAPHOST <ldap_host> LDAPPORT <ldap_port> LDAPSSL <on|off> SCHEMA <ad|openldap> REALM <realm> BIND_DN <bind_dn> BIND_PASSWORD <pass> BASE_DN <base_dn> USER_DN <user_dn> GROUP_DN <group_dn>
PUBLISH $(hostname -s):module.init ldapproxy0
EOF
```

Example:
```
cat >/dev/tcp/127.0.0.1/6379 <<EOF
HSET module/ldapproxy0/module.env EVENTS_IMAGE ghcr.io/nethserver/ldapproxy:latest PROXYPORT 3890 LDAPHOST 127.0.0.1 LDAPPORT 636 LDAPSSL on SCHEMA ad REALM AD.NETH.LOC BIND_DN administrator@AD.NETH.LOC BIND_PASSWORD Nethesis,1234 BASE_DN dc=ad,dc=neth,dc=loc USER_DN dc=ad,dc=neth,dc=loc GROUP_DN dc=ad,dc=neth,dc=loc
PUBLISH $(hostname -s):module.init ldapproxy0
EOF
```

### Nsdc

The Nsdc module runs a singleton and rootfull Samba 4 DC instance.

- *Rootfull* because Samba needs special privileges to store ACLs in the filesystem extended attributes
- *Singleton* because Samba services are bound to a host IP address to serve LAN clients, and 127.0.0.1

Initialize the Redis DB and start the installation with:

```
podman run -i --network host --rm docker.io/redis:6-alpine redis-cli <<EOF
HSET module/nsdc0/module.env EVENTS_IMAGE ghcr.io/nethserver/nsdc:latest NSDC_IMAGE ghcr.io/nethserver/nsdc:latest IPADDRESS 10.133.0.5 HOSTNAME nsdc1.$(hostname -d) NBDOMAIN AD REALM AD.$(hostname -d | tr a-z A-Z) ADMINPASS Nethesis,1234
PUBLISH $(hostname -s):module-rootfull.init nsdc0
EOF
```

The DC storage is persisted to the following Podman local volumes:

- nsdc0-data
- nsdc0-config

### VPN

Each node is connected to the master node using WireGuard VPN in a star network topology.
After the installation, the server will be configured as master node.

The VPN uses a `/24` private network, default is `10.5.4.2.0/24`.
The first node will be the master and has the default IP address set to `10.5.4.1`.
All other worker nodes will have IP address like `10.5.4.2`, `10.5.4.3`, etc.

## Applications installation

### Dokuwiki

To start a dokuwiki instance execute:
```
podman run -i --network host --rm docker.io/redis:6-alpine redis-cli <<EOF
HSET module/dokuwiki0/module.env EVENTS_IMAGE ghcr.io/nethserver/dokuwiki:latest
PUBLISH $(hostname -s):module.init dokuwiki0
EOF
```

Setup traefik routes:
```
N=dokuwiki HOST=dokuwiki.$(hostname -f); podman run -i --network host --rm docker.io/redis:6-alpine redis-cli <<EOF
SET traefik0/http/services/$N/loadbalancer/servers/0/url http://127.0.0.1:8080
SET traefik0/http/routers/$N-http/service $N
SET traefik0/http/routers/$N-http/entrypoints http,https
SET traefik0/http/routers/$N-http/rule "Host(\`$HOST\`)"
SET traefik0/http/routers/$N-https/entrypoints http,https
SET traefik0/http/routers/$N-https/rule "Host(\`$HOST\`)"
SET traefik0/http/routers/$N-https/tls true
SET traefik0/http/routers/$N-https/service $N
SET traefik0/http/routers/$N-https/tls/certresolver letsencrypt
SET traefik0/http/routers/$N-https/tls/domains/0/main $HOST
EOF
```

### Nextcloud

Nextcloud will use ldapproxy to configure the authentication, so make sure nsdc and ldappproxy are already running.

To start a nextcloud instance execute:
```
podman run -i --network host --rm docker.io/redis:6-alpine redis-cli <<EOF
HSET module/nextcloud0/module.env EVENTS_IMAGE ghcr.io/nethserver/nextcloud:latest
PUBLISH $(hostname -s):module.init nextcloud0
EOF
```

Setup traefik:
```
N=nextcloud HOST=nextcloud.$(hostname -f); podman run -i --network host --rm docker.io/redis:6-alpine redis-cli <<EOF
SET traefik0/http/services/$N/loadbalancer/servers/0/url http://127.0.0.1:8181
SET traefik0/http/routers/$N-http/service $N
SET traefik0/http/routers/$N-http/entrypoints http,https
SET traefik0/http/routers/$N-http/rule "Host(\`$HOST\`)"
SET traefik0/http/routers/$N-https/entrypoints http,https
SET traefik0/http/routers/$N-https/rule "Host(\`$HOST\`)"
SET traefik0/http/routers/$N-https/tls true
SET traefik0/http/routers/$N-https/service $N
SET traefik0/http/routers/$N-https/tls/certresolver letsencrypt
SET traefik0/http/routers/$N-https/tls/domains/0/main $HOST  
EOF
```

Configure LDAP authentication:
```
podman run -i --network host --rm docker.io/redis:6-alpine redis-cli <<EOF
PUBLISH nextcloud0:ldap-config ldapproxy0
EOF
```

Execute `occ` command:
```
ssh nextcloud0@localhost
podman exec -ti --user www-data nextcloud-app php occ
```

Setup nsdc account provider:
```
ssh nextcloud0@localhost
./scripts/setup_ad.sh
```

Note: the nsdc must have a user named `ldapservice` with password `Nethesis,1234`

### Mail

Installation

```
cat >/dev/tcp/127.0.0.1/6379 <<EOF
HSET module/mail0/module.env HOSTNAME mail.example.com BINDDN %n@AD.DP.NETHSERVER.NET EVENTS_IMAGE ghcr.io/nethserver/mail
PUBLISH $(hostname -s):module.init mail0
EOF
```

### Netdata

Install and configure netdata:
- port 19999 is not exposed
- application is accessibile only using Traefik at `netdata.<fqdn>` on master node
- the master node collects and display all metrics
- worker nodes send data to master node


Execute on the master node:
```
cat >/dev/tcp/127.0.0.1/6379 <<EOF
HSET module/netdata0/module.env EVENTS_IMAGE ghcr.io/nethserver/netdata:latest
PUBLISH $(hostname -s):module.init netdata0
EOF
```

Setup web access using traefik (execute only on master node):
```
N=netdata HOST=netdata.$(hostname -f); cat >/dev/tcp/127.0.0.1/6379 <<EOF
SET traefik0/http/services/$N/loadbalancer/servers/0/url http://127.0.0.1:19999
SET traefik0/http/routers/$N-http/service $N
SET traefik0/http/routers/$N-http/entrypoints http,https
SET traefik0/http/routers/$N-http/rule "Host(\`$HOST\`)"
SET traefik0/http/routers/$N-https/entrypoints http,https
SET traefik0/http/routers/$N-https/rule "Host(\`$HOST\`)"
SET traefik0/http/routers/$N-https/tls true
SET traefik0/http/routers/$N-https/service $N
SET traefik0/http/routers/$N-https/tls/certresolver letsencrypt
SET traefik0/http/routers/$N-https/tls/domains/0/main $HOST
EOF
```

To configure netdata on the worker node, execute on the master:
```
cat >/dev/tcp/127.0.0.1/6379 <<EOF
HSET module/netdata1/module.env EVENTS_IMAGE ghcr.io/nethserver/netdata:latest
PUBLISH <worker_short_hostname>:module.init netdata1
EOF
```

## Uninstall

The `core/uninstall.sh` script attempts to stop and erase core components and
additional modules. Handle it with care because it erases everything under `/home/*`!

    bash uninstall.sh

