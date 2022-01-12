*** Settings ***
Library           SSHLibrary

*** Variables ***
${SSH_KEYFILE}    %{HOME}/.ssh/id_ecdsa

*** Keywords ***
Connect to the node
    Open Connection   ${NODE_ADDR}
    Login With Public Key    root    ${SSH_KEYFILE}
    ${output} =    Execute Command    systemctl is-system-running  --wait
    Should Be Equal As Strings    ${output}    running

Uninstall core and modules
    Put File    imageroot/var/lib/nethserver/node/uninstall.sh    .    mode=0644
    Execute Command    bash uninstall.sh
    Close Connection

*** Settings ***
Suite Setup       Connect to the Node
Suite Teardown    Uninstall core and modules
