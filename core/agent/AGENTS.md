# Agent

The agent daemon executes **actions** and **events** on cluster nodes. It reads
tasks from Redis queues and runs step scripts found in the filesystem.

## Build & test

```bash
go build .                   # build the agent binary
bash test-agent.sh           # run Robot Framework tests (requires podman)
```

`test-agent.sh` spins up a Redis container and runs the test suite in
`test/suite/`. Mock actions live in `test/actions/` and `test/events/`.

## Module layout

| File | Purpose |
|------|---------|
| `agent.go` | Main entry, CLI flags, Redis subscription loop |
| `htask.go` | Task handler — `BRPOP` from `<agent-id>/tasks`, runs action steps |
| `hevent.go` | Event handler — subscribes to `*/event/*` patterns |
| `hbuiltin.go` | Built-in actions: `list-actions`, `cancel-task` |
| `envstate.go` | Reads the `./environment` state file (KEY=VALUE) |
| `worklimit.go` | Atomic concurrency limiter for goroutines |
| `models/` | Data structs: `Task`, `Processor`, `Step`, `Event` |
| `validation/` | JSON Schema validation via `gojsonschema` |

## How actions execute

1. Agent pops a task from Redis (`BRPOP <agent-id>/tasks`).
2. Locates the action directory and discovers step executables (sorted
   alphabetically).
3. Each step runs as a subprocess receiving:
   - **stdin (fd 0)**: task data JSON
   - **stdout (fd 1)**: captured as output
   - **stderr (fd 2)**: logged
   - **fd 3 (`AGENT_COMFD`)**: IPC channel for `set-progress`, `set-status`,
     `set-weight` commands
4. Validates input/output against `validate-input.json` /
   `validate-output.json` if present.
5. Writes results to Redis: `task/<agent-id>/<task-id>/{output,error,exit_code}`.

## Task lifecycle

`pending → running → completed | aborted | validation-failed`

## Exit codes

- **0**: success
- **8**: action not defined
- **9**: OS exec error
- **10**: validation failed
- **11**: agent overloaded (max concurrency reached)
- **1–7, 32–255**: action-specific

## Environment variables

Steps inherit the agent `./environment` contents plus:

- `AGENT_COMFD` — IPC file descriptor
- `AGENT_TASK_ID`, `AGENT_TASK_ACTION`, `AGENT_TASK_USER` — task context
- `AGENT_POLLING_INTERVAL`, `MAX_CONCURRENCY`, `OVERLOAD_SLEEP` — tuning

## Conventions

- Action steps are standalone executables (Python 3 or Bash), numbered for
  ordering (`00validate`, `50update`, `60post`).
- Secrets in task context are obscured with `"XXX"` in Redis.
- Logging uses systemd severity prefixes (`SD_ERR`, `SD_DEBUG`, etc.).
- Concurrency is bounded by `MAX_CONCURRENCY` (default 32); excess tasks get
  exit code 11.
