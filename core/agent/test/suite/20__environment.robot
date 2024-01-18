*** Settings ***
Resource         taskrun.resource
Suite Setup      Start the agent
Suite Teardown   Stop the agent and cleanup
Test Setup       Start command monitoring
Test Teardown    Stop command monitoring

*** Test Cases ***
Read a variable
    Given The task is submitted    read-environment
    When The command is received    set    /exit_code    0
    And The task output should be equal to    VAL1admin
