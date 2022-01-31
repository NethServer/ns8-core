*** Settings ***
Force Tags  promtail
Library    SSHLibrary

*** Test Cases ***
Service has not failed
    ${rc} =    Execute Command    systemctl is-failed -q promtail1
    ...    return_rc=True  return_stdout=False
    Should Be Equal As Integers    ${rc}    1
