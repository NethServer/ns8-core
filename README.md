# ns8-scratchpad

NethServer 8 experiments using containers on Fedora 33

- Podman running in rootless mode (every container has its own user)
- Traefik reads the dynamic configuration from Redis (control-plane)
- Local public services are reachable using the host network

## Initialize the control plane

1. Retrieve credentials for DigitalOcean registry and save the `docker-config.json` file

2. Execute as root:

       # export REGISTRY_AUTH_FILE=docker-config.json
       # bash init.sh

## Control plane components

1. Redis instance running as rootless container of the `cplane` user, TCP port 6379 - host network. Access to redis with:

       $ podman run -it --network host --rm redis redis-cli

2. `node-agent.service` unit, running as root. The events are defined in `/usr/local/share/agent/node-events` and `/var/lib/agent/node-events` (local sysadmin overrides).

Still missing:

3. VPN

4. ...


## Data plane services

...TODO