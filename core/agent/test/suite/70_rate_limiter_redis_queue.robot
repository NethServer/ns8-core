*** Settings ***
Resource            taskrun.resource
Suite Setup         Start the agent
Suite Teardown      Stop the agent and cleanup
Test Setup          Start command monitoring   
Test Teardown       Stop command monitoring

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
Rate Limiter Test
    Execute First Ten BRPOPs
    Execute Eleventh BRPOP

*** Keywords ***
Now Ms
    ${ms} =    Evaluate    int(__import__('time').perf_counter_ns() // 1_000_000)
    RETURN    ${ms}

Execute First Ten BRPOPs
    ${start} =     Now Ms
    FOR    ${i}    IN RANGE    10
        Given The Task is submitted    fail-validation    "test-rate-limiter"
    END

    FOR    ${i}    IN RANGE    10
        Then The Command is received    set    fail-validation    test-rate-limiter
    END
    ${end} =    Now Ms
    ${elapsed} =    Evaluate    ${end} - ${start}
    Should Be True    ${elapsed} < ${FILL_MS}
    # 
    # Should Be True    ${elapsed} < ${FILL_MS} + ${TOLERANCE_MS}

Execute Eleventh BRPOP
    FOR    ${i}    IN RANGE    10
        Given The Task is submitted    fail-validation    "test-rate-limiter"
    END

    FOR    ${i}    IN RANGE    10
        Then The Command is received    set    fail-validation    test-rate-limiter
    END

    ${start} =    Now Ms
        The Task is submitted    fail-validation    "test-rate-limiter"
        The Command is received    set    fail-validation    test-rate-limiter
    ${end} =    Now Ms
    ${elapsed} =    Evaluate    ${end} - ${start}

    Should Be True    ${elapsed} >= ${MIN_WAIT}
    Should Be True    ${elapsed} <= ${FILL_MS} + ${TOLERANCE_MS}