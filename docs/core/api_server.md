---
layout: default
title: API server
nav_order: 4
parent: Core
---

# API server

The API server is a daemon implemented using [Go](https://golang.org) and
listens on TCP port `9311`.

This component is used to send command from UI to Redis, using HTTP Rest API and Redis Pub/Sub protocol.

* TOC
{:toc}

## API Paths

APIs are documented using a [swagger.json file](https://raw.githubusercontent.com/NethServer/ns8-core/swagdoc/swagger.json).
To browse all existing APIs, use the [Swagger UI](https://petstore.swagger.io/?url=https://raw.githubusercontent.com/NethServer/ns8-core/swagdoc/swagger.json).

API flow:
- authenticate the user using the `login` API
- extract the JWT token from the login API response
- use the token to execute all next API calls

Authentication:
```bash
curl -X 'POST' \
  'http://localhost:9311/api/login' \
  -H 'accept: application/json' \
  -H 'Content-Type: application/json' \
  -d '{"username": "admin", "password": "Nethesis,12345"}' | jq
```

Response, note the JWT token saved inside the `token` field:
```json
{
  "code": 200,
  "expire": "2022-02-01T13:54:50Z",
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhY3Rpb25zIjpbXSwiZXhwIjoxNjQzNzIzNjkwLCJpZCI6ImFkbWluIiwib3JpZ19pYXQiOjE2NDMxMTg4OTAsInJvbGUiOiIifQ.BzzmleO6ln40DQUZr4FUyFTFEle6PkK-ar-vqwXJ5uo"
}
```

Use the obtained JWT `token` to list configured nodes:
```bash
curl -X 'GET' \ 
  'http://localhost:9311/api/nodes'
   -H 'accept: application/json' \
   -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhY3Rpb25zIjpbXSwiZXhwIjoxNjQzNzIzNjkwLCJpZCI6ImFkbWluIiwib3JpZ19pYXQiOjE2NDMxMTg4OTAsInJvbGUiOiIifQ.BzzmleO6ln40DQUZr4FUyFTFEle6PkK-ar-vqwXJ5uo"
```

Response:
```json
{
  "code": 200,
  "data": {
    "list": [
      "1"
    ],
    "queue": "node/*"
  },
  "message": "success"
}
```

You can also use the shortcut command called `api-cli`.
Usage example:
```
# api-cli list-actions
node/1/get-node-status
...
```

## Audit log

Every request made to the server, using its APIs or WebSocket, is logged inside an audit db.
The audit db is store in a file using a SQLite database schema. Each record is composed by:
- `ID`: the unique id of the record, autoincrement
- `User`: the name of the user that made the specific action
- `Action`: the name of the action made by the user
- `Data`: the payload of the action (if present)
- `Timestamp`: the time when the specific action is executed

Audit database is saved inside `AUDIT_FILE` variable, default path is `/var/lib/nethserver/api-server/audit.db`.

The audit logs can be access using the `audit` API or directly from command line.
Example:
```
sqlite3 /var/lib/nethserver/api-server/audit.db "SELECT * FROM audit LIMIT 10;" ".exit"
```

## api-cli

You can also use the shortcut command called api-cli. Usage example:
```
# api-cli login
username: admin
password:
```

If authentication is successful the `login` subcommand stores a JWT token in `$XDG_RUNTIME_DIR/.api-cli.cache`. With such token the `run` subcommand sends API calls to `api-server`. To destroy the token run
```
api-cli logout
```

If the token is missing `api-cli` tries to talk directly with Redis. In any case to run an action use the `run` subcommand, e.g.:
```
api-cli run list-actions
```

Shell auto-complete is available. It is possible to type `TAB` key to automatically complete the action name:
```
api-cli run <TAB><TAB>
```

