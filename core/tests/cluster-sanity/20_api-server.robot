*** Settings ***
Library           OperatingSystem

*** Test Cases ***
Service has not failed
    ${rc} =    Run And Return Rc    systemctl is-failed -q api-server
    Should Be Equal As Integers    ${rc}    1

Service is reachable
    ${output} =    Run    curl http://127.0.0.1:8080/api/nodes
    Should Contain    ${output}    "code":401
