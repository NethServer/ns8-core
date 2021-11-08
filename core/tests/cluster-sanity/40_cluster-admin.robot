*** Settings ***
Library           RequestsLibrary

*** Test Cases ***
HTTP endpoint is reachable
    ${request} =    GET    http://${NODE_ADDR}/cluster-admin/
    Should Contain    ${request.text}    <title>ns8-ui</title>
HTTPS endpoint is reachable
    ${request} =    GET    https://${NODE_ADDR}/cluster-admin/    verify=${FALSE}
    Should Contain    ${request.text}    <title>ns8-ui</title>
