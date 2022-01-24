# API server

The API server is a daemon implemented using [Go](https://golang.org).

This component is used to send command from UI to Redis, using HTTP Rest API and Redis Pub/Sub protocol.

## API Paths

To view the API paths and documentation follow these steps:
- Access the site https://petstore.swagger.io
- In the top bar paste https://raw.githubusercontent.com/NethServer/ns8-scratchpad/swagdoc/swagger.json
- Click on "Explore"
- Browse accordions to view API structure and data

## Audit

Every request made to the server, using its APIs or WebSocket, is logged inside an audit db. The audit db is store in a file using a SQLite database schema. Each record is composed by:
- `ID`: the unique id of the record, autoincrement
- `User`: the name of the user that made the specific action
- `Action`: the name of the action made by the user
- `Data`: the payload of the action (if present)
- `Timestamp`: the time when the specific action is executed

