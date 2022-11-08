# NethServer 8 UI

To develop NethServer 8 UI you need a working NS8 installation. Edit `public/config/config.development.js` by configuring `API_ENDPOINT` and `WS_ENDPOINT` with the IP address of your NS8 leader node.

You can develop NS8 UI on your workstation or inside a container.

## Develop NS8 UI on your workstation

- `yarn install`: project setup
- `yarn serve`: start development server with hot-reload
- `yarn build`: compiles and minifies for production
- `yarn storybook`: start Storybook webapp with hot-reload

## Develop NS8 UI inside a container

You have two options:
- Build and start a Podman container
- Use VS Code Dev Containers

### Build and start a Podman container

Build the container defined by `Containerfile`:

```
podman build -t ns8-core-dev .
```

Start development server (`--network=host` is required for hot-reload):
```
podman run -ti -v $(pwd):/app:Z --network=host --name ns8-core --replace ns8-core-dev serve
```

Compiles and minifies for production:
```
podman run -ti -v $(pwd):/app:Z --name ns8-core --replace ns8-core-dev build
```

Start Storybook webapp (`--network=host` is required for hot-reload):
```
podman run -ti -v $(pwd):/app:Z --network=host --name ns8-core --replace ns8-core-dev storybook
```

Note: if you want to run development server AND run Storybook webapp at the same time you can't use above commands; you would get a `yarn` error.
To run development server and run Storybook simultaneously:
```
podman run -ti -v $(pwd):/app:Z --network=host --name ns8-core --replace ns8-core-dev serve
```
and then
```
podman exec -ti ns8-core yarn storybook
```

### Use VS Code Dev Containers

- Install **Dev Containers** VS Code extension
- Dev Containers uses Docker by default but you can configure it to use Podman: go to `File > Preferences > Settings`, search `dev containers docker path` and type `podman` as `Docker path`
- Open `ns8-core` directory (the repository root) in VS Code, if you haven't already
- Open Command Palette (`CTRL+SHIFT+P`) and type `Reopen in Container`
- Open VS Code integrated terminal: `View > Terminal`
- Change directory to `core/ui`
- Enter one of the commands listed in [Develop NS8 UI on your workstation](#develop-ns8-ui-on-your-workstation), e.g. `yarn serve`

Container configuration is contained inside `.devcontainer/devcontainer.json`.
