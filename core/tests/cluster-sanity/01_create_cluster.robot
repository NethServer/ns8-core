*** Settings ***
Library           OperatingSystem

*** Variables ***
${COREBRANCH}     main
${ADMINPASS}      Nethesis,1234
${EP_PORT}        55820
${WG_NETWORK}     10.5.4.0/24

*** Test Cases ***
Create the cluster
    ${ep_host} =   Run  hostname -f
    ${rc}    ${output} =    Run And Return Rc And Output    bash -l -c "create-cluster ${ep_host}:${EP_PORT} ${WG_NETWORK} ${ADMINPASS}"
    Should Be Equal As Integers    ${rc}    0
    Should Contain    ${output}    join-cluster
