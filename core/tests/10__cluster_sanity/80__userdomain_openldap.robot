*** Settings ***
Library            SSHLibrary
Resource           api.resource
Resource           userdomain.resource
Suite Setup        Start the client tool container
Suite Teardown     Stop the client tool container

*** Variables ***
${domain}    openldap.nethserver.test
${domsuffix}    dc\=openldap,dc\=nethserver,dc\=test

*** Test Cases ***
Add domain
    ${response} =     Run task    cluster/add-internal-provider    {"image":"openldap","node":1}
    Set Suite Variable    ${mid1}    ${response['module_id']}
    Run task    module/${mid1}/configure-module    {"domain":"${domain}","admuser":"admin","admpass":"Nethesis,1234","provision":"new-domain"}

Reach directory directly
    ${val} =    Get server url    ${mid1}
    Set Suite Variable    ${surl}    ${val}
    Wait Until Keyword Succeeds    10s    1s    RootDSE is correct    ${surl}

Reach directory through proxy
    ${out} =    Execute Command    runagent python3 -c 'import agent.ldapproxy ; print(agent.ldapproxy.Ldapproxy().get_domain("${domain}"))'
    &{srv} =    Evaluate    ${out}
    Set Suite Variable    &{srv}    &{srv}
    Wait Until Keyword Succeeds    10s    1s    RootDSE is correct    ldap://${srv.host}:${srv.port}

Bind through proxy
    ${out}    ${err}    ${rc} =    Execute command    podman exec -i ldapclient ldapsearch -LLL -x -D '${srv.bind_dn}' -w '${srv.bind_password}' -H ldap://${srv.host}:${srv.port} -b '${srv.base_dn}' -s base
    ...    return_rc=${TRUE}    return_stderr=${TRUE}
    Should Be Equal As Integers    ${rc}    0
    Should Contain    ${out}    ${srv.base_dn}

Add provider
    ${response} =     Run task    cluster/add-internal-provider    {"domain":"${domain}","image":"openldap","node":1}
    Set Suite Variable    ${mid2}    ${response['module_id']}
    Run task    module/${mid2}/configure-module    {"domain":"${domain}","admuser":"admin","admpass":"Nethesis,1234","provision":"join-domain"}

Run user and group listing tests
    Run Keywords
    ...    Domain is ready
    ...    Create users and groups
    ...    List users
    ...    List groups
    ...    Check First User is unlocked
    ...    Check First User is member of g1
    #...    Check Second User is locked
    ...    Check Group One members
    ...    Get a non-existing user
    ...    Get a non-existing group

Remove provider
    Run task    cluster/remove-internal-provider    {"module_id":"${mid2}"}

Remove domain
    Run task    cluster/remove-internal-domain    {"domain":"${domain}"}

Directory is unreachable
    ${out}    ${err}    ${rc} =    Execute command    podman exec -i ldapclient ldapsearch -x -H ${surl} -b '' -s base
    ...    return_rc=${TRUE}    return_stderr=${TRUE}
    Should Be Equal As Integers    ${rc}    255


*** Keywords ***
RootDSE is correct
    [Arguments]    ${hurl}
    ${out}  ${err}  ${rc} =    Execute command    podman exec -i ldapclient ldapsearch -LLL -x -D '' -w '' -H ${hurl} -b '' -s base namingContexts
    ...    return_rc=${TRUE}    return_stderr=${TRUE}
    Should Be Equal As Integers    ${rc}    0
    Should Contain    ${out}    namingContexts: ${domsuffix}    ignore_case=${TRUE}

Get server url
    [Arguments]    ${mid}
    ${out} =    Execute Command    runagent redis-exec HGETALL module/${mid}/srv/tcp/ldap
    &{srv} =    Evaluate    ${out}
    [Return]    ldap://${srv.host}:${srv.port}


Start the client tool container
    Execute Command    podman run --rm --replace --detach --network host --name ldapclient alpine sh -c 'sleep 100000' ; podman exec -i ldapclient apk add openldap-clients

Stop the client tool container
    Execute Command    podman kill ldapclient
