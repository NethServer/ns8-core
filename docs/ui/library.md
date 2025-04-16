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

- [Build and start a Podman container](#build-and-start-a-podman-container), or
- [Use VS Code Dev Containers](#use-vs-code-dev-containers)

#### Build and start a Podman container

Build the container defined by `Containerfile`:

```
podman build -t ns8-ui-lib-dev .
```

Project setup:

```
podman run -ti -v $(pwd):/app:Z --name ns8-ui-lib --replace ns8-ui-lib-dev npm install
```

Compile and minify for production:

```
podman run -ti -v $(pwd):/app:Z --name ns8-ui-lib --replace ns8-ui-lib-dev npm run build
```

Create a tarball:
```
podman run -ti -v $(pwd):/app:Z --name ns8-ui-lib --replace ns8-ui-lib-dev npm run build-pack
```
Tarball generation is useful for testing a development or testing version of NS8 UI library. To import the generated tarball into another Node project (e.g. `ns8-core` or a NS8 module), see [Import development or testing version](#import-development-or-testing-version).

#### Use VS Code Dev Containers

- Install Visual Studio Code extension [Dev Containers](https://marketplace.visualstudio.com/items?itemName=ms-vscode-remote.remote-containers) (beware: this procedure may not work on [VSCodium](https://vscodium.com/))
- Dev Containers uses Docker by default but you can configure it to use Podman: go to `File > Preferences > Settings`, search `dev containers docker path` and type `podman` as `Docker path`
- Open `ns8-ui-lib` directory (the repository root) in VS Code, if you haven't already
- Open Command Palette (`CTRL+SHIFT+P`) and type `Reopen in Container`
- Open VS Code integrated terminal: `View > Terminal`
- Enter one of the following commands:
  - `npm install`: project setup, needed only the first time
  - `npm run build`: compile and minify for production
  - `npm run build-pack`: create a tarball

Tarball generation is useful for testing a development or pre-release version of ns8-ui-lib. To import the generated tarball into another Node project (e.g. `ns8-core` or a NS8 module), see [Import development or testing version](#import-development-or-testing-version).

Dev container configuration is contained inside `.devcontainer/devcontainer.json`.

### Develop NS8 UI library on your workstation

Developing NS8 UI library inside a container is the recommended way, but if you want to do it on your workstation:

- Install Node.js 18 and npm
- `npm install`: project setup, needed only the first time
- `npm run build`: compile and minify for production
- `npm run build-pack`: create a tarball

Tarball generation is useful for testing a development or pre-release version of ns8-ui-lib. To import the generated tarball into another Node project (e.g. `ns8-core` or a NS8 module), see [Import development or testing version](#import-development-or-testing-version).

## Release a new version

Releasing a new version of `ns8-ui-lib` involves publishing a new NPM package version, creating a GitHub tag and generating a GitHub release. Note that the release command should be executed on your workstation, as your public SSH key is required to push the new tag to GitHub.

Run one of the following commands, depending on the type of version bump. The command will trigger a GitHub action named "Publish to NPM and create release", which will handle the actual release of the new version.

Release a patch version:

```
npm run publish:patch
```

Release a minor version:

```
npm run publish:minor
```

Release a major version:

```
npm run publish:major
```

## Import NS8 UI library components

### Import latest version

To import the latest version of NS8 UI library:

```
yarn add @nethserver/ns8-ui-lib
```

### Import development or testing version

To import a development or testing version of `ns8-ui-lib` from a tarball (generated with `npm run build-pack`) into another Node project (e.g. `ns8-core` or a NS8 module):

```
cd other/node/project
yarn add /path/to/ns8-ui-lib/nethserver-ns8-ui-lib-x.y.z.tgz
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
