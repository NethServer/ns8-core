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
- `GET /api/cluster/tasks`
- `GET /api/node/<node_id>/tasks`
- `GET /api/module/<module_id>/tasks`

  Return the task information for each entity (Cluster, Node, Module).

  ```bash
  curl -s http://localhost:8080/api/cluster/tasks/ | jq
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
