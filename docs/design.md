---
layout: default
title: Design & Architecture
nav_order: 4
---

# Design & Architecture

* TOC
{:toc}

## Goals

The design of NethServer 8 (NS8) should fulfill these goals:

- solve sysadmin daily problems, like:
  - no-brain software installation
  - migrate software to another machine
  - easy and effective data backup and restore
  - prepare a machine in a lab, deploy at customer's office
- limit upstream dependency, increase upstream independence
- NS 7 migration (better, upgrade) path
- centralized management portal, installable on premise and as SaaS service

## Assumptions

While Kubernetes is great for managing applications inside big clusters, NS8 is designed
for small and medium business, with limited resources.

Assumptions:

- containers are the new standard
- built for cheap hardware or entry-level Virtual Private Server (VPS)
- few nodes: 1 node mostly, 2 or 3 nodes sometimes
- nodes can be on different geographical areas and the network could have poor latency

## Glossary

NS8 is composed by many parts, let's explain the main terminology.

- *Node*: a physical or virtual installation of NS8, it can be a *leader*, *worker* or both for all-in-one installations
- *Cluster*: a set of one or more nodes
- *Leader*: the node which executes the control logic and has read-write access to the configuration.
- *Worker*: a node which execute one or more modules, it has read-only access to the configuration. A leader node can be  a worker too.
- *Core*: set of fundamental software which is usually installed on all cluster nodes
- *Module*: additional software that can be installed on one or more cluster nodes, like Mail server, Nextcloud, Webtop.
- *Instance*: a running instance of a module with a specific configuration. Each instance runs inside an isolated environment and constituted by one or more Podman containers.
- *Application*: speaking in UI terms, it corresponds to one or more instances.

## Design choices

(almost) Everything runs inside a container.

Design goals:

- store the configuration on a reliable cluster-ready database: one read-write Redis master node for the whole cluster, multiple read-only replica nodes
- use TLS terminator whenever possible to simplify module configuration
- the configuration should be readable for any module inside the system, sensitive data (like passwords) should be store encrypted
- security as first class citizen: rootless container modules with SELinux / AppArmor, rootfull allowed too

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
- Authenticated Redis access for write operations
- Public Redis read only access
- Edge proxies with TLS termination
- Centralized certificate management
- FHS compliant
- API server
- UI structure

### Operating System (OS) requirements

NS8 should run on a OS with the following features:

- small footprint, x86_64 (and aarch64)
- systemd to keeps things running
- podman ≥ 2.2
- support WireGuard VPN
- cgroups v2
- firewalld to manage ports for the public network

### Users

As design choice, NS8 will use two different user bases for the core and the applications.

The core, including the administrator UI, will authenticate only against Redis.
Each module instance will use local Unix accounts.

Applications will be able to access users from the same account provider available on NS7:
- local and remote OpenLDAP
- local and remote AD

Account providers will be accessible using a LDAP proxy.
Each module instance can eventually select a different account provider.

## Components

The system is composed by two main components:
- core
- modules

### Core

The core purpose is managing the applications, providing the basics for their entire life cycle (install, upgrade, reconfigure, uninstall...).
It runs the following components:

- Redis [database and message bus](core/database.md)
- Node, cluster and module [agents](core/agents.md) written in Golang
- Traefik as [edge proxy](core/proxy_certificates.md), for TLS termination and centralized certificates management
- [LDAP proxy](core/user_domains.md), a rootless module listening on 127.0.0.1. It helps other
  modules to connect to account provider LDAP servers, with a clear text connection for local containers
- LDAP local account provider: [Samba DC](https://github.com/NethServer/ns8-core/blob/main/samba/README.md), OpenLDAP (not implemented yet)
- [VPN](core/vpn.md), each node is connected to the leader using WireGuard in a star network topology
- [API server](core/api_server.md), it handles authentication and authorization for UI and cli requests, it also audits executed tasks
- [UI](ui/index.md), it allows configuration of the cluster and applications


### Modules

The core instantiates a set of *modules*. Each module instance runs an
application (e.g. Webtop, Nextcloud) or it can be just a part of it (e.g.
Ldapproxy OpenLDAP and Samba are parts/components of User domains).

In any case a module instance is made of one or more Podman rootless containers.
In exceptional cases a module can run rootfull containers (i.e node_exporter, Crowdsec).

See [Modules architecture](modules/index.md) for more info.

