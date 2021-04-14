# Agent

* Go version 1.16+
* Run `go build` in this directory

## Startup

Start the agent passing the directories where actions are defined as additional command line arguments, e.g.:

    ./agent ~/actions ~/.config/actions

The agent expects some environment variables, so a complete invocation from Bash becomes:

    AGENT_PREFIX=cluster REDIS_ADDRESS=127.0.0.1:6379 REDIS_PASSWORD= ./agent ~/actions ~/.config/actions
