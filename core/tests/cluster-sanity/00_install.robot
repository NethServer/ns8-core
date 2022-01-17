*** Settings ***
Library    SSHLibrary
Library    String

*** Variables ***
${COREBRANCH}
${COREMODULES}

*** Test Cases ***
Install the core
    ${COREMODULES} =    Replace String    ${COREMODULES}    ,    ${SPACE}
    Put File    install.sh    .    mode=0644
    ${rc} =    Execute Command    bash install.sh ${COREBRANCH} ${COREMODULES}
    ...            return_stdout=False
    ...            return_stderr=False
    ...            return_rc=True
    Should Be Equal As Integers    ${rc}    0

Basic post install sanity checks
    File Should Exist    /etc/nethserver/wg0.key
    File Should Exist    /etc/nethserver/wg0.pub
    File Should Exist    /etc/nethserver/api-server.env
