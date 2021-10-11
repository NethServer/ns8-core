# Agent

* Go version 1.16+
* Run `go build .` in this directory

## Startup

Start the agent passing two or more arguments:

1. The agent prefix identifier (e.g. "cluster", "node/1", "module/mail1"...)
2. One or more root action directories where actions are defined

For instance:

    ./agent module/mail1 . ~/.config/actions

The agent accepts also some environment variables, so a complete invocation from Bash could be:

    REDIS_ADDRESS=127.0.0.1:6379 REDIS_PASSWORD= ./agent module/mail1 . ~/.config/actions

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

- `<agent-id>/task/<task-id>/output`, collects data sent to FD 1 by all steps
- `<agent-id>/task/<task-id>/error`, collects data sent to FD 2 by all steps
- `<agent-id>/task/<task-id>/exit_code`, 0 if all steps were successful, otherwise the exit code of the failing step
- `<agent-id>/environment`, copy of additional environment variables passed
  to the action steps. The values are stored in the local module
  filesystem and copied to Redis if the action has completed successfully
  or the `dump-env` command was called.

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

Each action step receives the `task.data` string from FD 0 (INPUT). Any data sent to FD 1 (OUTPUT) constitutes the action output data, and any data sent to FD 2 (ERROR) is immediately copied to the agent error stream and appended to the action error data.

An action step receives an additional file descriptor where it can write special command strings. The file descriptor number can be obtained from the `AGENT_COMFD` environment variable.

## Environment and working directory

A running action step receives the **current working directory** value and its operating system process environment from the agent. This environment consists of

- variables inherited from the agent environment (e.g. PATH, PYTHONPATH...)
- variables loaded from the local filesystem (file `./environment`) and
  copied to the Redis DB, under the hash key `<AGENT_ID>/environment`
  (e.g. MODULE_ID, IMAGE_URL...)
- runtime agent vars:

    * `AGENT_COMFD` (integer) The file descriptor number where **action commands** (see below) can be writte to
    * `AGENT_TASK_ID` (string) The unique identifier of the Task that started the Action
    * `AGENT_TASK_ACTION` (string) The Action name

## Action commands

A command sent to the `AGENT_COMFD` file descriptor is an UTF-8 encoded string with a trailing new-line character.
The string represents an array of space delimited fields. It must be composed at least by one field that identifies the
*command name*. Further fields are *arguments*:

    CommandName [ SPACE Argument ]... NEWLINE

Fields are separated by the space character. If a field value contains
a space it can be wrapped by double quotes.

Available commands are:

- `set-env`
- `unset-env`
- `dump-env`
- `set-status`
- `set-progress`
- `set-weight`

For example, to send a `set-progress` command in a Bash script:

    echo 'set-progress 73' >&${AGENT_COMFD}

The same command in Python 3

    import os
    fd = int(os.environ['AGENT_COMFD'])
    os.write(fd, 'set-progress 73\n'.encode('utf-8'))

### set-env

The `set-env` command modifies an environment variable for subsequent
steps. If the action is successful, it persists the value in the local
filesystem `./environment` file for future action invocations. For example

    set-env FULLNAME "First User"

The environment vars are also copied to the Redis DB under the
`{agent_id}/environment` hash key.

### unset-env

The `unset-env` command remove an environment variable for subsequent
steps. If the action is successful, it persists the value in the local
filesystem `./environment` file for future action invocations.
For example

    unset-env FULLNAME

### dump-env

The `dump-env` command writes to a special file all variables set using
`set-env`. The file can be used for subsequent steps. The `dump-env`
command is automatically called when the action is `completed`. The
generated file is written to the agent current working directory as
`./environment`, and it is also copied to the Redis DB under the
`{agent_id}/environment` hash key.

For example

    dump-env

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
