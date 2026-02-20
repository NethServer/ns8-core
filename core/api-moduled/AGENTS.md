# API Moduled

Lightweight, filesystem-driven HTTP API server for individual NS8 modules.
It wraps external executables (shell scripts, Python) with JWT auth and
JSON Schema validation. Each module instance runs its own `api-moduled`
process.

Known applications using `api-moduled` are ns8-samba and ns8-openldap.
They use it to implement a user self-service portal (create/edit/delete
users and groups of the LDAP DB). 

## Build

```bash
go build .                   # build with CGO_ENABLED=0 (no C deps)
```

## Module layout

The entire server is two files:

| File | Purpose |
|------|---------|
| `api-moduled.go` | Main entry — routing, JWT setup, handler dispatch (~370 LOC) |
| `validation/validation.go` | JSON Schema validation with external definitions support |

## How it works

Handlers are **directories on disk**, not Go code. The URL path maps directly
to a handler directory:

```
POST /api/myhandler  →  executes  handlers/myhandler/post
```

Each handler directory can contain:
- `post` — executable to run (receives JSON on stdin, returns JSON on stdout)
- `validate-input.json` — optional input schema
- `validate-output.json` — optional output schema

### Request flow

```
POST /api/:handler
  → validate input against validate-input.json
  → execute handlers/:handler/post with JSON stdin
  → validate output against validate-output.json
  → return JSON response
```

### Environment injected into handlers

- `JWT_ID` — authenticated user identity
- `JWT_CLAIMS` — full JWT claims as JSON
- `AMLD_EXPORT_ENV` — custom variables from config

## Authentication

- **Login**: `POST /api/login` calls `handlers/login/post` which must return
  JSON claims (including a `uid` field). Exit code 2–7 = bad credentials.
- **JWT**: Configurable timeout (default 4h). Scope claim is an optional array
  of allowed handler names (empty = unrestricted).
- **Basic auth**: `GET /api/auth` validates credentials without issuing a token.

## Configuration

All via `AMLD_` prefixed environment variables (uses Viper):

| Variable | Default | Purpose |
|----------|---------|---------|
| `AMLD_HANDLER_DIR` | `./handlers/` | Handler directory root |
| `AMLD_BIND_ADDRESS` | `:9313` | Listen address |
| `AMLD_JWT_TIMEOUT` | `4h` | Token expiry |
| `AMLD_ID_KEY` | `uid` | JWT identity claim key |
| `AMLD_EXPORT_ENV` | — | Extra env vars for handlers |

## Key differences from api-server

- **No Redis dependency** — auth is delegated to handler scripts
- **No SQLite** — no audit log
- **No WebSocket** — request/response only
- **Generic** — knows nothing about NethServer; any executable can be a handler
- **Stateless** — single static binary, trivial to deploy
