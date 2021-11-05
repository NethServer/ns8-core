# Design & Architecture

## Design choices

(almost) Everything runs inside a container.

Design goals:

- store the configuration on a reliable cluster-ready database: one read-write Redis master node for the whole cluster, multiple read-only nodes
- use TLS terminator whenever possible to simplify module configuration
- the configuration should be readable for any module inside the system, sensitive data (like passwords) should be store encrypted
- security as first class citizen: rootless container modules with SELinux / AppArmor, rootfull allowed too

### Operating System (OS) requirements

NS8 should run on a OS with the following features:

- small footprint, x86_64 (and aarch64)
- systemd to keeps things running
- podman ≥ 2.2
- support WireGuard VPN
- cgroups v2

### Users

As design choice, NS8 will use two different user bases for the core and the applications.

The core, including the administrator UI, will authenticate only against Redis.
Each module instance will use local Unix accounts.

Applications will be able to access users from the same account provider available on NS7:
- local and remote OpenLDAP
- local and remote AD

Account providers will be accessible using a LDAP proxy.
Each module instance can eventually select a different account provider.

### Prototype validation

List of things considered almost stable, with or without an existing prototype implementation:

- Systemd & Podman foundations
- Unix users for rootless modules
- WireGuard VPN among nodes
- Account providers:
  - Samba AD and OpenLDAP account providers (both are LDAP)
  - Remote LDAP account provider
  - No Unix accounts for domain users
- Cluster agent, Node agent, Module agents
- Events and Actions
- Container Registry as software repository for everything
- Store (a copy of) environment variables for actions and containers in Redis
- Authenticated Redis access for write operations
- Public Redis read only access
- Encrypted secrets in Redis DB
- Edge proxies with TLS termination
- Centralized certificate management
- FHS compliant
- API server
- UI structure

## Components

The system is composed by two main components:
- core
- modules

### Core components

The core purpose is managing the applications, providing the basics for their entire life cycle (install, upgrade, reconfigure, uninstall...). It runs the following components:

- Redis instance running as rootfull container, bound to TCP port 6379. The Redis DB
  * exchanges messages among agents and the API server
  * helps modules to discover information about the core and other modules
  * provides a signaling bus based on its PUB/SUB feature
  * stores (a copy of) the system and modules configuration
  * is an authentication backend for the management UI and the agents

- `agent@cluster.service` Systemd system unit, running as root. Its
  actions are defined in `/var/lib/nethserver/cluster/actions`

- `agent@node.service` Systemd system unit, running as root. Its actions
  are defined in `/var/lib/nethserver/node/actions`

- `agent.service` Systemd user unit, running as non-privileged Unix user
  for each rootless module instance. See the "Additional modules" section
  below for more details

- Edge proxy, for TLS termination and centralized certificates management (Traefik)

- LDAP proxy, a rootless module listening on 127.0.0.1. It helps other
  modules to connect to account provider LDAP servers, with a clear text
  connection

- LDAP local account provider (Samba DC, OpenLDAP)

- VPN, each node is connected to the leader using WireGuard in a star network topology

- API server, it handles authentication and authorization for all UI
  requests (and more...). It audits executed tasks

- UI, it allows configuration of the cluster and applications


### Modules

The core instantiates a set of *modules*. Each module instance
runs an application (e.g. Webtop, Nextcloud) as a set of one or more Podman **rootless containers**. In exceptional
cases a module can run rootfull containers (i.e Samba, Promtail)

See [Modules architecture](details.md#module-architecture) for more info.

---
Next: [Installation](installation.md)
