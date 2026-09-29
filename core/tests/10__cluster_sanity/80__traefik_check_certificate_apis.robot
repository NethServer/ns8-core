*** Settings ***
Library    RequestsLibrary
Library    SSHLibrary
Resource   api.resource

*** Test Cases ***
Use the HTTP-01 challenge
    # The CI reaches the node through a tunnel that terminates TLS, so TLS-ALPN-01 cannot reach traefik.
    Run task    module/traefik1/set-acme-server    {"url":"https://acme-v02.api.letsencrypt.org/directory","challenge":"HTTP-01"}

Set node FQDN certificate
    Run task    module/traefik1/set-certificate    {"fqdn":"${NODE_ADDR}"}

Check that admin interface is accessible with a valid certificate
    [Tags]    unstable
    Wait Until Keyword Succeeds    20 times    1 seconds    GET    https://${NODE_ADDR}/cluster-admin/
