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

-   `POST /login`

```bash
curl -s -H "Content-Type: application/json" -X POST http://localhost:8080/api/login --data '{"username": "<username>", "password": "<password>"}' | jq
```

```json
{
  "code": 200,
  "expire": "2021-04-21T14:05:35+02:00",
  "token":  "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhY3Rpb25zIjpbImNyZWF0ZS11c2VycyIsImRlbGV0ZS11c2VycyJdLCJleHAiOjE2MTkwMDY3MzUsImlkIjoiZWRvYXJkbyIsIm9yaWdfaWF0IjoxNjE4NDAxOTM1LCJyb2xlIjoiYWRtaW4ifQ.ru8CbqduPTBI4G9R3zINC_-L38Thggg_9ExFV3Grf18"
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

### Basic authentication

There are some cases where it is necessary to validate requests coming from other servers (Caddy, Nginx, ...) using the HTTP Basic authentication method. To achieve this functionality some routes verify authentication using basic auth.

The api-server could check that the action passed from route is among the user (obtained from Basic-Auth credentials) role authorized actions.

#### Basic authentication API list

-   `GET /module/:module_id/http-basic/:action` Use Basic HTTP auth for module

    Where:
    - `:module_id`: is the name of the module to check for action
    - `:action`: is the action to validate in the authorized actions

### 2FA (Two-Factor authentication)

If you feel comfortable and more secure to use 2FA, you can enable it inside the user setting page o via API. To enable these are the steps:

-   Login and get a valid JWT token
-   Enable 2FA for a specific user
-   Scan auto-generated QR code and add the secret inside your Google Authenticator (or similar protocol app)
-   Copy the generated OTP from the app and call the API to verifiy the OTP
-   After the OTP verification the JWT token can be used to call APIs

#### 2FA API list

-   `GET /api/2FA` Read current 2FA status
    ```json
    RESPONSE
    {
      "code": 200,
      "data": true,
      "message": "2FA set for this user"
    }

    {
      "code": 200,
      "data": false,
      "message": "2FA not set for this user"
    }
    ```

-   `DELETE /api/2FA` Revoke 2FA secret
    ```json
    BODY
    {
      "code": 200,
      "data": "",
      "message": "2FA revocate successfully"
    }
    ```

-   `GET /api/2FA/qr-code` Get URL with Authenticator setup (can be used to generate a QR Code)
    ```json
    RESPONSE
    {
      "code": 200,
      "data": {
        "key": "OQPZYGNFP4E37NJUCBR7ELEONC7WSXZW",
        "url": "otpauth://totp/NethServer:admin2fa?algorithm=SHA1\u0026digits=6\u0026issuer=NethServer\u0026period=30\u0026secret=OQPZYGNFP4E37NJUCBR7ELEONC7WSXZW"
      },
      "message": "QR code string"
    }
    ```

- `POST /api/2FA` Used to validate the OTP
   ```json
   BODY
   {
     "otp": "157720",
     "secret": "4W5XXQDHH3ROJBUOHEKUEL6M44IYUOMO"
   }

   RESPONSE
   {
    "code": 200,
    "data": null,
    "message": "2FA set successfully"
  }
   ```

### API list

The API doc is generated using swagger. The `swagger.json` file is built
automatically and published at:

<https://raw.githubusercontent.com/NethServer/ns8-core/swagdoc/swagger.json>

## Redis

  Each API is mapped to a specific command on Redis:

-   `GET /api/<cluster|node|module>/tasks` ⟶ `LRANGE <cluster|node|module>/tasks 0 -1`

-   `POST /api/<cluster|node|module>/tasks` ⟶ `LPUSH <cluster|node|module>/tasks <payload>`

    `PAYLOAD` contains the task information:

    -   `id` ⟶ task id is an uuid string generated from server
    -   `action` ⟶ action to execute
    -   `data` ⟶ data used from action to execute the task

## Audit

Every request made to the server, using its APIs or WebSocket, is logged inside an audit db. The audit db is store in a file using a SQLite database schema. Each record is composed by:

-   `ID`: the unique id of the record, autoincrement
-   `User`: the name of the user that made the specific action
-   `Action`: the name of the action made by the user
-   `Data`: the payload of the action (if present). To avoid storing
    passwords in clear-text, any attribute in `Data` matching any of the
    substrings "password", "secret", "token" is forced to value "XXX".
-   `Timestamp`: the time when the specific action is executed

## Websockets

The API server provides a WebSocket connection under the `/ws` path. Through websocket it is possible to send commands and receive data asynchronously via the socket channel.

### Commands

At the moment the supported commands are:

-   `logs-start`: used to get the logs for a specific entity (cluster, node or module). Use loki cli (`logcli`) to search and retrieve logs
-   `logs-stop`: used to stop the logs you are following (normally used in tail mode)

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
      "entity_name": "traefik1",
      "timezone": "Europe/Rome"
   }
}
```

-   `action`: fixed to `logs-start`
-   `payload`: payload of the action
    -   `id`: used to track multiple stream output (used in UI to get different logs stream in different panes)
    -   `mode`: must be `tail` or `dump` - `string` (choose how to retrieve logs)
    -   `lines`: could be empty or a number - `string` (could be empty in `tail` mode)
    -   `filter`: could be empty or string - `string` (used to search specific words inside logs)
    -   `from`: could be empty or iso8601 date string - `string` (used to search specific logs in a date range)
    -   `to`: could be empty or iso8601 date string - `string` (used to search specific logs in a date range)
    -   `entity`: must be `cluster` or `node` or `module` - `string`
    -   `entity_name`: could be empty (`cluster` case) or name of the entity - `string` (ex. hostname of the node or module id like `traefik1`)
    -   `timezone`: could be empty (default UTC) or a specific valid timezone (eg. Europe/Rome)

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

-   `name`: is the name of the action
-   `payload`: contains the response logs
    -   `message`: contains the log message
    -   `pid`: is the pid of the process that actually reads log
-   `timestamp`: timestamp of the action
-   `type`: used to identify the websocket outputs

#### logs-stop

```json
INPUT
{
   "action": "logs-stop",
   "payload": {
      "id": "2fd5b630-3d27-4b88-ac11-ad2d09f4bbd7",
      "pid": "14088"
   }
}
```

-   `action`: fixed to `logs-start`
-   `payload`: payload of the action
-   `pid`: contains the pid of the process that reads log

```json
OUTPUT
{
   "name": "logs-stop",
   "payload": {
      "id": "2fd5b630-3d27-4b88-ac11-ad2d09f4bbd7",
      "pid": "14431",
      "message": "logs follow stopped"
   },
   "timestamp":"2022-03-01T09:01:41.698943069Z",
   "type":"action"
}
```

-   `name`: is the name of the action
-   `payload`: contains the response logs
    -   `id`: id of the trackes log stream
    -   `pid`: is the pid of the process that actually reads log
-   `timestamp`: timestamp of the action
-   `type`: used to identify the websocket outputs

### Task events

The websocket also listens for every event of every task launched within the cluster and reports the payload through the socket channel.

# api-server-logs

This component is a wrapper of api-server websocket logs

## Usage

```bash
# ./api-server-logs

api-server-logs is a wrapper of api-server websocket logs

Usage:
  api-server-logs [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  help        Help about any command
  logs        get logs of a specific entity. default: cluster
  version     show api-server-logs version

Flags:
  -h, --help   help for api-server-logs

Use "api-server-logs [command] --help" for more information about a command.
```

### Commands

-   `completion`:

    ```bash
    Generate the autocompletion script for api-server-logs for the specified shell.
    See each sub-command's help for details on how to use the generated script.

    Usage:
      api-server-logs completion [command]

    Available Commands:
      bash        Generate the autocompletion script for bash
      fish        Generate the autocompletion script for fish
      powershell  Generate the autocompletion script for powershell
      zsh         Generate the autocompletion script for zsh

    Flags:
      -h, --help   help for completion
    ```

-   `help`: prints the Usage output
-   `logs`: the command to retrieve logs in dump or follow mode

    ```bash
    get logs of a specific entity. default: cluster

    Usage:
      api-server-logs logs [flags]

    Flags:
      -e, --entity string   get logs for a specific entity: cluster, node, module (default "cluster")
      -f, --from string     get logs from a specific date. ISO8601 format
      -h, --help            help for logs
      -l, --lines string    get logs for a specific lines in dump mode (default "25")
      -m, --mode string     get logs in a specific mode: tail or dump (default "tail")
      -n, --name string     get logs for a specific entity name. used in node or module
      -s, --search string   get logs for a specific search string
      -t, --to string       get logs to a specific date. ISO8601 format
      -z, --timezone string get logs in a specific timezone
    ```

-   `version`: prints the command version
    ```bash
    api-server-logs 0.0.1
    ```
