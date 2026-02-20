# API Server

Gin-based REST API serving the admin UI. Provides JWT authentication, WebSocket
real-time updates, task management, and SQLite audit logging.

## Build

```bash
go build .                   # build api-server binary
go build -tags api_server_logs -o api-server-logs api-server-logs.go  # logs CLI tool
```

Built with `CGO_ENABLED=1` (required by `go-sqlite3`). The `build-image.sh`
script uses `-ldflags='-extldflags=-static'` for static linking and adds
`-tags sqlite_omit_load_extension`.

## Module layout

| Path | Purpose |
|------|---------|
| `api-server.go` | Main entry — route definitions, middleware wiring |
| `api-server-logs.go` | Cobra CLI tool for sending WebSocket log actions |
| `methods/tasks.go` | Task CRUD: create, list, get output/status/context |
| `methods/auth.go` | Redis ACL auth, 2FA/OTP setup and validation |
| `methods/audit.go` | Audit log queries (filter by user, action, date range) |
| `middleware/` | JWT lifecycle: authenticator, authorizer, payload builder |
| `models/` | Data structs: Audit, Event, Task, UserAuthorizations, SocketAction |
| `redis/` | Redis client singleton, leader host lookup, health checks |
| `audit/` | SQLite audit store (WAL mode, transactions) |
| `socket/` | Melody WebSocket — session management, Redis event relay |
| `response/` | Standardized HTTP response envelopes (StatusOK, StatusBadRequest, etc.) |
| `configuration/` | All config from environment variables (no config files) |

## Routing structure

```
POST /api/login                           # no auth
POST /api/logout                          # no auth
GET  /api/module/:module_id/http-basic/:action  # basic auth, for app integration
─── JWT middleware applied to /api routes below ───
POST /api/cluster/tasks                   # cluster task operations
GET  /api/cluster/tasks                   # list cluster tasks
POST /api/node/:node_id/tasks             # node task operations
POST /api/module/:module_id/tasks         # module task operations
GET  /api/nodes                           # list nodes
GET  /api/modules                         # list modules
GET  /api/audit                           # audit log queries
GET  /api/audit/users                     # audit log users
GET  /api/audit/actions                   # audit log actions
*    /api/2FA                             # 2FA management

/ws                                       # WebSocket (Melody) at root; not behind Gin JWT middleware (session validated in WebSocket flow)
```

## Authentication

- **Primary**: Redis ACL credentials validated on login
- **JWT**: 14-day expiry; claims contain username, role, actions array, 2FA flag
- **Authorization**: Action-based (not role-based) using `filepath.Match`
  wildcard patterns against the user's authorized actions list
- **2FA**: Optional TOTP via `github.com/pquerna/otp`; enforced at login. Middleware bypasses authorization checks (not 2FA) for all GET requests, and `/api/2FA` POST/DELETE endpoints are explicitly exempted in the authorizer.

## Configuration

All via environment variables: `LISTEN_ADDRESS`, `REDIS_ADDRESS`, `REDIS_USER`,
`REDIS_PASSWORD`, `SECRET` (JWT key), `STATIC_PATH`, `AUDIT_FILE`,
`SENSITIVE_LIST`.

## Conventions

- Response envelopes use `structs.Map()` for struct-to-JSON conversion.
- Redis keys follow `task/<entity>/<id>/<task-id>/{context,output,error,exit_code}` (e.g. `task/node/<node-id>/<task-id>/output`).
- Audit events: `login-ok`, `login-fail`, `auth-fail` stored in SQLite.
- WebSocket sessions are validated on each ping-pong (JWT expiry check).
- No unit tests in this package; testing is via integration tests in
  `core/tests/`.
