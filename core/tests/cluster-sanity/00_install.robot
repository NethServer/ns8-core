*** Settings ***
Library    SSHLibrary

*** Variables ***
${COREBRANCH}     main

*** Test Cases ***
Install the core
    Put File    install.sh    .    mode=0644
    ${rc} =    Execute Command    bash install.sh ${COREBRANCH}
    ...            return_stdout=False
    ...            return_stderr=False
    ...            return_rc=True
    Should Be Equal As Integers    ${rc}    0

Basic post install sanity checks
    File Should Exist    /etc/nethserver/wg0.key
    File Should Exist    /etc/nethserver/wg0.pub
    File Should Exist    /etc/nethserver/api-server.env
    ${output} =     Execute Command    systemctl --failed
    Should Contain    ${output}    0 loaded units listed.
