*** Settings ***
Documentation     This is a resource file, that can contain variables and keywords.
...               Keywords defined here can be used where this Keywords.resource in loaded.
Library           SSHLibrary
Library           DateTime
Library           String

*** Variables ***
${SSH_KEYFILE}    %{HOME}/.ssh/id_ecdsa
${HEADLESS}       true

*** Keywords ***
Connect to the node
    [Arguments]    ${NODE_ADDR}    ${SSH_KEYFILE}
    Open Connection    ${NODE_ADDR}
    Login With Public Key    root    ${SSH_KEYFILE}    keep_alive_interval=30 seconds

System boot is complete
    ${output} =    Execute Command    systemctl is-system-running    --wait
    Should Be True    '${output}' == 'running' or '${output}' == 'degraded'

Install the core
    [Arguments]    ${COREMODULES}
    Save the journal begin timestamp
    ${COREMODULES} =    Replace String    ${COREMODULES}    ,    ${SPACE}
    Put File    install.sh    .    mode=0644
    ${rc} =    Execute Command    bash install.sh ${COREMODULES}
    ...    return_stdout=False
    ...    return_stderr=False
    ...    return_rc=True
    Should Be Equal As Integers    ${rc}    0

Save the journal begin timestamp
    ${tsnow} =    Get Current Date    result_format=epoch
    Set Global Variable    ${JOURNAL_SINCE}    ${tsnow}

Basic post install sanity checks
    SSHLibrary.File Should Exist    /etc/nethserver/wg0.key
    SSHLibrary.File Should Exist    /etc/nethserver/wg0.pub
    SSHLibrary.File Should Exist    /etc/nethserver/api-server.env

Uninstall core and modules
    Put File    imageroot/var/lib/nethserver/node/uninstall.sh    .    mode=0644
    Execute Command    bash uninstall.sh
    Collect the suite journal

Collect the suite journal
    Execute Command    journalctl -S @${JOURNAL_SINCE} >journal-dump.log
    Get File    journal-dump.log    ${OUTPUT DIR}/journal-${SUITE NAME}.log
