---
layout: default
title: User Interface
nav_order: 6
has_children: true
---

# User Interface

NS8 user interface is a web application developed with [VueJs](https://vuejs.org/) and based on [Carbon Design System](https://www.carbondesignsystem.com/).

NS8 UI can be accessed at `https://leader_node/cluster-admin/`. Default username is `admin` and default password is `Nethesis,1234`.

The web application source is composed by multiple VueJs projects:
- Core UI
- A distinct UI project for every module (e.g. Dokuwiki)
- A UI library that includes a set of reusable UI components and functions (VueJs mixins) used by core and modules UI

Source code of core UI is provided here: [https://github.com/NethServer/ns8-core/tree/main/core/ui](https://github.com/NethServer/ns8-core/tree/main/core/ui)

## UI design

[Carbon grid system](https://www.carbondesignsystem.com/guidelines/2x-grid/implementation/) promotes responsive design. A simple way to develop a responsive layout is to organize content by placing [tiles](https://www.carbondesignsystem.com/components/tile/usage) inside grid columns, using:

- Carbon `CvTile` component
- UI library `NsTile`
- UI library `Ns*Card` components (e.g. `NsInfoCard`, `NsStatusCard`, ...)

