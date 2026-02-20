# Copilot Instructions for ns8-core

> **Source of truth:** The canonical version of these instructions lives in
> `.github/copilot-instructions.md`. This `AGENTS.md` file is a duplicated/
> generated copy. To avoid drift, **edit only** `.github/copilot-instructions.md`
> and regenerate or sync this file as needed.
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
| `core/imageroot/` | Python 3 / Bash | Action scripts and the Python `agent` package, for action and event handler implementation |
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
cd core/api-server && go build .
```

UI development:

```bash
cd core/ui
yarn install
yarn serve          # dev server
yarn build          # production build
```

## Lint

```bash
cd core/ui && yarn lint
```

Bash scripts should pass `shellcheck`.

## Test

**Go agent tests** (runs in containers via podman):

```bash
cd core/agent && bash test-agent.sh
```

**Integration tests** (Robot Framework, requires a live node):

```bash
export LEADER_NODE=<ssh-address>
cd core && bash test-module.sh
```

Robot Framework tests are in `core/tests/` with numbered directories
(`10__cluster_sanity/`, `20__cluster_ui/`) and shared `keywords.resource`.

## Conventions

- **Action script numbering**: `00` = validation, `50` = main logic, `60+` =
  post-processing. Each step is a separate executable (Python 3 or Bash).
- **JSON Schema validation**: Actions define `validate-input.json` /
  `validate-output.json` for automatic request/response validation.
- **Redis as state store**: All persistent state is in Redis, not on the
  filesystem. Actions use `agent.redis_connect()` to read/write state.
- **UI design system**: Carbon Design System (`@carbon/vue`). Use Carbon
  components, not custom HTML widgets.
- **UI mixins**: Cross-cutting concerns (WebSocket, notifications, tasks, login,
  2FA, audit) are Vue mixins in `core/ui/src/mixins/`.
- **Event ID pattern**: Always generate a unique `eventId` per task creation to
  prevent event handler collisions:
  ```js
  const eventId = this.getUuid();
  // use `${taskAction}-completed-${eventId}` for event names
  ```
- **Go modules**: `core/agent` and `core/api-server` are independent Go modules
  with separate `go.mod` files.
- **Translations**: Managed via Weblate at
  `https://hosted.weblate.org/projects/ns8/`. Do not manually edit translation
  files.
- **`/etc/hosts` handling**: The `set-fqdn` action must not comment out
  localhost entries; instead remove old hostnames from those lines and normalize
  whitespace.
