# NethServer 8: cloud in a box

NethServer 8 (NS8), is the evolution of NethServer 7 (NS7).

## Goals

The design of NS 8 should fulfill these goals:

- solve sysadmin daily problems, like:
  - no-brain software installation
  - migrate software to another machine
  - easy and effective data backup and restore
  - prepare a machine in a lab, deploy at customer's office
- limit upstream dependency, increase upstream independence
- NS7 migration (better, upgrade) path
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
- *Application*: speaking in UI terms, it corresponds to a module and its instances.
