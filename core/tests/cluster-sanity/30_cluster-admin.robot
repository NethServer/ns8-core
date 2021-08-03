*** Settings ***
Library           OperatingSystem

*** Test Cases ***
HTTP endpoint is reachable
    ${output} =    Run    curl http://127.0.0.1/cluster-admin/
    Should Contain    ${output}    <title>ns8-ui</title>
HTTPS endpoint is reachable    
    ${output} =    Run    curl -k https://127.0.0.1/cluster-admin/
    Should Contain    ${output}    <title>ns8-ui</title>
