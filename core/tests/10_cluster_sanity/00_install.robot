*** Settings ***
Library           SSHLibrary
Library           String
Resource          ../keywords.resource

*** Variables ***
${COREBRANCH}
${COREMODULES}

*** Test Cases ***
NS8 Installation
    [Tags]    install
    Install the core    ${COREBRANCH}    ${COREMODULES}
    Basic post install sanity checks
