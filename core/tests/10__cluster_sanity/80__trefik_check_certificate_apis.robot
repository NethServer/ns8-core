*** Settings ***
Library    RequestsLibrary
Library    SSHLibrary
Resource   api.resource

*** Test Cases ***
Set node FQDN certificate
    Run task    module/traefik1/set-certificate    {"fqdn":"${NODE_ADDR}"}

Check that admin interface is accessible with a valid certificate
    [Tags]    unstable
    Wait Until Keyword Succeeds    20 times    1 seconds    GET    https://${NODE_ADDR}/cluster-admin/
