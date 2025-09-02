*** Settings ***
Resource            taskrun.resource
Suite Setup         Start the agent
Suite Teardown      Stop the agent and cleanup
Test Setup          Start command monitoring   
Test Teardown       Stop command monitoring

*** Test Cases ***
Rate Limiter Test
    # push phase
    FOR    ${i}     IN RANGE     11
        Given The task is submitted     fail-validation
    END

    # count the number of SET operations performed by Agent Module
    ${setSum}=     Set variable    0
    Sleep     1s
    FOR      ${result_cnt}      IN RANGE     11
        When The command is received   set
        ${setSum}=   Then The Counter must be incremented  ${setSum}
    END

    And Check The Counter    ${setSum}

*** Keywords ***
The Counter must be incremented
    [Arguments]                 ${i}
    ${i}=        Convert to Integer    ${i} 
    ${i}=       Set Variable    ${i + 1}
    RETURN  ${i}

Check The Counter
    [Arguments]                 ${i}
    ${const}=       Convert to Integer   11
    IF  ${i} != ${const}
        Fail        The Rate Limiter is not Working
    END