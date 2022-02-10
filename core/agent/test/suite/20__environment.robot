*** Settings ***
Resource         taskrun.resource
Suite Setup      Start the agent
Suite Teardown   Stop the agent and cleanup
Test Setup       Start command monitoring
Test Teardown    Stop command monitoring

*** Test Cases ***
Test the dump-env command
    Given The task is submitted    dump-environment    "ABC"

    When The command is received    hset    VAR1    ABC
    Then The agent environment contains    VAR1\=ABC

    When The command is received    set    /exit_code    1
    Then The task output should be equal to    ABC
    And The agent environment does not contain    XXX
    And The agent environment contains    ABC

Set a variable
    Given The task is submitted    set-environment

    When The command is received    set    /exit_code    0
    Then The agent environment contains    VAR1\=VAL1
    And The agent environment contains     VAR2\=VAL2
    And The agent environment does not contain    ABC

Read a variable
    Given The task is submitted    read-environment
    When The command is received    set    /exit_code    0
    And The task output should be equal to    VAL1

Unset a variable
    Given The task is submitted    unset-environment
    When The command is received    set    /exit_code    0
    And The agent environment contains    VAR2\=VAL2
    And The agent environment does not contain    VAR1\=VAL1
