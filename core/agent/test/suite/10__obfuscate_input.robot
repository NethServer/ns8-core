*** Settings ***
Resource         taskrun.resource
Suite Setup      Start the agent
Suite Teardown   Stop the agent and cleanup
Test Setup       Start command monitoring
Test Teardown    Stop command monitoring

*** Test Cases ***
Passwords are obfuscated in Redis context keys 
    Given The task is submitted    obfuscate-test    {"account_info":{"password":"Nethesis,1234"},"password_not_replaced":"PRE-SERVED"}
    When The command is received   set  obfuscate-test/context
    Then The task context should contain    XXX
    And The task context should contain    PRE-SERVED
