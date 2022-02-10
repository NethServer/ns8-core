*** Settings ***
Resource         taskrun.resource
Suite Setup      Set Environment Variable    AGENT_POLLING_INTERVAL    75ms
Suite Teardown   Set Environment Variable    AGENT_POLLING_INTERVAL    5s
Test Setup       Start the command monitoring and the agent
Test Teardown    Stop command monitoring the agent and flush the database

*** Test Cases ***
Long-running action sends progress heartbeat
    Given The task is submitted    run-longlasting
    Then The heartbeat is received every    %{AGENT_POLLING_INTERVAL}

*** Keywords ***
Start the command monitoring and the agent
    Start command monitoring
    Start the agent

The heartbeat is received every
    [Arguments]    ${polling_interval_ms}
    FOR    ${counter}    IN RANGE    0    10
        Sleep    ${polling_interval_ms}
        The command is received    publish    status    running
    END

Stop command monitoring the agent and flush the database
    Close All Connections
    Terminate All Processes
    Flush the Redis database