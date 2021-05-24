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

## Actions implementation

An Action is defined by one or more directories with the same name under the root directories
(those passed as arguments to the agent executable). Each executable file under the action
directories is an **Action Step**.

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

Each Action Step receives the `task.data` string from FD 0 (stdin). Any data sent to FD 1 (stdout) constitutes the Action output data, and any data sent to FD 2 (stderr) is immediately copied to the agent stderr and appended to the Action error data.

An Action Step receives an open file descriptor where it can write special command strings. The file descriptor number can be obtained from the `AGENT_COMFD` environment variable. Each command is composed at least by one field:

1. command name
2. first argument
...
n. nth argument

Fields are separated by the space character. If a field value contains
a space it can be wrapped by double quotes.

Available commands:

- `set-env`
- `dump-env`
- `set-progress`
- `set-weight`

### set-env

The `set-env` command modifies an environment variable for subsequent steps.
If the action is successful it persists the value in the Redis DB for future action invocations.
For example

    set-env FULLNAME "First User"

If an action step completes with a non-zero exit code, the whole action is considered failed and no more steps are
invoked.

When all steps are completed, the following keys are set in Redis DB

- `<agent-id>/task/<task-id>/output`, collects FD 1
- `<agent-id>/task/<task-id>/error`, collects FD 2
- `<agent-id>/task/<task-id>/exit_code`, 0 if all steps were successful, otherwise the exit code of the failing step
- `<agent-id>/environment`, additional environent variables passed to the action steps. The values are persisted to Redis if the action is successful

### dump-env

The `dump-env` command writes to a special file all variables set using `set-env`. The file can be used for subsequent steps.
If the action is successful, `dump-env` is automatically called at the end.
The generated file is saved inside the state directory and named `environment`.

For example

    dump-env

### set-progress

A step *progress* is a number from 0 to 100. The overall action progress value is given by the sum of all completed steps plus
the current step progress value. The `set-progress` step command changes the progress value for the current step. 

When the current step execution terminates successfully its progress value is always set to 100. For example

     set-progress 73

The above command set the current step progress to 73 (out of 100).

### set-weight

Each action step is assigned 1 as default *weight*. The step weight multiplies its progress value to obtain the overall action progress value. For example

    set-weight 50update 8

The above command sets the `50update` step weight to 8.
