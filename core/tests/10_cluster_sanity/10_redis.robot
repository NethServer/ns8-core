*** Settings ***
Library    SSHLibrary

*** Test Cases ***
Service has not failed
    ${rc} =    Execute Command    systemctl is-failed -q redis
    ...    return_rc=True  return_stdout=False
    Should Be Equal As Integers    ${rc}    1

Check redis-cli
    ${output} =    Execute Command    redis-cli PING
    Should Be Equal     ${output}    PONG

Check redis-exec
    ${output} =    Execute Command    runagent redis-exec PING
    Should Be Equal     ${output}    True
