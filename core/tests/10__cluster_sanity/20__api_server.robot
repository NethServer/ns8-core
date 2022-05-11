*** Settings ***
Library    SSHLibrary

*** Test Cases ***
Service has not failed
    ${rc} =    Execute Command    systemctl is-failed -q api-server    return_rc=True  return_stdout=False
    Should Be Equal As Integers    ${rc}    1

Service is reachable
    ${output} =    Execute Command    curl http://127.0.0.1:9311/api/nodes
    Should Contain    ${output}    "code":401
