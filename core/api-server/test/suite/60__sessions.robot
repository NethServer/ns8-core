*** Settings ***
Documentation      Administrator sessions: token refresh, logout, revocation
...                and join tokens. Session data stay out of Redis.
Resource           api.resource
Suite Setup        Start Fake Agent    cluster
Suite Teardown     Stop Fake Agents
Test Teardown      Close WebSocket

*** Keywords ***
Start a session
    [Arguments]    ${username}=admin
    ${response} =    Login    ${username}
    RETURN    ${response.json()}

Refresh
    [Arguments]    ${refresh_token}    ${expected_status}=200
    ${body} =    Create Dictionary    refresh_token=${refresh_token}
    ${response} =    API POST    /refresh    ${body}    expected_status=${expected_status}
    RETURN    ${response.json()}

Session is refused
    [Arguments]    ${token}
    ${response} =    API GET    /cluster/tasks    ${token}    expected_status=401
    Response message should be    ${response}    session is not valid

Publish credentials change
    [Arguments]    ${user}
    ${rdb} =    Evaluate    __import__('redis').Redis(username='cluster', password=$CLUSTER_PASS)
    Evaluate    $rdb.publish('cluster/event/user-credentials-changed', __import__('json').dumps({'user': $user, 'not_before': __import__('time').time()}))
    Evaluate    $rdb.close()

*** Test Cases ***
Refresh returns a new token pair of the same session
    ${login} =    Start a session
    ${refreshed} =    Refresh    ${login}[refresh_token]
    Should Not Be Equal    ${refreshed}[refresh_token]    ${login}[refresh_token]
    ${old} =    Decode Token Claims    ${login}[token]
    ${new} =    Decode Token Claims    ${refreshed}[token]
    Should Be Equal    ${new}[sid]    ${old}[sid]
    API GET    /cluster/tasks    ${refreshed}[token]
    API GET    /cluster/tasks    ${login}[token]

Replaced refresh token works only for a short time
    ${login} =    Start a session
    Refresh    ${login}[refresh_token]
    ${again} =    Refresh    ${login}[refresh_token]
    Should Not Be Empty    ${again}[token]

Unknown or missing refresh token is refused
    ${response} =    Refresh    not-a-refresh-token    expected_status=401
    Should Be Equal    ${response}[message]    session is not valid
    ${response} =    API POST    /refresh    ${{ {} }}    expected_status=401
    Should Be Equal    ${response.json()}[message]    missing refresh_token

An access token is not a refresh token
    ${login} =    Start a session
    Refresh    ${login}[token]    expected_status=401

Logout ends the session
    ${login} =    Start a session
    ${refreshed} =    Refresh    ${login}[refresh_token]
    ${body} =    Create Dictionary    refresh_token=${refreshed}[refresh_token]
    API POST    /logout    ${body}    ${refreshed}[token]
    Session is refused    ${refreshed}[token]
    Session is refused    ${login}[token]
    Refresh    ${refreshed}[refresh_token]    expected_status=401

Logout with an access token only ends the session
    ${login} =    Start a session
    API POST    /logout    token=${login}[token]
    Session is refused    ${login}[token]

Logout-others keeps only the current session
    ${first} =    Start a session
    ${second} =    Start a session
    ${response} =    API POST    /logout-others    token=${second}[token]
    Should Be True    ${response.json()}[data] >= 1
    Session is refused    ${first}[token]
    API GET    /cluster/tasks    ${second}[token]

Credentials change ends the sessions of the user
    ${admin} =    Start a session
    ${bob} =    Start a session    bob
    Publish credentials change    admin
    Wait Until Keyword Succeeds    3s    200ms    Session is refused    ${admin}[token]
    Refresh    ${admin}[refresh_token]    expected_status=401
    API GET    /cluster/tasks    ${bob}[token]
    ${later} =    Start a session
    API GET    /cluster/tasks    ${later}[token]

WebSocket refuses an ended session
    ${login} =    Start a session
    API POST    /logout    token=${login}[token]
    Open WebSocket
    Websocket Authorize    ${login}[token]
    ${message} =    Websocket Receive Json
    Should Be Equal    ${message}    ${{ {"type": "authorize-error", "payload": "session is not valid"} }}
    Websocket Should Be Closed

Join token can only add a node
    ${login} =    Start a session
    ${response} =    API POST    /cluster/join-token    token=${login}[token]
    ${join} =    Set Variable    ${response.json()}[token]
    ${claims} =    Decode Token Claims    ${join}
    Should Be Equal    ${claims}[scope]    join
    Dictionary Should Not Contain Key    ${claims}    sid
    ${lifetime} =    Evaluate    $claims['exp'] - $claims['orig_iat']
    Should Be Equal As Integers    ${lifetime}    3600
    ${body} =    Task request    list-actions
    ${response} =    API POST    /cluster/tasks    ${body}    ${join}    expected_status=201
    ${task} =    Fake Agent Receives Task    cluster
    Fake Agent Completes Task    cluster    ${task}[id]    output="ok"
    API GET    /cluster/task/${task}[id]/status    ${join}
    API GET    /cluster/task/${task}[id]/context    ${join}    expected_status=403
    ${body} =    Task request    list-modules
    API POST    /cluster/tasks    ${body}    ${join}    expected_status=403
    API GET    /nodes    ${join}    expected_status=403
    API POST    /cluster/join-token    token=${join}    expected_status=403
    Refresh    ${join}    expected_status=401

Join token needs the add-node permission
    ${login} =    Start a session    bob
    API POST    /cluster/join-token    token=${login}[token]    expected_status=403

Session data are not stored in Redis
    Start a session
    ${rdb} =    Evaluate    __import__('redis').Redis(username='cluster', password=$CLUSTER_PASS, decode_responses=True)
    ${keys} =    Evaluate    [k for k in $rdb.scan_iter() if 'session' in k or 'refresh' in k]
    Evaluate    $rdb.close()
    Should Be Empty    ${keys}
