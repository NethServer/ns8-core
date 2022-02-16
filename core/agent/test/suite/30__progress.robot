*** Settings ***
Resource         taskrun.resource
Suite Setup      Start the agent
Suite Teardown   Stop the agent and cleanup
Test Setup       Start command monitoring
Test Teardown    Stop command monitoring

*** Test Cases ***
Observe step weights
    Given The task is submitted    set-weight
    When The command is received   publish    running    :0
    And The command is received   publish    running    :0
    And The command is received   publish    running    :12
    And The command is received   publish    running    :25
    And The command is received   publish    running    :37
    And The command is received   publish    running    :49
    And The command is received   publish    running    :50
    And The command is received   publish    running    :70
    And The command is received   publish    running    :80
    And The command is received   publish    running    :100
    Then The command is received   set    /exit_code    "0"
    And The task output should be equal to    DONE\n
