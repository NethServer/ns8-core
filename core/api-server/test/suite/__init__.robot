*** Settings ***
Documentation      Black-box tests of the API server: it runs as a process
...                connected to a Redis server prepared like a cluster
...                leader node. Tests use HTTP and WebSocket clients, a fake
...                agent and the Redis MONITOR command.
Resource           api.resource
Suite Setup        Run Keywords    Prepare Redis    AND    Start the API server
Suite Teardown     Stop the API server
