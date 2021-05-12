# Architecture

## Core components

The core purpose is managing the applications, providing the basics for their entire lifecycle (install, upgrade, reconfigure, uninstall...). It runs the following components:

1. Redis instance running as rootfull container, bound to TCP port 6379. The Redis DB
   stores the system and modules configuration and provides a signaling bus based on its PUB/SUB feature.

2. `node-agent.service` Systemd unit, running as root. The events are defined in `/usr/local/share/agent/node-events`
   and `/var/local/node-events` (for local Sysadmin overrides).

3. `module-agent.service` Systemd units, running in each module as non-privileged users. See the "Additional modules" section below for more details.

4. Edge proxy, for TLS termination and centralized certificates management (Traefik)

5. LDAP proxy, listening on 127.0.0.1 port 3890. It helps other modules to connect to the account provider LDAP service,
   providing a fixed address and clear text connection.

6. LDAP local account provider (OpenLDAP, Samba DC)

7. VPN connects all cluster node

Further components will be added in the future (e.g. API Server).

## Modules

- The core instantiates a set of *modules*. Each module instance
  runs an application (e.g. Webtop, Nextcloud) as a set of one or more Podman **rootless containers**. In exceptional
  cases a module can run also rootfull containers; the only known exception by now is Samba DC.

- An exclusive unix user account is created for each module instance (e.g. `webtop0`, `nextcloud0`, `openldap0`...).

- The unix user account has session lingering enabled: it automatically starts a persistent Systemd user manager instance.

- A module is installed by the core `node-agent` service, when it receives a "module.init" event.
  The next sections show some examples of the Redis HSET/PUB commands that signal that event.

- A module provides a bunch of event handlers and Systemd unit definitions.
  An event is handled by one or more *action* scripts, stored under `$HOME/.config/module-events`. 
  The local Sysadmin can extend and/or override them by putting their action scripts under `$HOME/module-events`.
  Module-local systemd units are installed under `$HOME/.config/systemd/user`. System-wide units are installed under
  `/etc/systemd/user`.

- A module must provide a `module.init` event implementation that installs and configures
  the Podman containers. Installed services must be enabled to start on boot.

- The `module-agent.service` Systemd unit executes the event handlers. The agent daemon runs in the Python virtual
  environment installed in `/usr/local/share/agent/`: action scripts inherit the same environment. Additional binaries
  can be installed under `/usr/local/share/agent/bin/`.

- Each module has granted full read-only access to the Redis database.

- Each module has a public/private key pair to encrypt passwords and other secrets in the Redis database.

- `<worker_vpn_ip>` is the VPN worker IP set in the previous command

## Prototype validation

List of things considered almost stable, with or without an existing prototype implementation:

- Systemd & Podman foundations
- Unix users for rootless modules
- WireGuard VPN among nodes
- Account providers:
  - Samba AD and OpenLDAP account providers (both are LDAP)
  - Remote LDAP account provider
  - No Unix accounts for domain users
- Node agent / Module agents
- Events and Actions
- Container Registry as software repository for everything
- Store environment variables for actions and containers in Redis
- Authenticated Redis access for write operations
- Public Redis read only access
- Encrypted secrets in Redis DB
- Edge proxies with TLS termination
- Centralized certificate management
- FHS compliancy

## Backup

See [backup](backup.md).
