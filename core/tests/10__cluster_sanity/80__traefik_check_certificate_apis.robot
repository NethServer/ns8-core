*** Settings ***
Library    RequestsLibrary
Library    SSHLibrary
Resource   api.resource
# Needs a public name for the node. Excluded where the node has none, like the QEMU CI.
Force Tags    letsencrypt

*** Test Cases ***
Set node FQDN certificate
    Run task    module/traefik1/set-certificate    {"fqdn":"${NODE_ADDR}"}

Check that admin interface is accessible with a valid certificate
    [Tags]    unstable
    Wait Until Keyword Succeeds    20 times    1 seconds    GET    https://${NODE_ADDR}/cluster-admin/
