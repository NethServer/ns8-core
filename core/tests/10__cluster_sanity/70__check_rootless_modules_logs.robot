*** Settings ***
Library           SSHLibrary

*** Test Cases ***
Check logs for ldapproxy1
    ${output}  ${rc} =    Execute Command  logcli query -q --no-labels '{job="systemd-journal"} | json | SYSLOG_IDENTIFIER="ldapproxy1" | line_format "{{.nodename}} --> {{.MESSAGE}}"'
    ...    return_rc=True
    Should Be Equal As Integers    ${rc}    0
    Should Not Be Empty    ${output}

Check logs for loki1
    ${output}  ${rc} =    Execute Command  logcli query -q --no-labels '{job="systemd-journal"} | json | SYSLOG_IDENTIFIER="loki1" | line_format "{{.nodename}} --> {{.MESSAGE}}"'
    ...    return_rc=True
    Should Be Equal As Integers    ${rc}    0
    Should Not Be Empty    ${output}

Check logs for traefik1
    ${output}  ${rc} =    Execute Command  logcli query -q --no-labels '{job="systemd-journal"} | json | SYSLOG_IDENTIFIER="traefik1" | line_format "{{.nodename}} --> {{.MESSAGE}}"'
    ...    return_rc=True
    Should Be Equal As Integers    ${rc}    0
    Should Not Be Empty    ${output}
