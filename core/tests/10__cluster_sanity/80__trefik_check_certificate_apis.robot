*** Settings ***
Library    RequestsLibrary
Library    SSHLibrary
Resource   api.resource

*** Test Cases ***
Set Let's Encrypt staging ACME servr
    Run task    module/traefik1/set-acme-server    {"url":"https://acme-staging-v02.api.letsencrypt.org/directory"}

Check ACME server
    ${output} =    Run task    module/traefik1/get-acme-server    {}
    Should Be Equal As Strings    ${output["url"]}    https://acme-staging-v02.api.letsencrypt.org/directory

Set node FQDN certificate
    Run task    module/traefik1/set-certificate    {"fqdn":"${NODE_ADDR}"}

Check that admin interface is accessible with a valid certificate
    [Tags]    unstable
    Wait Until Keyword Succeeds    20 times    1 seconds    GET    https://${NODE_ADDR}/cluster-admin/
