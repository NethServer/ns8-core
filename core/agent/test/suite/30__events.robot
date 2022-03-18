*** Settings ***
Resource         taskrun.resource
Suite Setup      Start the agent
Suite Teardown   Stop the agent and cleanup
Test Setup       Start command monitoring
Test Teardown    Stop command monitoring
Force Tags       unstable

*** Test Cases ***
Output is sent to log
    Fail    not implemented

Error is sent to log
    Fail    not implemented

Environment is defined
    Fail    not implemented

Abort stops execution
    Fail    not implemented
