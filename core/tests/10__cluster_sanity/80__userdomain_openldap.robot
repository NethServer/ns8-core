*** Settings ***
Library            SSHLibrary
Resource           api.resource

*** Variables ***
${domain}    nethserver.test
${domsuffix}    dc\=nethserver,dc\=test

*** Test Cases ***
Add domain
    ${response} =     Run task    cluster/add-internal-provider    {"image":"openldap","node":1}
    Set Suite Variable    ${mid1}    ${response['module_id']}
    Run task    module/${mid1}/configure-module    {"domain":"${domain}","admuser":"admin","admpass":"Nethesis,1234","provision":"new-domain"}

Reach directory directly
    [Tags]    unstable
    Fail

Reach directory through proxy
    [Tags]    unstable
    Fail

Add provider
    ${response} =     Run task    cluster/add-internal-provider    {"domain":"${domain}","image":"openldap","node":1}
    Set Suite Variable    ${mid2}    ${response['module_id']}
    Run task    module/${mid2}/configure-module    {"domain":"${domain}","admuser":"admin","admpass":"Nethesis,1234","provision":"join-domain"}

Remove provider
    Run task    cluster/remove-internal-provider    {"module_id":"${mid2}"}

Remove domain
    Run task    cluster/remove-internal-domain    {"domain":"${domain}"}

Directory is unreachable
    [Tags]    unstable
    Fail
