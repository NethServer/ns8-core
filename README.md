# ns8-scratchpad

NethServer 8 experiments using containers on Fedora 33

## Initialize the control plane

1. Retrieve credentials for DigitalOcean registry and save the `docker-config.json` file

2. Execute as root:

       # export REGISTRY_AUTH_FILE=docker-config.json
       # curl https://raw.githubusercontent.com/DavidePrincipi/ns8-scratchpad/main/control-plane/init.sh | bash init.sh

## Control plane components

The control plane runs the following components:

1. Redis instance running as rootless container of the `cplane` user, TCP port 6379 - host network. Access to redis with:

       $ podman run -it --network host --rm redis redis-cli

2. `node-agent.service` Systemd unit, running as root. The events are defined in `/usr/local/share/agent/node-events` and `/var/lib/agent/node-events` (for local Sysadmin overrides).

Further components will be added in the future (e.g. API Server, VPN, ...).


## Data plane components

- The control plane instantiates a set of *modules* (e.g. Webtop, Nextcloud, OpenLDAP, Traefik...) which constitute
  the data plane.

- An exclusive unix user account is created for each module instance (e.g. `webtop0`, `nextcloud0`, `openldap0`...).

- The unix user account has session lingering enabled: 
  it automatically starts a persistent Systemd user manager.
  The Systemd user process runs the `module-agent.service` unit.

- Each module provides a bunch of event handlers. An event is handled by one or more *action* scripts, stored under `$HOME/.config/module-events`. The local Sysadmin can extend and/or override them by putting their action scripts under `$HOME/module-events`.

- The `module-agent.service` Systemd unit executes the event handlers.

- Each module has granted full read-only access to the Redis database.

- Each module has a public/private key pair to encrypt passwords and other secrets in the Redis database.