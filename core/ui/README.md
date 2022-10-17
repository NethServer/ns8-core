# NethServer 8 UI

To develop NethServer 8 UI you need a working NS8 installation. Edit `public/config/config.development.js` by configuring `API_ENDPOINT` and `WS_ENDPOINT` with the IP address of your NS8 leader node.

You can develop NS8 UI on your workstation or inside a container.

## Develop NS8 UI on your workstation

- `yarn install`: project setup
- `yarn serve`: start development server

## Develop NS8 UI inside a container

You have two options:
- Build and start a Podman container
- Use VS Code Dev Containers

### Build and start a Podman container

Build and start the container defined by `Containerfile`:

- `podman build -t ns8-core-dev .`
- `podman run -ti -p 8080:8080 -v $(pwd):/app:Z ns8-core-dev`

### Use VS Code Dev Containers

- Install **Dev Containers** VS Code extension
- Dev Containers uses Docker by default but you can configure it to use Podman: go to `File > Preferences > Settings`, search `dev containers docker path` and type `podman` as `Docker path`
- Open `ns8-core` directory (the repository root) in VS Code, if you haven't already
- Open Command Palette (`CTRL+SHIFT+P`) and type `Reopen in Container`
- Open VS Code integrated terminal, change directory to `core/ui` and then type `yarn serve` to start development server

Container configuration is contained inside `.devcontainer/devcontainer.json`.
