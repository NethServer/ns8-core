*** Settings ***
Documentation      Task submission to agent queues, task status, lists,
...                authorization of actions and audit records
Resource           api.resource
Suite Setup        Start the fake agents
Suite Teardown     Stop Fake Agents

*** Keywords ***
Start the fake agents
    Start Fake Agent    cluster
    Start Fake Agent    node/1
    Start Fake Agent    module/t1
    ${token} =    Get token of    admin
    Set Suite Variable    ${ADMIN_TOKEN}    ${token}
    ${token} =    Get token of    bob
    Set Suite Variable    ${BOB_TOKEN}    ${token}

*** Test Cases ***
Cluster task is queued for the cluster agent
    ${data} =    Create Dictionary    name=value
    ${body} =    Task request    list-actions    ${data}
    Start Redis Monitor
    ${response} =    API POST    /cluster/tasks    ${body}    ${ADMIN_TOKEN}    expected_status=201
    Redis Monitor Should Contain    (?i)^lpush cluster/tasks\\s
    Stop Redis Monitor
    Should Be Equal As Strings    ${response.json()}[message]    task queued successfully
    ${task} =    Fake Agent Receives Task    cluster
    Should Be Equal    ${task}[id]    ${response.json()}[data][id]
    Should Be Equal    ${task}[action]    list-actions
    Should Be Equal    ${task}[data]    ${data}
    Should Be Equal    ${task}[user]    admin
    Should Be Equal    ${task}[queue]    cluster/tasks
    Should Be Equal    ${task}[extra][title]    test list-actions
    Dictionary Should Contain Key    ${task}    timestamp
    Dictionary Should Contain Key    ${task}    parent

Node and module tasks are queued for their agents
    ${body} =    Task request    get-node-status
    API POST    /node/1/tasks    ${body}    ${ADMIN_TOKEN}    expected_status=201
    ${task} =    Fake Agent Receives Task    node/1
    Should Be Equal    ${task}[queue]    node/1/tasks
    ${body} =    Task request    get-configuration
    API POST    /module/t1/tasks    ${body}    ${ADMIN_TOKEN}    expected_status=201
    ${task} =    Fake Agent Receives Task    module/t1
    Should Be Equal    ${task}[queue]    module/t1/tasks

Task for an agent that is not running is rejected
    [Documentation]    module/t3 has roles but no running agent
    ${body} =    Task request    get-configuration
    ${response} =    API POST    /module/t3/tasks    ${body}    ${ADMIN_TOKEN}    expected_status=404
    Response message should be    ${response}    client idle check failed

Malformed task request is rejected
    ${response} =    POST On Session    api    /cluster/tasks    data=not-json
    ...    headers=${{ {"Authorization": "Bearer " + $ADMIN_TOKEN, "Content-Type": "application/json"} }}
    ...    expected_status=anything
    Should Be True    ${response.status_code} in [400, 403]

Task status and context are read from Redis
    ${body} =    Task request    list-actions
    ${response} =    API POST    /cluster/tasks    ${body}    ${ADMIN_TOKEN}    expected_status=201
    ${task} =    Fake Agent Receives Task    cluster
    ${task_id} =    Set Variable    ${task}[id]
    Fake Agent Completes Task    cluster    ${task_id}    exit_code=0    output={"result": [1, 2]}    error=some stderr
    ${response} =    API GET    /cluster/task/${task_id}/status    ${ADMIN_TOKEN}
    ${status} =    Set Variable    ${response.json()}[data]
    Should Be Equal As Integers    ${status}[exit_code]    0
    Should Be Equal    ${status}[output]    ${{ {"result": [1, 2]} }}
    Should Be Equal    ${status}[error]    some stderr
    Should Be Equal    ${status}[file]    task/cluster/${task_id}
    Evaluate    __import__('redis').Redis(username='cluster', password=$CLUSTER_PASS).set('task/cluster/' + $task_id + '/context', '{"action": "list-actions"}')
    ${response} =    API GET    /cluster/task/${task_id}/context    ${ADMIN_TOKEN}
    Should Be Equal    ${response.json()}[data][context]    ${{ {"action": "list-actions"} }}

Task output that is not JSON is returned as a string
    Fake Agent Completes Task    module/t1    plain1    exit_code=2    output=plain text
    ${response} =    API GET    /module/t1/task/plain1/status    ${ADMIN_TOKEN}
    Should Be Equal    ${response.json()}[data][output]    plain text
    Should Be Equal As Integers    ${response.json()}[data][exit_code]    2

Status of an unknown task is an error
    ${response} =    API GET    /cluster/task/unknown-task/status    ${ADMIN_TOKEN}    expected_status=400

Task lists
    Fake Agent Completes Task    node/1    listed1    output="ok"
    ${response} =    API GET    /node/1/tasks    ${ADMIN_TOKEN}
    Dictionary Should Contain Key    ${response.json()}[data][tasks]    task/node/1/listed1
    ${response} =    API GET    /nodes    ${ADMIN_TOKEN}
    List Should Contain Value    ${response.json()}[data][list]    1
    ${response} =    API GET    /modules    ${ADMIN_TOKEN}
    List Should Contain Value    ${response.json()}[data][list]    t1

Reader can read and run allowed actions
    API GET    /cluster/tasks    ${BOB_TOKEN}
    ${body} =    Task request    list-modules
    API POST    /cluster/tasks    ${body}    ${BOB_TOKEN}    expected_status=201
    Fake Agent Receives Task    cluster

Reader cannot run other actions
    ${body} =    Task request    add-module
    ${response} =    API POST    /cluster/tasks    ${body}    ${BOB_TOKEN}    expected_status=403
    Response message should be    ${response}    you don't have permission to access this resource
    ${body} =    Task request    get-configuration
    API POST    /module/t1/tasks    ${body}    ${BOB_TOKEN}    expected_status=201
    Fake Agent Receives Task    module/t1
    ${body} =    Task request    configure-module
    API POST    /module/t1/tasks    ${body}    ${BOB_TOKEN}    expected_status=403

User without roles on an agent cannot run its actions
    ${body} =    Task request    get-node-status
    API POST    /node/1/tasks    ${body}    ${BOB_TOKEN}    expected_status=403

Sensitive task data is hidden in the audit log
    ${data} =    Create Dictionary    user=alice    password=VerySecret!
    ${body} =    Task request    alter-user    ${data}
    API POST    /cluster/tasks    ${body}    ${ADMIN_TOKEN}    expected_status=201
    ${task} =    Fake Agent Receives Task    cluster
    Should Be Equal    ${task}[data][password]    VerySecret!
    ${response} =    GET On Session    api    /audit    params=action=create-task&user=admin
    ...    headers=${{ {"Authorization": "Bearer " + $ADMIN_TOKEN} }}
    ${text} =    Convert To String    ${response.content}
    Should Contain    ${text}    alter-user
    Should Not Contain    ${text}    VerySecret!

Login is recorded in the audit log
    ${response} =    GET On Session    api    /audit    params=action=login-ok&user=admin
    ...    headers=${{ {"Authorization": "Bearer " + $ADMIN_TOKEN} }}
    ${text} =    Convert To String    ${response.content}
    Should Contain    ${text}    login-ok
