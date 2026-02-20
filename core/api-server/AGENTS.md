# API Server

Gin-based REST API serving the admin UI. Provides JWT authentication, WebSocket
real-time updates, task management, and SQLite audit logging.

## Build

```bash
go build .                   # build api-server binary
go build -tags api_server_logs -o api-server-logs api-server-logs.go  # logs CLI tool
```

Built with `CGO_ENABLED=1` (required by `go-sqlite3`). The `build-image.sh`
script uses `-ldflags '-linkmode external -extldflags "-static"'` for static
linking.

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
─── JWT middleware applied below ───
POST /api/cluster/tasks                   # cluster task operations
POST /api/node/:node_id/tasks             # node task operations
POST /api/module/:module_id/tasks         # module task operations
GET  /api/audit                           # audit log queries
*    /api/2FA                             # 2FA management
/ws                                       # WebSocket (Melody)
```

## Authentication

- **Primary**: Redis ACL credentials validated on login
- **JWT**: 14-day expiry; claims contain username, role, actions array, 2FA flag
- **Authorization**: Action-based (not role-based) using `filepath.Match`
  wildcard patterns against the user's authorized actions list
- **2FA**: Optional TOTP via `github.com/pquerna/otp`; bypass for GET requests

## Configuration

All via environment variables: `LISTEN_ADDRESS`, `REDIS_ADDRESS`, `REDIS_USER`,
`REDIS_PASSWORD`, `SECRET` (JWT key), `STATIC_PATH`, `AUDIT_FILE`,
`SENSITIVE_LIST`.

## Conventions

- Response envelopes use `structs.Map()` for struct-to-JSON conversion.
- Redis keys follow `{cluster,node,module}/<id>/tasks/<task-id>` hierarchy.
- Audit events: `login-ok`, `login-fail`, `auth-fail` stored in SQLite.
- WebSocket sessions are validated on each ping-pong (JWT expiry check).
- No unit tests in this package; testing is via integration tests in
  `core/tests/`.
