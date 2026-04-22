# ns8-core

This is production-grade project with thousands of installations. Breaking
changes are not allowed.

## Architecture

NethServer 8 is a container-based Linux platform orchestrated via Redis. This
repository contains the **core** — cluster/node management, API server, agent
framework, and admin UI. External application modules live in separate
`ns8-<app>` repositories.

### Components

| Path | Language | Purpose |
|------|----------|---------|
| `core/agent/` | Go | Agent daemon — reads tasks from Redis, handles events, executes actions, on behalf of some system component (cluster, node, application) |
| `core/api-server/` | Go | Gin-based REST API with JWT auth, WebSocket, SQLite audit log |
| `core/api-moduled/` | Go | Generic HTTP API for individual modules (JWT, JSON Schema validation) |
| `core/ui/` | Vue 2 + JS | Admin web UI — Carbon Design System, Vuex, Vue Router, WebSocket |
| `core/imageroot/` | Python 3 / Bash | Action scripts and the Python `agent` package (at `usr/local/agent/pypkg/agent/`), for action and event handler implementation |
| `core/support/` | Shell | Remote support session via OpenVPN |
| `core/restic/`, `core/rsync/` | Shell | Backup and replication container wrappers |
| `core/tests/` | Robot Framework | Integration tests (cluster sanity, UI) |
| `docs/` | Markdown (Jekyll) | Developer manual at nethserver.github.io/ns8-core |

### How actions work

Actions are the primary unit of work. They live under
`core/imageroot/var/lib/nethserver/{cluster,node}/actions/<action-name>/` and
consist of numbered step scripts (`00validate_input`, `50update`, etc.) executed
sequentially by the Go agent. Each action directory may include
`validate-input.json` and `validate-output.json` for JSON Schema validation.

Python action scripts follow this pattern:

```python
#!/usr/bin/env python3
import json, sys, agent, agent.tasks

request = json.load(sys.stdin)
rdb = agent.redis_connect(privileged=True)
# ... business logic ...
json.dump(result, sys.stdout)
```

Key `agent` module APIs: `redis_connect()`, `set_env()`, `read_envfile()`,
`set_progress()`, `run_helper()`, `agent.tasks.run()`, `agent.tasks.runp()`.

### UI task/event model

The Vue UI communicates via WebSocket. Components invoke server-side actions
through `TaskService` mixin methods (e.g., `createModuleTaskForApp`). Each task
gets a unique `eventId` (via `getUuid()`) to avoid event handler collisions.
Progress updates arrive over the WebSocket and flow through the notification
mixin into the Vuex store.

## Build

Container images are built with **buildah** (not Docker):

```bash
cd core && bash build-image.sh
```

Go binaries are compiled with `CGO_ENABLED=0` (agent, api-moduled) or
`CGO_ENABLED=1 -static` (api-server, for SQLite). There are no Makefiles; the
build script is the source of truth.

For local Go development:

```bash
cd core/agent && go build .
cd core/api-server && go build api-server.go
```

UI development:

```bash
cd core/ui
yarn install
yarn serve          # dev server
yarn build          # production build
```

## Lint

For UI/Frontend:

```bash
cd core/ui && yarn lint
```

For Backend:

Bash scripts should pass `shellcheck`.

## Test

**Go agent tests** (runs in containers via podman):

```bash
cd core/agent && bash test-agent.sh
```

**Integration tests** (Robot Framework, requires a live node):

```bash
cd core
run-ns8-tests <LEADER_NODE> [robot_options...]
```

Robot Framework tests are in `core/tests/` with numbered directories
(`10__cluster_sanity/`) and shared `keywords.resource`. The
`run-ns8-tests` helper is provided by the `ns8-github-actions` repository;
you also need Podman in `PATH` and an SSH private key that can access the
target leader node.

## Conventions

- **Action script numbering**: By convention, `00` = validation, `50` = main
  logic, `60+` = post-processing, but intermediate numbers (e.g. `01`, `05`,
  `10`, `20`, `30`) are commonly used. Each step is a separate executable
  (Python 3 or Bash).
- **JSON Schema validation**: Actions define `validate-input.json` /
  `validate-output.json` for automatic request/response validation.
- **Redis as state store**: All persistent state is in Redis, not on the
  filesystem. Actions use `agent.redis_connect()` to read/write state.
- **UI design system**: Carbon Design System (`@carbon/vue`). Use Carbon
  components, not custom HTML widgets.
- **UI mixins**: Cross-cutting concerns (WebSocket, notifications, login,
  2FA, audit, node) are Vue mixins in `core/ui/src/mixins/`. Task-related
  services (`TaskService`, `UtilService`, etc.) come from the external
  `@nethserver/ns8-ui-lib` package.
- **Event ID pattern**: Always generate a unique `eventId` per task creation to
  prevent event handler collisions:
  ```js
  const eventId = this.getUuid();
  // use `${taskAction}-completed-${eventId}` for event names
  ```
- **Go modules**: `core/agent`, `core/api-server`, and `core/api-moduled` are
  independent Go modules with separate `go.mod` files. Use `go fmt` before
  commits.
- **Translations**: Managed via Weblate at
  `https://hosted.weblate.org/projects/ns8/`. Do not manually edit translation
  files.
- **`/etc/hosts` handling**: The `set-fqdn` action comments out lines
  containing the old hostname (prefixed with `# commented by set-fqdn #`)
  and appends a new `127.0.1.1` entry for the new hostname.
- **Branch names**: never use "/" in branch names. Use only chars allowed
  by container registry tags, like "-" and alphanumeric chars. This is a
  requirement for container image uploads.
- **Commits**: Use conventional commit style. Short title line, 50 chars max.
  Wrap body text at 72.
