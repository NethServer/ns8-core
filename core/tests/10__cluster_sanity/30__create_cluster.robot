*** Settings ***
Library           SSHLibrary

*** Variables ***
${ADMINPASS}      Nethesis,1234
${EP_PORT}        55820
${WG_NETWORK}     10.5.4.0/24

*** Test Cases ***
Create the cluster
    [Tags]    install
    ${ep_host} =    Execute Command    hostname -f
    ${output}    ${rc} =    Execute Command    create-cluster ${ep_host}:${EP_PORT} ${WG_NETWORK} ${ADMINPASS}
    ...    return_rc=True
    Should Be Equal As Integers    ${rc}    0
    Should Contain    ${output}    join-cluster
