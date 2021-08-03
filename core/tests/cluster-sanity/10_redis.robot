*** Settings ***
Library           OperatingSystem

*** Test Cases ***
Service has not failed
    ${rc} =    Run And Return Rc    systemctl is-failed -q redis
    Should Be Equal As Integers    ${rc}    1

Check redis-cli
    ${output} =    Run    bash -l -c "redis-cli PING"
    Should Be Equal     ${output}    PONG

Check redis-exec
    ${output} =    Run    bash -l -c "redis-exec PING"
    Should Be Equal     ${output}    True
