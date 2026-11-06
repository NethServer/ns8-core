*** Settings ***
Library     SSHLibrary
Library     RequestsLibrary
Library     String
Resource    api.resource
Suite Setup    Setup sessions suite
Suite Teardown    Run Keyword And Ignore Error    Run task    cluster/remove-user    {"user":"testsession"}

*** Variables ***
${API}    https://${NODE_ADDR}/cluster-admin/api

*** Test Cases ***
Login returns a refresh token
    ${login} =    Login testsession
    Should Not Be Empty    ${login}[token]
    Should Not Be Empty    ${login}[refresh_token]
    Set Suite Variable    ${LOGIN}    ${login}

Token is accepted
    Get nodes    ${LOGIN}[token]    200

Refresh returns a new token pair
    ${refreshed} =    Refresh    ${LOGIN}[refresh_token]    200
    Should Not Be Equal    ${refreshed}[refresh_token]    ${LOGIN}[refresh_token]
    Get nodes    ${refreshed}[token]    200
    Set Suite Variable    ${REFRESHED}    ${refreshed}

Old access token stays valid until it expires
    Get nodes    ${LOGIN}[token]    200

Logout ends the session
    ${headers} =    Create Dictionary    Authorization=Bearer ${REFRESHED}[token]
    ${body} =    Create Dictionary    refresh_token=${REFRESHED}[refresh_token]
    POST    ${API}/logout    json=${body}    headers=${headers}    verify=${FALSE}    expected_status=200
    Get nodes    ${REFRESHED}[token]    401
    Get nodes    ${LOGIN}[token]    401
    Refresh    ${REFRESHED}[refresh_token]    401

Logout-others keeps only the current session
    ${first} =    Login testsession
    ${second} =    Login testsession
    ${headers} =    Create Dictionary    Authorization=Bearer ${second}[token]
    POST    ${API}/logout-others    headers=${headers}    verify=${FALSE}    expected_status=200
    Get nodes    ${first}[token]    401
    Get nodes    ${second}[token]    200

Password change by an administrator ends the sessions
    ${login} =    Login testsession
    Get nodes    ${login}[token]    200
    Run task    cluster/alter-user    {"user":"testsession","set":{"password":"Nethesis,4321"}}
    Wait Until Keyword Succeeds    10s    1s    Get nodes    ${login}[token]    401
    Refresh    ${login}[refresh_token]    401
    Run task    cluster/alter-user    {"user":"testsession","set":{"password":"Nethesis,1234"}}

Join token can only add a node
    ${login} =    Login testsession
    ${headers} =    Create Dictionary    Authorization=Bearer ${login}[token]
    ${response} =    POST    ${API}/cluster/join-token    headers=${headers}    verify=${FALSE}    expected_status=200
    ${join_token} =    Set Variable    ${response.json()}[token]
    ${join_headers} =    Create Dictionary    Authorization=Bearer ${join_token}
    ${task} =    Evaluate    {"action": "list-actions", "data": None, "extra": {"title": "join token probe", "isNotificationHidden": True}}
    POST    ${API}/cluster/tasks    json=${task}    headers=${join_headers}    verify=${FALSE}    expected_status=201
    ${task} =    Evaluate    {"action": "list-modules", "data": None, "extra": {"title": "join token probe", "isNotificationHidden": True}}
    POST    ${API}/cluster/tasks    json=${task}    headers=${join_headers}    verify=${FALSE}    expected_status=403
    Get nodes    ${join_token}    403
    POST    ${API}/cluster/join-token    headers=${join_headers}    verify=${FALSE}    expected_status=403

Removing the user ends the sessions
    ${login} =    Login testsession
    Run task    cluster/remove-user    {"user":"testsession"}
    Wait Until Keyword Succeeds    10s    1s    Get nodes    ${login}[token]    401

*** Keywords ***
Setup sessions suite
    Evaluate    __import__('urllib3').disable_warnings(__import__('urllib3').exceptions.InsecureRequestWarning)
    ${out}  ${rc} =    Execute Command    add-user --role owner --password Nethesis,1234 testsession    return_rc=True
    Should Be Equal As Integers    ${rc}    0

Login testsession
    ${jsondata} =    Create Dictionary    username=testsession    password=Nethesis,1234
    ${response} =    POST    ${API}/login    json=${jsondata}    verify=${FALSE}    expected_status=200
    RETURN    ${response.json()}

Refresh
    [Arguments]    ${refresh_token}    ${expected_status}
    ${body} =    Create Dictionary    refresh_token=${refresh_token}
    ${response} =    POST    ${API}/refresh    json=${body}    verify=${FALSE}    expected_status=${expected_status}
    RETURN    ${response.json()}

Get nodes
    [Arguments]    ${token}    ${expected_status}
    ${headers} =    Create Dictionary    Authorization=Bearer ${token}
    GET    ${API}/nodes    headers=${headers}    verify=${FALSE}    expected_status=${expected_status}
