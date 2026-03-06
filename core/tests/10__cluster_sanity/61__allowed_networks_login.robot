*** Settings ***
Library     SSHLibrary
Library     RequestsLibrary
Resource    api.resource
Suite Setup    Evaluate    __import__('urllib3').disable_warnings(__import__('urllib3').exceptions.InsecureRequestWarning)

*** Test Cases ***
Login is denied when source IP not in allowed_networks
    ${out}  ${rc} =     Execute Command    add-user --allowed-networks 127.0.0.1/32 --role owner --password Nethesis,1234 testblock    return_rc=True
    Should Be Equal As Integers    ${rc}    0
    ${jsondata} =    Evaluate    json.loads('{"username":"testblock","password":"Nethesis,1234"}')    modules=json
    ${request} =    POST    https://${NODE_ADDR}/cluster-admin/api/login
    ...    json=${jsondata}
    ...    verify=${FALSE}
    ...    expected_status=401
    Should Contain    ${request.text}    incorrect Username or Password    Unexpected response contents

Login succeeds after removing the restriction
    Run task    cluster/alter-user    {"user":"testblock","set":{"allowed_networks":[]}}
    ${jsondata} =    Evaluate    json.loads('{"username":"testblock","password":"Nethesis,1234"}')    modules=json
    ${request} =    POST    https://${NODE_ADDR}/cluster-admin/api/login
    ...    json=${jsondata}
    ...    verify=${FALSE}
    ...    expected_status=200
    Set Suite Variable    ${TESTBLOCK_TOKEN}    ${request.json()}[token]

JWT is rejected after restriction is re-applied mid-session
    Run task    cluster/alter-user    {"user":"testblock","set":{"allowed_networks":["127.0.0.1/32"]}}
    ${headers} =    Create Dictionary    Authorization=Bearer ${TESTBLOCK_TOKEN}
    ${request} =    GET    https://${NODE_ADDR}/cluster-admin/api/nodes
    ...    headers=${headers}
    ...    verify=${FALSE}
    ...    expected_status=403

cluster agent credentials are blocked from external IP
    ${password} =    Execute Command    runagent grep REDIS_PASSWORD agent.env | cut -d= -f2
    ${jsondata} =    Evaluate    json.loads('{"username":"cluster","password":"${password}"}')    modules=json
    ${request} =    POST    https://${NODE_ADDR}/cluster-admin/api/login
    ...    json=${jsondata}
    ...    verify=${FALSE}
    ...    expected_status=401
    Should Contain    ${request.text}    incorrect Username or Password    Unexpected response contents

cluster agent credentials are allowed from cluster network
    ${password} =    Execute Command    runagent grep REDIS_PASSWORD agent.env | cut -d= -f2
    ${jsondata} =    Evaluate    json.loads('{"username":"cluster","password":"${password}"}')    modules=json
    ${request} =    POST    https://10.5.4.1/cluster-admin/api/login
    ...    json=${jsondata}
    ...    verify=${FALSE}
    ...    expected_status=200

node/1 agent credentials are blocked from external IP
    ${password} =    Execute Command    runagent -m node grep REDIS_PASSWORD agent.env | cut -d= -f2
    ${jsondata} =    Evaluate    json.loads('{"username":"node/1","password":"${password}"}')    modules=json
    ${request} =    POST    https://${NODE_ADDR}/cluster-admin/api/login
    ...    json=${jsondata}
    ...    verify=${FALSE}
    ...    expected_status=401
    Should Contain    ${request.text}    incorrect Username or Password    Unexpected response contents

node/1 agent credentials are allowed from cluster network
    ${password} =    Execute Command    runagent -m node grep REDIS_PASSWORD agent.env | cut -d= -f2
    ${jsondata} =    Evaluate    json.loads('{"username":"node/1","password":"${password}"}')    modules=json
    ${request} =    POST    https://10.5.4.1/cluster-admin/api/login
    ...    json=${jsondata}
    ...    verify=${FALSE}
    ...    expected_status=200

module/traefik1 credentials are blocked from external IP (login)
    ${password} =    Execute Command    runagent -m traefik1 grep REDIS_PASSWORD agent.env | cut -d= -f2
    Set Suite Variable    ${TRAEFIK1_PASSWORD}    ${password}
    ${jsondata} =    Evaluate    json.loads('{"username":"module/traefik1","password":"${password}"}')    modules=json
    ${request} =    POST    https://${NODE_ADDR}/cluster-admin/api/login
    ...    json=${jsondata}
    ...    verify=${FALSE}
    ...    expected_status=401
    Should Contain    ${request.text}    incorrect Username or Password    Unexpected response contents

module/traefik1 credentials are allowed from cluster network
    ${password} =    Execute Command    runagent -m traefik1 grep REDIS_PASSWORD agent.env | cut -d= -f2
    ${jsondata} =    Evaluate    json.loads('{"username":"module/traefik1","password":"${TRAEFIK1_PASSWORD}"}')    modules=json
    ${request} =    POST    https://10.5.4.1/cluster-admin/api/login
    ...    json=${jsondata}
    ...    verify=${FALSE}
    ...    expected_status=200

HTTP-Basic with module/traefik1 credentials is denied from external IP
    ${basic_auth} =    Evaluate    ('module/traefik1', '${TRAEFIK1_PASSWORD}')
    ${request} =    GET    https://${NODE_ADDR}/cluster-admin/api/module/traefik1/http-basic/set-route
    ...    auth=${basic_auth}
    ...    verify=${FALSE}
    ...    expected_status=403
    Should Contain    ${request.text}    IP not allowed    Unexpected response contents

HTTP-Basic with module/traefik1 credentials is allowed from cluster network
    ${basic_auth} =    Evaluate    ('module/traefik1', '${TRAEFIK1_PASSWORD}')
    ${request} =    GET    https://10.5.4.1/cluster-admin/api/module/traefik1/http-basic/set-route
    ...    auth=${basic_auth}
    ...    verify=${FALSE}
    ...    expected_status=200

Cleanup testblock user
    Run task    cluster/remove-user    {"user":"testblock"}
