*** Settings ***
Library            RequestsLibrary
Library            Process
Library            OperatingSystem
Library            Collections
Library            ApiServerLibrary.py

*** Variables ***
${API}                 http://127.0.0.1:9311/api
${SECRET}              test-jwt-secret-0123456789
${API_SERVER_PASS}     api-server-test-password
${CLUSTER_PASS}        cluster-test-password
${ADMIN_PASS}          Nethesis,1234
${OTP_SECRET}          JBSWY3DPEHPK3PXPJBSWY3DPEHPK3PXP
${TOKEN_DAYS}          ${14}

*** Keywords ***
Prepare Redis
    [Documentation]    Create the ACL users and keys of a cluster leader
    ...    node, as install-finalize.sh and the user actions do
    ${rdb} =    Evaluate    __import__('redis').Redis(host='127.0.0.1', decode_responses=True)
    Evaluate    $rdb.flushall()
    Evaluate    $rdb.execute_command('ACL', 'SETUSER', 'cluster', 'ON', '>' + $CLUSTER_PASS, '~*', '&*', '+@all')
    Evaluate    $rdb.execute_command('ACL', 'SETUSER', 'api-server', 'ON', '>' + $API_SERVER_PASS, '~*', '&*', 'nocommands', '+@read', '+@pubsub', '+lpush', '+@transaction', '+@connection', '+role', '+hset')
    Evaluate    $rdb.execute_command('ACL', 'SETUSER', 'module/t1', 'ON', '>t1-test-password', 'resetkeys', '~module/t1/*', '~task/module/t1/*', 'resetchannels', '&progress/module/t1/*', 'nocommands', '+@read', '+@write', '+@transaction', '+@connection', '+publish', '-@admin')
    FOR    ${user}    IN    admin    bob    otpuser    netuser
        Evaluate    $rdb.execute_command('ACL', 'SETUSER', $user, 'ON', '>' + $ADMIN_PASS, 'resetkeys', 'resetchannels', 'nocommands')
    END
    # roles and allowed actions
    Evaluate    $rdb.sadd('cluster/roles/owner', '*')
    Evaluate    $rdb.sadd('node/1/roles/owner', '*')
    Evaluate    $rdb.sadd('module/t1/roles/owner', '*')
    Evaluate    $rdb.sadd('module/t3/roles/owner', '*')
    Evaluate    $rdb.sadd('cluster/roles/reader', 'list-*', 'get-*')
    Evaluate    $rdb.sadd('module/t1/roles/reader', 'get-*')
    Evaluate    $rdb.hset('roles/admin', mapping={'cluster': 'owner', 'node/1': 'owner', 'module/t1': 'owner', 'module/t3': 'owner'})
    Evaluate    $rdb.hset('roles/bob', mapping={'cluster': 'reader', 'module/t1': 'reader'})
    Evaluate    $rdb.hset('roles/otpuser', mapping={'cluster': 'owner'})
    Evaluate    $rdb.hset('roles/netuser', mapping={'cluster': 'owner'})
    Evaluate    $rdb.hset('user/otpuser', '2fa', $OTP_SECRET)
    Evaluate    $rdb.hset('user/netuser', 'allowed_networks', '192.0.2.0/24')
    Evaluate    $rdb.set('cluster/network', '10.5.4.0/24')
    Evaluate    $rdb.hset('node/1/environment', 'NODE_ID', '1')
    Evaluate    $rdb.hset('module/t1/environment', 'MODULE_ID', 't1')
    # restrict the default user, last
    Evaluate    $rdb.execute_command('ACL', 'SETUSER', 'default', 'ON', 'nopass', 'resetkeys', '~cluster/*', '~node/*', '~module/*', 'resetchannels', '&*', 'nocommands', '+@read', '+@connection', '+subscribe', '+psubscribe', '-@admin')
    Evaluate    $rdb.close()
    Set Redis Credentials    cluster    ${CLUSTER_PASS}

Start the API server
    Create Directory    /srv/state
    Remove Files    /srv/state/audit.db*    /srv/state/sessions.db*
    Start Process    /usr/src/api-server
    ...    alias=api-server
    ...    stdout=/srv/log/api-server.log    stderr=STDOUT
    ...    env:GIN_MODE=release
    ...    env:LISTEN_ADDRESS=127.0.0.1:9311
    ...    env:REDIS_ADDRESS=127.0.0.1:6379
    ...    env:REDIS_USER=api-server
    ...    env:REDIS_PASSWORD=${API_SERVER_PASS}
    ...    env:SECRET=${SECRET}
    ...    env:AUDIT_FILE=/srv/state/audit.db
    ...    env:SESSIONS_FILE=/srv/state/sessions.db
    ...    env:STATIC_PATH=/srv/state
    ...    env:GLOBAL_RATE_LIMIT_AVERAGE=1000
    ...    env:GLOBAL_RATE_LIMIT_BURST=1000
    Wait Until Keyword Succeeds    10s    250ms    The API server answers
    Create Session    api    ${API}

The API server answers
    ${response} =    POST    ${API}/login    json=${{ {} }}    expected_status=anything
    Should Be Equal As Integers    ${response.status_code}    401

Stop the API server
    Close WebSocket
    Stop Redis Monitor
    Stop Fake Agents
    Terminate Process    api-server
    Delete All Sessions

Login
    [Arguments]    ${username}    ${password}=${ADMIN_PASS}    ${otp}=${None}
    ...    ${expected_status}=200
    ${body} =    Create Dictionary    username=${username}    password=${password}
    IF    $otp is not None
        Set To Dictionary    ${body}    otp=${otp}
    END
    ${response} =    POST On Session    api    /login    json=${body}
    ...    expected_status=${expected_status}
    RETURN    ${response}

Get token of
    [Arguments]    ${username}    ${password}=${ADMIN_PASS}
    ${response} =    Login    ${username}    ${password}
    RETURN    ${response.json()}[token]

Auth headers
    [Arguments]    ${token}
    ${headers} =    Create Dictionary    Authorization=Bearer ${token}
    RETURN    ${headers}

API GET
    [Arguments]    ${path}    ${token}=${None}    ${expected_status}=200
    ${headers} =    Create Dictionary
    IF    $token is not None
        ${headers} =    Auth headers    ${token}
    END
    ${response} =    GET On Session    api    ${path}    headers=${headers}
    ...    expected_status=${expected_status}
    RETURN    ${response}

API POST
    [Arguments]    ${path}    ${body}=${None}    ${token}=${None}    ${expected_status}=200
    ${headers} =    Create Dictionary
    IF    $token is not None
        ${headers} =    Auth headers    ${token}
    END
    ${response} =    POST On Session    api    ${path}    json=${body}    headers=${headers}
    ...    expected_status=${expected_status}
    RETURN    ${response}

Response message should be
    [Arguments]    ${response}    ${message}
    Should Be Equal As Strings    ${response.json()}[message]    ${message}

Task request
    [Arguments]    ${action}    ${data}=${None}
    ${extra} =    Create Dictionary    title=test ${action}    isNotificationHidden=${TRUE}
    ${body} =    Create Dictionary    action=${action}    data=${data}    extra=${extra}
    RETURN    ${body}
