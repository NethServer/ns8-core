# Design

(almost) Everything runs inside a container.

Design goals:

- store the configuration on a reliable cluster-ready database: one read-write Redis master node for the whole cluster, multiple read-only nodes
- use TLS terminator whenever possibile to simplify module configuration
- the configuration should be readable for any module inside the system, sensitive data (like passwords) should be store encrypted
- security as first class citizen: rootless container modules with SELinux / AppArmor, rootfull allowed too

## Operating System (OS) requirements

NS8 should run on a OS with the following features:

- small footprint, x86_64 and aarch64
- systemd to keeps things running
- podman ≥ 2.2
- support WireGuard VPN

## Account providers

As design choice, NS8 will use two different user bases for the core and the applications.

The core, including the administrator UI, will authenticate only against Redis.
Each module instance will use local unix accounts.

Applications will be able to access users from the same account provider available on NS7:
- local and remote OpenLDAP
- local and remote AD

Account providers will be accessibile using a LDAP proxy.
Each module instance can eventually select a different account provider.
