*** Settings ***
Documentation    Check different task exit codes
Resource         taskrun.resource
Suite Setup      Start the agent
Suite Teardown   Stop the agent and cleanup
Test Setup       Start command monitoring
Test Teardown    Stop command monitoring

*** Test Cases ***
Action completes succesfully
    Given The task is submitted    run-success
    When The command is received   set  run-success/context
    And The command is received    publish  pending
    And The command is received    publish  running  progress  0
    And The command is received    publish  running  progress  100
    And The command is received    set    /exit_code    0
    And The command is received    publish  status  completed  progress  100
    Then The task has exit code    0

Action step exits with error
    Given The task is submitted    run-failure
    When The command is received   set  run-failure/context
    And The command is received    publish  pending
    And The command is received    publish  running  progress  0
    And The command is received    set    /exit_code    1
    And The command is received    publish  status  aborted  progress  0
    Then The task has exit code    1
    And The task output should not contain    NEVER REACH THIS POINT

Action does not exist
    Given The task is submitted    non-existent
    When The command is received   set  non-existent/context
    And The command is received    publish  pending  progress  0
    And The command is received    set    /error        is not defined
    And The command is received    set    /exit_code    8
    And The command is received    publish  status  aborted  progress  0
    Then The task has exit code    8

Validation failed
    Given The task is submitted    fail-validation    "FAILVALUE"
    When The command is received   set  fail-validation  context
    And The command is received    publish  pending  progress  0
    And The command is received    publish  running  progress  0
    And The command is received    set    /exit_code    2
    And The command is received    publish  status  validation-failed  progress  0
    Then The task has exit code    2
    And The task output should contain    FAILVALUE

Run list actions builtin
    Given The task is submitted    list-actions
    When The command is received    set    /exit_code    0
    Then The task output should contain    fail-validation
    And The task output should contain     run-success
    And The task output should contain     run-failure

A missing step does not abort the whole action
    Given The task is submitted    simulate-update
    When The command is received   set  simulate-update/context
    And The command is received    publish  pending
    And The command is received    publish  running  progress  100
    Then The task has exit code    0
    And The task output should not contain    NEVER REACH THIS POINT1
    And The task output should not contain    NEVER REACH THIS POINT2

A missing step does not abort a failing action
    Given The task is submitted    simulate-update-fail
    When The command is received    set    /exit_code    1
    Then The task output should not contain    NEVER REACH THIS POINT1
    And The task output should not contain    NEVER REACH THIS POINT2
