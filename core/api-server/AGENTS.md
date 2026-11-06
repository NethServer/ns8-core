# API Server

Gin-based REST API serving the admin UI. Provides JWT authentication, WebSocket
real-time updates, task management, and SQLite audit logging.

## Build

```bash
go build api-server.go       # build api-server binary
go build api-server-logs.go  # logs CLI tool
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
| `utils/` | Error logging, string helpers, mutex-guarded clock for session expiry |
| `redis/` | Redis client factory (new client per call), leader host lookup, health checks |
| `audit/` | SQLite audit store (WAL mode, transactions) |
| `socket/` | Melody WebSocket — session management, Redis event relay |
| `response/` | Standardized HTTP response envelopes (StatusOK, StatusBadRequest, etc.) |
| `configuration/` | All config from environment variables (no config files) |

## Routing structure

```
POST /api/login                                        # no auth
POST /api/logout                                       # no auth; ends the session
POST /api/refresh                                      # no auth; refresh token -> new token pair
GET  /api/module/:module_id/http-basic/:action         # basic auth, for app integration
─── JWT and session check middlewares applied to /api routes below ───
POST /api/logout-others                                # end the other sessions of the user
POST /api/cluster/join-token                           # token for a node join code (needs add-node)
GET  /api/cluster/tasks                                # list cluster tasks
GET  /api/cluster/task/:task_id/status                 # task status
GET  /api/cluster/task/:task_id/context                # task context
POST /api/cluster/tasks                                # create cluster task
GET  /api/nodes                                        # list nodes
GET  /api/node/:node_id/tasks                          # list node tasks
GET  /api/node/:node_id/task/:task_id/status           # node task status
GET  /api/node/:node_id/task/:task_id/context          # node task context
POST /api/node/:node_id/tasks                          # create node task
GET  /api/modules                                      # list modules
GET  /api/module/:module_id/tasks                      # list module tasks
GET  /api/module/:module_id/task/:task_id/status       # module task status
GET  /api/module/:module_id/task/:task_id/context      # module task context
POST /api/module/:module_id/tasks                      # create module task
GET  /api/audit                                        # audit log queries
GET  /api/audit/users                                  # audit log users
GET  /api/audit/actions                                # audit log actions
POST /api/2FA                                          # enable 2FA
GET  /api/2FA                                          # get 2FA status
DELETE /api/2FA                                        # disable 2FA
GET  /api/2FA/qr-code                                  # generate QR code

/ws                                                    # WebSocket (Melody); auth via "authorize" action with JWT payload
```

## Authentication

- **Primary**: Redis ACL credentials validated on login
- **JWT**: claims contain username, 2FA flag, and placeholder
  `role`/`actions` fields (always empty). Actual role and actions are fetched
  fresh from Redis on every request via `IdentityHandler`. Lifetime: 15
  minutes for administrators, 14 days for agent accounts (`cluster`,
  `node/*`, `module/*`), 1 hour for join tokens (`"scope": "join"`, limited
  to `add-node`/`list-actions` by `joinTokenAllowed`).
- **Sessions** (`sessions/`): an administrator login starts a session in a
  local SQLite file (`SESSIONS_FILE`), never in Redis because it is
  replicated to worker nodes. Tokens carry its ID in the `sid` claim, checked
  on every request by `SessionCheck` and on WebSocket ping-pong. Refresh
  tokens are stored as SHA-256 hashes and replaced at each use; sessions end
  after 12 hours without refresh or 14 days after login. The
  `cluster/event/user-credentials-changed` event ends the sessions of a user.
- **Authorization**: Action-based (not role-based) using `filepath.Match`
  wildcard patterns against the user's authorized actions list (from Redis).
  All GET requests bypass authorization checks.
- **2FA**: Optional TOTP via `github.com/dgryski/dgoogauth`; enforced at
  login. `/api/2FA` POST/DELETE endpoints are explicitly exempted in the
  authorizer.

## Configuration

All via environment variables: `LISTEN_ADDRESS`, `REDIS_ADDRESS`, `REDIS_USER`,
`REDIS_PASSWORD`, `SECRET` (JWT key), `STATIC_PATH`, `AUDIT_FILE`,
`SESSIONS_FILE` (default `/var/lib/nethserver/api-server/sessions.db`),
`SENSITIVE_LIST`, `ISSUER` (TOTP issuer name, default `NethServer`),
`GLOBAL_RATE_LIMIT_AVERAGE` (per-IP requests/sec across all routes, default
`25`; `0` disables), and `GLOBAL_RATE_LIMIT_BURST` (burst allowance, default
`300`).

The sustained average is what bounds a flood; the burst only sets how many
requests may arrive at once, and is sized for a legitimate UI page load. A burst
below `1` is clamped to `1` by `RateLimiter`; use `GLOBAL_RATE_LIMIT_AVERAGE=0`
to disable the limiter.

## Conventions

- Response envelopes use `structs.Map()` for struct-to-JSON conversion.
- Redis keys follow `task/<entity>/<id>/<task-id>/{context,output,error,exit_code}` (e.g. `task/node/<node-id>/<task-id>/output`).
- Audit events: `login-ok`, `auth-fail`, `create-task`, `join-token` stored in SQLite.
- WebSocket sessions are validated on each ping-pong (JWT expiry and login
  session check); clients send `authorize` again after a token refresh.
- Unit tests only for `sessions/` (`go test -tags sqlite_omit_load_extension
  ./sessions/`); other testing is via integration tests in `core/tests/`.
