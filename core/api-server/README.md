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
Each API is authenticated and authorized through a `JWT (JSON Web Token)`. In order to get access to the APIs you have to login before and get a JWT `token`. `SECRET` is used to get the JWT signature.

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
