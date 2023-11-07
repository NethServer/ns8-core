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

Obfuscation does not alter the input structure
    Given The task is submitted    obfuscate-test    {"tags":["t1","t2"],"account_info":{"password":"Nethesis,1234","claims":["c1","c2"],"ratio":1.2},"password_not_replaced":"PRE-SERVED"}
    When The command is received   set  obfuscate-test/context
    Then The task context should contain    XXX
    And The task context should contain    PRE-SERVED
    And The task context should contain    ["t1","t2"]
    And The task context should contain    ["c1","c2"]
    And The task context should contain    1.2
