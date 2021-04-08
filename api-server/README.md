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
`/opt/nethserver/api-server/conf.json`
```json
{
    "listen_address": "0.0.0.0:8080",
    "redis_type": "tcp", // could be `tcp` or `unix`
    "redis_address": "127.0.0.1:6379", // could be an tcp address or unix socket
    "secret": "YourSecret,11"
}
```

## Running
To execute the binary run
```bash
./api-server
```

You can also specify a different config file location with
```bash
./api-server -c /your/config/path/file.json
```

## APIs
- `GET /api/tasks/:task_id`

  Return the task information.

  ```bash
  curl -s http://localhost:8080/api/tasks/8de678c2-5b00-4199-b6d0-52b8ab6580a0 | jq
  ```

  ```json
  {
    "id": "8de678c2-5b00-4199-b6d0-52b8ab6580a0",
    "module": "mail",
    "node_id": 1,
    "params": "MY_VAR1=2;MY_VAR2=4"
  }
  ```

- `POST /api/tasks`

  Create a new task and publish to `tasks-channel`

  ```bash
  curl -s -X POST http://localhost:8080/api/tasks --data '{"id": "", "module": "mail", "node_id": 1, "params": "MY_VAR1=2;MY_VAR2=4"}'
  ```

  ```json
  {
    "message": "task created successfully",
    "task_id": "18dfe562-2a12-4e46-9b42-1403a1460cf7"
  }
  ```
  
- `DELETE /api/tasks/:task_id`

  Delete a specific task

  ```bash
  curl -s -X DELETE http://localhost:8080/api/tasks/91a5b85f-9ce7-42cf-a18e-ede1c6d64671 | jq
  ```

  ```json
  {
    "message": "task deleted successfully"
  }
  ```

  ## Redis
  Each API is mapped to a specific command on Redis:
  - `GET /api/tasks/:task_id` ⟶ `HGET tasks/:task_id PAYLOAD`

    `PAYLOAD` contains the task information.

  - `POST /api/tasks` ⟶ `HSET tasks/:task_id PAYLOAD <payload>` and `PUBLISH tasks-channel tasks/:task_id`

    `PAYLOAD` contains the task information and `tasks-channel` is the Pub/Sub channel for communications.

  - `DELETE /api/tasks/:task_id` ⟶ `DEL tasks/:task_id`