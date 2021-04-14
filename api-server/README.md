# api-server
This component is used to send command from UI to Redis, using the Pub/Sub protocol.

## Building
To build the binary just execute
```bash
# get dependencies
go get

# build binary
go build
```

## Configuration
The configuration is read from ENV
```bash
LISTEN_ADDRESS=0.0.0.0:8080
REDIS_ADDRESS=127.0.0.1:6379
SECRET=MY_SECRET,11
```

## Running
To execute the binary run
```bash
./api-server
```

You can also specify configurations using env vars
```bash
LISTEN_ADDRESS=0.0.0.0:8080 REDIS_ADDRESS=127.0.0.1:6379 SECRET=MY_SECRET,11./api-server
```

## APIs
Each API is authenticated and authorized through a `JWT (JSON Web Token)`. In order to get access to the APIs you have to login before and get a JWT `token`.

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

### API list (not authenticated)
- `POST /login`
- `POST /logout`
- `GET /refresh_token` used to refresh the access token when is expired

### API list (authenticated)

- `GET /api/cluster/tasks`
- `GET /api/node/tasks`
- `GET /api/module/tasks`
- `GET /api/node/<node_id>/tasks`
- `GET /api/module/<module_id>/tasks`

  Return the task information for each entity (Cluster, Node, Module).

  ```bash
  curl -s http://localhost:8080/api/cluster/tasks/ | jq
  curl -s http://localhost:8080/api/node/tasks/ | jq
  curl -s http://localhost:8080/api/module/tasks/ | jq
  curl -s http://localhost:8080/api/node/<node_id>/tasks/ | jq
  curl -s http://localhost:8080/api/module/<module_id>/tasks/ | jq
  ```

  ```json
  {
    "queue": "cluster/tasks",
    "tasks": [
      {
        "id": "30f49b7f-2f9d-46c1-b743-111b09229797",
        "action": "create-task",
        "data": "MY_VAR=2;MY_CIAO=4"
      },
      ...
    ]
  }
  ```
  ```json
  [
    {
      "queue": "node/1/tasks",
      "tasks": [
        {
          "id": "361cc217-328a-44fb-9b56-d8051f2d539d",
          "action": "create-task",
          "data": "MY_VAR=2;MY_CIAO=4"
        },
       ...
      ]
    },
    {
      "queue": "node/2/tasks",
      "tasks": [
        {
          "id": "da73d079-e23a-4dff-9fd2-5666a794c8cd",
          "action": "create-task",
          "data": "MY_VAR=2;MY_CIAO=4"
        },
        ...
      ]
    }
  ]
  ```
  ```json
  [
    {
      "queue": "module/1/tasks",
      "tasks": [
        {
          "id": "361cc217-328a-44fb-9b56-d8051f2d539d",
          "action": "create-task",
          "data": "MY_VAR=2;MY_CIAO=4"
        },
       ...
      ]
    },
    {
      "queue": "module/2/tasks",
      "tasks": [
        {
          "id": "da73d079-e23a-4dff-9fd2-5666a794c8cd",
          "action": "create-task",
          "data": "MY_VAR=2;MY_CIAO=4"
        },
        ...
      ]
    }
  ]
  ```
  ```json
  {
    "queue": "node/<node_id>/tasks",
    "tasks": [
      {
        "id": "30f49b7f-2f9d-46c1-b743-111b09229797",
        "action": "create-task",
        "data": "MY_VAR=2;MY_CIAO=4"
      },
      ...
    ]
  }
  ```
  ```json
  {
    "queue": "module/<module_id>/tasks",
    "tasks": [
      {
        "id": "30f49b7f-2f9d-46c1-b743-111b09229797",
        "action": "create-task",
        "data": "MY_VAR=2;MY_CIAO=4"
      },
      ...
    ]
  }

  ```

- `POST /api/cluster/tasks`
- `POST /api/node/<node_id>/tasks`
- `POST /api/module/<module_id>/tasks`

  Create a new task and add to specific entity queue

  ```bash
  curl -s -X POST http://localhost:8080/api/cluster/tasks --data '{"id": "", "action": "create-task", "data": "MY_VAR=2;MY_CIAO=4"}' | jq
  curl -s -X POST http://localhost:8080/api/node/<node_id>/tasks --data '{"id": "", "action": "create-task", "data": "MY_VAR=2;MY_CIAO=4"}' | jq
  curl -s -X POST http://localhost:8080/api/module/<module_id>/tasks --data '{"id": "", "action": "create-task", "data": "MY_VAR=2;MY_CIAO=4"}' | jq
  ```

  ```json
  {
    "message": "task queued successfully",
    "queue": "cluster/tasks",
    "task": {
      "id": "ff77b29c-8a5f-4775-8858-8dcfc8f421bf",
      "action": "create-task",
      "data": "MY_VAR=2;MY_CIAO=4"
    }
  }
  {
    "message": "task queued successfully",
    "queue": "node/<node_id>/tasks",
    "task": {
      "id": "ff77b29c-8a5f-4775-8858-8dcfc8f421bf",
      "action": "create-task",
      "data": "MY_VAR=2;MY_CIAO=4"
    }
  }
  {
    "message": "task queued successfully",
    "queue": "module/<module_id>/tasks",
    "task": {
      "id": "ff77b29c-8a5f-4775-8858-8dcfc8f421bf",
      "action": "create-task",
      "data": "MY_VAR=2;MY_CIAO=4"
    }
  }
  ```

  ## Redis
  Each API is mapped to a specific command on Redis:
  - `GET /api/<cluster|node|module>/tasks` ⟶ `LRANGE <cluster|node|module>/tasks 0 -1`

  - `POST /api/<cluster|node|module>/tasks` ⟶ `LPUSH <cluster|node|module>/tasks <payload>`

    `PAYLOAD` contains the task information:
    - `id` ⟶ task id is an uuid string generated from server
    - `action` ⟶ action to execute
    - `data` ⟶ data used from action to execute the task
