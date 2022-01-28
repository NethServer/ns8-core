*** Settings ***
Library           RequestsLibrary

*** Test Cases ***

HTTP endpoint is not reachable
    ${request} =    GET    http://${NODE_ADDR}/cluster-admin/     expected_status=403
    Status Should Be    403    ${request}
HTTPS endpoint is reachable
    ${request} =    GET    https://${NODE_ADDR}/cluster-admin/    verify=${FALSE}
    Should Contain    ${request.text}    <title>ns8-ui</title>
