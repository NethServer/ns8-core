---
layout: default
title: Modules UI
nav_order: 4
parent: User Interface
---

# Modules User Interface

* TOC
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

- Install NethServer 8 on a development machine
- Creating a new NS8 module? Start from [NS8 kickstart](https://github.com/NethServer/ns8-kickstart)
- Install the module on your development machine
- Replace `MODULE` with the actual module name, e.g. `mail1` in the following steps
- You can develop NS8 modules UI inside a container (recommended) or on your workstation.

### Develop NS8 modules UI inside a container

You have two options:

- Build and start a Podman container, or
- Use VS Code Dev Containers

#### Build and start a Podman container

- Build the container defined by `Containerfile`:

```
cd ui

podman build -t ns8-MODULE-dev .
```

- Compile sources and watch for changes:

```
podman run -ti -v $(pwd):/app:Z --name ns8-MODULE --replace ns8-MODULE-dev watch
```

- [Sync module sources to your NS8 development machine](#sync-module-sources-to-ns8-development-machine)

#### Use VS Code Dev Containers

- Install Visual Studio Code extension [Dev Containers](https://marketplace.visualstudio.com/items?itemName=ms-vscode-remote.remote-containers) (beware: this procedure may not work on [VSCodium](https://vscodium.com/))
- Dev Containers uses Docker by default but you can configure it to use Podman: go to `File > Preferences > Settings`, search `dev containers docker path` and type `podman` as `Docker path`
- Open `ns8-MODULE` directory (the repository root) in VS Code, if you haven't already
- Open Command Palette (`CTRL+SHIFT+P`) and type `Reopen in Container`
- Open VS Code integrated terminal: `View > Terminal`
- Change directory to `ui`
- Enter the following commands:
  - `yarn install`: project setup, needed only the first time
  - `yarn watch`: compile sources and watch for changes
- [Sync module sources to your NS8 development machine](#sync-module-sources-to-ns8-development-machine)

Container configuration is contained inside `.devcontainer/devcontainer.json`.

### Develop NS8 modules UI on your workstation

Developing NS8 modules UI inside a container is the recommended way, but if you want to do it on your workstation:

- Install Node.js and npm (LTS version, currently v18)
- Change directory to `ui`
- `yarn install`: project setup, needed only the first time
- `yarn watch`: compile sources and watch for changes
- [Sync module sources to your NS8 development machine](#sync-module-sources-to-ns8-development-machine)

### Sync module sources to NS8 development machine

To sync module sources from your workstation to your NS8 development machine:

- Enter the following commands on your workstation (outside any container):

```
cd PATH_TO_MODULE/ui

watch -n 1 rsync -avz --delete dist/* root@YOUR_NS8_MACHINE:/var/lib/nethserver/cluster/ui/apps/MODULE/
```

- Open your browser at `https://YOUR_NS8_MACHINE/cluster-admin/#/apps/MODULE`
- Edit module sources
- Reload browser page to see your changes (disable browser cache if needed)
