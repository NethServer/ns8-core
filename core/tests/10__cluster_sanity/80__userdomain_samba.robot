*** Settings ***
Library            SSHLibrary
Resource           api.resource
Resource           userdomain.resource
Suite Setup        Start the client tool container
Suite Teardown     Stop the client tool container

*** Variables ***
${domain}    samba.nethserver.test
${domsuffix}    dc\=samba,dc\=nethserver,dc\=test

*** Test Cases ***

Add domain
    Add the first samba module instance
    Prepare the suite configuration for the first samba module instance
    Run task    module/${mid1}/configure-module    {"hostname":"${hostname}","nbdomain":"${nbdomain}","ipaddress":"${ipaddress}","realm":"${domain}","adminuser":"administrator","adminpass":"Nethesis,1234","provision":"new-domain"}

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

Run user and group listing tests
    Run Keywords
    ...    Domain is ready
    ...    Create users and groups
    ...    List users
    ...    List groups
    ...    Check First User is unlocked
    ...    Check First User is member of g1
    ...    Check Second User is locked
    ...    Check Group One members
    ...    Get a non-existing user
    ...    Get a non-existing group

Remove domain
    Run task    cluster/remove-internal-domain    {"domain":"${domain}"}

Directory is unreachable
    ${out}    ${err}    ${rc} =    Execute command    podman exec -i ldapclient ldapsearch -x -H ${surl} -b '' -s base
    ...    return_rc=${TRUE}    return_stderr=${TRUE}
    Should Be Equal As Integers    ${rc}    255


*** Keywords ***
Add the first samba module instance
    ${response} =     Run task    cluster/add-internal-provider    {"image":"samba","node":1}
    Set Suite Variable    ${mid1}    ${response['module_id']}

Prepare the suite configuration for the first samba module instance
    ${response} =    Run task    module/${mid1}/get-defaults    {"provision":"new-domain"}
    Set Suite Variable    ${ipaddress}    ${response['ipaddress_list'][0]['ipaddress']}
    Should Not Be Empty    ${ipaddress}
    Set Suite Variable    ${nbdomain}    ${response['nbdomain']}
    Set Suite Variable    ${hostname}    ${response['hostname']}

Get server url
    [Arguments]    ${mid}
    ${out} =    Execute Command    runagent redis-exec HGETALL module/${mid}/srv/tcp/ldap
    &{srv} =    Evaluate    ${out}
    [Return]    ldaps://${srv.host}:${srv.port}

RootDSE is correct
    [Arguments]    ${hurl}
    ${out}  ${err}  ${rc} =    Execute Command    podman exec -i ldapclient ldapsearch -LLL -x -D '' -w '' -H ${hurl} -b '' -s base namingContexts
    ...    return_stderr=${TRUE}    return_rc=${TRUE}
    Should Contain    ${out}    namingContexts: ${domsuffix}    ignore_case=${TRUE}

Start the client tool container
    Execute Command    podman run --env\=LDAPTLS_REQCERT\=never --rm --replace --detach --network host --name ldapclient alpine sh -c 'sleep 100000' ; podman exec -i ldapclient apk add openldap-clients

Stop the client tool container
    Execute Command    podman kill ldapclient
