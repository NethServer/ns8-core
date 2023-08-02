*** Settings ***
Library           SSHLibrary

*** Test Cases ***
Log of ldapproxy1 is not empty
    ${output}  ${rc} =    Execute Command  logcli query -q --no-labels '{module_id="ldapproxy1"} | json | line_format "{{.MESSAGE}}"'
    ...    return_rc=True
    Should Be Equal As Integers    ${rc}    0
    Should Not Be Empty    ${output}

Log of loki1 is not empty
    ${output}  ${rc} =    Execute Command  logcli query -q --no-labels '{module_id="loki1"} | json | line_format "{{.MESSAGE}}"'
    ...    return_rc=True
    Should Be Equal As Integers    ${rc}    0
    Should Not Be Empty    ${output}

Log of traefik1 is not empty
    ${output}  ${rc} =    Execute Command  logcli query -q --no-labels '{module_id="traefik1"} | json | line_format "{{.MESSAGE}}"'
    ...    return_rc=True
    Should Be Equal As Integers    ${rc}    0
    Should Not Be Empty    ${output}

Log of node 1 is not empty
    ${output}  ${rc} =    Execute Command  logcli query -q --no-labels '{node_id="1"} | json | line_format "{{.MESSAGE}}"'
    ...    return_rc=True
    Should Be Equal As Integers    ${rc}    0
    Should Not Be Empty    ${output}
