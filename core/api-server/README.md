# api-server
This component is used to send command from UI to Redis, using HTTP Rest API and Redis Pub/Sub protocol.

## Building
To build the binary just execute
```bash
# get dependencies
go get

# build binary
go build

# build docs
swag init
```

## Configuration
The configuration is read from ENV
```bash
LISTEN_ADDRESS=0.0.0.0:8080
REDIS_ADDRESS=127.0.0.1:6379
REDIS_USER=default
REDIS_PASSWORD=""
SECRET=MY_SECRET,11
```

## Running
To execute the binary run
```bash
./api-server
```

You can also specify configurations using env vars
```bash
LISTEN_ADDRESS=0.0.0.0:8080 REDIS_ADDRESS=127.0.0.1:6379 REDIS_USER=my-user REDIS_PASSWORD=MyPass SECRET=MY_SECRET,11./api-server
```

## APIs
Each API is authenticated and authorized through a `JWT (JSON Web Token)`. In order to get access to the APIs you have to login before and get a JWT `token`. `SECRET` is used to set the JWT signature.

- `POST /login`
```bash
curl -s -H "Content-Type: application/json" -X POST http://localhost:8080/api/login --data '{"username": "<username>", "password": "<password>"}' | jq
```

```json
{
  "code": 200,
  "expire": "2021-04-21T14:05:35+02:00",
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhY3Rpb25zIjpbImNyZWF0ZS11c2VycyIsImRlbGV0ZS11c2VycyJdLCJleHAiOjE2MTkwMDY3MzUsImlkIjoiZWRvYXJkbyIsIm9yaWdfaWF0IjoxNjE4NDAxOTM1LCJyb2xlIjoiYWRtaW4ifQ.ru8CbqduPTBI4G9R3zINC_-L38Thggg_9ExFV3Grf18"
}
```

Once you have the `token`, each request must include this `token` in the header.

```bash
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhY3Rpb25zIjpbImNyZWF0ZS11c2VycyIsImRlbGV0ZS11c2VycyJdLCJleHAiOjE2MTg5MDc2NTIsImlkIjoiZWRvYXJkbyIsIm9yaWdfaWF0IjoxNjE4MzAyODUyLCJyb2xlIjoiYWRtaW4ifQ.dyqSFWi7L0aKAe7mujBJ9eN2nFnC9PcnFPhSQOZc2Nc"
```

The finally API request should be:
```bash
curl -s -H "Authorization: Bearer <your_token>" -X <VERB> http://localhost:8080/api/<api_name> --data '{<your_json_data}' | jq
```

### API list
The API doc is generated using swagger. The `swagger.json` file is built
automatically and published at:

https://raw.githubusercontent.com/NethServer/ns8-scratchpad/swagdoc/swagger.json

## Redis
  Each API is mapped to a specific command on Redis:
  - `GET /api/<cluster|node|module>/tasks` ⟶ `LRANGE <cluster|node|module>/tasks 0 -1`

  - `POST /api/<cluster|node|module>/tasks` ⟶ `LPUSH <cluster|node|module>/tasks <payload>`

    `PAYLOAD` contains the task information:
    - `id` ⟶ task id is an uuid string generated from server
    - `action` ⟶ action to execute
    - `data` ⟶ data used from action to execute the task

## Audit
Every request made to the server, using its APIs or WebSocket, is logged inside an audit db. The audit db is store in a file using a SQLite database schema. Each record is composed by:
- `ID`: the unique id of the record, autoincrement
- `User`: the name of the user that made the specific action
- `Action`: the name of the action made by the user
- `Data`: the payload of the action (if present)
- `Timestamp`: the time when the specific action is executed

## Websockets
The API server provides a WebSocket connection under the `/ws` path. Through websocket it is possible to send commands and receive data asynchronously via the socket channel.

### Commands
At the moment the supported commands are:
- `logs-start`: used to get the logs for a specific entity (cluster, node or module). Use loki cli (`logcli`) to search and retrieve logs
- `logs-stop`: used to stop the logs you are following (normally used in tail mode)

To execute commands via websocket, you need to send a payload to the websocket and listen to the results.

#### logs-start
```json
INPUT
{
   "action": "logs-start",
   "payload": {
      "id": "2fd5b630-3d27-4b88-ac11-ad2d09f4bbd7",
      "mode": "dump",
      "lines": "25",
      "filter": "",
      "from": "2021-01-19T10:00:00Z",
      "to": "2021-01-19T20:00:00Z",
      "entity" :"module",
      "entity_name": "traefik1"
   }
}
```

- `action`: fixed to `logs-start`
- `payload`: payload of the action
  - `id`: used to track multiple stream output (used in UI to get different logs stream in different panes)
  - `mode`: must be `tail` or `dump` - `string` (choose how to retrieve logs)
  - `lines`: could be empty or a number - `string` (could be empty in `tail` mode)
  - `filter`: could be empty or string - `string` (used to search specific words inside logs)
  - `from`: could be empty or iso8601 date string - `string` (used to search specific logs in a date range)
  - `to`: could be empty or iso8601 date string - `string` (used to search specific logs in a date range)
  - `entity`: must be `cluster` or `node` or `module` - `string`
  - `entity_name`: could be empty (`cluster` case) or name of the entity - `string` (ex. hostname of the node or module id like `traefik1`)


```json
OUTPUT
{
   "name": "logs-start",
   "payload": {
      "message": "2022-02-17T20:47:24+01:00 fedora33.n2.edoardo --\u003e [GIN-debug] Listening and serving HTTP on 127.0.0.1:8080",
      "pid": "14431"
   },
   "timestamp":"2022-03-01T09:01:41.698943069Z",
   "type":"action"
}
```

- `name`: is the name of the action
- `payload`: contains the response logs
  - `message`: contains the log message
  - `pid`: is the pid of the process that actually reads log
- `timestamp`: timestamp of the action
- `type`: used to identify the websocket outputs

#### logs-stop
```json
{
   "action":"logs-stop",
   "payload":{
      "pid":"14088"
   }
}
```

- `action`: fixed to `logs-start`
- `payload`: payload of the action
  - `pid`: contains the pid of the process that reads log

### Task events
The websocket also listens for every event of every task launched within the cluster and reports the payload through the socket channel.
