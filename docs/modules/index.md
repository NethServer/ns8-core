---
layout: default
title: Modules
nav_order: 5
has_children: true
---

# Modules: implementation details

A module is an NS8 application which is usually composed by a backend and
a frontend.

Every time a module instance is added to the cluster, the new instance
is named as the module itself followed by a progressive number
starting from 1. Given a module named `myapp`, instances will be named
`myapp1`, `myapp2`, etc. This uniquely identifies the instance inside the cluster.

Modules can be managed using these commands:

- `add-module [--force] <module> <node_id>`: install a module on the given node with ID `node_id`; search for `module` inside enabled repositories and install
  latest available version. If `module` is a image registry URL, just install the module straight from it; this method can be used to install
  customized images. Like the Software Center, the command refuses a node that is not eligible, for example because of
  the `org.nethserver.max-per-node` or `org.nethserver.min-core` image labels; `--force` skips this check.
- `remove-module [--no-preserve] <module>`: remove an installed module; if `--no-preserve` is given, erase also all module data

