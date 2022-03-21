*** Settings ***
Resource         taskrun.resource
Test Setup      Start the agent and wait psubscribe
Test Teardown   Stop the agent and cleanup

*** Test Cases ***
Payload is echoed across runs
    When The event is raised    echo-received    payload0000
    Then Wait until the agent log contains    payload0000
    When The event is raised    echo-received    payload0001
    Then Wait until the agent log contains    payload0001
    When The event is raised    echo-received    payload0002
    Then Wait until the agent log contains    payload0002

Output is sent to log
    When The event is raised    echo-received    payload0003
    Then Wait until the agent log contains    Environment variables:
    And The event exit code is    0

Error is sent to log
    When The event is raised    echo-received    payload0004
    Then Wait until the agent log contains    Message sent to stderr
    And The event exit code is    0

Commands are sent to log
    When The event is raised    echo-received    payload0007
    Then Wait until the agent log contains    set-env IGNORED DURING_EVENTS
    And The event exit code is    0

Environment is defined
    Given Create File    /srv/state/environment    MYVARIABLE\=MYVAL\n
    When The event is raised    echo-received    payload0005
    Then Wait until the agent log contains    AGENT_EVENT_NAME
    And Wait until the agent log contains    AGENT_EVENT_SOURCE
    And Wait until the agent log contains    MYVARIABLE

Abort stops execution
    Given The event is raised    handler-errored    payload0006
    When The event exit code is    2
    Then The agent log does not contain    NEVER REACH THIS POINT

*** Keywords ***
The event is raised
    [Arguments]    ${event_name}    ${event_payload}
    Set Test Variable    ${LAST_EVENT}    ${event_name}
    Evaluate    $RDB.publish("source/event/"+$event_name, $event_payload)

The agent log does not contain
    [Arguments]    ${pattern}
    Run Keyword And Expect Error    "${pattern}" not found in the agent log
    ...    Grep in the agent log    ${pattern}

The event exit code is
    [Arguments]    ${expected_exit_code}
    IF    ${expected_exit_code} == 0
        Wait until the agent log contains    Handler of source/event/${LAST_EVENT} exited with status "completed" (0) at step
    ELSE
        Wait until the agent log contains    Handler of source/event/${LAST_EVENT} exited with status "aborted" (${expected_exit_code}) at step
    END

Start the agent and wait psubscribe
    Start command monitoring
    Start the agent
    The command is received    psubscribe
    Stop command monitoring
