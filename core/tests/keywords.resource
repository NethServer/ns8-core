# //// use this file?

*** Settings ***
Documentation     This is a resource file, that can contain variables and keywords.
...               Keywords defined here can be used where this Keywords.resource in loaded.
Library           SSHLibrary

*** Keywords ***
Connect to the node
    Open Connection    ${NODE_ADDR}
    Login With Public Key    root    ${SSH_KEYFILE}
    ${output} =    Execute Command    systemctl is-system-running    --wait
    Should Be True    '${output}' == 'running' or '${output}' == 'degraded'

Uninstall core and modules
    Put File    imageroot/var/lib/nethserver/node/uninstall.sh    .    mode=0644
    Execute Command    bash uninstall.sh
    Close Connection

Install the core
    ${COREMODULES} =    Replace String    ${COREMODULES}    ,    ${SPACE}
    Put File    install.sh    .    mode=0644
    ${rc} =    Execute Command    bash install.sh ${COREBRANCH} ${COREMODULES}
    ...    return_stdout=False
    ...    return_stderr=False
    ...    return_rc=True
    Should Be Equal As Integers    ${rc}    0

Basic post install sanity checks
    File Should Exist    /etc/nethserver/wg0.key
    File Should Exist    /etc/nethserver/wg0.pub
    File Should Exist    /etc/nethserver/api-server.env
