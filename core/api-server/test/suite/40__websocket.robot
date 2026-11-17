*** Settings ***
Documentation      WebSocket authorization and relay of task progress
Resource           api.resource
Test Teardown      Close WebSocket

*** Test Cases ***
Authorized client receives task progress
    ${token} =    Get token of    admin
    Open WebSocket
    Websocket Authorize    ${token}
    ${message} =    Websocket Receive Json    timeout=1
    Should Be Equal    ${message}    ${None}    msg=no reply expected on successful authorization
    Fake Agent Completes Task    cluster    ws-task-1    output="ok"
    ${message} =    Websocket Receive Json
    Should Be Equal    ${message}[type]    task
    Should Be Equal    ${message}[name]    progress/cluster/task/ws-task-1
    Should Be Equal    ${message}[payload][status]    completed
    Should Be Equal As Integers    ${message}[payload][progress]    100
    Dictionary Should Contain Key    ${message}    timestamp

Client without authorization receives nothing
    Open WebSocket
    Fake Agent Completes Task    cluster    ws-task-2    output="ok"
    ${message} =    Websocket Receive Json    timeout=1
    Should Be Equal    ${message}    ${None}

Expired token is refused
    ${claims} =    Admin Claims    admin    expires_in=-60
    ${token} =    Forge Token    ${SECRET}    ${claims}
    Open WebSocket
    Websocket Authorize    ${token}
    ${message} =    Websocket Receive Json
    Should Be Equal    ${message}    ${{ {"type": "authorize-error", "payload": "Token is expired"} }}
    Websocket Should Be Closed

Token with a wrong signature is refused
    ${claims} =    Admin Claims    admin
    ${token} =    Forge Token    another-secret    ${claims}
    Open WebSocket
    Websocket Authorize    ${token}
    ${message} =    Websocket Receive Json
    Should Be Equal    ${message}[type]    authorize-error
    Websocket Should Be Closed

Authorization can be renewed with a new token
    ${token} =    Get token of    admin
    Open WebSocket
    Websocket Authorize    ${token}
    ${token} =    Get token of    admin
    Websocket Authorize    ${token}
    Fake Agent Completes Task    cluster    ws-task-3    output="ok"
    ${message} =    Websocket Receive Json
    Should Be Equal    ${message}[name]    progress/cluster/task/ws-task-3
