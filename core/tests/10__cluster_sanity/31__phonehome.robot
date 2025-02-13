*** Settings ***
Library     SSHLibrary
Resource    api.resource


*** Test Cases ***
Check if phonehome.timer is enabled
    SSHLibrary.File Should Exist    /etc/systemd/system/timers.target.wants/phonehome.timer

Check phonehome data format
    ${output}    ${error}    ${rc} =    Execute Command    runagent print-phonehome    return_stderr=${True}    return_rc=${True}
    Should Be Equal As Integers    ${rc}    0
    Create Binary File    ${TEMPDIR}/phonehome.json    ${output}
    ${response} =    Evaluate    json.load(open("${TEMPDIR}/phonehome.json"))    modules=json
    Remove File    ${TEMPDIR}/phonehome.json
    Should Not Be Empty    ${response['$schema']}
    Should Not Be Empty    ${response['uuid']}
    Should Be Equal    ${response['installation']}    nethserver
    Should Not Be Empty    ${response['facts']}
