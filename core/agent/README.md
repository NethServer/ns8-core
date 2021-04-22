# Agent

* Go version 1.16+
* Run `go build .` in this directory

## Startup

Start the agent passing two or more arguments:

1. The agent prefix identifier (e.g. "cluster", "node/1", "module/mail1"...)
2. One or more directories where actions are defined

For instance:

    ./agent module/mail1 ~/actions ~/.config/actions

The agent accepts also some environment variables, so a complete invocation from Bash could be:

    REDIS_ADDRESS=127.0.0.1:6379 REDIS_PASSWORD= ./agent module/mail1 ~/actions ~/.config/actions
