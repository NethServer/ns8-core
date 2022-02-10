*** Settings ***
Resource         taskrun.resource
Suite Setup      Start the agent
Suite Teardown   Stop the agent and cleanup
Test Setup       Start command monitoring
Test Teardown    Stop command monitoring

*** Test Cases ***
Input JSON decode error
    Given The task is submitted    run-success    INVALID_JSON_STRING
    When The command is received   LPUSH  INVALID_JSON_STRING
    Then Wait until the agent log contains    Task ignored for decoding error: invalid character 'I'

JSON Schema check success
    Given The task is submitted    validate-schema    "OUTPUT-OK"
    When The command is received   publish    status    running    progress    100
    And The command is received    set    /exit_code    "0"
    And The command is received    publish    status    completed
    And The command is received    exec
    Then The task output should be equal to    "OUTPUT-OK"

Input JSON schema check failure
    Given The task is submitted    validate-schema    "INVALID-INPUT"
    When The command is received   publish    status    pending
    Then The command is received   set    /exit_code    "10"
    And The command is received    publish    status    validation-failed
    And The command is received    exec
    And The task output should contain    parameter
    And The task error should contain    Validation errors

Output JSON schema check failure
    Given The task is submitted    validate-schema    "INPUT-OK"
    When The command is received   publish    status    running    progress    100
    Then The command is received   set    /exit_code    "1"
    And The command is received    publish    status    aborted
    And The command is received    exec
    And The task output should be equal to    "INPUT-OK"
    And The task error should contain    Validation errors
