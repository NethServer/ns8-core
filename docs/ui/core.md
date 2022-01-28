---
layout: default
title: Core UI
nav_order: 2
parent: User Interface
---

# Core User Interface

* TOC
{:toc}

## Components

Core UI includes the following components:
- Login page
- Shell
  - Side menu
  - Top header
  - Global search
  - Notification drawer
  - App launcher
- Cluster status page
- Software center
- Cluster logs (auditing)
- Cluster settings

Source code of core UI is provided [here](https://github.com/NethServer/ns8-core/tree/main/core/ui).

Core UI also includes a [Storybook](https://storybook.js.org/) to explore and test the reusable components included in the UI library.
To launch Storybook web app:

```bash
cd core/ui
yarn storybook
```

Application routing is implemented with [Vue Router](https://router.vuejs.org/). See `core/ui/src/router/index.js`

Core UI uses [Vuex](https://vuex.vuejs.org/) to handle the global state of the application. See `core/ui/src/store/index.js`.

Source code of core UI is provided here: [https://github.com/NethServer/ns8-core/tree/main/core/ui](https://github.com/NethServer/ns8-core/tree/main/core/ui)

## Shortcuts

The core and the applications can add quickly accessible actions called shortcuts.
Shortcuts will be available from the search bar inside the UI.
Examples of shortcuts could be something like "Add a user" or "Show application X logs".

Shortcuts are described inside a `shortcuts.json` file which must be inside the `ui` path:
- core: `/var/lib/nethserver/cluster/ui/shortcuts.json`
- applications: `/var/lib/nethserver/cluster/ui/apps/<app>/shortcuts.json`

A `shorcuts.json` must always contain the following fields:
- `name`: the action name
- `description`: a text which describes the action
- `tags`: a list of tags to ease the research
- `path`: the UI path to open when the shortcut has been selected. The path **must** be relative.

The `list-shortcuts` API will read all `shortcuts.json` files and combine them on a single array and it will generate absolute UI paths.

## Tasks and notifications

Core and modules can submit tasks to the cluster, to a node or to a module.
NS8 provides automatic toast notifications for task creation, progress, completion and abortion; task events are sent from API server to UI through web-socket.
Validation errors can be intercepted as well, in order to highlight invalid input fields or do something after validation is successful (e.g. close a modal).

To execute a task silently and prevent automatic notifications on task events, set the following attribute in task object:

```
extra.isNotificationHidden: true
```

Tasks that only read data should always be silent.

## Error notifications

Errors are notified to the user using inline notifications (`NsInlineNotification` component) or toast notifications (`NsToastNotification` component). Use inline notifications whenever possible, showing the error message in the exact place of the UI where the error occurred. Use toast notifications when there is no proper place to display the error (e.g. background or long running tasks).

## Query string parameters

You can sync the state of some components inside page URL. This is useful in the following cases:

- You share the URL you are visiting with someone (they will see the same UI state)
- You want to link to a specific point and state of the web app, so that you can browse to a page and prefill some input fields, or open a modal automatically. Global shortcuts rely on this feature

Do not use query parameters for data that is not under user control, e.g. data retrieved by the backend.

Modules must use query parameters for their internal navigation, using view data `q.page`.

## Core UI development

Follow the steps below to prepare the development environment for coding core UI:

- Install NethServer 8 on a development machine
- Edit `core/ui/public/config/config.development.js`:
  - Configure `API_ENDPOINT` and `WS_ENDPOINT` with the IP address of the NethServer 8 leader node
- Install development tools on your workstation:
  - Node.js and npm (LTS version, currently v16)
  - Yarn
- Run a web server on your workstation (hot reloading enabled):
  - `cd core/ui`
  - `yarn install` (needed only the first time)
  - `yarn serve`

