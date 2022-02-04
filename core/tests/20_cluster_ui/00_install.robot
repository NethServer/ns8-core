*** Settings ***
Library           SSHLibrary
Library           String
Resource          ../keywords.resource

*** Variables ***
${COREBRANCH}     main
${COREMODULES}

*** Test Cases ***
NS8 Installation
    [Tags]    install
    Connect to the Node    ${NODE_ADDR}    ${SSH_KEYFILE}
    Install the core    ${COREBRANCH}    ${COREMODULES}
    Basic post install sanity checks
