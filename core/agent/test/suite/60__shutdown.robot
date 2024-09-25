*** Settings ***
Resource         taskrun.resource
Test Setup       Start the command monitoring and the agent
Test Teardown    Stop command monitoring and flush the database

*** Test Cases ***
Stop an idle agent gracefully
    Given The command is received    brpop    module/t1000/tasks
    When Send Signal To Process    SIGUSR1
    Then The agent completes successfully

Stop a busy agent gracefully
    [Teardown]    Run Keyword If Test Failed    Terminate all processes
    Given The task is submitted    run-longlasting
    And The command is received    publish    run-longlasting    status    running    progress  33
    And Send Signal To Process    SIGUSR1
    And Process should be running
    Then The agent completes successfully

Terminate a busy agent
    [Teardown]    Run Keyword If Test Failed    Terminate all processes
    Given The task is submitted    run-longlasting
    When Send Signal To Process    SIGTERM
    Then The agent quits within one second

*** Keywords ***
Start the command monitoring and the agent
    Start command monitoring    timeout=${2.0}
    Start the agent

The agent completes successfully
    [Arguments]    ${timeout}=${2000}
    ${res} =    Wait For Process    timeout=${timeout}ms
    Process Should Be Stopped
    Should Be Equal As Integers    ${res.rc}    ${0}

The agent quits within one second
    ${res} =    Wait For Process    timeout=1s
    Should Be True    ${res}    msg=The agent process is still running
    Should Be Equal As Integers    ${res.rc}    ${-15}
