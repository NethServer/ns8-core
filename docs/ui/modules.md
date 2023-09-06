---
layout: default
title: Modules UI
nav_order: 4
parent: User Interface
---

# Modules User Interface

- TOC
  {:toc}

## Guidelines

Users can open modules UI from the App drawer or using Global search.

Every NS8 module follows the same UI guidelines in order to provide a uniform user experience. Almost every module has at least these standard pages:

- Status
- Settings
- About

Status page is the landing page of the module, it should provide a dashboard displaying the current status of the module, including instance name, installation node and information about module systemd services.

Settings page should contain a form to set and review module configuration.

About page should provide module meta-information such as documentation URL, source code and author information.

NS8 modules make use of UI library components and functions, e.g. by including the shared `NsButton` component in a module form, or requesting the execution of a cluster task.

The UI of a module lives inside an &lt;iframe&gt; of core UI (see `core/ui/src/views/Applications.vue`). It has a reference to core VueJs instance, named `ns8Core`, used to:

- Events management (e.g. register to task completion)
- Access to core API URL and other core configurations
- Access to core i18n strings

Modules use [Vue Router](https://router.vuejs.org/) to implement page routing and [Vuex](https://vuex.vuejs.org/) to handle global state.

## Module UI development

Follow the steps below to prepare the development environment for coding modules UI:

- Install NethServer 8 on a development machine
- Install development tools on your workstation:
  - Node.js and npm (LTS version, currently v16)
  - Yarn
- Creating a new module? Start from NS8 kickstart: [ns8-kickstart](https://github.com/NethServer/ns8-kickstart)
- Sync module sources on your development machine:
  - `cd PATH_TO_MODULE/ui`
  - `yarn install` (needed only the first time)
  - `yarn watch`
  - on you NS8 machine: `sshfs USER@WORKSTATION:PATH_TO_MODULE/ui/dist/ /var/lib/nethserver/cluster/ui/apps/your_module1/ -o allow_other,default_permissions`
