# Install

**NOTE**: this is just a prototype, *do not* use in production.

## Core installation

When installing on Debian 10 Buster, first make sure to have the latest running kernel and curl installed:
```
apt-get update && apt-get --with-new-pkgs upgrade -y && apt-get install curl -y && reboot
```

The above commands are not required on Debian 11 Bullseye.

Start the installation procedure as `root`:
```
curl https://raw.githubusercontent.com/NethServer/ns8-scratchpad/main/core/install.sh | bash
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

Once the core has been initialized, you can access Redis with one of the following commands:

    redis-cli <<EOF
    PING
    EOF

The `redis-cli` command is faster when invoked by the root user, because it runs a command in
the existing `redis` rootfull container. For non-root users a temporary container is created and destroyed.

The `redis-cli` command attempts to connect with the higher available Redis privileges, by reading `agent.env` files.

An experimental Python-based helper script is available too and it better suits
Bash action scripts as it does not start any container.
This client is synchronous as it waits for the server response.

    redis-exec <<EOF
    PING
    EOF

An alternative, synchronous invocation relies on the `nc` command, provided by the `nmap-ncat` RPM (good for development environments):

    nc 127.0.0.1 6379 <<EOF
    PING
    EOF

For completeness here is a pure Bash invocation that does not wait the server response:

    cat >/dev/tcp/127.0.0.1/6379 <<EOF
    PING
    EOF

### Traefik

To inspect and modify the module start a SSH session. SSH is preferred to `su - traefik0` because the latter
does not properly initialize the Systemd session environment. Check the services are running with:

    # ssh traefik1@localhost
    # systemctl --user status

To uninstall the `traefik1` module run

    # systemctl stop user@$(id -u traefik1)
    # loginctl disable-linger traefik1
    # userdel -r traefik1

See traefik README for more info.

### Ldapproxy

The LDAP account provider service can be local or remote, and can require TLS or not. To help modules using the LDAP
service, a LDAP proxy is always available at 127.0.0.1:3890 without TLS.

Rootless containers can establish a connection from their private network to the loopback interface with the
following Podman arguments:

    --network=slirp4netns:allow_host_loopback=true --add-host=accountprovider:10.0.2.2

To create a Ldapproxy instance run the following commands, adjusting the LDAPHOST value if needed.

XXX: does not work

### Samba

The Samba module runs a rootfull Samba 4 DC instance.

Rootfull mode is required because Samba needs special privileges
to store ACLs in the filesystem extended attributes.

Only one instance of Samba can run on a node because Samba services are bound
to a host IP address to serve LAN clients, and 127.0.0.1

Install with:

```
add-module samba 1
```

If the module gets id `samba1`, start the provisioning procedure like:

```
IPADDRESS=10.133.0.2 PROVISION_TYPE=new-domain ADMINUSER=administrator ADMINPASS=Nethesis,1234 REALM=AD.DP.NETHSERVER.NET NBDOMAIN=AD HOSTNAME=dc1.ad.dp.nethserver.net /var/lib/nethserver/samba1/bin/provision-domain
```

The DC storage is persisted to the following Podman local volumes:

- samba1-data
- samba1-config


### VPN

Each node is connected to the leader node using WireGuard VPN in a star network topology.

The VPN uses private network, default is `10.5.4.2.0/24`.
The leader node gets the first IP address, by default `10.5.4.1`.
Worker nodes will have IP address like `10.5.4.2`, `10.5.4.3`, etc.

## Applications installation

### Dokuwiki

See Dokuwiki README.

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
additional modules. Handle it with care because it erases everything under `/home/*` and `/var/lib/nethserver/*`!

    bash /var/lib/nethserver/node/uninstall.sh

