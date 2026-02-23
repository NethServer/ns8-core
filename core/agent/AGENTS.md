# Agent

The agent daemon executes **actions** and **events** on cluster nodes. It reads
tasks from Redis queues and runs step scripts found in the filesystem.

## Build & test

```bash
go build .                   # build the agent binary
bash test-agent.sh           # run Robot Framework tests (requires buildah and podman)
```

`test-agent.sh` spins up a Redis container and runs the test suite in
`test/suite/`. Mock actions live in `test/actions/` and `test/events/`.

## Module layout

| File | Purpose |
|------|---------|
| `agent.go` | Main entry, CLI flags, signal handling (SIGUSR1), goroutine orchestration |
| `htask.go` | Task handler — `BRPOP` from `<agent-id>/tasks`, runs action steps, publishes progress, obscures secrets |
| `hevent.go` | Event handler — subscribes to `*/event/*` patterns |
| `hbuiltin.go` | Built-in actions: `list-actions`, `cancel-task`, `rejectAction` (overload handling) |
| `envstate.go` | Reads the `./environment` state file (KEY=VALUE) |
| `worklimit.go` | Atomic concurrency limiter for goroutines |
| `models/` | Data types and factory functions: `Task`, `Processor`, `Step`, `Event`, `CreateTaskProcessor()`, `CreateEventHandler()`, `ListProcessors()` |
| `validation/` | JSON Schema validation via `gojsonschema` |

## How actions execute

1. Agent pops a task from Redis (`BRPOP <agent-id>/tasks`).
2. Locates the action directory and discovers step executables (sorted
   alphabetically).
3. Input is validated against `validate-input.json` if present (inserted as
   first step). Output is validated against `validate-output.json` if present
   (inserted as last step).
4. Each step runs as a subprocess receiving:
   - **stdin (fd 0)**: task data JSON
   - **stdout (fd 1)**: captured as output
   - **stderr (fd 2)**: logged
   - **fd 3 (`AGENT_COMFD`)**: IPC channel for `set-progress`,
     `set-status` (only `validation-failed` is accepted), `set-weight` commands
5. Writes results to Redis: `task/<agent-id>/<task-id>/{context,output,error,exit_code}`.
   Task results expire after 8 hours.

## Task lifecycle

`pending → running → completed | aborted | validation-failed`

## Exit codes

- **0**: success
- **1**: validation infrastructure fault / output validation error
- **8**: action not defined
- **10**: validation failed
- **11**: agent overloaded (max concurrency reached)
- **2–7, 32–255**: action-specific

## Environment variables

**Action steps** inherit the agent `./environment` contents plus:

- `AGENT_COMFD=3` — IPC file descriptor for action commands
- `AGENT_TASK_ID`, `AGENT_TASK_ACTION`, `AGENT_TASK_USER` — task context

**Event handler steps** inherit the agent `./environment` contents plus:

- `AGENT_COMFD=2` — commands are redirected to stderr (ignored)
- `AGENT_EVENT_NAME`, `AGENT_EVENT_SOURCE` — event context

**Agent-level tuning** (consumed by the agent process, not passed to steps):

- `AGENT_POLLING_INTERVAL` — BRPOP timeout in ms (default 5000)
- `MAX_CONCURRENCY` — max concurrent tasks (default 32)
- `OVERLOAD_SLEEP` — grace period before rejecting (default varies)

## Conventions

- Action steps are standalone executables (Python 3 or Bash), numbered for
  ordering (`00validate`, `50update`, `60post`).
- Secrets in task context are obscured with `"XXX"` in Redis.
- Logging uses systemd severity prefixes (`SD_ERR`, `SD_DEBUG`, etc.).
- Concurrency is bounded by `MAX_CONCURRENCY` (default 32); excess tasks get
  exit code 11.
