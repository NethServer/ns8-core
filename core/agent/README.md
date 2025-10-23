# Agent

* Go version 1.18+
* Run `go build .` in this directory

## Startup

Start the agent passing the following arguments

1. The agent prefix identifier (e.g. "cluster", "node/1", "module/mail1"...)
2. One or more root action/event directories

For instance:

    ./agent --agentid module/mail1 . --actionsdir ~/.config/actions --eventsdir ~/.config/events

The agent accepts also some environment variables, so a complete invocation from Bash could be:

    REDIS_ADDRESS=127.0.0.1:6379 REDIS_PASSWORD= ./agent module/mail1 . ~/.config/actions

## Agent behavior test suite

During development it is useful to test code changes quickly. The
`test-agent.sh` build the Golang `agent` binary and starts a
Robotframework test suite:

    bash test-agent.sh

The test suite aims to describe the commands sent to Redis by the agent
and the expected results in the agent log and environment file.

- The `test-agent.sh` script starts a detached Podman container that runs
  a Redis DB instance: `test-agent-redis`. Run the following command on a
  separate terminal to receive a detailed Redis command trace:

      podman exec -ti test-agent-redis nc 127.0.0.1 6379 <<<$'MONITOR\r\n'

- The test suite runs in a foreground container, `test-agent-suite`. On
  each run, the test suite starts and stops the `agent` process and deals
  with Redis DB setup/cleanup. Additional script arguments are forwarded
  to the Robotframework command. For example to invoke `robot` help:

      bash test-agent.sh --help

- Test results are written to the persistent `tstate` volume. To export
  its contents run:

      mkdir tstate
      podman volume export tstate | tar -C tstate -v -x -f -

When testing is finished, stop and destroy the Redis container with:

    podman stop test-agent-redis


## Task processing

After startup, the agent process blocks and waits for incoming tasks on a
Redis list (BRPOP). The list key is given by `<agent_prefix>/tasks`, e.g.
`module/mail1/tasks`.

Each list item represents a **task** to be processed by an **action**. It
must be a string representing a valid JSON-encoded object with the
following structure:

```json
{
    "id": "UUID-STRING",
    "action": "action-name",
    "data": {},
    "extra": {}
}
```

When a task is popped the agent

1. write the task payload to Redis key
   `task/<agent-id>/<task-id>/context`. Any attribute of the task `data`
   that ends with `password`, `secret`, or `token` has value replaced with
   the string `XXX`.

2. starts the action in background (more details about Action in the next
   section)

3. immediately returns to wait for more tasks on the Redis list while the
   action is running

## Actions execution model

An **action** is defined by one or more directories with the same name under the root directories
(those passed as arguments to the agent executable). Each executable file under the action
directories is an **action step**.

    dir1
     +--first-action
         +--step1
         +--step2
         +--step3

    dir2
     +--first-action
         +--step3
         +--step4

If two action steps have the same file name the last root action directory argument wins,
and previous executables with the same name are ignored.

Action steps are invoked in alphabetical order, no matter the directory where the file is placed.

For instance if agent is started with `agent id dir1 dir2` and the directories of the figure above
when `first-action` is invoked the excuted steps are:

- `dir1/first-action/step1`
- `dir1/first-action/step2`
- `dir2/first-action/step3`
- `dir2/first-action/step4`

## Action status and outcomes

The possible action status are:

- `pending`
- `running`
- `completed`
- `aborted`
- `validation-failed`

An action status is initially set to `pending` then switches to `running` during normal steps execution.

If the execution of a step cannot be started successfully by the OS `exec()` system call the whole action
is `aborted` and the exit code 9 is set.

If an action step completes with a non-zero exit code, the whole action is considered `aborted` and no more steps are
invoked.

When all steps are executed successfully the action is considered `completed`. Then the following keys are set in Redis DB

- `task/<agent-id>/<task-id>/output`, collects data sent to FD 1 by all steps
- `task/<agent-id>/<task-id>/error`, collects data sent to FD 2 by all steps
- `task/<agent-id>/<task-id>/exit_code`, 0 if all steps were successful, otherwise the exit code of the failing step
- `<agent-id>/environment`, copy of additional environment variables passed
  to the action steps. The values are stored in the local module
  filesystem and copied to Redis if the action has completed successfully.

The `validation-failed` status must be set manually with the `set-status` command, described below. When the `validation-failed`
status is set, no more steps are invoked and the action exit code is set to 10 unless another non-zero exit code is returned by
the last step.

Exit codes and their meaning:

- `0` - success
- `1` - **available to use** for generic errors
- `2-7` - **available to use** for action-specific error numbers
- `8` - the invoked action is not defined or has 0 steps (the directory is empty)
- `9` - OS `exec()` error, like "file is not executable" and similar
- `10` - default exit code once `validation-failed` status is set. It
  can be overridden by terminating the step with a non-zero exit code.
- `11-31` - reserved to the agent implementation
- `32-255` - **available to use** for action-specific error numbers

## File descriptors

Each action step receives the `task.data` contents in JSON representation
from FD 0 (INPUT). Any data sent to FD 1 (OUTPUT) constitutes the action
output data, and any data sent to FD 2 (ERROR) is immediately copied to
the agent error stream and appended to the action error data.

An action step receives an additional file descriptor where it can write special command strings. The file descriptor number can be obtained from the `AGENT_COMFD` environment variable.

## Validation

Validation of input and output can be enforced using [JSON
Schema](https://json-schema.org/).

To validate the input just drop a `validate-input.json` file inside the
action directory, to validate the output just create a
`validate-output.json` file.

As general rule, schema files are loaded from special paths:
- `${AGENT_INSTALL_DIR}/validator-definitions.json`: this is always picked
  up by the agent from all actions. It can contain the common data format
  definitions
- `${AGENT_INSTALL_DIR}/actions/<action-name>/validate-input.json`: schema
  applied to validate the input data
- `${AGENT_INSTALL_DIR}/actions/<action-name>/validate-output.json`:
  schema applied to validate the output data

Additional validation logic can be implemented in an early step script
like `01validation`. If the validation fails, the step must execute the
[set-status](#set-status) command and return a non-zero exit code.


## Environment and working directory

A running action step receives the **current working directory** value and its operating system process environment from the agent. This environment consists of

- variables inherited from the agent environment (e.g. PATH, PYTHONPATH...)
- variables loaded from the local filesystem (state file `./environment`) and
  copied to the Redis DB, under the hash key `<AGENT_ID>/environment`
  (e.g. MODULE_ID, IMAGE_URL...)
- runtime agent vars:

    * `AGENT_COMFD` (integer) The file descriptor number where **action commands** (see below) can be writte to
    * `AGENT_TASK_ID` (string) The unique identifier of the Task that started the Action
    * `AGENT_TASK_ACTION` (string) The Action name
    * `AGENT_TASK_USER` (string) The user invoking the action, if available

## Action commands

A command sent to the `AGENT_COMFD` file descriptor is an UTF-8 encoded string with a trailing new-line character.
The string represents an array of space delimited fields. It must be composed at least by one field that identifies the
*command name*. Further fields are *arguments*:

    CommandName [ SPACE Argument ]... NEWLINE

Fields are separated by the space character. If a field value contains
a space it can be wrapped by double quotes.

Available commands are:

- `set-status`
- `set-progress`
- `set-weight`

For example, to send a `set-progress` command in a Bash script:

    echo 'set-progress 73' >&${AGENT_COMFD}

The same command in Python 3

    import os
    fd = int(os.environ['AGENT_COMFD'])
    os.write(fd, 'set-progress 73\n'.encode('utf-8'))

### set-status

The `set-status` command manually changes the status of the running action. The command requires the new status name
as argument. Accepted status names are:

- `validation-failed` - see the **Action status** section for more information

For example

    set-status validation-failed

### set-progress

A step *progress* is a number from 0 to 100. It helps the UI to show the user a smooth progress bar.

When a step execution completes successfully its progress value is automatically set to 100.

For instance the following command is sent during the 3rd step of a running action. The action has
5 steps. The command sets the current (3rd) step progress to 73 (out of 100).

    set-progress 73

The overall action progress value is given by the sum of all completed steps plus
the current step progress value. From the previous example the sum is

    100 + 100 + 73 = 273

As the action has 5 steps (100 units each) the overall action progress is

    273 * 100 / 500 = 54%

The UI can show the user that the running task progress is at 54%.

### set-weight

Each action step is assigned 1 as default *weight*. The step weight multiplies its progress value to obtain the overall action progress value. For example

    set-weight 50update 8

The above command sets the `50update` step weight to 8. A step that downloads data from the network, or performs a long filesystem
operation can be assigned a bigger weight to balance the progress increments of the action.

## Builtin actions

The `agent` binary implements some builtin actions. They are not read from
the directories passed as command arguments and cannot be overridden.

### list-actions

The `list-actions` action scans the action directories on the filesystems
and returns an array of action names. Builtin actions are included, too.

Action input is ignored.

Action output (JSON format):

    ["action1", "action2", "list-actions", "cancel-task"]

### cancel-task

The `cancel-task` action cancels a running task, by sending a TERM signal
to the currently running step process and any other process forked from
it.  Depending on how the TERM signal is handled by the step process and
its exit code, the action can be subsequently aborted or continue
execution.

If the given task ID does not match any running task, the action aborts
with code `2`.

Action input (JSON format):

    {"task":"fbee04e5-f251-4cda-a835-d7598449bf16"}

Action output is undefined.

## Events

Actions are started when a task object is received: it is designed to be a
one-to-one communication channel.  On the opposite, Events are a form of
one-to-many communication.  When the agent starts it listens for events
with a Redis `PSUBSCRIBE */event/*` command.

If the TCP connection is closed, the underlying `go-redis` library handles
the server reconnection.

When an event matching the PSUBSCRIBE pattern is received, the agent
starts to process it.

The Redis channel matches the pattern `<Source>/event/<Name>`. The prefix
is the event **source**. The suffix is the event **name**.

The event processing model resembles the actions one: a matching event
name is searched in the directory list defined by ``--eventsdir`` command
arguments. Then a list of executable event handler steps is built with the
same rules of action steps. The steps are executed in alphabetical order
and if one of them has non-zero exit code, the event processing is
aborted.

Action execution similarities are:

- The event payload can be read from the step stdin
- The step stdout and stderr file descriptors are redirected to the agent stderr

Differences with the action execution are:

- The final event exit code is written to the agent log only
- Commands sent to AGENT_COMFD are ignored (and redirected to the agent stderr)
