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
${MIN_WAIT}        240
${CAPACITY}        10
${TOLERANCE_MS}    100

*** Test Cases ***
Tasks executed within rate limit boundaries
    The First Ten BRPOPs run quickly
    But The Eleventh BRPOP run with delay due to rate limiter

*** Keywords ***
Give me a Timestamp in Ms
    ${ms} =  Get Timestamp Ms
    RETURN    ${ms}

The First Ten BRPOPs run quickly
    ${start} =     Give me a Timestamp in Ms
    
    Execute a battery of ten BRPOPs
    
    ${end} =       Give me a Timestamp in Ms
    ${elapsed} =    Evaluate    ${end} - ${start}
    Should Be True    ${elapsed} < ${FILL_MS}
    
    # run the following statement only if the previous one
    # fails
    # Should Be True    ${elapsed} < ${FILL_MS} + ${TOLERANCE_MS}

But The Eleventh BRPOP run with delay due to rate limiter
    Execute a battery of eleven BRPOPs

Execute a battery of ten BRPOPs
    FOR    ${i}    IN RANGE    10
        Execute a single BRPOP test
    END

Execute a single BRPOP test
    The Task is submitted    fail-validation    "test-rate-limiter"
    The Command is received    set    fail-validation    test-rate-limiter

Execute a battery of eleven BRPOPs
    FOR    ${i}    IN RANGE    11 
        ${start} =    Give me a Timestamp in Ms
        Execute a single BRPOP test
        ${end} =      Give me a Timestamp in Ms
        ${elapsed} =    Evaluate    ${end} - ${start}
        IF    ${i} < 11
            Should Be True    ${elapsed} < ${FILL_MS}
        END
    END
    Should Be True    ${elapsed} >= ${MIN_WAIT}
    Should Be True    ${elapsed} <= ${FILL_MS} + ${TOLERANCE_MS}