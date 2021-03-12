# ns8-scratchpad

NethServer 8 experiments using containers on Fedora 33.

The install procedure should work also on Debian 10.


## Core components

The core purpose is managing the applications, providing the basics for their entire lifecycle (install, upgrade, reconfigure, uninstall...). It runs the following components:

1. Redis instance running as rootless container of the `redis0` user,
   bound to TCP port 6379 in host network namespace. The Redis DB
   stores the system and modules configuration and provides a signaling bus based on its PUB/SUB feature.

2. `node-agent.service` Systemd unit, running as root. The events are defined in `/usr/local/share/agent/node-events`
   and `/var/local/node-events` (for local Sysadmin overrides).

3. `module-agent.service` Systemd units, running in each module as non-privileged users. See the "Additional modules" section below for more details.

Some modules are considered core components, as they provide basic services for additional modules:

- LDAP local account provider (OpenLDAP, Samba DC)
- Edge proxy, for TLS termination and centralized certificates management (Traefik)

Further components will be added in the future (e.g. API Server, VPN, ...).

## Applications

- The core instantiates a set of *modules*. Each module instance
  runs an application (e.g. Webtop, Nextcloud) as one or more Podman **rootless containers**. In exceptional
  cases a module can run also a set of rootfull containers; the known exception is Samba DC.

- An exclusive unix user account is created for each module instance (e.g. `webtop0`, `nextcloud0`, `openldap0`...).

- The unix user account has session lingering enabled: it automatically starts a persistent Systemd user manager instance.

- A module is installed by the core `node-agent` service, when it receives a "module.init" event.
  The next sections show some examples of the Redis HSET/PUB commands that signal that event.

- A module provides a bunch of event handlers and Systemd unit definitions.
  An event is handled by one or more *action* scripts, stored under `$HOME/.config/module-events`. 
  The local Sysadmin can extend and/or override them by putting their action scripts under `$HOME/module-events`.
  Module-local systemd units are installed under `$HOME/.config/systemd/user`. System-wide units are installed under
  `/etc/systemd/user`.

- A module must provides a `module.init` event implementation that installs and configures
  the Podman containers. Installed services must be enabled to start on boot.

- The `module-agent.service` Systemd unit executes the event handlers. The agent daemon runs in the Python virtual
  environment installed in `/usr/local/share/agent/`: action scripts inherit the same environment. Additional binaries
  can be installed under `/usr/local/share/agent/bin/`.

- Each module has granted full read-only access to the Redis database.

- Each module has a public/private key pair to encrypt passwords and other secrets in the Redis database.

## Core installation

Execute as root:
```
# curl https://raw.githubusercontent.com/DavidePrincipi/ns8-scratchpad/main/core/install.sh | bash
```

If you're a developer and you need to push images to the registry, you must configure the authentication.
Create a [GitHub PAT](https://docs.github.com/en/github/authenticating-to-github/creating-a-personal-access-token)
for the **ghcr.io** registry (for read-only access `read:packages private` scope should be enough) then run the following command, specifying
your GitHub user name and providing the generated PAT as password:
```
# podman login --authfile /usr/local/etc/registry.json ghcr.io
```

### Traefik

Traefik will proxy all HTTP/HTTPs connections to web applications.

Once the core has been initialized run this Redis command to initialize the Traefik
module instance and start it:

    podman run -i --network host --rm docker.io/redis:6-alpine redis-cli <<EOF
    SET traefik ''
    HSET traefik0/module.env LE_EMAIL root@$(hostname -f) EVENTS_IMAGE ghcr.io/nethserver/traefik:latest
    PUBLISH $(hostname -s):module.init traefik0
    EOF

As alternative, use `nc` command:

    # nc 127.0.0.1 6379 <<EOF
    ...
    EOF

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
SET traefik/http/routers/$N-http/service $N
SET traefik/http/routers/$N-http/entrypoints http,https
SET traefik.http/routers/$N-http/rule "Host(\`$HOST\`)"
SET traefik/http/routers/$N-https/entrypoints http,https
SET traefik/http/routers/$N-https/rule "Host(\`$HOST\`)"
SET traefik/http/routers/$N-https/tls true
SET traefik/http/routers/$N-https/service $N
SET traefik/http/routers/$N-https/tls/certresolver letsencrypt
SET traefik/http/routers/$N-https/tls/domains/0/main $HOST
EOF
```

Traefik will generate the certificate without exposing any new service.

### Nsdc

The Nsdc module runs a singleton and rootfull Samba 4 DC instance.

- *Rootfull* because Samba needs special privileges to store ACLs in the filesystem extended attributes
- *Singleton* because Samba services are bound to a host IP address, to serve LAN clients

Initialize the Redis DB and start the installation with:

```
podman run -ti --network host --rm docker.io/redis:6-alpine redis-cli <<EOF
HSET nsdc0/module.env EVENTS_IMAGE ghcr.io/nethserver/nsdc:latest NSDC_IMAGE ghcr.io/nethserver/nsdc:latest IPADDRESS 10.133.0.5 HOSTNAME nsdc1.$(hostname -d) NBDOMAIN AD REALM AD.$(hostname -d | tr a-z A-Z) ADMINPASS Nethesis,1234
PUBLISH $(hostname -s):module-rootfull.init nsdc0
EOF
```

The DC storage is persisted to the following Podman local volumes:

- nsdc0-data
- nsdc0-config

## Applications installation

### Dokuwiki

To start a dokuwiki instance execute:
```
podman run -i --network host --rm docker.io/redis:6-alpine redis-cli <<EOF
HSET dokuwiki0/module.env EVENTS_IMAGE ghcr.io/nethserver/dokuwiki:latest
PUBLISH $(hostname -s):module.init dokuwiki0
EOF
```

Setup traefik routes:
```
N=dokuwiki HOST=dokuwiki.$(hostname -f); podman run -i --network host --rm docker.io/redis:6-alpine redis-cli <<EOF
SET traefik/http/services/$N/loadbalancer/servers/0/url http://127.0.0.1:8080
SET traefik/http/routers/$N-http/service $N
SET traefik/http/routers/$N-http/entrypoints http,https
SET traefik.http/routers/$N-http/rule "Host(\`$HOST\`)"
SET traefik/http/routers/$N-https/entrypoints http,https
SET traefik/http/routers/$N-https/rule "Host(\`$HOST\`)"
SET traefik/http/routers/$N-https/tls true
SET traefik/http/routers/$N-https/service $N
SET traefik/http/routers/$N-https/tls/certresolver letsencrypt
SET traefik/http/routers/$N-https/tls/domains/0/main $HOST
EOF
```

### Nextcloud

To start a nextcloud instance execute:
```
podman run -i --network host --rm docker.io/redis:6-alpine redis-cli <<EOF
HSET nextcloud0/module.env EVENTS_IMAGE ghcr.io/nethserver/nextcloud:latest
PUBLISH $(hostname -s):module.init nextcloud0
EOF
```

Setup traefik:
```
N=nextcloud HOST=nextcloud.$(hostname -f); podman run -i --network host --rm docker.io/redis:6-alpine redis-cli <<EOF
SET traefik/http/services/$N/loadbalancer/servers/0/url http://127.0.0.1:8181
SET traefik/http/routers/$N-http/service $N
SET traefik/http/routers/$N-http/entrypoints http,https
SET traefik.http/routers/$N-http/rule "Host(\`$HOST\`)"
SET traefik/http/routers/$N-https/entrypoints http,https
SET traefik/http/routers/$N-https/rule "Host(\`$HOST\`)"
SET traefik/http/routers/$N-https/tls true
SET traefik/http/routers/$N-https/service $N
SET traefik/http/routers/$N-https/tls/certresolver letsencrypt
SET traefik/http/routers/$N-https/tls/domains/0/main $HOST  
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

## Backup & restore

### First configuration

The core and each instance can implement the `backup` event. The event takes the name of the backup as first argument.

First, generate a password for the new backup and set restic repository base URL:
```
podman run -i --network host --rm docker.io/redis:6-alpine redis-cli SETNX backup/<backup_name>/password <password>
podman run -i --network host --rm docker.io/redis:6-alpine redis-cli SET backup/<backup_name>/base_repository rest:http://127.0.0.1:8383
```

Example:
```
podman run -i --network host --rm docker.io/redis:6-alpine redis-cli SETNX backup/backup1/password $(cat /dev/urandom | tr -dc 'a-zA-Z0-9' | fold -w 32 | head -n 1)
podman run -i --network host --rm docker.io/redis:6-alpine redis-cli SET backup/backup1/base_repository rest:http://127.0.0.1:8383
```

### Backup the core

To execute the backup of the core execute:
```
podman run -i --network host --rm docker.io/redis:6-alpine redis-cli PUBLISH <hostname>:backup <backup_name>
```

Example:
```
podman run -i --network host --rm docker.io/redis:6-alpine redis-cli PUBLISH $(hostname -s):backup backup1
```

### Backup an instance

To execute the backup of an instance execute:
```
podman run -i --network host --rm docker.io/redis:6-alpine redis-cli PUBLISH <user_instance>:backup <backup_name>
```

Example:
```
podman run -i --network host --rm docker.io/redis:6-alpine redis-cli PUBLISH nextcloud0:backup backup1
```

### Restore an instance

Each instance implementing the `backup` event, must implement also the `restore` event. The event takes the name of the backup as first argument.

To execute the restore of an instance execute:
```
podman run -i --network host --rm docker.io/redis:6-alpine redis-cli PUBLISH <user_instance>:restore <backup_name>
```

Example:
```
podman run -i --network host --rm docker.io/redis:6-alpine redis-cli PUBLISH nextcloud0:restore backup1
```


## Uninstall

The `core/uninstall.sh` script attempts to stop and erase core components and
additional modules. Handle it with care because it erases everything under `/home/*`!

    bash uninstall.sh
