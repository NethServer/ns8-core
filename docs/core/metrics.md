---
layout: default
title: Metrics and alerting
nav_order: 16
parent: Core
---

# Metrics and alerting

The core provides a module named `metrics` that collects metrics from the nodes
and sends alerts to the Nethesis portal. The module is installed by the cluster
initialization script.

The module includes the following services:

- [Prometheus](https://prometheus.io/)
- [Alertmanager](https://prometheus.io/docs/alerting/alertmanager/)
- [Grafana](https://grafana.com/)
- [alert-proxy](alert-proxy/README.md)

Key points:

- Single instance running on the leader node
- Monitors all cluster nodes automatically
- Removed if the leader node becomes a worker
- Prometheus: port 9091
- Alertmanager: port 9093
- alert-proxy: port 9095
- Grafana: port 3000 (disabled by default, enabled with Traefik route)

Configuration:

- Prometheus and Alertmanager configurations are recreated on Prometheus restart
- Module restarts on node addition/removal
- alert-proxy restarts on subscription change to send alerts to Nethesis portals

If a subscription is enabled, alerts are sent to Nethesis portals by default.

For more information, please refer to the module [README](https://github.com/NethServer/ns8-metrics) for details on alert configuration and customization.