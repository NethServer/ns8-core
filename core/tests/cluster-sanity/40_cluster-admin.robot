*** Settings ***
Library           RequestsLibrary

*** Test Cases ***
HTTP endpoint is reachable
    ${request} =    GET    http://127.0.0.1/cluster-admin/
    Should Contain    ${request.text}    <title>ns8-ui</title>
HTTPS endpoint is reachable
    ${request} =    GET    https://127.0.0.1/cluster-admin/    verify=${FALSE}
    Should Contain    ${request.text}    <title>ns8-ui</title>
