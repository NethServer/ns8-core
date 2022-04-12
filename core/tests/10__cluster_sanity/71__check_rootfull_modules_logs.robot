*** Settings ***
Library           SSHLibrary

*** Test Cases ***
Check logs for promtail1
    ${output}  ${rc} =    Execute Command  logcli query -q --no-labels '{job="systemd-journal"} | json | CONTAINER_TAG="promtail1" | line_format "{{.nodename}} --> {{.MESSAGE}}"'
    ...    return_rc=True
    Should Be Equal As Integers    ${rc}    0
    Should Not Be Empty    ${output}
