*** Settings ***
Library    RequestsLibrary
Library    SSHLibrary

*** Test Cases ***
Set node FQDN certificate
    ${output}  ${rc} =    Execute Command    api-cli run set-certificate --agent module/traefik1 --data '{"fqdn":"${NODE_ADDR}"}'
    ...    return_rc=True
    Should Be Equal As Integers    ${rc}    0

Check that admin iterface is accessible with a valid certificate
    Wait Until Keyword Succeeds    20 times    1 seconds    GET    https://${NODE_ADDR}/cluster-admin/
