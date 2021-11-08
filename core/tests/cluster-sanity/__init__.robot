*** Settings ***
Library           SSHLibrary

*** Variables ***
${SSH_KEYFILE}    %{HOME}/.ssh/id_ecdsa

*** Keywords ***
Connect to the node
    Open Connection    127.0.0.1
    Login With Public Key    root    ${SSH_KEYFILE}

Uninstall core and modules
    Put File    imageroot/var/lib/nethserver/node/uninstall.sh    .    mode=0644
    Execute Command    bash uninstall.sh
    Close Connection

*** Settings ***
Suite Setup       Connect to the Node
Suite Teardown    Uninstall core and modules
