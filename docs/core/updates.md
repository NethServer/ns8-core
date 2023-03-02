---
layout: default
title: Updates
nav_order: 12
parent: Core
---

# Core updates

* TOC
{:toc}

## NS8 updates

Core updates include enhancements and bug fixes for the `core` image. Such image contains: agents, UI and common actions.

To execute the core update, use:

    api-cli run update-core --data '{"core_url":"ghcr.io/nethserver/core:latest","nodes":[1, 2]}'

The cluster agent forwards the request to selected nodes.
Each node agent downloads the new image and restart core services: agents, redis and api-server.
The core services restart occurs only if the core version update has at least a minor release increment.

The following applications are part of the core but can be updated as normal modules:

- ldapproxy
- traefik
- samba
- loki
- promtail

See [module updates]({{site.baseurl}}/modules/updates) for more info.

## Distribution updates

Each distribution has its own life-cycle and should be kept updated to handle security fixes.

Simple and dirty commands for Debian:

    apt-get update && apt-get upgrade

Simple and dirty command for CentOS Stream:

    dnf update

Please refer to [Debian](https://www.debian.org/doc/manuals/debian-faq/uptodate.en.html) and CentOS official documentation for more info.
