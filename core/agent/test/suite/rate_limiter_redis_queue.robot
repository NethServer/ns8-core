*** Settings ***
Resource       taskrun.resource
Suite Setup    Start the agent
Suite Teardown Stop the agent and cleanup
Test Setup     Start command monitoring   
Test Teardown  Stop command monitoring

***  Test Cases ***
Rate Limiter Test
    FOR ${i} IN RANGE 200
        Given The Task Is Submitted fail-validation
        IF ${i} > 10
            # TODO -> test agent pop rate
    END