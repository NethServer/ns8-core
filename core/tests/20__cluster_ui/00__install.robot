*** Settings ***
Library           SSHLibrary
Library           String
Resource          ../keywords.resource

*** Variables ***
${COREMODULES}

*** Test Cases ***
NS8 Installation
    [Tags]    install
    Connect to the Node    ${NODE_ADDR}    ${SSH_KEYFILE}
    Install the core    ${COREMODULES}
    Basic post install sanity checks
