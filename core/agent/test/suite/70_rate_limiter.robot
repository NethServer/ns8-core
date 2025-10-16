*** Settings ***
Resource            taskrun.resource
Suite Setup         Start the agent
Suite Teardown      Stop the agent and cleanup
Test Setup          Start command monitoring   
Test Teardown       Stop command monitoring
Library             timeutils.py            

# IMPORTANT NOTE
# This test is based on the following assumptions:
# 1. The filler sleeps around 300 milliseconds.
# 2. The maximum capacity is fixed at 10 tokens.
# Because the sleep time is not precise the variable
# MIN_WAIT plays a crucial role.

*** Variables ***
${FILL_MS}         300
${MIN_WAIT}        250
${CAPACITY}        10
${TOLERANCE_MS}    100

*** Test Cases ***
Tasks executed within rate limit boundaries
    The First Ten BRPOPs run quickly
    But The eleventh runs Slowly due to rate limiting

Events executed wihtin rate limit boundaries
    These Ten Events run very Slowly due to rate limiting

*** Keywords ***
The First Ten BRPOPs run quickly
    ${start} =     Get Timestamp Ms
    
    Execute a battery of ten BRPOPs
    
    ${end} =       Get Timestamp Ms
    ${elapsed} =    Evaluate    ${end} - ${start}
    Should Be True    ${elapsed} < ${FILL_MS}
    
    # run the following statement only if the previous one
    # fails
    # Should Be True    ${elapsed} < ${FILL_MS} + ${TOLERANCE_MS}

But The eleventh runs Slowly due to rate limiting
    ${start} =      Get Timestamp Ms
    Execute a single BRPOP test
    ${end} =        Get Timestamp Ms
    ${elapsed} =      Evaluate     ${end} - ${start}
    Should Be True    ${elapsed} >= ${MIN_WAIT}
    Should Be True    ${elapsed} <= ${FILL_MS} + ${TOLERANCE_MS}

Execute a battery of ten BRPOPs
    FOR    ${i}    IN RANGE    10
        Execute a single BRPOP test
    END

Execute a single BRPOP test
    Given The task is submitted    run-success
    When The command is received   set  run-success/context

Execute a single PUBSUB test
    Given The task is submitted    run-success
    When The command is received    publish  pending

These Ten Events run very Slowly due to rate limiting
    FOR     ${i}    IN RANGE    10
        ${start} =          Get Timestamp Ms
        Execute a single PUBSUB test
        ${end} =            Get Timestamp Ms

        ${elapsed} =      Evaluate    ${end} - ${start}
        Should Be True    ${elapsed} >= ${MIN_WAIT}
        Should Be True    ${elapsed} <= ${FILL_MS} + ${TOLERANCE_MS}
    END