# Install

**NOTE**: this is just a prototype, *do not* use in production.

## Core installation

When installing on Debian 10 Buster, first make sure to have the latest running kernel and curl installed:
```
apt-get update && apt-get --with-new-pkgs upgrade -y && apt-get install curl -y
reboot
```

Start the installation procedure:
```
# curl https://raw.githubusercontent.com/NethServer/ns8-scratchpad/main/core/install.sh | bash
```

The procedure configures a single node cluster. It prints how to invoke additional commands
to initialize a multi node cluster (`create-cluster`) or to join an existing cluster (`join-cluster`).

### Developer configuration

If you're a developer and you need to push images to the registry, you must configure the authentication.
Create a [GitHub PAT](https://docs.github.com/en/github/authenticating-to-github/creating-a-personal-access-token)
for the **ghcr.io** registry (for read-only access `read:packages private` scope should be enough) then run the following command, specifying
your GitHub user name and providing the generated PAT as password:
```
buildah login ghcr.io
```

The core is composed also by the following components:

- traefik, running with `traefik1` user
- restic rest-server, running with user `restic1`

### Redis

Once the core has been initialized, you can access Redis with one of the following command:

    # podman exec -ti redis redis-cli <<EOF
    PING
    EOF

The above command works only for the root user. An experimental Python based helper script
is available for every user. This client is synchronous: it waits for the server response.

    $ redis-exec <<EOF
    ...
    EOF

An alternative, synchronous invocation relies on the `nc` command, provided by the `nmap-ncat` RPM:

    $ nc 127.0.0.1 6379 <<EOF
    ...
    EOF

Or even shorter in Bash (asynchronous):

    $ cat >/dev/tcp/127.0.0.1/6379 <<EOF
    ...
    EOF

### Traefik

To inspect and modify the module start a SSH session. SSH is preferred to `su - traefik0` because the latter
does not properly initialize the Systemd session environment. Check the services are running with:

    # ssh traefik1@localhost
    # systemctl --user status

To uninstall the `traefik1` module run

    # loginctl disable-linger traefik1
    # userdel -r traefik1

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

### Samba

The Samba module runs a rootfull Samba 4 DC instance.

Rootfull mode is required because Samba needs special privileges
to store ACLs in the filesystem extended attributes.

Only one instance of Samba can run on a node because Samba services are bound
to a host IP address to serve LAN clients, and 127.0.0.1

Install with:

```
install-module samba
```

If the module gets id `samba1`, start the provisioning procedure like:

```
IPADDRESS=10.133.0.2 PROVISION_TYPE=new-domain ADMINUSER=administrator ADMINPASS=Nethesis,1234 REALM=AD.DP.NETHSERVER.NET NBDOMAIN=AD HOSTNAME=dc1.ad.dp.nethserver.net /var/lib/nethserver/samba1/root/bin/provision-domain
```

The DC storage is persisted to the following Podman local volumes:

- samba1-data
- samba1-config


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

The `uninstall.sh` script attempts to stop and erase core components and
additional modules. Handle it with care because it erases everything under `/home/*`!

    bash /var/lib/nethserver/node/uninstall.sh

