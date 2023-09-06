---
layout: default
title: UI library
nav_order: 3
parent: User Interface
---

# UI library

Since core and modules UI share multiple components and features, these facilities have been aggregated into NPM package [@nethserver/ns8-ui-lib](https://www.npmjs.com/package/@nethserver/ns8-ui-lib).
The reusable UI components included in the library are easily recognizable since their name starts with `Ns` prefix, e.g. `NsButton`, `NsInlineNotification`.
The library also includes a set of VueJs mixins used by core and modules to access utility functions and perform common tasks.

Source code of UI library is provided [here](https://github.com/NethServer/ns8-ui-lib).

## UI library development

You can develop NS8 UI library inside a container (recommended) or on your workstation.

### Develop NS8 UI library inside a container

You have two options:

- Build and start a Podman container, or
- Use VS Code Dev Containers

#### Build and start a Podman container

Build the container defined by `Containerfile`:

```
podman build -t ns8-ui-lib-dev .
```

Compile and minify for production:

```
podman run -ti -v $(pwd):/app:Z --name ns8-ui-lib --replace ns8-ui-lib-dev build
```

Publish NPM package:

```
podman run -ti -v $(pwd):/app:Z --name ns8-ui-lib --replace ns8-ui-lib-dev publish
```

#### Use VS Code Dev Containers

- Install Visual Studio Code extension [Dev Containers](https://marketplace.visualstudio.com/items?itemName=ms-vscode-remote.remote-containers) (beware: this procedure may not work on [VSCodium](https://vscodium.com/))
- Dev Containers uses Docker by default but you can configure it to use Podman: go to `File > Preferences > Settings`, search `dev containers docker path` and type `podman` as `Docker path`
- Open `ns8-ui-lib` directory (the repository root) in VS Code, if you haven't already
- Open Command Palette (`CTRL+SHIFT+P`) and type `Reopen in Container`
- Open VS Code integrated terminal: `View > Terminal`
- Enter one of the following commands:
  - `npm install`: project setup, needed only the first time
  - `npm run build`: compile and minify for production
  - `npm publish`: publish NPM package

Container configuration is contained inside `.devcontainer/devcontainer.json`.

### Develop NS8 UI library on your workstation

Developing NS8 UI library inside a container is the recommended way, but if you want to do it on your workstation:

- Install Node.js and npm (LTS version, currently v18)
- `npm install`: project setup, needed only the first time
- `npm run build`: compile and minify for production
- `npm publish`: publish NPM package

## Import NS8 UI library components

To import the latest version of NS8 UI library:

```
yarn add @nethserver/ns8-ui-lib
```

### Import library components and mixins

To globally include NS8 components, edit `main.js` of the project importing the library:

```js
// main.js

import ns8Lib from "@nethserver/ns8-ui-lib";
Vue.use(ns8Lib);
```

To import a specific mixin inside a component (e.g. `TaskService` mixin) of the project importing the library:

```js
// ComponentName.vue

import { TaskService } from "@nethserver/ns8-ui-lib";

export default {
  name: "ComponentName",
  mixins: [TaskService],
  ...
}
```
