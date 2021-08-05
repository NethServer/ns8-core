*** Settings ***
Library           SSHLibrary

*** Variables ***
${COREBRANCH}     main
${SSH_KEYFILE}    %{HOME}/.ssh/id_ecdsa

*** Keywords ***
Install the core
    Open Connection    127.0.0.1
    Login With Public Key    root    ${SSH_KEYFILE}
    Execute Command    rm -rf install.sh uninstall.sh /etc/nethserver /var/lib/nethserver /usr/local/nethserver /home/*[0-9]
    Put File    install.sh    .    mode=0644
    ${rc} =    Execute Command    bash install.sh ${COREBRANCH}
    ...            return_stdout=False
    ...            return_stderr=False
    ...            return_rc=True
    Should Be Equal As Integers    ${rc}    0
    # Basic post install sanity checks:
    File Should Exist    /etc/nethserver/wg0.key
    File Should Exist    /etc/nethserver/wg0.pub
    File Should Exist    /etc/nethserver/api-server.env
    ${output} =     Execute Command    systemctl --failed
    Should Contain    ${output}    0 loaded units listed.

Uninstall core and modules
    Put File    imageroot/var/lib/nethserver/node/uninstall.sh    .    mode=0644
    Execute Command    bash uninstall.sh
    Close Connection

*** Settings ***
Suite Setup       Install the core
Suite Teardown    Uninstall core and modules
